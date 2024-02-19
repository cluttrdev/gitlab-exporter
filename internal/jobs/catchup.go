package jobs

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
	"github.com/cluttrdev/gitlab-exporter/internal/exporter"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
	"github.com/cluttrdev/gitlab-exporter/internal/tasks"
	"github.com/cluttrdev/gitlab-exporter/pkg/worker"
)

type ProjectCatchUpJob struct {
	Config   config.Project
	GitLab   *gitlab.Client
	Exporter *exporter.Exporter

	WorkerPool *worker.Pool
}

func (j *ProjectCatchUpJob) Run(ctx context.Context) {
	opt := gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},

		Scope: &[]string{"finished"}[0],
	}
	if j.Config.CatchUp.UpdatedAfter != "" {
		after, err := time.Parse("2006-01-02T15:04:05Z", j.Config.CatchUp.UpdatedAfter)
		if err != nil {
			slog.Error("error parsing catchup update_after", "error", err)
		} else {
			opt.UpdatedAfter = &after
		}
	}
	if j.Config.CatchUp.UpdatedBefore != "" {
		before, err := time.Parse("2006-01-02T15:04:05Z", j.Config.CatchUp.UpdatedBefore)
		if err != nil {
			slog.Error("error parsing catchup update_before", "error", err)
		} else {
			opt.UpdatedBefore = &before
		}
	}

	pipelines := j.GitLab.ListProjectPipelines(ctx, j.Config.Id, opt)
	var wg sync.WaitGroup
	for r := range pipelines {
		if r.Error != nil && !errors.Is(r.Error, context.Canceled) {
			slog.Error("error listing project pipelines", "error", r.Error)
			continue
		}

		pipelineID := r.Pipeline.Id
		wg.Add(1)
		j.WorkerPool.Submit(func(ctx context.Context) {
			defer wg.Done()
			if err := j.export(ctx, pipelineID); err != nil {
				slog.Error(err.Error())
			}
		})
	}

	wg.Wait()
}

func (j *ProjectCatchUpJob) export(ctx context.Context, pipelineID int64) error {
	options := tasks.ExportPipelineHierarchyOptions{
		ProjectID:  j.Config.Id,
		PipelineID: pipelineID,

		ExportSections:           j.Config.Export.Sections.Enabled,
		ExportTestReports:        j.Config.Export.TestReports.Enabled,
		ExportTraces:             j.Config.Export.Traces.Enabled,
		ExportLogEmbeddedMetrics: j.Config.Export.LogEmbeddedMetrics.Enabled,
	}

	return tasks.ExportPipelineHierarchy(ctx, j.GitLab, j.Exporter, options)
}
