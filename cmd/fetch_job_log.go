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

type FetchJobLogConfig struct {
	FetchConfig

	printSections bool
	printMetrics  bool
}

func NewFetchJobLogCmd(out io.Writer) *cli.Command {
	cfg := FetchJobLogConfig{
		FetchConfig: FetchConfig{
			RootConfig: RootConfig{
				out:   out,
				flags: flag.NewFlagSet(fmt.Sprintf("%s fetch joblog", exeName), flag.ContinueOnError),
			},
		},
	}

	cfg.RegisterFlags(cfg.flags)

	return &cli.Command{
		Name:       "joblog",
		ShortUsage: fmt.Sprintf("%s fetch joblog [option]... project_id job_id", exeName),
		ShortHelp:  "Fetch job log",
		Flags:      cfg.flags,
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

	jobID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing `job_id` argument: %w", err)
	}

	trace, err := glc.GetJobLog(ctx, projectID, jobID)
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
