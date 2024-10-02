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

	"github.com/alitto/pond"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/cli"

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
	slog.Info("Resolving projects to export...")
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
	slog.Info("Resolving projects to export... done")

	g := &run.Group{}

	var pool *pond.WorkerPool
	if len(projects) > 0 { // jobs
		slog.Info(fmt.Sprintf("Found %d projects to export", len(projects)))
		ctxJobs, cancelJobs := context.WithCancel(context.Background())

		slog.Info("Starting worker pool")
		pool = pond.New(42, 1024, pond.Context(ctxJobs))

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
						job.Run(ctxJobs)
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
					job.Run(ctxJobs)
				}()
			}

			wg.Wait()
			return nil
		}, func(err error) { // interrupt
			slog.Info("Cancelling jobs...")
			cancelJobs()
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

		if pool != nil {
			registerWorkerPoolMetrics(reg, pool)
		}

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
	projectConfigs := make(map[int64]config.Project)

	opt := gitlab.ListNamespaceProjectsOptions{}
	for _, namespace := range cfg.Namespaces {
		opt.Kind = namespace.Kind
		if namespace.Visibility != "" {
			opt.Visibility = (*_gitlab.VisibilityValue)(&namespace.Visibility)
		}

		opt.WithShared = _gitlab.Ptr(namespace.WithShared)
		opt.IncludeSubgroups = _gitlab.Ptr(namespace.IncludeSubgroups)

		err := glab.ListNamespaceProjects(ctx, namespace.Id, opt, func(projects []*_gitlab.Project) bool {
			for _, project := range projects {
				projectID := int64(project.ID)
				projectConfigs[projectID] = config.Project{
					ProjectSettings: namespace.ProjectSettings,
					Id:              projectID,
				}
			}

			for _, pid := range namespace.ExcludeProjects {
				p, _, err := glab.Client().Projects.GetProject(pid, nil, _gitlab.WithContext(ctx))
				if err != nil {
					slog.Error("error getting namespace project", "namespace_id", namespace.Id, "project", pid, "error", err)
					return false
				}
				delete(projectConfigs, int64(p.ID))
			}

			return true
		})
		if err != nil {
			return nil, err
		}
	}

	// overwrite with explicitly configured projects
	for _, p := range cfg.Projects {
		projectConfigs[p.Id] = p
	}

	projects := make([]config.Project, 0, len(projectConfigs))
	for _, p := range projectConfigs {
		projects = append(projects, p)
	}

	return projects, nil
}

func registerWorkerPoolMetrics(reg *prometheus.Registry, pool *pond.WorkerPool) {
	// Worker pool metrics
	reg.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "pool_workers_running",
			Help: "Number of running worker goroutines",
		},
		func() float64 {
			return float64(pool.RunningWorkers())
		}))
	reg.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "pool_workers_idle",
			Help: "Number of idle worker goroutines",
		},
		func() float64 {
			return float64(pool.IdleWorkers())
		}))

	// Task metrics
	reg.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "pool_tasks_submitted_total",
			Help: "Number of tasks submitted",
		},
		func() float64 {
			return float64(pool.SubmittedTasks())
		}))
	reg.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "pool_tasks_waiting_total",
			Help: "Number of tasks waiting in the queue",
		},
		func() float64 {
			return float64(pool.WaitingTasks())
		}))
	reg.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "pool_tasks_successful_total",
			Help: "Number of tasks that completed successfully",
		},
		func() float64 {
			return float64(pool.SuccessfulTasks())
		}))
	reg.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "pool_tasks_failed_total",
			Help: "Number of tasks that completed with panic",
		},
		func() float64 {
			return float64(pool.FailedTasks())
		}))
	reg.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "pool_tasks_completed_total",
			Help: "Number of tasks that completed either successfully or with panic",
		},
		func() float64 {
			return float64(pool.CompletedTasks())
		}))
}
