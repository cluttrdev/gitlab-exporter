package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/controller"
	gitlab "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/gitlab"
)

type exportProjectWorker struct {
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}

	// ensure the worker can only be started once
	start sync.Once
	// ensure the worker can only be stopped once
	stop sync.Once

	ctl *controller.Controller
	cfg *config.Project
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

	ctx := w.ctx
	ctl := w.ctl
	cfg := w.cfg

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
			for r := range ctl.GitLab.ListProjectPipelines(ctx, cfg.Id, opt) {
				if r.Error != nil {
					log.Println(r.Error)
					continue
				}

				wg.Add(1)
				go func(pid int64) {
					defer wg.Done()

					if err := ctl.ExportPipeline(ctx, cfg.Id, pid); err != nil {
						log.Printf("error exporting pipeline: %s\n", err)
					}
					log.Printf("Exported projects/%d/pipelines/%d\n", cfg.Id, pid)
				}(r.Pipeline.ID)
			}
			wg.Wait()
		}
	}
}

func NewExportProjectWorker(ctl *controller.Controller, cfg *config.Project) Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &exportProjectWorker{
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),

		ctl: ctl,
		cfg: cfg,
	}
}
