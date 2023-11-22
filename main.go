package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/cluttrdev/gitlab-exporter/cmd"
)

var (
	version = "(devel)"
)

func main() {
	var (
		out                 = os.Stdout
		rootCmd, rootConfig = cmd.NewRootCmd()
		versionCmd          = cmd.NewVersionCmd(out, version)
		runCmd              = cmd.NewRunCmd(rootConfig, out)
		fetchCmd            = cmd.NewFetchCmd(rootConfig, out)
		exportCmd           = cmd.NewExportCmd(rootConfig, out)
		deduplicateCmd      = cmd.NewDeduplicateCmd(rootConfig)
	)

	rootCmd.Subcommands = []*ffcli.Command{
		versionCmd,
		runCmd,
		fetchCmd,
		exportCmd,
		deduplicateCmd,
	}

	if len(os.Args[1:]) == 0 {
		fmt.Fprintln(out, rootCmd.UsageFunc(rootCmd))
		os.Exit(0)
	}

	if err := rootCmd.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing args: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	if err := rootCmd.Run(ctx); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		os.Exit(1)
	}
}
