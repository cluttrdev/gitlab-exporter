package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"golang.org/x/exp/slices"

	"github.com/cluttrdev/cli"

	"github.com/cluttrdev/gitlab-exporter/pkg/healthz"
	"github.com/cluttrdev/gitlab-exporter/pkg/worker"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
	"github.com/cluttrdev/gitlab-exporter/internal/exporter"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
	"github.com/cluttrdev/gitlab-exporter/internal/jobs"
	"github.com/cluttrdev/gitlab-exporter/internal/server"
)

type RunConfig struct {
	RootConfig

	projects projectList

	flags *flag.FlagSet
}

type projectList []int64

func (f *projectList) String() string {
	return fmt.Sprintf("%v", []int64(*f))
}

func (f *projectList) Set(value string) error {
	values := strings.Split(value, ",")
	for _, s := range values {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*f = append(*f, v)
	}
	return nil
}

func NewRunCmd(out io.Writer) *cli.Command {
	cfg := RunConfig{
		RootConfig: RootConfig{
			out: out,
		},

		flags: flag.NewFlagSet("run", flag.ExitOnError),
	}

	cfg.RegisterFlags(cfg.flags)

	return &cli.Command{
		Name:       "run",
		ShortUsage: fmt.Sprintf("%s run [option]...", exeName),
		ShortHelp:  "Run in daemon mode",
		Flags:      cfg.flags,
		Exec:       cfg.Exec,
	}
}

func (c *RunConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootConfig.RegisterFlags(fs)

	fs.Var(&c.projects, "projects", "Comma separated list of project ids.")

	_ = fs.String("log-level", "info", "The logging level, one of 'debug', 'info', 'warn', 'error'. (default: 'info')")
	_ = fs.String("log-format", "text", "The logging format, either 'text' or 'json'. (default: 'text')")
}

func (c *RunConfig) Exec(ctx context.Context, _ []string) error {
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

	// add projects passed to run command
	for _, pid := range c.projects {
		exists := slices.ContainsFunc(cfg.Projects, func(p config.Project) bool {
			return p.Id == pid
		})

		if !exists {
			cfg.Projects = append(cfg.Projects, config.Project{
				ProjectSettings: config.DefaultProjectSettings(),
				Id:              pid,
			})
		}
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

	// configure workers
	pool := worker.NewWorkerPool(42)
	var wg sync.WaitGroup
	for _, p := range cfg.Projects {
		if p.CatchUp.Enabled {
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

		job := jobs.ProjectExportJob{
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

	go startServer(ctx, cfg.Server, func() error {
		return gitlabclient.CheckReadiness(ctx)
	})

	slog.Info("Starting workers")
	pool.Start(ctx)

	<-ctx.Done()
	return nil
}

func startServer(ctx context.Context, cfg config.Server, ready healthz.Check) {
	srv := server.New(server.ServerConfig{
		Host:  cfg.Host,
		Port:  cfg.Port,
		Debug: false,

		ReadinessCheck: ready,
	})

	if err := srv.Serve(ctx); err != nil {
		slog.Error("error during server shutdown", "error", err)
	}
}
