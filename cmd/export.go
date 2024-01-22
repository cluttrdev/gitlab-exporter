package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/cluttrdev/cli"
)

type ExportConfig struct {
	RootConfig
}

func NewExportCmd(out io.Writer) *cli.Command {
	fs := flag.NewFlagSet(fmt.Sprintf("%s export", exeName), flag.ContinueOnError)

	cfg := ExportConfig{
		RootConfig: RootConfig{
			out: out,
		},
	}

	cfg.RegisterFlags(fs)

	var (
		exportPipelineCmd = NewExportPipelineCmd(out)
	)

	return &cli.Command{
		Name:       "export",
		ShortUsage: fmt.Sprintf("%s export <subcommand> [option]... [args]...", exeName),
		ShortHelp:  "Export data from the GitLab API",
		Flags:      fs,
		Subcommands: []*cli.Command{
			exportPipelineCmd,
		},
		Exec: cfg.Exec,
	}
}

func (c *ExportConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootConfig.RegisterFlags(fs)
}

func (c *ExportConfig) Exec(ctx context.Context, _ []string) error {
	return flag.ErrHelp
}
