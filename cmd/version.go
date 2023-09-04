package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/peterbourgon/ff/v3/ffcli"
)

type VersionConfig struct {
	out io.Writer

	version string
}

func NewVersionCmd(out io.Writer, version string) *ffcli.Command {
	config := VersionConfig{
		out:     out,
		version: version,
	}

	fs := flag.NewFlagSet(fmt.Sprintf("%s version", exeName), flag.ContinueOnError)

	return &ffcli.Command{
		Name:       "version",
		ShortUsage: fmt.Sprintf("%s version", exeName),
		ShortHelp:  "Display version information",
		UsageFunc:  usageFunc,
		FlagSet:    fs,
		Exec:       config.Exec,
	}
}

func (c *VersionConfig) Exec(ctx context.Context, _ []string) error {
	fmt.Fprintf(c.out, "%s\n", c.version)
	return nil
}
