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

    fs := flag.NewFlagSet(fmt.Sprintf("%s fetch pipeline", exeName), flag.ExitOnError)
    config.RegisterFlags(fs)
    config.fetchConfig.rootConfig.RegisterFlags(fs)

    return &ffcli.Command{
        Name: "pipeline",
        ShortUsage: fmt.Sprintf("%s [flags] fetch pipeline [flags] project_id pipeline_id", exeName),
        ShortHelp: "Fetch pipeline data",
        FlagSet: fs,
        Exec: config.Exec,
    }
}

func (c *FetchPipelineConfig) RegisterFlags(fs *flag.FlagSet) {
    fs.BoolVar(&c.all, "all", false, "Fetch pipeline hierarchy.")
}

func (c *FetchPipelineConfig) Exec(ctx context.Context, args []string) error {
    if len(args) != 2 {
        return fmt.Errorf("invalid number of positional arguments: %v", args)
    }

    log.SetOutput(c.fetchConfig.out)

    ctl := c.fetchConfig.rootConfig.Controller

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
        ph, err := ctl.GitLab.GetPipelineHierarchy(ctx, projectID, pipelineID)
        if err != nil {
            return fmt.Errorf("error fetching pipeline hierarchy: %w", err)
        }

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
