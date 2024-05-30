package jobs

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	_gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
	"github.com/cluttrdev/gitlab-exporter/internal/exporter"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
	"github.com/cluttrdev/gitlab-exporter/internal/tasks"
	"github.com/cluttrdev/gitlab-exporter/pkg/worker"
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type ProjectCatchUpJob struct {
	Config   config.Project
	GitLab   *gitlab.Client
	Exporter *exporter.Exporter

	WorkerPool *worker.Pool
}

func (j *ProjectCatchUpJob) Run(ctx context.Context) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		j.exportProjectPipelines(ctx, &wg)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		j.exportProjectMergeRequests(ctx)
	}()

	wg.Wait()
}

func (j *ProjectCatchUpJob) exportProjectPipelines(ctx context.Context, wg *sync.WaitGroup) {
	projectID := j.Config.Id

	opt := gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},

		Scope: gitlab.Ptr("finished"),
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
	for r := range pipelines {
		if r.Error != nil {
			if errors.Is(r.Error, context.Canceled) {
				return
			} else {
				slog.Error("error listing project pipelines", "error", r.Error)
			}
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

func (j *ProjectCatchUpJob) exportProjectMergeRequests(ctx context.Context) {
	projectID := j.Config.Id

	opt := _gitlab.ListProjectMergeRequestsOptions{
		ListOptions: _gitlab.ListOptions{
			Pagination: "keyset",
			PerPage:    100,
			OrderBy:    "updated_at",
			Sort:       "desc",
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

	options := []_gitlab.RequestOptionFunc{
		_gitlab.WithContext(ctx),
	}

	var wg sync.WaitGroup
	for {
		mrs, resp, err := j.GitLab.Client().MergeRequests.ListProjectMergeRequests(int(projectID), &opt, options...)
		if err != nil {
			slog.Error("error fetching project merge requests", "project", projectID, "error", err)
			break
		}

		iids := make([]int, 0, len(mrs))
		for _, mr := range mrs {
			iids = append(iids, mr.IID)
		}

		wg.Add(1)
		j.WorkerPool.Submit(func(ctx context.Context) {
			defer wg.Done()
			mergerequests := make([]*typespb.MergeRequest, 0, len(iids))

			opt := _gitlab.GetMergeRequestsOptions{}
			for _, iid := range iids {
				mr, _, err := j.GitLab.Client().MergeRequests.GetMergeRequest(int(projectID), iid, &opt, _gitlab.WithContext(ctx))
				if err != nil {
					if errors.Is(err, context.Canceled) {
						break
					}
					slog.Error(err.Error())
					continue
				}

				mergerequests = append(mergerequests, gitlab.ConvertMergeRequest(mr))
			}

			if len(mergerequests) == 0 {
				return
			}

			if err := j.Exporter.ExportMergeRequests(ctx, mergerequests); err != nil {
				slog.Error(err.Error())
			}
		})

		if resp.NextLink == "" {
			break
		}

		options = []_gitlab.RequestOptionFunc{
			_gitlab.WithContext(ctx),
			_gitlab.WithKeysetPaginationParameters(resp.NextLink),
		}
	}
	wg.Wait()
}
