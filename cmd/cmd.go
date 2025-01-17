package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/cluttrdev/cli"
)

// Version is the version to be overriden when building the binary
var Version string

func Execute() error {
	var (
		out        = os.Stdout
		rootCmd    = NewRootCmd(out)
		versionCmd = cli.NewVersionCommand(cli.NewBuildInfo(Version), out)
		runCmd     = NewRunCmd(out)
		fetchCmd   = NewFetchCmd(out)
		exportCmd  = NewExportCmd(out)
		catchupCmd = NewCatchUpCmd(out)
	)

	rootCmd.Subcommands = []*cli.Command{
		versionCmd,
		runCmd,
		fetchCmd,
		exportCmd,
		catchupCmd,
		NewOAuthCmd(out),
	}

	opts := []cli.ParseOption{
		cli.WithEnvVarPrefix("GLE"),
	}

	if err := rootCmd.Parse(os.Args[1:], opts...); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		} else {
			return fmt.Errorf("error parsing args: %v\n", err)
		}
	}

	return rootCmd.Run(context.Background())
}
