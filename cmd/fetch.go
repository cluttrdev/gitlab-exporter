package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/peterbourgon/ff/v3/ffcli"
)

type FetchConfig struct {
	rootConfig *RootConfig

	out io.Writer

	flags *flag.FlagSet
}

func NewFetchCmd(rootConfig *RootConfig, out io.Writer) *ffcli.Command {
	fs := flag.NewFlagSet(fmt.Sprintf("%s fetch", exeName), flag.ContinueOnError)

	cfg := FetchConfig{
		rootConfig: rootConfig,
		out:        out,

		flags: fs,
	}

	cfg.RegisterFlags(fs)

	var (
		fetchPipelineCmd   = NewFetchPipelineCmd(&cfg)
		fetchTestReportCmd = NewFetchTestReportCmd(&cfg)
	)

	return &ffcli.Command{
		Name:       "fetch",
		ShortUsage: fmt.Sprintf("%s fetch <subcommand> [flags] [<args>...]", exeName),
		ShortHelp:  "Fetch data from the GitLab API",
		UsageFunc:  usageFunc,
		FlagSet:    fs,
		Subcommands: []*ffcli.Command{
			fetchPipelineCmd,
			fetchTestReportCmd,
		},
		Exec: cfg.Exec,
	}
}

func (c *FetchConfig) RegisterFlags(fs *flag.FlagSet) {
	c.rootConfig.RegisterFlags(fs)
}

func (c *FetchConfig) Exec(ctx context.Context, _ []string) error {
	return flag.ErrHelp
}
