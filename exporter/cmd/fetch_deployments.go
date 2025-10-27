package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/cluttrdev/cli"

	"go.cluttr.dev/gitlab-exporter/internal/config"
	"go.cluttr.dev/gitlab-exporter/internal/gitlab/rest"
	"go.cluttr.dev/gitlab-exporter/internal/types"
)

type FetchDeploymentsConfig struct {
	FetchConfig

	projectId     int64
	updatedAfter  *time.Time
	updatedBefore *time.Time
}

func NewFetchDeploymentsCmd(out io.Writer) *cli.Command {
	cfg := FetchDeploymentsConfig{
		FetchConfig: FetchConfig{
			RootConfig: RootConfig{
				out:   out,
				flags: flag.NewFlagSet("fetch-deployments", flag.ContinueOnError),
			},
		},
	}

	cfg.RegisterFlags(cfg.flags)

	return &cli.Command{
		Name:       "deployments",
		ShortUsage: fmt.Sprintf("%s fetch deployments [option]...", exeName),
		ShortHelp:  "Fetch project deployments.",
		Flags:      cfg.flags,
		Exec:       cfg.Exec,
	}
}

func (c *FetchDeploymentsConfig) RegisterFlags(fs *flag.FlagSet) {
	c.FetchConfig.RegisterFlags(fs)

	fs.Int64Var(&c.projectId, "project-id", 0, "The project id.")

	fs.Func("updated-after", "", func(s string) error {
		t, err := parseTimeISO8601(s)
		if err != nil {
			return err
		}
		c.updatedAfter = &t
		return nil
	})
	fs.Func("updated-before", "", func(s string) error {
		t, err := parseTimeISO8601(s)
		if err != nil {
			return err
		}
		c.updatedBefore = &t
		return nil
	})
}

func (c *FetchDeploymentsConfig) Exec(ctx context.Context, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("too many arguments: %v", args)
	}

	if c.projectId == 0 {
		return fmt.Errorf("missing required option: --project-id")
	}

	cfg := config.Default()
	if err := loadConfig(c.FetchConfig.RootConfig.filename, c.flags, &cfg); err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	glab, err := createGitLabClient(cfg)
	if err != nil {
		return fmt.Errorf("create gitlab client: %w", err)
	}

	opt := rest.GetProjectDeploymentsOptions{
		UpdatedAfter:  c.updatedAfter,
		UpdatedBefore: c.updatedBefore,
	}

	ds, err := glab.Rest.GetProjectDeployments(ctx, c.projectId, opt)
	if err != nil {
		return fmt.Errorf("fetch deployments: %w", err)
	}

	deployments := make([]types.Deployment, 0, len(ds))
	for _, d := range ds {
		deployments = append(deployments, rest.ConvertDeployment(d))
	}

	out, err := json.Marshal(deployments)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(c.FetchConfig.RootConfig.out, string(out))
	if err != nil {
		return err
	}

	return nil
}
