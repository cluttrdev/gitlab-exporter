package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/controller"
)

type FetchTestReportConfig struct {
	fetchConfig *FetchConfig

	flags *flag.FlagSet
}

func NewFetchTestReportCmd(fetchConfig *FetchConfig) *ffcli.Command {
	fs := flag.NewFlagSet(fmt.Sprintf("%s fetch testreport", exeName), flag.ContinueOnError)

	config := FetchTestReportConfig{
		fetchConfig: fetchConfig,

		flags: fs,
	}

	return &ffcli.Command{
		Name:       "testreport",
		ShortUsage: fmt.Sprintf("%s fetch testreport [flags] project_id pipeline_id", exeName),
		ShortHelp:  "Fetch pipeline testreport",
		UsageFunc:  usageFunc,
		FlagSet:    fs,
		Options:    rootCmdOptions,
		Exec:       config.Exec,
	}
}

func (c *FetchTestReportConfig) RegisterFlags(fs *flag.FlagSet) {
	c.fetchConfig.RegisterFlags(fs)
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

	log.SetOutput(c.fetchConfig.out)

	cfg := config.Default()
	if err := loadConfig(c.fetchConfig.rootConfig.filename, c.flags, &cfg); err != nil {
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

	fmt.Fprint(c.fetchConfig.out, string(b))

	return nil
}
