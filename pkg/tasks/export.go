package tasks

import (
	"context"
	"fmt"

	clickhouse "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/clickhouse"
	gitlab "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/gitlab"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/models"
)

type ExportPipelineHierarchyOptions struct {
	ProjectID  int64
	PipelineID int64

	ExportSections    bool
	ExportTestReports bool
	ExportTraces      bool
}

func ExportPipelineHierarchy(ctx context.Context, opts *ExportPipelineHierarchyOptions, gl *gitlab.Client, ch *clickhouse.Client) error {
	return <-exportPipelineHierarchy(ctx, opts, gl, ch)
}

func exportPipelineHierarchy(ctx context.Context, opts *ExportPipelineHierarchyOptions, gl *gitlab.Client, ch *clickhouse.Client) <-chan error {
	out := make(chan error)

	go func() {
		defer close(out)

		opt := &gitlab.GetPipelineHierarchyOptions{
			FetchSections: opts.ExportSections,
		}

		phr := <-gl.GetPipelineHierarchy(ctx, opts.ProjectID, opts.PipelineID, opt)
		if err := phr.Error; err != nil {
			out <- fmt.Errorf("error getting pipeline hierarchy: %w", err)
			return
		}
		ph := phr.PipelineHierarchy

		if err := clickhouse.InsertPipelineHierarchy(ctx, ph, ch); err != nil {
			out <- fmt.Errorf("error inserting pipeline hierarchy: %w", err)
			return
		}

		if opts.ExportTraces {
			pts := ph.GetAllTraces()
			if err := clickhouse.InsertTraces(ctx, pts, ch); err != nil {
				out <- fmt.Errorf("error inserting traces: %w", err)
				return
			}
		}

		if opts.ExportTestReports {
			trs, err := gl.GetPipelineHierarchyTestReports(ctx, ph)
			if err != nil {
				out <- fmt.Errorf("error getting testreports: %w", err)
				return
			}
			tss := []*models.PipelineTestSuite{}
			tcs := []*models.PipelineTestCase{}
			for _, tr := range trs {
				tss = append(tss, tr.TestSuites...)
				for _, ts := range tr.TestSuites {
					tcs = append(tcs, ts.TestCases...)
				}
			}
			if err = clickhouse.InsertTestReports(ctx, trs, ch); err != nil {
				out <- fmt.Errorf("error inserting testreports: %w", err)
				return
			}
			if err = clickhouse.InsertTestSuites(ctx, tss, ch); err != nil {
				out <- fmt.Errorf("error inserting testsuites: %w", err)
				return
			}
			if err = clickhouse.InsertTestCases(ctx, tcs, ch); err != nil {
				out <- fmt.Errorf("error inserting testcases: %w", err)
				return
			}
		}
	}()

	return out
}
