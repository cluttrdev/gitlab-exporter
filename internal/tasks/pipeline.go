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

	ExportSections    bool
	ExportTestReports bool
	ExportTraces      bool
	ExportMetrics     bool
}

func ExportPipelineHierarchy(ctx context.Context, glc *gitlab.Client, exp *exporter.Exporter, opts ExportPipelineHierarchyOptions) error {
	opt := &gitlab.GetPipelineHierarchyOptions{
		FetchSections:   opts.ExportSections,
		FetchJobMetrics: opts.ExportSections,
	}

	phr := <-glc.GetPipelineHierarchy(ctx, opts.ProjectID, opts.PipelineID, opt)
	if err := phr.Error; err != nil {
		return fmt.Errorf("error getting pipeline hierarchy (project=%d pipeline=%d): %w", opts.ProjectID, opts.PipelineID, err)
	}
	ph := phr.PipelineHierarchy

	if err := exp.ExportPipelineHierarchy(ctx, ph); err != nil {
		return fmt.Errorf("error exporting pipeline hierarchy (project=%d pipeline=%d): %w", opts.ProjectID, opts.PipelineID, err)
	}

	if opts.ExportTraces {
		traces := ph.GetAllTraces()
		if err := exp.ExportTraces(ctx, traces); err != nil {
			return fmt.Errorf("error exporting traces (project=%d pipeline=%d): %w", opts.ProjectID, opts.PipelineID, err)
		}
	}

	if opts.ExportMetrics {
		if err := exp.ExportMetrics(ctx, phr.Metrics); err != nil {
			return fmt.Errorf("error exporting metrics (project=%d pipeline=%d): %w", opts.ProjectID, opts.PipelineID, err)
		}
	}

	if opts.ExportTestReports {
		results, err := glc.GetPipelineHierarchyTestReports(ctx, ph)
		if err != nil {
			return fmt.Errorf("error getting testreports (project=%d pipeline=%d): %w", opts.ProjectID, opts.PipelineID, err)
		}
		if err := exp.ExportTestReports(ctx, results.TestReports); err != nil {
			return fmt.Errorf("error exporting testreports (project=%d pipeline=%d): %w", opts.ProjectID, opts.PipelineID, err)
		}
		if err := exp.ExportTestSuites(ctx, results.TestSuites); err != nil {
			return fmt.Errorf("error exporting testsuites (project=%d pipeline=%d): %w", opts.ProjectID, opts.PipelineID, err)
		}
		if err := exp.ExportTestCases(ctx, results.TestCases); err != nil {
			return fmt.Errorf("error exporting testcases (project=%d pipeline=%d): %w", opts.ProjectID, opts.PipelineID, err)
		}
	}

	return nil
}
