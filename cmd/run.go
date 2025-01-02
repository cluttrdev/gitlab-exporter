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
	"syscall"
	"time"

	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/cluttrdev/cli"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
	"github.com/cluttrdev/gitlab-exporter/internal/exporter"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
	"github.com/cluttrdev/gitlab-exporter/internal/healthz"
	"github.com/cluttrdev/gitlab-exporter/internal/tasks"
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
	var err error

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

	// add projects passed as arguments
	for _, pid := range c.projects {
		exists := slices.ContainsFunc(cfg.Projects, func(p config.Project) bool {
			return p.Id == pid
		})
		if exists {
			continue
		}

		cfg.Projects = append(cfg.Projects, config.Project{
			Id:              pid,
			ProjectSettings: config.DefaultProjectSettings(),
		})
	}

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

			ExportInterval:  5 * time.Minute,
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

			if c.catchup {
				go func() {
					if err := ctrl.CatchUp(ctx); err != nil {
						slog.Error("error catching up", "error", err)
					}
				}()
			}

			return ctrl.Run(ctx)
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

func createGitLabClient(cfg config.Config) (*gitlab.Client, error) {
	var (
		oauthConfig *gitlab.OAuthConfig
		err         error
	)

	oauthConfig, err = configureOAuth(cfg.GitLab.OAuth)
	if err != nil {
		return nil, fmt.Errorf("oauth config: %w", err)
	}

	glab, err := gitlab.NewGitLabClient(gitlab.ClientConfig{
		URL:   cfg.GitLab.Url,
		Token: cfg.GitLab.Token,

		OAuth: oauthConfig,

		RateLimit: cfg.GitLab.Client.Rate.Limit,
	})
	if err != nil {
		return nil, err
	}

	if config.IsOAuthRequired(cfg) {
		if err := glab.HTTP.ChechAuth(); err != nil {
			return nil, err
		}
	}

	return glab, err
}

func configureOAuth(cfg config.GitLabOAuth) (*gitlab.OAuthConfig, error) {
	var err error

	secrets := cfg.GitLabOAuthSecrets
	if cfg.SecretsFile != "" {
		secrets, err = config.LoadOAuthSecretsFile(cfg.SecretsFile)
		if err != nil {
			return nil, fmt.Errorf("load oauth secrets: %w", err)
		}
	}

	return &gitlab.OAuthConfig{
		GitLabOAuthSecrets: secrets,
		FlowType:           cfg.FlowType,
		Scopes:             []string{"openid"},
	}, nil
}
