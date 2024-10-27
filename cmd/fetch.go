package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/cluttrdev/cli"
)

type FetchConfig struct {
	RootConfig
}

func NewFetchCmd(out io.Writer) *cli.Command {
	cfg := FetchConfig{
		RootConfig: RootConfig{
			out:   out,
			flags: flag.NewFlagSet(fmt.Sprintf("%s fetch", exeName), flag.ContinueOnError),
		},
	}

	cfg.RegisterFlags(cfg.flags)

	var (
		fetchPipelineCmd   = NewFetchPipelineCmd(out)
		fetchJobLogCmd     = NewFetchJobLogCmd(out)
		fetchTestReportCmd = NewFetchTestReportCmd(out)

		fetchProjectUpdatesCmd = NewFetchProjectsCommand(out)
	)

	return &cli.Command{
		Name:       "fetch",
		ShortUsage: fmt.Sprintf("%s fetch <subcommand> [option]... [args]...", exeName),
		ShortHelp:  "Fetch data from the GitLab API",
		Flags:      cfg.flags,
		Subcommands: []*cli.Command{
			fetchPipelineCmd,
			fetchJobLogCmd,
			fetchTestReportCmd,
			fetchProjectUpdatesCmd,
		},
		Exec: cfg.Exec,
	}
}

func (c *FetchConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootConfig.RegisterFlags(fs)
}

func (c *FetchConfig) Exec(ctx context.Context, _ []string) error {
	return flag.ErrHelp
}
