package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/cluttrdev/gitlab-exporter/pkg/config"
	"github.com/cluttrdev/gitlab-exporter/pkg/controller"
)

type FetchJobLogConfig struct {
	fetchConfig *FetchConfig

	flags *flag.FlagSet
}

func NewFetchJobLogCmd(fetchConfig *FetchConfig) *ffcli.Command {
	fs := flag.NewFlagSet(fmt.Sprintf("%s fetch joblog", exeName), flag.ContinueOnError)

	cfg := FetchJobLogConfig{
		fetchConfig: fetchConfig,

		flags: fs,
	}

	cfg.RegisterFlags(fs)

	return &ffcli.Command{
		Name:       "joblog",
		ShortUsage: fmt.Sprintf("%s fetch joblog [flags] project_id job_id", exeName),
		ShortHelp:  "Fetch job log",
		UsageFunc:  usageFunc,
		FlagSet:    fs,
		Options:    rootCmdOptions,
		Exec:       cfg.Exec,
	}
}

func (c *FetchJobLogConfig) RegisterFlags(fs *flag.FlagSet) {
	c.fetchConfig.RegisterFlags(fs)
}

func (c *FetchJobLogConfig) Exec(ctx context.Context, args []string) error {
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

	jobID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing `job_id` argument: %w", err)
	}

	trace, err := ctl.GitLab.GetJobLog(ctx, projectID, jobID)
	if err != nil {
		return fmt.Errorf("error fetching job log: %w", err)
	}

	_, err = io.Copy(c.fetchConfig.out, trace)
	if err != nil {
		return err
	}

	return nil
}
