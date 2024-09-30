package tasks

import (
	"context"
	"fmt"

	"github.com/cluttrdev/gitlab-exporter/internal/exporter"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
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

	phr, err := glc.GetPipelineHierarchy(ctx, opts.ProjectID, opts.PipelineID, opt)
	if err != nil {
		return fmt.Errorf("get pipeline hierarchy: %w", err)
	}
	ph := phr.PipelineHierarchy

	if err := exp.ExportPipelineHierarchy(ctx, ph); err != nil {
		return fmt.Errorf("export pipeline hierarchy: %w", err)
	}

	if opts.ExportTraces {
		traces := ph.GetAllTraces()
		if err := exp.ExportTraces(ctx, traces); err != nil {
			return fmt.Errorf("export traces: %w", err)
		}
	}

	if opts.ExportMetrics {
		if err := exp.ExportMetrics(ctx, phr.Metrics); err != nil {
			return fmt.Errorf("export metrics: %w", err)
		}
	}

	if opts.ExportTestReports {
		if err := exportPipelineHierarchyTestReports(ctx, glc, exp, ph); err != nil {
			return fmt.Errorf("export testreports: %w", err)
		}
	}

	return nil
}

func exportPipelineHierarchyTestReports(ctx context.Context, glab *gitlab.Client, exp *exporter.Exporter, ph *gitlab.PipelineHierarchy) error {
	var (
		projectID  = ph.Pipeline.ProjectId
		pipelineID = ph.Pipeline.Id
	)

	r, err := glab.GetPipelineTestReport(ctx, projectID, pipelineID)
	if err != nil {
		return fmt.Errorf("error getting testreports (project=%d pipeline=%d): %w", projectID, pipelineID, err)
	}

	if err := exp.ExportTestReports(ctx, []*typespb.TestReport{r.TestReport}); err != nil {
		return fmt.Errorf("error exporting testreports (project=%d pipeline=%d): %w", projectID, pipelineID, err)
	}
	if err := exp.ExportTestSuites(ctx, r.TestSuites); err != nil {
		return fmt.Errorf("error exporting testsuites (project=%d pipeline=%d): %w", projectID, pipelineID, err)
	}
	if err := exp.ExportTestCases(ctx, r.TestCases); err != nil {
		return fmt.Errorf("error exporting testcases (project=%d pipeline=%d): %w", projectID, pipelineID, err)
	}

	for _, dph := range ph.DownstreamPipelines {
		if err := exportPipelineHierarchyTestReports(ctx, glab, exp, dph); err != nil {
			return err
		}
	}

	return nil
}
