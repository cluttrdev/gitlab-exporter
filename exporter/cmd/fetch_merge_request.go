package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"github.com/cluttrdev/cli"

	"go.cluttr.dev/gitlab-exporter/exporter/internal/config"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab/graphql"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/types"
)

type FetchMergeRequestConfig struct {
	FetchConfig

	projectPath     string
	mergeRequestIid int64
}

func NewFetchMergeRequestCmd(out io.Writer) *cli.Command {
	cfg := FetchMergeRequestConfig{
		FetchConfig: FetchConfig{
			RootConfig: RootConfig{
				out:   out,
				flags: flag.NewFlagSet("fetch merge_request", flag.ContinueOnError),
			},
		},
	}

	cfg.RegisterFlags(cfg.flags)

	return &cli.Command{
		Name:       "merge_request",
		ShortUsage: fmt.Sprintf("%s fetch merge_request [option]...", exeName),
		ShortHelp:  "Fetch merge request data",
		Flags:      cfg.flags,
		Exec:       cfg.Exec,
	}
}

func (c *FetchMergeRequestConfig) RegisterFlags(fs *flag.FlagSet) {
	c.FetchConfig.RegisterFlags(fs)

	fs.StringVar(&c.projectPath, "project-path", "", "The merge request project's full path.")
	fs.Int64Var(&c.mergeRequestIid, "merge-request-iid", 0, "The merge request's iid.")
}

func (c *FetchMergeRequestConfig) Exec(ctx context.Context, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("no positional arguments expected, got %d: %v", len(args), args)
	}

	cfg := config.Default()
	if err := loadConfig(c.FetchConfig.RootConfig.filename, c.flags, &cfg); err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	// create gitlab client
	glab, err := createGitLabClient(cfg)
	if err != nil {
		return fmt.Errorf("create gitlab client: %w", err)
	}

	mrf, err := glab.GraphQL.GetProjectMergeRequest(ctx, c.projectPath, c.mergeRequestIid)
	if err != nil {
		return fmt.Errorf("get merge request fields: %w", err)
	}

	mr, err := graphql.ConvertMergeRequest(mrf)
	if err != nil {
		return fmt.Errorf("convert merge request fields: %w", err)
	}

	var mrCommits []types.MergeRequestCommit
	if len(mrf.Commits) > 0 {
		mrRef := types.MergeRequestReference{
			Id:      mr.Id,
			Iid:     mr.Iid,
			Project: mr.Project,
		}
		for _, cf := range mrf.Commits {
			mrc, err := graphql.ConvertMergeRequestCommit(mrRef, cf)
			if err != nil {
				return fmt.Errorf("convert merge request commit: %w", err)
			}
			mrCommits = append(mrCommits, mrc)
		}
	}

	v := struct {
		types.MergeRequest
		Commits []types.MergeRequestCommit
	}{
		MergeRequest: mr,
		Commits:      mrCommits,
	}
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal merge request: %w", err)
	}

	fmt.Fprint(c.FetchConfig.RootConfig.out, string(b))

	return nil
}
