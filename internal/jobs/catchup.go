package jobs

import (
	"context"
	"log/slog"
	"sync"
	"time"

	_gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
	"github.com/cluttrdev/gitlab-exporter/internal/exporter"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
	"github.com/cluttrdev/gitlab-exporter/internal/tasks"
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type ProjectCatchUpJob struct {
	Config   config.Project
	GitLab   *gitlab.Client
	Exporter *exporter.Exporter

	WorkerPool WorkerPool
}

func (j *ProjectCatchUpJob) Run(ctx context.Context) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		j.exportProject(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		j.exportProjectPipelines(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		j.exportProjectMergeRequests(ctx)
	}()

	wg.Wait()
}

func (j *ProjectCatchUpJob) submit(task func()) bool {
	if j.WorkerPool.Stopped() {
		return false
	}

	j.WorkerPool.Submit(task)
	return true
}

func (j *ProjectCatchUpJob) exportProject(ctx context.Context) {
	projectID := j.Config.Id

	_ = j.submit(func() {
		project, err := j.GitLab.GetProject(ctx, projectID)
		if err != nil {
			slog.Error("error getting project", "project_id", projectID, "error", err)
			return
		}

		if err := j.Exporter.ExportProjects(ctx, []*typespb.Project{project}); err != nil {
			slog.Error("error exporting project", "project_id", projectID, "error", err)
		}
	})
}

func (j *ProjectCatchUpJob) exportProjectPipelines(ctx context.Context) {
	projectID := j.Config.Id
	glab := j.GitLab.Client()

	opt := _gitlab.ListProjectPipelinesOptions{
		ListOptions: _gitlab.ListOptions{
			PerPage: 100,
			OrderBy: "updated_at",
			Sort:    "desc",
		},
		Scope: _gitlab.Ptr("finished"),
	}
	if j.Config.CatchUp.UpdatedAfter != "" {
		after, err := time.Parse("2006-01-02T15:04:05Z", j.Config.CatchUp.UpdatedAfter)
		if err != nil {
			slog.Error("error parsing catchup update_after", "error", err)
			return
		} else {
			opt.UpdatedAfter = &after
		}
	}
	if j.Config.CatchUp.UpdatedBefore != "" {
		before, err := time.Parse("2006-01-02T15:04:05Z", j.Config.CatchUp.UpdatedBefore)
		if err != nil {
			slog.Error("error parsing catchup update_before", "error", err)
			return
		} else {
			opt.UpdatedBefore = &before
		}
	}

	err := gitlab.ListProjectPipelines(ctx, glab, projectID, opt, func(pipelines []*_gitlab.PipelineInfo) bool {
		for _, pipeline := range pipelines {
			pipelineID := int64(pipeline.ID)
			submitted := j.submit(func() {
				err := tasks.ExportPipelineHierarchy(ctx, j.GitLab, j.Exporter, tasks.ExportPipelineHierarchyOptions{
					ProjectID:  projectID,
					PipelineID: pipelineID,

					ExportSections:    j.Config.Export.Sections.Enabled,
					ExportTestReports: j.Config.Export.TestReports.Enabled,
					ExportTraces:      j.Config.Export.Traces.Enabled,
					ExportMetrics:     j.Config.Export.Metrics.Enabled,
				})
				if err != nil {
					slog.Error("error exporting pipeline hierarchy", "project_id", projectID, "pipeline_id", pipelineID, "error", err)
				}
			})

			if !submitted {
				return false
			}
		}

		return true
	})

	if err != nil {
		slog.Error("error listing project pipelines", "project_id", projectID, "error", err)
	}
}

func (j *ProjectCatchUpJob) exportProjectMergeRequests(ctx context.Context) {
	projectID := int(j.Config.Id)
	glab := j.GitLab.Client()
	exp := j.Exporter

	opt := _gitlab.ListProjectMergeRequestsOptions{
		ListOptions: _gitlab.ListOptions{
			PerPage: 100,
			OrderBy: "updated_at",
			Sort:    "desc",
		},
		View: _gitlab.Ptr("simple"),
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

	err := gitlab.ListProjectMergeRequests(ctx, glab, int64(projectID), opt, func(mrs []*_gitlab.MergeRequest) bool {
		iids := make([]int, 0, len(mrs))
		for _, mr := range mrs {
			iids = append(iids, mr.IID)
		}

		submitted := j.submit(func() {
			opt := tasks.ExportProjectMergeRequestsOptions{
				ProjectID:        projectID,
				MergeRequestIIDs: iids,

				ExportNoteEvents: j.Config.Export.MergeRequests.NoteEvents,
			}
			if err := tasks.ExportProjectMergeRequests(ctx, glab, exp, opt); err != nil {
				slog.Error("error exporting project merge requests", "project_id", projectID, "error", err)
			}
		})

		return submitted
	})

	if err != nil {
		slog.Error("error listing project merge requests", "project_id", projectID, "error", err)
	}
}
