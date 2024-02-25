package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cluttrdev/cli"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/cluttrdev/gitlab-exporter/pkg/worker"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
	"github.com/cluttrdev/gitlab-exporter/internal/exporter"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
	"github.com/cluttrdev/gitlab-exporter/internal/jobs"
)

type CatchUpConfig struct {
	RootConfig
}

func NewCatchUpCmd(out io.Writer) *cli.Command {
	cfg := CatchUpConfig{
		RootConfig: RootConfig{
			out:   out,
			flags: flag.NewFlagSet("catchup", flag.ExitOnError),
		},
	}

	cfg.RegisterFlags(cfg.flags)

	cmd := &cli.Command{
		Name:       "catchup",
		ShortUsage: fmt.Sprintf("%s catchup [option]...", exeName),
		ShortHelp:  "Catch up on project history",
		Flags:      cfg.flags,
		Exec:       cfg.Exec,
	}

	return cmd
}

func (c *CatchUpConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootConfig.RegisterFlags(fs)

	_ = fs.String("log-level", "info", "The logging level, one of 'debug', 'info', 'warn', 'error'. (default: 'info')")
	_ = fs.String("log-format", "text", "The logging format, either 'text' or 'json'. (default: 'text')")
}

func (c *CatchUpConfig) Exec(ctx context.Context, args []string) error {
	// load configuration
	cfg := config.Default()
	if err := loadConfig(c.RootConfig.filename, c.flags, &cfg); err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	// override values passed as env vars or flags
	c.flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "log-level":
			cfg.Log.Level = f.Value.String()
		case "log-format":
			cfg.Log.Format = f.Value.String()
		}
	})

	if c.debug {
		cfg.HTTP.Enabled = true
		cfg.HTTP.Debug = true
		cfg.Log.Level = "debug"
	}

	if cfg.Log.Level == "debug" {
		writeConfig(c.out, cfg)
	}
	initLogging(c.out, cfg.Log)

	// create gitlab client
	gitlabclient, err := gitlab.NewGitLabClient(gitlab.ClientConfig{
		URL:   cfg.GitLab.Api.URL,
		Token: cfg.GitLab.Api.Token,

		RateLimit: cfg.GitLab.Client.Rate.Limit,
	})
	if err != nil {
		return fmt.Errorf("error creating gitlab client: %w", err)
	}

	// create exporter
	endpoints := exporter.CreateEndpointConfigs(cfg.Endpoints)
	exp, err := exporter.New(endpoints)
	if err != nil {
		return err
	}

	g := &run.Group{}

	pool := worker.NewWorkerPool(42)
	{ // worker pool
		ctx, cancel := context.WithCancel(context.Background())

		g.Add(func() error { // execute
			slog.Info("Starting worker pool")
			pool.Start(ctx)
			<-ctx.Done()
			return ctx.Err()
		}, func(err error) { // interrupt
			defer cancel()
			slog.Info("Stopping worker pool...")
			pool.Stop()
			slog.Info("Stopping worker pool... done")
		})
	}

	{ // jobs
		ctx, cancel := context.WithCancel(context.Background())

		g.Add(func() error { // execute
			var wg sync.WaitGroup
			for _, p := range cfg.Projects {
				if !p.CatchUp.Enabled {
					continue
				}

				job := jobs.ProjectCatchUpJob{
					Config:   p,
					GitLab:   gitlabclient,
					Exporter: exp,

					WorkerPool: pool,
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					job.Run(ctx)
				}()
			}

			wg.Wait()
			return nil
		}, func(err error) { // interrupt
			slog.Info("Cancelling jobs...")
			cancel()
			<-ctx.Done()
			slog.Info("Cancelling jobs... done")
		})
	}

	if cfg.HTTP.Enabled {
		colls := []prometheus.Collector{
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		}
		for _, endpoint := range cfg.Endpoints {
			if mc := exp.MetricsCollectorFor(endpoint.Address); mc != nil {
				colls = append(colls, mc)
			}
		}
		reg := prometheus.NewRegistry()
		reg.MustRegister(colls...)

		g.Add(serveHTTP(cfg.HTTP, reg))
	}

	{ // signal handler
		ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
		g.Add(func() error { // execute
			<-ctx.Done()
            err := ctx.Err()
            if !errors.Is(err, context.Canceled) {
                slog.Info("Got SIGINT/SIGTERM, exiting")
            }
			return err
		}, func(err error) { // interrupt
			cancel()
		})
	}

	return g.Run()
}
