package tasks

import (
	"context"
	"fmt"

	"github.com/cluttrdev/gitlab-exporter/internal/exporter"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
)

type ExportPipelineHierarchyOptions struct {
	ProjectID  int64
	PipelineID int64

	ExportSections           bool
	ExportTestReports        bool
	ExportTraces             bool
	ExportLogEmbeddedMetrics bool
}

func ExportPipelineHierarchy(ctx context.Context, glc *gitlab.Client, exp *exporter.Exporter, opts ExportPipelineHierarchyOptions) error {
	opt := &gitlab.GetPipelineHierarchyOptions{
		FetchSections:   opts.ExportSections,
		FetchJobMetrics: opts.ExportSections,
	}

	phr := <-glc.GetPipelineHierarchy(ctx, opts.ProjectID, opts.PipelineID, opt)
	if err := phr.Error; err != nil {
		return fmt.Errorf("error getting pipeline hierarchy: %w", err)
	}
	ph := phr.PipelineHierarchy

	if err := exp.ExportPipelineHierarchy(ctx, ph); err != nil {
		return fmt.Errorf("error exporting pipeline hierarchy: %w", err)
	}

	if opts.ExportTraces {
		traces := ph.GetAllTraces()
		if err := exp.ExportTraces(ctx, traces); err != nil {
			return fmt.Errorf("error exporting traces: %w", err)
		}
	}

	if opts.ExportLogEmbeddedMetrics {
		if err := exp.ExportLogEmbeddedMetrics(ctx, phr.LogEmbeddedMetrics); err != nil {
			return fmt.Errorf("error exporting log embedded metrics: %w", err)
		}
	}

	if opts.ExportTestReports {
		results, err := glc.GetPipelineHierarchyTestReports(ctx, ph)
		if err != nil {
			return fmt.Errorf("error getting testreports: %w", err)
		}
		if err := exp.ExportTestReports(ctx, results.TestReports); err != nil {
			return fmt.Errorf("error exporting testreports: %w", err)
		}
		if err := exp.ExportTestSuites(ctx, results.TestSuites); err != nil {
			return fmt.Errorf("error exporting testsuites: %w", err)
		}
		if err := exp.ExportTestCases(ctx, results.TestCases); err != nil {
			return fmt.Errorf("error exporting testcases: %w", err)
		}
	}

	return nil
}
