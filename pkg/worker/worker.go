package worker

import (
	"context"
	"sync"
)

type Worker interface {
	Start(context.Context)
	Stop()
	Done() <-chan struct{}
}

type worker struct {
	cancel context.CancelFunc
	done   chan struct{}

	// ensure the worker can only be started once
	start sync.Once
	// ensure the worker can only be stopped once
	stop sync.Once

	run func(context.Context)
}

func (w *worker) Start(ctx context.Context) {
	ctx, w.cancel = context.WithCancel(ctx)
	go func() {
		w.start.Do(func() {
			defer w.cancel()
			w.run(ctx)
		})
	}()
}

func (w *worker) Stop() {
	w.stop.Do(func() {
		w.cancel()
		close(w.done)
	})
}

func (w *worker) Done() <-chan struct{} {
	return w.done
}

func NewWorker(run func(context.Context)) Worker {
	return &worker{
		done: make(chan struct{}),
		run:  run,
	}
}
