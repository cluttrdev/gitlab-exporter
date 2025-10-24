package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"github.com/cluttrdev/cli"

	"go.cluttr.dev/gitlab-exporter/internal/config"
	"go.cluttr.dev/gitlab-exporter/internal/tasks"
)

type FetchRunnersConfig struct {
	FetchConfig
}

func NewFetchRunnerCmd(out io.Writer) *cli.Command {
	cfg := FetchRunnersConfig{
		FetchConfig: FetchConfig{
			RootConfig: RootConfig{
				out:   out,
				flags: flag.NewFlagSet("fetch-runners", flag.ContinueOnError),
			},
		},
	}

	cfg.RegisterFlags(cfg.flags)

	return &cli.Command{
		Name:       "runners",
		ShortUsage: fmt.Sprintf("%s fetch runneris [option]...", exeName),
		ShortHelp:  "Fetch all runners (requires admin access)",
		Flags:      cfg.flags,
		Exec:       cfg.Exec,
	}
}

func (c *FetchRunnersConfig) RegisterFlags(fs *flag.FlagSet) {
	c.FetchConfig.RegisterFlags(fs)
}

func (c *FetchRunnersConfig) Exec(ctx context.Context, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("too many arguments: %v", args)
	}

	cfg := config.Default()
	if err := loadConfig(c.FetchConfig.RootConfig.filename, c.flags, &cfg); err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	glab, err := createGitLabClient(cfg)
	if err != nil {
		return fmt.Errorf("create gitlab client: %w", err)
	}

	runners, err := tasks.FetchRunners(ctx, glab)
	if err != nil {
		return fmt.Errorf("fetch runners: %w", err)
	}

	out, err := json.Marshal(runners)
	if err != nil {
		return fmt.Errorf("marshal runners: %w", err)
	}

	_, err = fmt.Fprintln(c.FetchConfig.RootConfig.out, string(out))
	if err != nil {
		return err
	}

	return nil
}
