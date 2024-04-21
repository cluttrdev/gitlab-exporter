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
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type ProjectExportJob struct {
	Config   config.Project
	GitLab   *gitlab.Client
	Exporter *exporter.Exporter

	WorkerPool *worker.Pool

	lastUpdate time.Time
}

func (j *ProjectExportJob) Run(ctx context.Context) {
	period := 1 * time.Minute

	j.lastUpdate = time.Now().UTC()

	projectID := j.Config.Id

	ticker := time.NewTicker(period)
	var iteration int = 0
	var wg sync.WaitGroup
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			iteration++
			slog.Debug("Exporting project", "id", projectID, "iteration", iteration)

			wg.Add(1)
			j.WorkerPool.Submit(func(ctx context.Context) {
				defer wg.Done()
				j.exportProject(ctx)
			})

			j.exportProjectPipelines(ctx, &wg)

			wg.Wait()
			j.lastUpdate = time.Now().UTC()
		}
	}
}

func (j *ProjectExportJob) exportProject(ctx context.Context) {
	projectID := j.Config.Id

	project, err := j.GitLab.GetProject(ctx, projectID)
	if err != nil {
		slog.Error("error fetching project", "project", projectID, "error", err)
		return
	} else if !project.LastActivityAt.AsTime().After(j.lastUpdate) {
		return
	}

	if err := j.Exporter.ExportProjects(ctx, []*typespb.Project{project}); err != nil {
		slog.Error(err.Error())
	}
}

func (j *ProjectExportJob) exportProjectPipelines(ctx context.Context, wg *sync.WaitGroup) {
	projectID := j.Config.Id

	opt := gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},

		Scope:        &[]string{"finished"}[0],
		UpdatedAfter: &j.lastUpdate,
	}

	pipelines := j.GitLab.ListProjectPipelines(ctx, projectID, opt)
	for r := range pipelines {
		if r.Error != nil {
			if errors.Is(r.Error, context.Canceled) {
				return
			}
			slog.Error("error listing project pipelines", "error", r.Error)
			continue
		}

		pipelineID := r.Pipeline.Id
		wg.Add(1)
		j.WorkerPool.Submit(func(ctx context.Context) {
			defer wg.Done()
			err := tasks.ExportPipelineHierarchy(ctx, j.GitLab, j.Exporter, tasks.ExportPipelineHierarchyOptions{
				ProjectID:  projectID,
				PipelineID: pipelineID,

				ExportSections:    j.Config.Export.Sections.Enabled,
				ExportTestReports: j.Config.Export.TestReports.Enabled,
				ExportTraces:      j.Config.Export.Traces.Enabled,
				ExportMetrics:     j.Config.Export.Metrics.Enabled,
			})
			if err != nil {
				slog.Error(err.Error())
			}
		})
	}
}
