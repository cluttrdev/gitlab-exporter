package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/cluttrdev/cli"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
	"github.com/cluttrdev/gitlab-exporter/internal/junitxml"
	"github.com/cluttrdev/gitlab-exporter/internal/types"
)

type FetchReportConfig struct {
	FetchConfig

	projectPath   string
	jobId         int64
	fileType      string
	artifactPaths stringList
}

type stringList []string

func (f *stringList) String() string {
	return fmt.Sprintf("%v", []string(*f))
}

func (f *stringList) Set(value string) error {
	values := strings.Split(value, ",")
	*f = append(*f, values...)
	return nil
}

func NewFetchReportCmd(out io.Writer) *cli.Command {
	cfg := FetchReportConfig{
		FetchConfig: FetchConfig{
			RootConfig: RootConfig{
				out:   out,
				flags: flag.NewFlagSet("fetch report", flag.ExitOnError),
			},
		},
	}

	cfg.RegisterFlags(cfg.flags)

	return &cli.Command{
		Name:       "report",
		ShortUsage: fmt.Sprintf("%s fetch report [option]...", exeName),
		ShortHelp:  "Fetch job report",
		Flags:      cfg.flags,
		Exec:       cfg.Exec,
	}
}

func (c *FetchReportConfig) RegisterFlags(fs *flag.FlagSet) {
	c.FetchConfig.RegisterFlags(fs)

	fs.StringVar(&c.projectPath, "project-path", "", "The full project path.")
	fs.Int64Var(&c.jobId, "job-id", 0, "The job id.")
	fs.StringVar(&c.fileType, "file-type", "", "The report file type.")
	fs.Var(&c.artifactPaths, "artifact-paths", "Comma separated list of artifact paths.")
}

func (c *FetchReportConfig) Exec(ctx context.Context, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("too many arguments: %v", args)
	}

	if c.projectPath == "" {
		return fmt.Errorf("missing required option: --project-path")
	}
	if c.jobId == 0 {
		return fmt.Errorf("missing required option: --job-id")
	}
	if c.fileType != "junit" {
		return fmt.Errorf("file type not supported, yet: %s", c.fileType)
	}
	if len(c.artifactPaths) == 0 {
		return fmt.Errorf("missing report artifact paths")
	}

	cfg := config.Default()
	if err := loadConfig(c.FetchConfig.RootConfig.filename, c.flags, &cfg); err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	// create gitlab client
	glab, err := createGitLabClient(cfg)
	if err != nil {
		return fmt.Errorf("create gitlab client: %w", err)
	}

	var report junitxml.TestReport
	for _, path := range c.artifactPaths {
		reader, err := glab.Rest.GetProjectJobArtifact(ctx, c.projectPath, c.jobId, path)
		if errors.Is(err, gitlab.ErrNotFound) {
			continue
		} else if err != nil {
			return fmt.Errorf("download file: %w", err)
		}

		r, err := junitxml.Parse(reader)
		if err != nil {
			return fmt.Errorf("parse file: %w", err)
		}

		report.Tests += r.Tests
		report.Failures += r.Failures
		report.Errors += r.Errors
		report.Skipped += r.Skipped
		report.Time += r.Time
		report.Timestamp = r.Timestamp
		report.TestSuites = append(report.TestSuites, r.TestSuites...)
	}

	tr, ts, tc := junitxml.ConvertTestReport(report, types.JobReference{
		Id: c.jobId,
		Pipeline: types.PipelineReference{
			Project: types.ProjectReference{
				FullPath: c.projectPath,
			},
		},
	})

	b, err := json.Marshal(map[string]any{
		"testreport": tr,
		"testsuites": ts,
		"testcases":  tc,
	})
	if err != nil {
		return fmt.Errorf("error marshalling pipeline testreport: %w", err)
	}

	fmt.Fprint(c.FetchConfig.RootConfig.out, string(b))

	return nil
}
