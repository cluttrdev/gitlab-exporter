package tasks

import (
	"context"
	"fmt"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/gitlab"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/models"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/internal/datastore"
)

type ExportPipelineHierarchyOptions struct {
	ProjectID  int64
	PipelineID int64

	ExportSections    bool
	ExportTestReports bool
	ExportTraces      bool
}

func ExportPipelineHierarchy(ctx context.Context, opts ExportPipelineHierarchyOptions, gl *gitlab.Client, ds datastore.DataStore) error {
	return <-exportPipelineHierarchy(ctx, opts, gl, ds)
}

func exportPipelineHierarchy(ctx context.Context, opts ExportPipelineHierarchyOptions, gl *gitlab.Client, ds datastore.DataStore) <-chan error {
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

		if err := ds.InsertPipelineHierarchy(ctx, ph); err != nil {
			out <- fmt.Errorf("error inserting pipeline hierarchy: %w", err)
			return
		}

		if opts.ExportTraces {
			pts := ph.GetAllTraces()
			if err := ds.InsertTraces(ctx, pts); err != nil {
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
			if err = ds.InsertTestReports(ctx, trs); err != nil {
				out <- fmt.Errorf("error inserting testreports: %w", err)
				return
			}
			if err = ds.InsertTestSuites(ctx, tss); err != nil {
				out <- fmt.Errorf("error inserting testsuites: %w", err)
				return
			}
			if err = ds.InsertTestCases(ctx, tcs); err != nil {
				out <- fmt.Errorf("error inserting testcases: %w", err)
				return
			}
		}
	}()

	return out
}
