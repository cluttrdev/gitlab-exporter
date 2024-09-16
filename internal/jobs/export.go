package jobs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	_gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
	"github.com/cluttrdev/gitlab-exporter/internal/exporter"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
	"github.com/cluttrdev/gitlab-exporter/internal/tasks"
	"github.com/cluttrdev/gitlab-exporter/internal/types"
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
				j.exportProject(ctx, iteration == 1)
			})

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
			j.lastUpdate = time.Now().UTC()
		}
	}
}

func (j *ProjectExportJob) exportProject(ctx context.Context, first bool) {
	projectID := j.Config.Id

	project, err := j.GitLab.GetProject(ctx, projectID)
	if err != nil {
		slog.Error("error fetching project", "project", projectID, "error", err)
		return
	} else if !project.LastActivityAt.AsTime().After(j.lastUpdate) && !first {
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

		Scope:        gitlab.Ptr("finished"),
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

func (j *ProjectExportJob) exportProjectMergeRequests(ctx context.Context) {
	projectID := int(j.Config.Id)
	glab := j.GitLab.Client()

	opt := _gitlab.ListProjectMergeRequestsOptions{
		ListOptions: _gitlab.ListOptions{
			Pagination: "keyset",
			PerPage:    100,
			OrderBy:    "updated_at",
			Sort:       "desc",
		},
		View: _gitlab.Ptr("simple"),

		UpdatedAfter: &j.lastUpdate,
	}

	options := []_gitlab.RequestOptionFunc{
		_gitlab.WithContext(ctx),
	}

	var wg sync.WaitGroup
	for {
		// get iids of updated merge requests
		mrs, resp, err := glab.MergeRequests.ListProjectMergeRequests(projectID, &opt, options...)
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
			mrNoteEvents := make([]*typespb.MergeRequestNoteEvent, 0, len(iids))

			opt := _gitlab.GetMergeRequestsOptions{}
			for _, iid := range iids {
				mr, _, err := glab.MergeRequests.GetMergeRequest(projectID, iid, &opt, _gitlab.WithContext(ctx))
				if err != nil {
					if errors.Is(err, context.Canceled) {
						break
					}
					slog.Error("error fetching merge request", "project_id", projectID, "iid", iid)
					continue
				}

				mergerequests = append(mergerequests, types.ConvertMergeRequest(mr))

				if j.Config.Export.MergeRequests.NoteEvents {
					notes, err := tasks.FetchMergeRequestNotes(ctx, glab, projectID, iid)
					if err != nil {
						if errors.Is(err, context.Canceled) {
							break
						}
						slog.Error("error fetching merge request note events", "project_id", projectID, "iid", iid)
						continue
					}

					for _, note := range notes {
						if ev := types.ConvertToMergeRequestNoteEvent(note); ev != nil {
							mrNoteEvents = append(mrNoteEvents, ev)
						}
					}
				}
			}

			if len(mergerequests) > 0 {
				if err := j.Exporter.ExportMergeRequests(ctx, mergerequests); err != nil {
					slog.Error(fmt.Sprintf("error exporting merge requests: %v", err))
				}
			}

			if len(mrNoteEvents) > 0 {
				if err := j.Exporter.ExportMergeRequestNoteEvents(ctx, mrNoteEvents); err != nil {
					slog.Error(fmt.Sprintf("error exporting merge request note events: %v", err))
				}
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

	// wait for all paginated exports to finish
	wg.Wait()
}
