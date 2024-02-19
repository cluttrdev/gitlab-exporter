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

type FetchPipelineConfig struct {
	FetchConfig

	fetchHierarchy bool
	fetchSections  bool

	outputTrace bool
}

func NewFetchPipelineCmd(out io.Writer) *cli.Command {
	cfg := FetchPipelineConfig{
		FetchConfig: FetchConfig{
			RootConfig: RootConfig{
				out:   out,
				flags: flag.NewFlagSet(fmt.Sprintf("%s fetch pipeline", exeName), flag.ContinueOnError),
			},
		},
	}

	cfg.RegisterFlags(cfg.flags)

	return &cli.Command{
		Name:       "pipeline",
		ShortUsage: fmt.Sprintf("%s fetch pipeline [option]... project_id pipeline_id", exeName),
		ShortHelp:  "Fetch pipeline data",
		Flags:      cfg.flags,
		Exec:       cfg.Exec,
	}
}

func (c *FetchPipelineConfig) RegisterFlags(fs *flag.FlagSet) {
	c.FetchConfig.RegisterFlags(fs)

	fs.BoolVar(&c.fetchHierarchy, "hierarchy", false, "Fetch pipeline hierarchy. (default: false)")
	fs.BoolVar(&c.fetchSections, "fetch-sections", true, "Fetch job sections. (default: true)")
	fs.BoolVar(&c.outputTrace, "trace", false, "Output pipeline trace. (default: false)")
}

func (c *FetchPipelineConfig) Exec(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("invalid number of positional arguments: %v", args)
	}

	cfg := config.Default()
	if err := loadConfig(c.FetchConfig.RootConfig.filename, c.flags, &cfg); err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	// create gitlab client
	glc, err := gitlab.NewGitLabClient(gitlab.ClientConfig{
		URL:   cfg.GitLab.Api.URL,
		Token: cfg.GitLab.Api.Token,

		RateLimit: cfg.GitLab.Client.Rate.Limit,
	})
	if err != nil {
		return fmt.Errorf("error creating gitlab client: %w", err)
	}

	projectID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing `project_id` argument: %w", err)
	}

	pipelineID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing `pipeline_id` argument: %w", err)
	}

	var b []byte
	if c.fetchHierarchy || c.outputTrace {
		opt := &gitlab.GetPipelineHierarchyOptions{
			FetchSections: c.fetchSections,
		}

		phr := <-glc.GetPipelineHierarchy(ctx, projectID, pipelineID, opt)
		if err := phr.Error; err != nil {
			return fmt.Errorf("error fetching pipeline hierarchy: %w", err)
		}
		ph := phr.PipelineHierarchy

		if c.outputTrace {
			ts := ph.GetAllTraces()
			b, err = json.Marshal(ts)
			if err != nil {
				return fmt.Errorf("error marshalling pipeline traces: %w", err)
			}
		} else {
			b, err = json.Marshal(ph)
			if err != nil {
				return fmt.Errorf("error marshalling pipeline hierarchy: %w", err)
			}
		}
	} else {
		p, err := glc.GetPipeline(ctx, projectID, pipelineID)
		if err != nil {
			return fmt.Errorf("error fetching pipeline: %w", err)
		}

		b, err = json.Marshal(p)
		if err != nil {
			return fmt.Errorf("error marshalling pipeline %w", err)
		}
	}

	fmt.Fprint(c.FetchConfig.RootConfig.out, string(b))

	return nil
}
