package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/cluttrdev/cli"

	"go.cluttr.dev/gitlab-exporter-clickhouse-recorder/internal/clickhouse"
	"go.cluttr.dev/gitlab-exporter-clickhouse-recorder/internal/config"
)

type DeduplicateConfig struct {
	RootConfig

	final       bool
	by          columnList
	except      columnList
	throwIfNoop bool

	flags *flag.FlagSet
}

func NewDeduplicateCmd(out io.Writer) *cli.Command {
	fs := flag.NewFlagSet(fmt.Sprintf("%s deduplicate", exeName), flag.ContinueOnError)

	cfg := DeduplicateConfig{
		RootConfig: RootConfig{
			out: out,
		},
		flags: fs,
	}
	cfg.RegisterFlags(fs)

	return &cli.Command{
		Name:       "deduplicate",
		ShortUsage: fmt.Sprintf("%s deduplicate [option]... table", exeName),
		ShortHelp:  "Deduplicate database table",
		Flags:      fs,
		Exec:       cfg.Exec,
	}
}

func (c *DeduplicateConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootConfig.RegisterFlags(fs)

	fs.BoolVar(&c.final, "final", true, "Optimize even if all data is already in one part. (default: true)")
	fs.Var(&c.by, "by", "Comma separated list of columns to deduplicate by. (default: [])")
	fs.Var(&c.except, "except", "Comma separated list of columns to not deduplicate by. (default: [])")
	fs.BoolVar(&c.throwIfNoop, "throw-if-noop", true, "Notify if deduplication is not performed. (default: true)")
}

func (c *DeduplicateConfig) Exec(ctx context.Context, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("invalid number of positional arguments: %v", args)
	}

	// load configuration
	var cfg config.Config
	config.SetDefaults(&cfg)
	if err := loadConfig(c.RootConfig.filename, c.flags, &cfg); err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

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

	table := args[0]

	opt := clickhouse.DeduplicateTableOptions{
		Database:    cfg.ClickHouse.Database,
		Table:       table,
		Final:       &c.final,
		By:          c.by,
		Except:      c.except,
		ThrowIfNoop: &c.throwIfNoop,
	}

	return clickhouse.DeduplicateTable(ctx, opt, client)
}

type columnList []string

func (f *columnList) String() string {
	return fmt.Sprintf("%v", []string(*f))
}

func (f *columnList) Set(value string) error {
	values := strings.Split(value, ",")
	for _, v := range values {
		*f = append(*f, v)
	}
	return nil
}
