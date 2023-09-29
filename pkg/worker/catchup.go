package worker

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/controller"
	gitlab "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/gitlab"
)

type catchUpProjectWorker struct {
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

	ctx := w.ctx
	ctl := w.ctl
	cfg := w.cfg

	opt := &gitlab.ListProjectPipelineOptions{
		PerPage: 100,
		Page:    1,

		Scope: &[]string{"finished"}[0],
	}
	if cfg.CatchUp.UpdatedAfter != "" {
		after, err := time.Parse("2006-01-02T15:04:05Z", cfg.CatchUp.UpdatedAfter)
		if err != nil {
			log.Println(err)
		} else {
			opt.UpdatedAfter = &after
		}
	}
	if cfg.CatchUp.UpdatedBefore != "" {
		before, err := time.Parse("2006-01-02T15:04:05Z", cfg.CatchUp.UpdatedBefore)
		if err != nil {
			log.Println(err)
		} else {
			opt.UpdatedBefore = &before
		}
	}

	for i := 0; i < retries; i++ {
		latestUpdates, err := ctl.QueryLatestProjectPipelineUpdates(ctx, cfg.Id)
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Println(err)
			continue
		}

		for r := range ctl.GitLab.ListProjectPipelines(ctx, cfg.Id, opt) {
			if r.Error != nil {
				if errors.Is(r.Error, context.Canceled) {
					return
				}
				log.Println(r.Error)
				continue
			}
			pi := r.Pipeline

			lastUpdatedAt, ok := latestUpdates[pi.ID]
			if ok && pi.UpdatedAt.Compare(lastUpdatedAt) <= 0 {
				continue
			}

			if err := ctl.ExportPipeline(ctx, cfg.Id, pi.ID); err != nil {
				if errors.Is(r.Error, context.Canceled) {
					return
				}
				log.Printf("error exporting pipeline: %s\n", err)
			}
			log.Printf("Caught up on projects/%d/pipelines/%d\n", cfg.Id, pi.ID)
		}
	}
}

func NewCatchUpProjectWorker(ctl *controller.Controller, cfg *config.Project) Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &catchUpProjectWorker{
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),

		ctl: ctl,
		cfg: cfg,
	}
}
