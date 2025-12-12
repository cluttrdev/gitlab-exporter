package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log/slog"

	"github.com/cluttrdev/cli"
	"go.cluttr.dev/gitlab-exporter-clickhouse-recorder/internal/clickhouse"
	"go.cluttr.dev/gitlab-exporter-clickhouse-recorder/internal/config"
)

var (
	MigrationsFileSystem fs.FS
	MigrationsPath       string
)

type MigrateConfig struct {
	RootConfig

	flags *flag.FlagSet
}

func NewMigrateCommand(out io.Writer) *cli.Command {
	flags := flag.NewFlagSet("migrate", flag.ExitOnError)

	cfg := MigrateConfig{
		RootConfig: RootConfig{
			out: out,
		},

		flags: flags,
	}
	cfg.RegisterFlags(flags)

	return &cli.Command{
		Name:       "migrate",
		ShortUsage: fmt.Sprintf("%s migrate [option]...", exeName),
		ShortHelp:  "Migrate database schema",
		Flags:      flags,
		Exec:       cfg.Exec,
	}
}

func (c *MigrateConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootConfig.RegisterFlags(fs)
}

func (c *MigrateConfig) Exec(ctx context.Context, args []string) error {
	// load configuration
	var cfg config.Config
	config.SetDefaults(&cfg)
	if err := loadConfig(c.RootConfig.filename, c.flags, &cfg); err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	if c.debug {
		cfg.Log.Level = "debug"
	}
	initLogging(c.out, cfg.Log)

	opts := clickhouse.MigrationOptions{
		ClientConfig: clickhouse.ClientConfig{
			Host:     cfg.ClickHouse.Host,
			Port:     cfg.ClickHouse.Port,
			Database: cfg.ClickHouse.Database,
			User:     cfg.ClickHouse.User,
			Password: cfg.ClickHouse.Password,
		},

		FileSystem: MigrationsFileSystem,
		Path:       MigrationsPath,
	}
	if err := clickhouse.MigrateUp(opts); err != nil {
		if errors.Is(err, clickhouse.ErrMigrateNoChange) {
			slog.Info("no schema changes")
			return nil
		}
		return fmt.Errorf("error migrating database schema: %w", err)
	}
	return nil
}
