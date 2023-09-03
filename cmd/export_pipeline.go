package cmd

import (
	"context"
	"flag"
	"fmt"
	"strconv"

	clickhouse "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/clickhouse"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/models"
	"github.com/peterbourgon/ff/v3/ffcli"
)

type ExportPipelineConfig struct {
	exportConfig *ExportConfig

	exportTrace       bool
	exportTestReports bool
}

func NewExportPipelineCmd(exportConfig *ExportConfig) *ffcli.Command {
	config := &ExportPipelineConfig{
		exportConfig: exportConfig,
	}

	fs := flag.NewFlagSet(fmt.Sprintf("%s export pipeline", exeName), flag.ContinueOnError)
	config.RegisterFlags(fs)
	config.exportConfig.rootConfig.RegisterFlags(fs)

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
	fs.BoolVar(&c.exportTrace, "export-traces", true, "Export pipeline trace.")
	fs.BoolVar(&c.exportTestReports, "export-testreports", true, "Export pipeline test reports.")
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

	ctl := c.exportConfig.rootConfig.Controller

	phr := <-ctl.GitLab.GetPipelineHierarchy(ctx, projectID, pipelineID)
	if err := phr.Error; err != nil {
		return fmt.Errorf("error fetching pipeline hierarchy: %w", err)
	}
	ph := phr.PipelineHierarchy

	pts := [][]*models.Span{}
	if c.exportTrace {
		pts = ph.GetAllTraces()
	}

	trs := []*models.PipelineTestReport{}
	if c.exportTestReports {
		trs, err = ctl.GitLab.GetPipelineHierarchyTestReports(ctx, ph)
		if err != nil {
			return fmt.Errorf("error fetching pipeline hierarchy test reports: %w", err)
		}
	}
	tss := []*models.PipelineTestSuite{}
	tcs := []*models.PipelineTestCase{}
	for _, tr := range trs {
		tss = append(tss, tr.TestSuites...)
		for _, ts := range tr.TestSuites {
			tcs = append(tcs, ts.TestCases...)
		}
	}

	if err = clickhouse.InsertPipelineHierarchy(ctx, ph, ctl.ClickHouse); err != nil {
		return fmt.Errorf("error inserting pipeline hierarchy: %w", err)
	}

	if c.exportTrace {
		if err = clickhouse.InsertTraces(ctx, pts, ctl.ClickHouse); err != nil {
			return fmt.Errorf("error inserting pipeline trace: %w", err)
		}
	}

	if c.exportTestReports {
		if err = clickhouse.InsertTestReports(ctx, trs, ctl.ClickHouse); err != nil {
			return fmt.Errorf("error inserting testreports: %w", err)
		}
		if err = clickhouse.InsertTestSuites(ctx, tss, ctl.ClickHouse); err != nil {
			return fmt.Errorf("error inserting testsuites: %w", err)
		}
		if err = clickhouse.InsertTestCases(ctx, tcs, ctl.ClickHouse); err != nil {
			return fmt.Errorf("error inserting testcases: %w", err)
		}
	}

	return nil
}
