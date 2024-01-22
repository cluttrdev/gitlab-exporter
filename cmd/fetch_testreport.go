package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/cluttrdev/cli"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
	"github.com/cluttrdev/gitlab-exporter/internal/controller"
)

type FetchTestReportConfig struct {
	FetchConfig
}

func NewFetchTestReportCmd(out io.Writer) *cli.Command {
	fs := flag.NewFlagSet(fmt.Sprintf("%s fetch testreport", exeName), flag.ContinueOnError)

	config := FetchTestReportConfig{
		FetchConfig: FetchConfig{
			RootConfig: RootConfig{},
		},
	}

	return &cli.Command{
		Name:       "testreport",
		ShortUsage: fmt.Sprintf("%s fetch testreport [option]... project_id pipeline_id", exeName),
		ShortHelp:  "Fetch pipeline testreport",
		Flags:      fs,
		Exec:       config.Exec,
	}
}

func (c *FetchTestReportConfig) RegisterFlags(fs *flag.FlagSet) {
	c.FetchConfig.RegisterFlags(fs)
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

	log.SetOutput(c.FetchConfig.RootConfig.out)

	cfg := config.Default()
	if err := loadConfig(c.FetchConfig.RootConfig.filename, &c.flags, &cfg); err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	ctl, err := controller.NewController(cfg)
	if err != nil {
		return fmt.Errorf("error constructing controller: %w", err)
	}

	tr, err := ctl.GitLab.GetPipelineTestReport(ctx, projectID, pipelineID)
	if err != nil {
		return fmt.Errorf("error fetching pipeline testreport: %w", err)
	}

	b, err := json.Marshal(tr)
	if err != nil {
		return fmt.Errorf("error marshalling pipeline testreport %w", err)
	}

	fmt.Fprint(c.FetchConfig.RootConfig.out, string(b))

	return nil
}
