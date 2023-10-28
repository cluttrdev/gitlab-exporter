package cmd

import (
	"context"
	"flag"
	"fmt"
	"strconv"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/controller"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/tasks"
)

type ExportPipelineConfig struct {
	exportConfig *ExportConfig

	exportSections    bool
	exportTestReports bool
	exportTraces      bool

	flags *flag.FlagSet
}

func NewExportPipelineCmd(exportConfig *ExportConfig) *ffcli.Command {
	fs := flag.NewFlagSet(fmt.Sprintf("%s export pipeline", exeName), flag.ContinueOnError)

	config := &ExportPipelineConfig{
		exportConfig: exportConfig,

		flags: fs,
	}

	config.RegisterFlags(fs)

	return &ffcli.Command{
		Name:       "pipeline",
		ShortUsage: fmt.Sprintf("%s export pipeline [flags] project_id pipeline_id", exeName),
		ShortHelp:  "Export pipeline data",
		UsageFunc:  usageFunc,
		FlagSet:    fs,
		Options:    rootCmdOptions,
		Exec:       config.Exec,
	}
}

func (c *ExportPipelineConfig) RegisterFlags(fs *flag.FlagSet) {
	c.exportConfig.RegisterFlags(fs)

	fs.BoolVar(&c.exportSections, "export-sections", true, "Export job sections. (default: true)")
	fs.BoolVar(&c.exportTraces, "export-traces", true, "Export pipeline trace. (default: true)")
	fs.BoolVar(&c.exportTestReports, "export-testreports", true, "Export pipeline test reports. (default: true)")
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
	if err := loadConfig(c.exportConfig.rootConfig.filename, c.flags, &cfg); err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	ctl, err := controller.NewController(cfg)
	if err != nil {
		return fmt.Errorf("error constructing controller: %w", err)
	}

	opts := tasks.ExportPipelineHierarchyOptions{
		ProjectID:  projectID,
		PipelineID: pipelineID,

		ExportSections:    c.exportSections,
		ExportTestReports: c.exportTestReports,
		ExportTraces:      c.exportTraces,
	}

	return tasks.ExportPipelineHierarchy(ctx, opts, &ctl.GitLab, &ctl.ClickHouse)
}
