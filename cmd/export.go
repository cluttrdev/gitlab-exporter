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
}

func NewExportCmd(rootConfig *RootConfig, out io.Writer) *ffcli.Command {
	config := ExportConfig{
		rootConfig: rootConfig,
		out:        out,
	}

	fs := flag.NewFlagSet(fmt.Sprintf("%s export", exeName), flag.ContinueOnError)
	config.rootConfig.RegisterFlags(fs)

	var (
		exportPipelineCmd = NewExportPipelineCmd(&config)
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
		Exec: config.Exec,
	}
}

func (c *ExportConfig) Exec(ctx context.Context, _ []string) error {
	return flag.ErrHelp
}
