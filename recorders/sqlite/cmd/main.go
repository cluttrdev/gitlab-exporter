package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v3"

	"go.cluttr.dev/gitlab-exporter/grpc/server"
	"go.cluttr.dev/gitlab-exporter/recorders/sqlite"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Fatal error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		address    string
		configPath string
	)

	flag.StringVar(&address, "address", "", "Address to listen on (e.g., unix:///tmp/recorder.sock or :9090)")
	flag.StringVar(&configPath, "config", "", "Path to configuration file")
	flag.Parse()

	if address == "" {
		return fmt.Errorf("--address is required")
	}

	// Create context that cancels on interrupt signals
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Create recorder instance
	rec := sqlite.New(address)

	// Load and apply configuration
	configBytes, err := loadConfig(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	var settings sqlite.Settings
	if err := yaml.Unmarshal(configBytes, &settings); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	if err := rec.Initialize(ctx, settings); err != nil {
		return fmt.Errorf("initialize recorder: %w", err)
	}

	if err := rec.Start(ctx); err != nil {
		return fmt.Errorf("start recorder: %w", err)
	}
	defer func() {
		if err := rec.Stop(context.Background()); err != nil {
			slog.Error("Error stopping recorder", "error", err)
		}
	}()

	if err := rec.CheckHealth(ctx); err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	// Create and start gRPC server
	srv := server.New(rec)

	slog.Info("Starting SQLite recorder", "address", address)

	return srv.ListenAndServe(ctx, address)
}

// loadConfig reads the config file.
func loadConfig(configPath string) ([]byte, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	return data, nil
}
