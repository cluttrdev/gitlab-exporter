package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/cli"

	"github.com/cluttrdev/gitlab-exporter/pkg/worker"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
	"github.com/cluttrdev/gitlab-exporter/internal/exporter"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
	"github.com/cluttrdev/gitlab-exporter/internal/healthz"
	"github.com/cluttrdev/gitlab-exporter/internal/jobs"
)

type RunConfig struct {
	RootConfig

	projects projectList
	catchup  bool
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
			out:   out,
			flags: flag.NewFlagSet("run", flag.ExitOnError),
		},
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
	fs.BoolVar(&c.catchup, "catchup", false, "Whether to export historical data. (default: false)")

	_ = fs.String("log-level", "info", "The logging level, one of 'debug', 'info', 'warn', 'error'. (default: 'info')")
	_ = fs.String("log-format", "text", "The logging format, either 'text' or 'json'. (default: 'text')")
}

func (c *RunConfig) Exec(ctx context.Context, _ []string) error {
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

	// gather projects from config
	projects, err := resolveProjects(ctx, cfg, gitlabclient)
	if err != nil {
		return fmt.Errorf("error resolving projects: %w", err)
	}

	// add projects passed as arguments
	for _, pid := range c.projects {
		exists := slices.ContainsFunc(cfg.Projects, func(p config.Project) bool {
			return p.Id == pid
		})

		if !exists {
			projects = append(cfg.Projects, config.Project{
				ProjectSettings: config.DefaultProjectSettings(),
				Id:              pid,
			})
		}
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

	if len(projects) > 0 { // jobs
		ctx, cancel := context.WithCancel(context.Background())

		g.Add(func() error { // execute
			var wg sync.WaitGroup
			for _, p := range projects {
				if c.catchup && p.CatchUp.Enabled {
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

			wg.Wait()
			return nil
		}, func(err error) { // interrupt
			slog.Info("Cancelling jobs...")
			cancel()
			<-ctx.Done()
			slog.Info("Cancelling jobs... done")
		})
	} else {
		slog.Warn("There are no projects configured for export")
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

func serveHTTP(cfg config.HTTP, reg *prometheus.Registry) (func() error, func(error)) {
	m := http.NewServeMux()

	m.Handle("/healthz/", http.StripPrefix("/healthz", healthz.NewHandler()))

	m.Handle(
		"/metrics",
		promhttp.InstrumentMetricHandler(
			reg, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		),
	)

	if cfg.Debug {
		m.HandleFunc("/debug/pprof/", pprof.Index)
		m.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		m.HandleFunc("/debug/pprof/profile", pprof.Profile)
		m.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		m.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler: m,
	}

	execute := func() error {
		slog.Info("Starting http server", "addr", httpServer.Addr)
		return httpServer.ListenAndServe()
	}

	interrupt := func(error) {
		slog.Info("Stopping http server...")
		if err := httpServer.Shutdown(context.Background()); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				slog.Error("error shutting down http server", "error", err)
			}
		}
		slog.Info("Stopping http server... done")
	}

	return execute, interrupt
}

func resolveProjects(ctx context.Context, cfg config.Config, glab *gitlab.Client) ([]config.Project, error) {
	pm := make(map[int64]config.Project)

	opt := gitlab.ListNamespaceProjectsOptions{}
	for _, namespace := range cfg.Namespaces {
		opt.Kind = namespace.Kind
		opt.Visibility = (*gitlab.VisibilityValue)(&namespace.Visibility)
		opt.WithShared = namespace.WithShared
		opt.IncludeSubgroups = namespace.IncludeSubgroups

		ps, err := glab.ListNamespaceProjects(ctx, namespace.Id, opt)
		if err != nil {
			return nil, err
		}

		for _, p := range ps {
			pm[p.Id] = config.Project{
				ProjectSettings: namespace.ProjectSettings,
				Id:              p.Id,
			}
		}

		for _, pid := range namespace.ExcludeProjects {
			p, _, err := glab.Client().Projects.GetProject(pid, nil, _gitlab.WithContext(ctx))
			if err != nil {
				return nil, fmt.Errorf("error getting project %q: %w", pid, err)
			}
			delete(pm, int64(p.ID))
		}
	}

	// overwrite with explicitly configured projects
	for _, p := range cfg.Projects {
		pm[p.Id] = p
	}

	projects := make([]config.Project, 0, len(pm))
	for _, p := range pm {
		projects = append(projects, p)
	}

	return projects, nil
}
