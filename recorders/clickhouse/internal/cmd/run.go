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
	"syscall"

	"github.com/cluttrdev/cli"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"go.cluttr.dev/gitlab-exporter/grpc/server"

	"go.cluttr.dev/gitlab-exporter-clickhouse-recorder/internal/clickhouse"
	"go.cluttr.dev/gitlab-exporter-clickhouse-recorder/internal/config"
	"go.cluttr.dev/gitlab-exporter-clickhouse-recorder/internal/recorder"
)

type RunConfig struct {
	RootConfig

	ServerHost string
	ServerPort string

	LogLevel  string
	LogFormat string

	flags *flag.FlagSet
}

func NewRunCmd(out io.Writer) *cli.Command {
	fs := flag.NewFlagSet("run", flag.ExitOnError)

	cfg := RunConfig{
		RootConfig: RootConfig{
			out: out,
		},

		flags: fs,
	}
	cfg.RegisterFlags(fs)

	return &cli.Command{
		Name:       "run",
		ShortUsage: fmt.Sprintf("%s run [option]...", exeName),
		ShortHelp:  "Run gRPC server",
		Flags:      fs,
		Exec:       cfg.Exec,
	}
}

func (c *RunConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootConfig.RegisterFlags(fs)

	fs.StringVar(&c.ServerHost, "server-host", "127.0.0.1", "The gRPC server listen host. (default: '127.0.0.1')")
	fs.StringVar(&c.ServerPort, "server-port", "0", "The gRPC server listen port. (default: '0', random)")

	fs.StringVar(&c.LogLevel, "log-level", "info", "The logging level, one of 'debug', 'info', 'warn', 'error'. (default: 'info')")
	fs.StringVar(&c.LogFormat, "log-format", "text", "The logging format, either 'text' or 'json'. (default: 'text')")
}

func (c *RunConfig) Exec(ctx context.Context, args []string) error {
	// load configuration
	var cfg config.Config
	config.SetDefaults(&cfg)
	if err := loadConfig(c.RootConfig.filename, c.flags, &cfg); err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}
	// override values passed as env vars or flags
	c.flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "server-host":
			cfg.Server.Host = f.Value.String()
		case "server-port":
			cfg.Server.Port = f.Value.String()
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

	// create clickhouse client
	opts := clickhouse.ClientOptions(clickhouse.ClientConfig{
		Host:     cfg.ClickHouse.Host,
		Port:     cfg.ClickHouse.Port,
		Database: cfg.ClickHouse.Database,
		User:     cfg.ClickHouse.User,
		Password: cfg.ClickHouse.Password,
	})
	conn, err := clickhouse.Connect(&opts)
	if err != nil {
		return fmt.Errorf("error creating clickhouse connection")
	}
	client := clickhouse.NewClient(conn, cfg.ClickHouse.Database)

	if err := c.checkSchemaVersion(ctx, client); err != nil {
		return fmt.Errorf("error checking database schema: %w", err)
	}

	// create recorder
	rec := recorder.New(client)

	// create grpc server
	grpcServer := server.New(rec)

	// setup run group
	g := &run.Group{}

	{ // serve grpc
		ctx, cancel := context.WithCancel(ctx)
		g.Add(func() error { // execute
			slog.Info("Starting grpc server")
			addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
			return grpcServer.ListenAndServe(ctx, addr)
		}, func(err error) { // interrupt
			slog.Info("Stopping grpc server...")
			cancel()
			<-ctx.Done()
			slog.Info("Stopping grpc server... done")
		})

		{ // monitor health
			g.Add(func() error { // execute
				grpcServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
				if err := rec.GetReady(ctx); err != nil {
					return err
				}
				grpcServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

				errChan := rec.WatchReadiness(ctx)
				var latestErr error
				for {
					err, ok := <-errChan
					if !ok {
						// channel closed, return last error
						return latestErr
					}

					if errors.Is(err, context.Canceled) {
						return err
					}

					if err != nil {
						latestErr = err
						grpcServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
					} else if latestErr != nil {
						slog.Info("Readiness check successful")
						grpcServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
						latestErr = nil
					}
				}
			}, func(err error) { // interrupt
				cancel()
			})
		}
	}

	if cfg.HTTP.Enabled { // serve http
		reg := prometheus.NewRegistry()
		reg.MustRegister(
			grpcServer.MetricsCollector(),
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		)

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

func (c *RunConfig) checkSchemaVersion(ctx context.Context, ch *clickhouse.Client) error {
	schemaVersion, dirty, err := clickhouse.GetSchemaVersion(ch, ctx)
	if err != nil {
		return fmt.Errorf("error getting schema version: %w", err)
	} else if dirty {
		return fmt.Errorf("database schema is dirty")
	}

	fsys := MigrationsFileSystem
	path := MigrationsPath
	migrationsVersion, err := clickhouse.GetLatestMigrationVersion(fsys, path)
	if err != nil {
		return fmt.Errorf("error getting migrations version: %w", err)
	}

	if schemaVersion != migrationsVersion {
		slog.Error("Database schema version does not match migrations", "schema", schemaVersion, "migrations", migrationsVersion)
		return fmt.Errorf("database schema version mismatch")
	}
	slog.Debug("Database schema version matches migrations", "schema", schemaVersion, "migrations", migrationsVersion)
	return nil
}

func serveHTTP(cfg config.HTTP, reg *prometheus.Registry) (func() error, func(error)) {
	m := http.NewServeMux()

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
