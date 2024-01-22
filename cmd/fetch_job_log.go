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
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
)

type FetchJobLogConfig struct {
	FetchConfig

	printSections bool
	printMetrics  bool
}

func NewFetchJobLogCmd(out io.Writer) *cli.Command {
	fs := flag.NewFlagSet(fmt.Sprintf("%s fetch joblog", exeName), flag.ContinueOnError)

	cfg := FetchJobLogConfig{
		FetchConfig: FetchConfig{
			RootConfig: RootConfig{
				out: out,
			},
		},
	}

	cfg.RegisterFlags(fs)

	return &cli.Command{
		Name:       "joblog",
		ShortUsage: fmt.Sprintf("%s fetch joblog [option]... project_id job_id", exeName),
		ShortHelp:  "Fetch job log",
		Flags:      fs,
		Exec:       cfg.Exec,
	}
}

func (c *FetchJobLogConfig) RegisterFlags(fs *flag.FlagSet) {
	c.FetchConfig.RegisterFlags(fs)

	fs.BoolVar(&c.printSections, "sections", false, "Print parsed job log sections. (default: false)")
	fs.BoolVar(&c.printMetrics, "metrics", false, "Print parsed job log embedded metrics. (default: false)")
}

func (c *FetchJobLogConfig) Exec(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("invalid number of positional arguments: %v", args)
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

	if c.printSections || c.printMetrics {
		data, err := gitlab.ParseJobLog(trace)
		if err != nil {
			return err
		}

		var m []byte
		switch {
		case c.printSections && c.printMetrics:
			m, err = json.Marshal(data)
		case c.printSections:
			m, err = json.Marshal(data.Sections)
		case c.printMetrics:
			m, err = json.Marshal(data.Metrics)
		}
		if err != nil {
			return err
		}
		fmt.Println(string(m))
	} else {
		_, err = io.Copy(c.FetchConfig.RootConfig.out, trace)
		if err != nil {
			return err
		}
	}

	return nil
}
