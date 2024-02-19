package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strconv"

	"github.com/cluttrdev/cli"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
	"github.com/cluttrdev/gitlab-exporter/internal/exporter"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
	"github.com/cluttrdev/gitlab-exporter/internal/tasks"
)

type ExportPipelineConfig struct {
	ExportConfig

	exportSections    bool
	exportTestReports bool
	exportTraces      bool
	exportMetrics     bool
}

func NewExportPipelineCmd(out io.Writer) *cli.Command {
	fs := flag.NewFlagSet(fmt.Sprintf("%s export pipeline", exeName), flag.ContinueOnError)

	config := &ExportPipelineConfig{
		ExportConfig: ExportConfig{
			RootConfig: RootConfig{
				out: out,
			},
		},
	}

	config.RegisterFlags(fs)

	return &cli.Command{
		Name:       "pipeline",
		ShortUsage: fmt.Sprintf("%s export pipeline [option]... project_id pipeline_id", exeName),
		ShortHelp:  "Export pipeline data",
		Flags:      fs,
		Exec:       config.Exec,
	}
}

func (c *ExportPipelineConfig) RegisterFlags(fs *flag.FlagSet) {
	c.ExportConfig.RegisterFlags(fs)

	fs.BoolVar(&c.exportSections, "export-sections", true, "Export job sections. (default: true)")
	fs.BoolVar(&c.exportTraces, "export-traces", true, "Export pipeline trace. (default: true)")
	fs.BoolVar(&c.exportTestReports, "export-testreports", true, "Export pipeline test reports. (default: true)")
	fs.BoolVar(&c.exportMetrics, "export-metrics", true, "Export job log embedded metrics. (default: true)")
}

func (c *ExportPipelineConfig) Exec(ctx context.Context, args []string) error {
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

	cfg := config.Default()
	if err := loadConfig(c.ExportConfig.RootConfig.filename, &c.flags, &cfg); err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	// create gitlab client
	gitlabclient, err := gitlab.NewGitLabClient(gitlab.ClientConfig{
		URL:   cfg.GitLab.Api.URL,
		Token: cfg.GitLab.Api.Token,

		RateLimit: cfg.GitLab.Client.Rate.Limit,
	})
	if err != nil {
		return fmt.Errorf("error creating gitlab client: %w", err)
	}

	// create exporter
	endpoints := exporter.CreateEndpointConfigs(cfg.Endpoints)
	exp, err := exporter.New(endpoints)
	if err != nil {
		return fmt.Errorf("error creating exporter: %w", err)
	}

	opts := tasks.ExportPipelineHierarchyOptions{
		ProjectID:  projectID,
		PipelineID: pipelineID,

		ExportSections:           c.exportSections,
		ExportTestReports:        c.exportTestReports,
		ExportTraces:             c.exportTraces,
		ExportLogEmbeddedMetrics: c.exportMetrics,
	}

	return tasks.ExportPipelineHierarchy(ctx, gitlabclient, exp, opts)
}
