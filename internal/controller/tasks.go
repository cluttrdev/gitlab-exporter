package controller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
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

func ExportPipelineHierarchy(ctl *Controller, ctx context.Context, opts ExportPipelineHierarchyOptions) error {
	opt := &gitlab.GetPipelineHierarchyOptions{
		FetchSections:   opts.ExportSections,
		FetchJobMetrics: opts.ExportSections,
	}

	phr := <-ctl.GitLab.GetPipelineHierarchy(ctx, opts.ProjectID, opts.PipelineID, opt)
	if err := phr.Error; err != nil {
		return fmt.Errorf("error getting pipeline hierarchy: %w", err)
	}
	ph := phr.PipelineHierarchy

	if err := ctl.Exporter.ExportPipelineHierarchy(ctx, ph); err != nil {
		return fmt.Errorf("error exporting pipeline hierarchy: %w", err)
	}

	if opts.ExportTraces {
		traces := ph.GetAllTraces()
		if err := ctl.Exporter.ExportTraces(ctx, traces); err != nil {
			return fmt.Errorf("error exporting traces: %w", err)
		}
	}

	if opts.ExportLogEmbeddedMetrics {
		if err := ctl.Exporter.ExportLogEmbeddedMetrics(ctx, phr.LogEmbeddedMetrics); err != nil {
			return fmt.Errorf("error exporting log embedded metrics: %w", err)
		}
	}

	if opts.ExportTestReports {
		results, err := ctl.GitLab.GetPipelineHierarchyTestReports(ctx, ph)
		if err != nil {
			return fmt.Errorf("error getting testreports: %w", err)
		}
		if err := ctl.Exporter.ExportTestReports(ctx, results.TestReports); err != nil {
			return fmt.Errorf("error exporting testreports: %w", err)
		}
		if err := ctl.Exporter.ExportTestSuites(ctx, results.TestSuites); err != nil {
			return fmt.Errorf("error exporting testsuites: %w", err)
		}
		if err := ctl.Exporter.ExportTestCases(ctx, results.TestCases); err != nil {
			return fmt.Errorf("error exporting testcases: %w", err)
		}
	}

	return nil
}

// ===========================================================================

type ProjectExportTask struct {
	Config config.Project
}

func (t *ProjectExportTask) Run(ctl *Controller, ctx context.Context) {
	interval := 60 * time.Second

	opt := gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},

		Scope: &[]string{"finished"}[0],
	}

	before := time.Now().UTC().Add(-interval)
	opt.UpdatedBefore = &before

	var first bool = true
	ticker := time.NewTicker(1 * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if first {
				ticker.Stop()
				ticker = time.NewTicker(interval)
				first = false
			}

			now := time.Now().UTC()
			opt.UpdatedAfter = opt.UpdatedBefore
			opt.UpdatedBefore = &now

			var wg sync.WaitGroup
			for r := range ctl.GitLab.ListProjectPipelines(ctx, t.Config.Id, opt) {
				if r.Error != nil {
					slog.Error("error listing project pipelines", "error", r.Error)
					continue
				}

				wg.Add(1)
				go func(pid int64) {
					defer wg.Done()

					opts := ExportPipelineHierarchyOptions{
						ProjectID:  t.Config.Id,
						PipelineID: pid,

						ExportSections:           t.Config.Export.Sections.Enabled,
						ExportTestReports:        t.Config.Export.TestReports.Enabled,
						ExportTraces:             t.Config.Export.Traces.Enabled,
						ExportLogEmbeddedMetrics: t.Config.Export.LogEmbeddedMetrics.Enabled,
					}

					if err := ExportPipelineHierarchy(ctl, ctx, opts); err != nil {
						if !errors.Is(err, context.Canceled) {
							slog.Error("error exporting pipeline hierarchy", "error", err)
						}
					} else {
						slog.Debug("Exported project pipeline hierarchy", "project_id", opts.ProjectID, "pipeline_id", opts.PipelineID)
					}
				}(r.Pipeline.Id)
			}
			wg.Wait()
		}
	}
}

// ===========================================================================

type ProjectCatchUpTask struct {
	Config config.Project
}

func (t *ProjectCatchUpTask) Run(ctl *Controller, ctx context.Context) {
	opt := gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},

		Scope: &[]string{"finished"}[0],
	}
	if t.Config.CatchUp.UpdatedAfter != "" {
		after, err := time.Parse("2006-01-02T15:04:05Z", t.Config.CatchUp.UpdatedAfter)
		if err != nil {
			slog.Error("error parsing catchup update_after", "error", err)
		} else {
			opt.UpdatedAfter = &after
		}
	}
	if t.Config.CatchUp.UpdatedBefore != "" {
		before, err := time.Parse("2006-01-02T15:04:05Z", t.Config.CatchUp.UpdatedBefore)
		if err != nil {
			slog.Error("error parsing catchup update_before", "error", err)
		} else {
			opt.UpdatedBefore = &before
		}
	}

	ch := t.produce(ctl, ctx, opt)
	t.process(ctl, ctx, ch)
}

func (t *ProjectCatchUpTask) produce(ctl *Controller, ctx context.Context, opt gitlab.ListProjectPipelinesOptions) <-chan int64 {
	ch := make(chan int64)

	go func() {
		defer close(ch)

		resChan := ctl.GitLab.ListProjectPipelines(ctx, t.Config.Id, opt)
		for {
			select {
			case <-ctx.Done():
				return
			case r, ok := <-resChan:
				if !ok { // channel closed
					return
				}

				if r.Error != nil && !errors.Is(r.Error, context.Canceled) {
					slog.Error("error listing project pipelines", "error", r.Error)
					continue
				}

				ch <- r.Pipeline.Id
			}
		}
	}()

	return ch
}

func (t *ProjectCatchUpTask) process(ctl *Controller, ctx context.Context, pipelineChan <-chan int64) {
	numWorkers := 10
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for pipelineID := range pipelineChan {
				opts := ExportPipelineHierarchyOptions{
					ProjectID:  t.Config.Id,
					PipelineID: pipelineID,

					ExportSections:           t.Config.Export.Sections.Enabled,
					ExportTestReports:        t.Config.Export.TestReports.Enabled,
					ExportTraces:             t.Config.Export.Traces.Enabled,
					ExportLogEmbeddedMetrics: t.Config.Export.Sections.Enabled, // for now, export metrics if we fetch the logs for sections anyway
				}

				if err := ExportPipelineHierarchy(ctl, ctx, opts); err != nil {
					if !errors.Is(err, context.Canceled) {
						slog.Error("error exporting pipeline hierarchy", "error", err)
					}
				} else {
					slog.Debug("Caught up on project pipeline hierarchy", "project_id", opts.ProjectID, "pipeline_id", opts.PipelineID)
				}
			}
		}()
	}
	wg.Wait()
}
