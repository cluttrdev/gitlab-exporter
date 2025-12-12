package cmd

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"

	"github.com/cluttrdev/cli"
	"go.cluttr.dev/gitlab-exporter-clickhouse-recorder/internal/config"
)

func NewRootCmd(out io.Writer) *cli.Command {
	fs := flag.NewFlagSet(exeName, flag.ExitOnError)

	cfg := RootConfig{
		out:   out,
		flags: fs,
	}
	cfg.RegisterFlags(fs)

	return &cli.Command{
		Name:  exeName,
		Flags: fs,
		Exec:  cfg.Exec,
	}
}

type RootConfig struct {
	filename string
	out      io.Writer
	flags    *flag.FlagSet
	debug    bool
}

func (c *RootConfig) RegisterFlags(fs *flag.FlagSet) {
	_ = fs.String("clickhouse-host", "127.0.0.1", "The ClickHouse server name (default: '127.0.0.1').")
	_ = fs.String("clickhouse-port", "9000", "The ClickHouse port to connect to (default: '9000')")
	_ = fs.String("clickhouse-database", "default", "Select the current default ClickHouse database (default: 'default').")
	_ = fs.String("clickhouse-user", "default", "The ClickHouse username to connect with (default: 'default').")
	_ = fs.String("clickhouse-password", "", "The ClickHouse password (default: '').")

	_ = fs.Int64("clickhouse-client-max-concurrent-queries", 0, "The maximum number of concurrent queries the client sends to clickhouse (default: 0, unlimited).")

	fs.StringVar(&c.filename, "config", "", "The configuration file to use.")

	fs.BoolVar(&c.debug, "debug", false, "Run in debug mode.")
}

func (c *RootConfig) Exec(ctx context.Context, args []string) error {
	return flag.ErrHelp
}

func loadConfig(filename string, flags *flag.FlagSet, cfg *config.Config) error {
	// load configuration file first
	if filename != "" {
		if err := config.LoadFile(filename, cfg); err != nil {
			return err
		}
	}

	// override with values passed as env vars or flags
	flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "clickhouse-host":
			cfg.ClickHouse.Host = f.Value.String()
		case "clickhouse-port":
			cfg.ClickHouse.Port = f.Value.String()
		case "clickhouse-database":
			cfg.ClickHouse.Database = f.Value.String()
		case "clickhouse-user":
			cfg.ClickHouse.User = f.Value.String()
		case "clickhouse-password":
			cfg.ClickHouse.Password = f.Value.String()

		case "clickhouse-client-max-concurrent-queries":
			n, err := strconv.ParseInt(f.Value.String(), 10, 64)
			if err != nil {
				n = -1
			}
			cfg.ClickHouse.Client.MaxConcurrentQueries = n
		}
	})

	if cfg.ClickHouse.Client.MaxConcurrentQueries < 0 {
		return fmt.Errorf("invalid config: max_concurrent_queries")
	}

	return nil
}

func writeConfig(out io.Writer, cfg config.Config) {
	_cfg := cfg
	_cfg.ClickHouse.Password = fmt.Sprintf("%x", sha256String(cfg.ClickHouse.Password))

	b, err := json.MarshalIndent(_cfg, "", "  ")
	if err != nil {
		fmt.Fprintf(out, "error marshalling config: %v\n", err)
	}
	fmt.Fprint(out, string(b))
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
