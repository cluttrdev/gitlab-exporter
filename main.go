package main

import (
	"context"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/cmd"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/controller"
)

func main() {
	var (
		out                 = os.Stdout
		rootCmd, rootConfig = cmd.NewRootCmd()
		runCmd              = cmd.NewRunCmd(rootConfig, out)
		fetchCmd            = cmd.NewFetchCmd(rootConfig, out)
		exportCmd           = cmd.NewExportCmd(rootConfig, out)
	)

	rootCmd.Subcommands = []*ffcli.Command{
		runCmd,
		fetchCmd,
		exportCmd,
	}

	if err := rootCmd.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing args: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args[1:]) == 0 {
		fmt.Fprintf(out, rootCmd.UsageFunc(rootCmd))
		os.Exit(0)
	}

	ctl, err := controller.NewController(rootConfig.Config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error constructing controller: %v", err)
		os.Exit(1)
	}

	rootConfig.Controller = &ctl

	ctx := context.Background()
	if err := rootCmd.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
