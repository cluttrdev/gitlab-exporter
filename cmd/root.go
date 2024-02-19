package cmd

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/cluttrdev/cli"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
)

const (
	exeName      string = "gitlab-exporter"
	envVarPrefix string = "GLE"
)

type RootConfig struct {
	filename string

	out   io.Writer
	flags *flag.FlagSet
}

func NewRootCmd(out io.Writer) *cli.Command {
	cfg := RootConfig{
		filename: "",

		out:   out,
		flags: flag.NewFlagSet(exeName, flag.ContinueOnError),
	}

	cfg.RegisterFlags(cfg.flags)

	return &cli.Command{
		Name:       exeName,
		ShortUsage: fmt.Sprintf("%s <subcommand> [option]... [arg]...", exeName),
		Flags:      cfg.flags,
		Exec:       cfg.Exec,
	}
}

func (c *RootConfig) RegisterFlags(fs *flag.FlagSet) {
	defaults := config.Default()

	fs.String("gitlab-api-url", defaults.GitLab.Api.URL, fmt.Sprintf("The GitLab API URL (default: '%s').", defaults.GitLab.Api.URL))
	fs.String("gitlab-api-token", defaults.GitLab.Api.Token, fmt.Sprintf("The GitLab API Token (default: '%s').", defaults.GitLab.Api.Token))

	fs.StringVar(&c.filename, "config", "", "Configuration file to use.")
}

func (c *RootConfig) Exec(context.Context, []string) error {
	return flag.ErrHelp
}

func loadConfig(filename string, flags *flag.FlagSet, cfg *config.Config) error {
	if filename != "" {
		if err := config.LoadFile(filename, cfg); err != nil {
			return err
		}
	}

	flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "gitlab-api-url":
			cfg.GitLab.Api.URL = f.Value.String()
		case "gitlab-api-token":
			cfg.GitLab.Api.Token = f.Value.String()
		}
	})

	return nil
}

func writeConfig(out io.Writer, cfg config.Config) {
	fmt.Fprintln(out, "----")
	fmt.Fprintf(out, "GitLab URL: %s\n", cfg.GitLab.Api.URL)
	fmt.Fprintf(out, "GitLab Token: %x\n", sha256String(cfg.GitLab.Api.Token))
	fmt.Fprintln(out, "----")

	for _, s := range cfg.Endpoints {
		fmt.Fprintf(out, "Endpoint: %s\n", s.Address)
	}
	fmt.Fprintln(out, "----")

	projects := []int64{}
	for _, p := range cfg.Projects {
		projects = append(projects, p.Id)
	}
	fmt.Fprintf(out, "Projects: %v\n", projects)
	fmt.Fprintln(out, "----")
	fmt.Fprintf(out, "Log Level: %s\n", cfg.Log.Level)
	fmt.Fprintf(out, "Log Format: %s\n", cfg.Log.Format)
	fmt.Fprintln(out, "----")
}

func sha256String(s string) []byte {
	h := sha256.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

func initLogging(out io.Writer, cfg config.Log) {
	if out == nil {
		out = os.Stderr
	}

	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	switch cfg.Format {
	case "text":
		handler = slog.NewTextHandler(out, &opts)
	case "json":
		handler = slog.NewJSONHandler(out, &opts)
	default:
		handler = slog.NewTextHandler(out, &opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func setupDaemon(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-signalChan:
			slog.Debug("Got SIGINT/SIGTERM, exiting")
			signal.Stop(signalChan)
			cancel()
		case <-ctx.Done():
			slog.Debug("Done")
		}
	}()

	return ctx, cancel
}
