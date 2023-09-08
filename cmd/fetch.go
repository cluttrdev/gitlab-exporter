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
}

func NewFetchCmd(rootConfig *RootConfig, out io.Writer) *ffcli.Command {
	config := FetchConfig{
		rootConfig: rootConfig,
		out:        out,
	}

	fs := flag.NewFlagSet(fmt.Sprintf("%s fetch", exeName), flag.ContinueOnError)
	config.rootConfig.RegisterFlags(fs)

	var (
		fetchPipelineCmd   = NewFetchPipelineCmd(&config)
		fetchTestReportCmd = NewFetchTestReportCmd(&config)
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
		Exec: config.Exec,
	}
}

func (c *FetchConfig) Exec(ctx context.Context, _ []string) error {
	return flag.ErrHelp
}
