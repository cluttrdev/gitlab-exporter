package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strconv"

	"github.com/cluttrdev/cli"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
)

type FetchTestReportConfig struct {
	FetchConfig

	summary bool
}

func NewFetchTestReportCmd(out io.Writer) *cli.Command {
	cfg := FetchTestReportConfig{
		FetchConfig: FetchConfig{
			RootConfig: RootConfig{
				out:   out,
				flags: flag.NewFlagSet(fmt.Sprintf("%s fetch testreport", exeName), flag.ContinueOnError),
			},
		},
	}

	cfg.RegisterFlags(cfg.flags)

	return &cli.Command{
		Name:       "testreport",
		ShortUsage: fmt.Sprintf("%s fetch testreport [option]... project_id pipeline_id", exeName),
		ShortHelp:  "Fetch pipeline testreport",
		Flags:      cfg.flags,
		Exec:       cfg.Exec,
	}
}

func (c *FetchTestReportConfig) RegisterFlags(fs *flag.FlagSet) {
	c.FetchConfig.RegisterFlags(fs)

	fs.BoolVar(&c.summary, "summary", false, "Fetch testreport summary")
}

func (c *FetchTestReportConfig) Exec(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("invalid number of positional arguments: %v", args)
	}

	projectID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing `project_id` argument: %w", err)
	}

	pipelineID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing `pipeline_id` argument: %w", err)
	}

	cfg := config.Default()
	if err := loadConfig(c.FetchConfig.RootConfig.filename, c.flags, &cfg); err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	// create gitlab client
	glab, err := gitlab.NewGitLabClient(gitlab.ClientConfig{
		URL:   cfg.GitLab.Url,
		Token: cfg.GitLab.Token,

		RateLimit: cfg.GitLab.Client.Rate.Limit,
	})
	if err != nil {
		return fmt.Errorf("error creating gitlab client: %w", err)
	}

	var b []byte
	if c.summary {
		tr, err := glab.Rest.GetPipelineTestReportSummary(ctx, projectID, pipelineID)
		if err != nil {
			return fmt.Errorf("error fetching pipeline testreport summary: %w", err)
		}

		b, err = json.Marshal(tr)
		if err != nil {
			return fmt.Errorf("error marshalling pipeline testreport summary: %w", err)
		}

	} else {
		tr, _, err := glab.Rest.GetPipelineTestReport(ctx, projectID, pipelineID)
		if err != nil {
			return fmt.Errorf("error fetching pipeline testreport: %w", err)
		}

		b, err = json.Marshal(tr)
		if err != nil {
			return fmt.Errorf("error marshalling pipeline testreport: %w", err)
		}
	}

	fmt.Fprint(c.FetchConfig.RootConfig.out, string(b))

	return nil
}
