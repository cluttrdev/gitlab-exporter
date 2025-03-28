package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/cluttrdev/cli"

	"go.cluttr.dev/gitlab-exporter/internal/config"
	"go.cluttr.dev/gitlab-exporter/internal/exporter"
	"go.cluttr.dev/gitlab-exporter/internal/tasks"
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
	glab, err := createGitLabClient(cfg)
	if err != nil {
		return fmt.Errorf("create gitlab client: %w", err)
	}

	// create exporter
	endpoints := exporter.CreateEndpointConfigs(cfg.Endpoints)
	exp, err := exporter.New(endpoints)
	if err != nil {
		return err
	}

	g := &run.Group{}

	{ // controller
		ctrl := tasks.NewController(glab, exp, tasks.ControllerConfig{
			GitLab:     cfg.GitLab,
			Projects:   cfg.Projects,
			Namespaces: cfg.Namespaces,

			CatchUpInterval: 24 * time.Hour,
		})

		ctx, cancel := context.WithCancel(context.Background())

		g.Add(func() error { // execute
			slog.Info("Resolving projects...")
			count, err := ctrl.ResolveProjects(ctx)
			if err != nil {
				return fmt.Errorf("resolve projecst: %w", err)
			} else if count == 0 {
				return fmt.Errorf("No projects found")
			}
			slog.Info("Resolving projects... done", "found", count)

			return ctrl.CatchUp(ctx)
		}, func(err error) { // interrupt
			slog.Info("Stopping controller...")
			cancel()
			slog.Info("Stopping controller... done")
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
