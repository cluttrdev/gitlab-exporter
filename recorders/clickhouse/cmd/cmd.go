package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/cluttrdev/cli"
)

const (
	exeName      string = "gitlab-exporter-clickhouse-recorder"
	envVarPrefix string = "GLCHR"
)

var Version string

func Execute() error {
	out := os.Stderr

	root := NewRootCmd(out)
	root.Subcommands = []*cli.Command{
		NewRunCmd(out),
		NewDeduplicateCmd(out),
		NewMigrateCommand(out),
		cli.NewVersionCommand(cli.NewBuildInfo(Version), out),
	}

	args := os.Args[1:]
	opts := []cli.ParseOption{
		cli.WithEnvVarPrefix(envVarPrefix),
	}

	if err := root.Parse(args, opts...); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		} else {
			return fmt.Errorf("error parsing arguments: %w", err)
		}
	}

	return root.Run(context.Background())
}
