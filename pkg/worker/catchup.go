package worker

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	clickhouse "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/clickhouse"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
	gitlab "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/gitlab"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/tasks"
)

type catchUpProjectWorker struct {
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

func NewCatchUpProjectWorker(cfg *config.Project, gl *gitlab.Client, ch *clickhouse.Client) Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &catchUpProjectWorker{
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),

		project:    cfg,
		gitlab:     gl,
		clickhouse: ch,
	}
}

func (w *catchUpProjectWorker) Start() {
	go func() {
		w.start.Do(func() {
			w.run()
		})
	}()
}

func (w *catchUpProjectWorker) Stop() {
	w.stop.Do(func() {
		w.cancel()
		close(w.done)
	})
}

func (w *catchUpProjectWorker) Done() <-chan struct{} {
	return w.done
}

func (w *catchUpProjectWorker) run() {
	retries := 3

	opt := &gitlab.ListProjectPipelineOptions{
		PerPage: 100,
		Page:    1,

		Scope: &[]string{"finished"}[0],
	}
	if w.project.CatchUp.UpdatedAfter != "" {
		after, err := time.Parse("2006-01-02T15:04:05Z", w.project.CatchUp.UpdatedAfter)
		if err != nil {
			log.Println(err)
		} else {
			opt.UpdatedAfter = &after
		}
	}
	if w.project.CatchUp.UpdatedBefore != "" {
		before, err := time.Parse("2006-01-02T15:04:05Z", w.project.CatchUp.UpdatedBefore)
		if err != nil {
			log.Println(err)
		} else {
			opt.UpdatedBefore = &before
		}
	}

	for i := 0; i < retries; i++ {
		ch := w.produce(opt)

		w.process(ch)
	}
}

func (w *catchUpProjectWorker) produce(opt *gitlab.ListProjectPipelineOptions) <-chan int64 {
	ch := make(chan int64)

	go func() {
		defer close(ch)

		latestUpdates, err := w.clickhouse.QueryProjectPipelinesLatestUpdate(w.ctx, w.project.Id)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Println(err)
		}

		resChan := w.gitlab.ListProjectPipelines(w.ctx, w.project.Id, opt)
		for {
			select {
			case <-w.ctx.Done():
				return
			case r := <-resChan:
				if r.Error != nil && !errors.Is(r.Error, context.Canceled) {
					log.Println(r.Error)
					continue
				}

				lastUpdatedAt, ok := latestUpdates[r.Pipeline.ID]
				if ok && r.Pipeline.UpdatedAt.Compare(lastUpdatedAt) <= 0 {
					continue
				}

				ch <- r.Pipeline.ID
			}
		}
	}()

	return ch
}

func (w *catchUpProjectWorker) process(pipelineChan <-chan int64) {
	numWorkers := 10
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for pipelineID := range pipelineChan {
				opts := &tasks.ExportPipelineHierarchyOptions{
					ProjectID:  w.project.Id,
					PipelineID: pipelineID,

					ExportSections:    w.project.Export.Sections.Enabled,
					ExportTestReports: w.project.Export.TestReports.Enabled,
					ExportTraces:      w.project.Export.Traces.Enabled,
				}

				if err := tasks.ExportPipelineHierarchy(w.ctx, opts, w.gitlab, w.clickhouse); err != nil {
					if !errors.Is(err, context.Canceled) {
						log.Printf("error exporting pipeline hierarchy: %s\n", err)
					}
				}
				log.Printf("Caught up on projects/%d/pipelines/%d\n", opts.ProjectID, opts.PipelineID)
			}
		}()
	}
	wg.Wait()
}
