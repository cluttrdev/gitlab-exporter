package worker

import (
	"context"
	"log"
	"sync"
	"time"

	clickhouse "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/clickhouse"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
	gitlab "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/gitlab"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/tasks"
)

type exportProjectWorker struct {
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}

	// ensure the worker can only be started once
	start sync.Once
	// ensure the worker can only be stopped once
	stop sync.Once

	project    *config.Project
	gitlab     *gitlab.Client
	clickhouse *clickhouse.Client
}

func NewExportProjectWorker(cfg *config.Project, gl *gitlab.Client, ch *clickhouse.Client) Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &exportProjectWorker{
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),

		project:    cfg,
		gitlab:     gl,
		clickhouse: ch,
	}
}

func (w *exportProjectWorker) Start() {
	go func() {
		w.start.Do(func() {
			w.run()
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

func (w *exportProjectWorker) run() {
	interval := 60 * time.Second

	opt := &gitlab.ListProjectPipelineOptions{
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
		case <-w.ctx.Done():
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
			for r := range w.gitlab.ListProjectPipelines(w.ctx, w.project.Id, opt) {
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

					if err := tasks.ExportPipelineHierarchy(w.ctx, &opts, w.gitlab, w.clickhouse); err != nil {
						log.Printf("error exporting pipeline hierarchy: %s\n", err)
					}
					log.Printf("Exported projects/%d/pipelines/%d\n", opts.ProjectID, opts.PipelineID)
				}(r.Pipeline.ID)
			}
			wg.Wait()
		}
	}
}
