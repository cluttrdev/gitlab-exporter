package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"

	clickhouse "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/clickhouse"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/tasks"
)

type DeduplicateConfig struct {
	rootConfig *RootConfig
	out        io.Writer

	database    string
	table       string
	final       bool
	by          columnList
	except      columnList
	throwIfNoop bool

	flags *flag.FlagSet
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

func NewDeduplicateCmd(rootConfig *RootConfig) *ffcli.Command {
	fs := flag.NewFlagSet(fmt.Sprintf("%s deduplicate", exeName), flag.ContinueOnError)

	cfg := DeduplicateConfig{
		rootConfig: rootConfig,

		flags: fs,
	}

	cfg.RegisterFlags(fs)

	return &ffcli.Command{
		Name:       "deduplicate",
		ShortUsage: fmt.Sprintf("%s deduplicate [flags] table", exeName),
		ShortHelp:  "Deduplicate database table",
		UsageFunc:  usageFunc,
		FlagSet:    fs,
		Options:    rootCmdOptions,
		Exec:       cfg.Exec,
	}
}

func (c *DeduplicateConfig) RegisterFlags(fs *flag.FlagSet) {
	c.rootConfig.RegisterFlags(fs)

	fs.StringVar(&c.database, "database", "gitlab_ci", "The database name. (default: 'gitlab_ci')")
	fs.BoolVar(&c.final, "final", true, "Optimize even if all data is already in one part. (default: true)")
	fs.Var(&c.by, "by", "Comma separated list of columns to deduplicate by. (default: [])")
	fs.Var(&c.except, "except", "Comma separated list of columns to not deduplicate by. (default: [])")
	fs.BoolVar(&c.throwIfNoop, "throw-if-noop", true, "Notify if deduplication is not performed. (default: true)")
}

func (c *DeduplicateConfig) Exec(ctx context.Context, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("invalid number of positional arguments: %v", args)
	}

	table := args[0]

	cfg := config.Default()
	if err := loadConfig(c.rootConfig.filename, c.flags, &cfg); err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	ch, err := clickhouse.NewClickHouseClient(clickhouse.ClientConfig{
		Host:     cfg.ClickHouse.Host,
		Port:     cfg.ClickHouse.Port,
		Database: cfg.ClickHouse.Database,
		User:     cfg.ClickHouse.User,
		Password: cfg.ClickHouse.Password,
	})
	if err != nil {
		return fmt.Errorf("error creating clickhouse client: %w", err)
	}

	opt := tasks.DeduplicateTableOptions{
		Database:    c.database,
		Table:       table,
		Final:       &c.final,
		By:          c.by,
		Except:      c.except,
		ThrowIfNoop: &c.throwIfNoop,
	}

	return tasks.DeduplicateTable(ctx, opt, ch)
}
