package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/cluttrdev/cli"
)

func Execute() error {
	var (
		out        = os.Stdout
		rootCmd    = NewRootCmd(out)
		versionCmd = cli.DefaultVersionCommand(out)
		runCmd     = NewRunCmd(out)
		fetchCmd   = NewFetchCmd(out)
		exportCmd  = NewExportCmd(out)
	)

	rootCmd.Subcommands = []*cli.Command{
		versionCmd,
		runCmd,
		fetchCmd,
		exportCmd,
	}

	if err := rootCmd.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		} else {
			return fmt.Errorf("error parsing args: %v\n", err)
		}
	}

	return rootCmd.Run(context.Background())
}
