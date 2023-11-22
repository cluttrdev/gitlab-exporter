package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/gitlab"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/tasks"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/internal/datastore"
)

type exportProjectWorker struct {
	cancel context.CancelFunc
	done   chan struct{}

	// ensure the worker can only be started once
	start sync.Once
	// ensure the worker can only be stopped once
	stop sync.Once

	project   config.Project
	gitlab    *gitlab.Client
	datastore datastore.DataStore
}

func NewExportProjectWorker(cfg config.Project, gl *gitlab.Client, ds datastore.DataStore) Worker {
	return &exportProjectWorker{
		done: make(chan struct{}),

		project:   cfg,
		gitlab:    gl,
		datastore: ds,
	}
}

func (w *exportProjectWorker) Start(ctx context.Context) {
	ctx, w.cancel = context.WithCancel(ctx)
	go func() {
		w.start.Do(func() {
			w.run(ctx)
		})
	}()
}

func (w *exportProjectWorker) Stop() {
	w.stop.Do(func() {
		w.cancel()
		close(w.done)
	})
}

func (w *exportProjectWorker) Done() <-chan struct{} {
	return w.done
}

func (w *exportProjectWorker) run(ctx context.Context) {
	interval := 60 * time.Second

	opt := gitlab.ListProjectPipelineOptions{
		PerPage: 100,
		Page:    1,

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
			for r := range w.gitlab.ListProjectPipelines(ctx, w.project.Id, opt) {
				if r.Error != nil {
					log.Println(r.Error)
					continue
				}

				wg.Add(1)
				go func(pid int64) {
					defer wg.Done()

					opts := tasks.ExportPipelineHierarchyOptions{
						ProjectID:  w.project.Id,
						PipelineID: pid,

						ExportSections:    w.project.Export.Sections.Enabled,
						ExportTestReports: w.project.Export.TestReports.Enabled,
						ExportTraces:      w.project.Export.Traces.Enabled,
					}

					if err := tasks.ExportPipelineHierarchy(ctx, opts, w.gitlab, w.datastore); err != nil {
						log.Printf("error exporting pipeline hierarchy: %s\n", err)
					} else {
						log.Printf("Exported projects/%d/pipelines/%d\n", opts.ProjectID, opts.PipelineID)
					}
				}(r.Pipeline.ID)
			}
			wg.Wait()
		}
	}
}
