package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/peterbourgon/ff/v3/ffcli"
)

type FetchPipelineConfig struct {
	fetchConfig *FetchConfig

	all bool
}

func NewFetchPipelineCmd(fetchConfig *FetchConfig) *ffcli.Command {
	config := FetchPipelineConfig{
		fetchConfig: fetchConfig,
	}

	fs := flag.NewFlagSet(fmt.Sprintf("%s fetch pipeline", exeName), flag.ContinueOnError)
	config.registerFlags(fs)

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

func (c *FetchPipelineConfig) registerFlags(fs *flag.FlagSet) {
	fs.BoolVar(&c.all, "all", false, "Fetch pipeline hierarchy.")
}

func (c *FetchPipelineConfig) Exec(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("invalid number of positional arguments: %v", args)
	}

	log.SetOutput(c.fetchConfig.out)

	ctl, err := c.fetchConfig.rootConfig.newController()
	if err != nil {
		return err
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
	if c.all {
		phr := <-ctl.GitLab.GetPipelineHierarchy(ctx, projectID, pipelineID)
		if err := phr.Error; err != nil {
			return fmt.Errorf("error fetching pipeline hierarchy: %w", err)
		}
		ph := phr.PipelineHierarchy

		b, err = json.Marshal(ph)
		if err != nil {
			return fmt.Errorf("error marshalling pipeline hierarchy: %w", err)
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
