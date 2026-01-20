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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cluttrdev/cli"

	"go.cluttr.dev/gitlab-exporter/exporter/internal/config"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/exporter"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/healthz"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/subprocess"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/tasks"
	grpc_client "go.cluttr.dev/gitlab-exporter/grpc/client"
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

	// initialize grpc clients
	clients, launchers, err := initGrpcClients(cfg)
	if err != nil {
		return fmt.Errorf("initialize grpc clients: %w", err)
	}

	// setup exporter
	exp := exporter.New()
	for _, client := range clients {
		if err := exp.AddClient(client); err != nil {
			return fmt.Errorf("add grpc client: %w", err)
		}
	}

	g := &run.Group{}

	{ // controller
		ctrl := tasks.NewController(glab, exp, tasks.ControllerConfig{
			GitLab:     cfg.GitLab,
			Projects:   cfg.Projects,
			Namespaces: cfg.Namespaces,

			Export: cfg.Export,

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

	{ // recorders
		ctx, cancel := context.WithCancel(context.Background())

		g.Add(func() error { //execute
			slog.Info("Starting recorder subprocesses...")
			for _, launcher := range launchers {
				if err := launcher.Start(ctx); err != nil {
					return fmt.Errorf("start recorder subprocess: %w", err)
				}
			}
			slog.Info("Starting recorder subprocesses... done")
			<-ctx.Done()
			return ctx.Err()
		}, func(err error) { // interrupt
			slog.Info("Stopping recorder subprocesses...")
			cancel()
			for _, launcher := range launchers {
				if err := launcher.Stop(ctx); err != nil {
					slog.Error("error stopping subprocess recorder", "error", err)
				}
			}
			slog.Info("Stopping recorder subprocesses... done")
		})
	}

	if cfg.HTTP.Enabled {
		colls := []prometheus.Collector{
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		}
		for _, client := range clients {
			colls = append(colls, client.MetricsCollector())
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

func initGrpcClients(cfg config.Config) ([]*grpc_client.Client, []*subprocess.Launcher, error) {
	var clients []*grpc_client.Client
	var launchers []*subprocess.Launcher

	// for backwards compatibility with deprecated endpoints config
	var recorderConfigs []config.Recorder
	for _, endpoint := range cfg.Endpoints {
		recorderConfigs = append(recorderConfigs, config.Recorder{
			Address: endpoint.Address,
			Mode:    config.RecorderModeExternal,
			Enabled: true,
		})
	}

	recorderConfigs = append(recorderConfigs, cfg.Recorders...)

	for _, rec := range recorderConfigs {
		if !rec.Enabled {
			continue
		}

		switch rec.Mode {
		case config.RecorderModeExternal:
			if rec.Address == "" {
				return nil, nil, fmt.Errorf("external recorder %s: address is required", rec.Type)
			}

			client, err := grpc_client.NewCLient(rec.Address,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				return nil, nil, fmt.Errorf("connect to external recorder %s at %s: %w", rec.Type, rec.Address, err)
			}

			clients = append(clients, client)
		case config.RecorderModeSubprocess:
			// Extract settings for launcher configuration
			var command string
			var maxRestarts int = subprocess.DefaultMaxRestarts

			if cmd, ok := rec.Settings["command"].(string); ok {
				command = cmd
			}
			if mr, ok := rec.Settings["max_restarts"].(int); ok {
				maxRestarts = mr
			}

			launcher, err := subprocess.NewLauncher(subprocess.LauncherConfig{
				RecorderType: rec.Type,
				Command:      command,
				SocketPath:   rec.Address,
				MaxRestarts:  maxRestarts,
			})
			if err != nil {
				return nil, nil, fmt.Errorf("create launcher for %s: %w", rec.Type, err)
			}

			client, err := grpc_client.NewCLient("unix://"+launcher.SocketPath(),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				return nil, nil, fmt.Errorf("create client: %w", err)
			}

			clients = append(clients, client)
			launchers = append(launchers, launcher)

		default:
			return nil, nil, fmt.Errorf("recorder %s: invalid mode %q", rec.Type, rec.Mode)
		}
	}

	return clients, launchers, nil
}

func (c *RunConfig) startRecorders(ctx context.Context, launchers []*subprocess.Launcher) error {
	return nil
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
	clientConfig := gitlab.ClientConfig{
		URL:   cfg.GitLab.Url,
		Token: cfg.GitLab.Token,

		RateLimit: cfg.GitLab.Client.Rate.Limit,
	}

	if cfg.GitLab.Username != "" && cfg.GitLab.Password != "" {
		clientConfig.Auth = gitlab.AuthConfig{
			AuthType: gitlab.SessionAuth,
			Basic: gitlab.BasicAuthConfig{
				Username: cfg.GitLab.Username,
				Password: cfg.GitLab.Password,
			},
		}
	}

	glab, err := gitlab.NewGitLabClient(clientConfig)
	if err != nil {
		return nil, err
	}

	if config.IsAuthedHTTPRequired(cfg) {
		if err := glab.HTTP.CheckAuthed(); err != nil {
			return nil, fmt.Errorf("check http auth: %w", err)
		}
	}

	return glab, err
}
