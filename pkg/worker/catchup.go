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
	cancel context.CancelFunc
	done   chan struct{}

	// ensure the worker can only be started once
	start sync.Once
	// ensure the worker can only be stopped once
	stop sync.Once

	project    config.Project
	gitlab     *gitlab.Client
	clickhouse *clickhouse.Client
}

func NewCatchUpProjectWorker(cfg config.Project, gl *gitlab.Client, ch *clickhouse.Client) Worker {
	return &catchUpProjectWorker{
		done: make(chan struct{}),

		project:    cfg,
		gitlab:     gl,
		clickhouse: ch,
	}
}

func (w *catchUpProjectWorker) Start(ctx context.Context) {
	ctx, w.cancel = context.WithCancel(ctx)
	go func() {
		w.start.Do(func() {
			w.run(ctx)
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

func (w *catchUpProjectWorker) run(ctx context.Context) {
	opt := gitlab.ListProjectPipelineOptions{
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

	ch := w.produce(ctx, opt)
	w.process(ctx, ch)
}

func (w *catchUpProjectWorker) produce(ctx context.Context, opt gitlab.ListProjectPipelineOptions) <-chan int64 {
	ch := make(chan int64)

	go func() {
		defer close(ch)

		latestUpdates, err := w.clickhouse.QueryProjectPipelinesLatestUpdate(ctx, w.project.Id)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Println(err)
		}

		resChan := w.gitlab.ListProjectPipelines(ctx, w.project.Id, opt)
		for {
			select {
			case <-ctx.Done():
				return
			case r, ok := <-resChan:
				if !ok { // channel closed
					return
				}

				if r.Error != nil && !errors.Is(r.Error, context.Canceled) {
					log.Println(r.Error)
					continue
				}

				if !w.project.CatchUp.Forced {
					// if not forced, skip pipelines that have not been updated
					lastUpdatedAt, ok := latestUpdates[r.Pipeline.ID]
					if ok && r.Pipeline.UpdatedAt.Compare(lastUpdatedAt) <= 0 {
						continue
					}
				}

				ch <- r.Pipeline.ID
			}
		}
	}()

	return ch
}

func (w *catchUpProjectWorker) process(ctx context.Context, pipelineChan <-chan int64) {
	numWorkers := 10
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for pipelineID := range pipelineChan {
				opts := tasks.ExportPipelineHierarchyOptions{
					ProjectID:  w.project.Id,
					PipelineID: pipelineID,

					ExportSections:    w.project.Export.Sections.Enabled,
					ExportTestReports: w.project.Export.TestReports.Enabled,
					ExportTraces:      w.project.Export.Traces.Enabled,
				}

				if err := tasks.ExportPipelineHierarchy(ctx, opts, w.gitlab, w.clickhouse); err != nil {
					if !errors.Is(err, context.Canceled) {
						log.Printf("error exporting pipeline hierarchy: %s\n", err)
					}
				} else {
					log.Printf("Caught up on projects/%d/pipelines/%d\n", opts.ProjectID, opts.PipelineID)
				}
			}
		}()
	}
	wg.Wait()
}
