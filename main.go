package main

import (
	"context"
	"log"
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
		log.Fatalf("error parsing args: %v\n", err)
	}

	ctl, err := controller.NewController(rootConfig.Config)
	if err != nil {
		log.Fatalf("error constructing controller: %v", err)
	}

	rootConfig.Controller = &ctl

	ctx := context.Background()
	if err := rootCmd.Run(ctx); err != nil {
		log.Fatalf("%v", err)
	}
}
