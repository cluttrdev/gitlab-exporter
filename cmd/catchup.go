package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cluttrdev/cli"

	"github.com/cluttrdev/gitlab-exporter/pkg/worker"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
	"github.com/cluttrdev/gitlab-exporter/internal/exporter"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
	"github.com/cluttrdev/gitlab-exporter/internal/jobs"
)

type CatchUpConfig struct {
	RootConfig

	flags *flag.FlagSet
}

func NewCatchUpCmd(out io.Writer) *cli.Command {
	cfg := CatchUpConfig{
		RootConfig: RootConfig{
			out: out,
		},

		flags: flag.NewFlagSet("catchup", flag.ExitOnError),
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
	// setup daemon
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

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

	// configure workers
	pool := worker.NewWorkerPool(42)
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

	go func() {
		// cancel context when work is done to stop worker pool
		wg.Wait()
		cancel()
	}()

	slog.Info("Starting workers")
	pool.Start(ctx)

	<-ctx.Done()
	return nil
}
