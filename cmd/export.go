package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/peterbourgon/ff/v3/ffcli"
)

type ExportConfig struct {
	rootConfig *RootConfig

	out io.Writer

	flags *flag.FlagSet
}

func NewExportCmd(rootConfig *RootConfig, out io.Writer) *ffcli.Command {
	fs := flag.NewFlagSet(fmt.Sprintf("%s export", exeName), flag.ContinueOnError)

	cfg := ExportConfig{
		rootConfig: rootConfig,
		out:        out,

		flags: fs,
	}

	cfg.RegisterFlags(fs)

	var (
		exportPipelineCmd = NewExportPipelineCmd(&cfg)
	)

	return &ffcli.Command{
		Name:       "export",
		ShortUsage: fmt.Sprintf("%s export <subcommand> [flags] [<args>...]", exeName),
		ShortHelp:  "Export data from the GitLab API to ClickHouse",
		UsageFunc:  usageFunc,
		FlagSet:    fs,
		Subcommands: []*ffcli.Command{
			exportPipelineCmd,
		},
		Exec: cfg.Exec,
	}
}

func (c *ExportConfig) RegisterFlags(fs *flag.FlagSet) {
	c.rootConfig.RegisterFlags(fs)
}

func (c *ExportConfig) Exec(ctx context.Context, _ []string) error {
	return flag.ErrHelp
}
