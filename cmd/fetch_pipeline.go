package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/cluttrdev/gitlab-exporter/pkg/config"
	"github.com/cluttrdev/gitlab-exporter/pkg/controller"
	"github.com/cluttrdev/gitlab-exporter/pkg/gitlab"
)

type FetchPipelineConfig struct {
	fetchConfig *FetchConfig

	fetchHierarchy bool
	fetchSections  bool

	outputTrace bool

	flags *flag.FlagSet
}

func NewFetchPipelineCmd(fetchConfig *FetchConfig) *ffcli.Command {
	fs := flag.NewFlagSet(fmt.Sprintf("%s fetch pipeline", exeName), flag.ContinueOnError)

	config := FetchPipelineConfig{
		fetchConfig: fetchConfig,

		flags: fs,
	}

	config.RegisterFlags(fs)

	return &ffcli.Command{
		Name:       "pipeline",
		ShortUsage: fmt.Sprintf("%s fetch pipeline [flags] project_id pipeline_id", exeName),
		ShortHelp:  "Fetch pipeline data",
		UsageFunc:  usageFunc,
		FlagSet:    fs,
		Options:    rootCmdOptions,
		Exec:       config.Exec,
	}
}

func (c *FetchPipelineConfig) RegisterFlags(fs *flag.FlagSet) {
	c.fetchConfig.RegisterFlags(fs)

	fs.BoolVar(&c.fetchHierarchy, "hierarchy", false, "Fetch pipeline hierarchy. (default: false)")
	fs.BoolVar(&c.fetchSections, "fetch-sections", true, "Fetch job sections. (default: true)")
	fs.BoolVar(&c.outputTrace, "trace", false, "Output pipeline trace. (default: false)")
}

func (c *FetchPipelineConfig) Exec(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("invalid number of positional arguments: %v", args)
	}

	log.SetOutput(c.fetchConfig.out)

	cfg := config.Default()
	if err := loadConfig(c.fetchConfig.rootConfig.filename, c.flags, &cfg); err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	ctl, err := controller.NewController(cfg)
	if err != nil {
		return fmt.Errorf("error constructing controller: %w", err)
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

		phr := <-ctl.GitLab.GetPipelineHierarchy(ctx, projectID, pipelineID, opt)
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
		p, err := ctl.GitLab.GetPipeline(ctx, projectID, pipelineID)
		if err != nil {
			return fmt.Errorf("error fetching pipeline: %w", err)
		}

		b, err = json.Marshal(p)
		if err != nil {
			return fmt.Errorf("error marshalling pipeline %w", err)
		}
	}

	fmt.Fprint(c.fetchConfig.out, string(b))

	return nil
}
