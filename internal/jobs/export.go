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

type ProjectExportJob struct {
	Config   config.Project
	GitLab   *gitlab.Client
	Exporter *exporter.Exporter

	WorkerPool *worker.Pool
}

func (j *ProjectExportJob) Run(ctx context.Context) {
	period := 1 * time.Minute

	opt := gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},

		Scope: &[]string{"finished"}[0],
	}

	now := time.Now().UTC()
	opt.UpdatedBefore = &now

	ticker := time.NewTicker(period)
	var iteration int = 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			iteration++
			slog.Debug("Exporting updated pipelines", "project", j.Config.Id, "iteration", iteration)

			now := time.Now().UTC()
			opt.UpdatedAfter = opt.UpdatedBefore
			opt.UpdatedBefore = &now

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
	}
}

func (j *ProjectExportJob) export(ctx context.Context, pipelineID int64) error {
	options := tasks.ExportPipelineHierarchyOptions{
		ProjectID:  j.Config.Id,
		PipelineID: pipelineID,

		ExportSections:    j.Config.Export.Sections.Enabled,
		ExportTestReports: j.Config.Export.TestReports.Enabled,
		ExportTraces:      j.Config.Export.Traces.Enabled,
		ExportMetrics:     j.Config.Export.Metrics.Enabled,
	}

	return tasks.ExportPipelineHierarchy(ctx, j.GitLab, j.Exporter, options)
}
