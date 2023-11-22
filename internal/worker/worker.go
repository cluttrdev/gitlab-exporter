package worker

import (
	"context"
)

type Worker interface {
	Start(context.Context)
	Stop()
	Done() <-chan struct{}
}

type worker struct {
	cancel context.CancelFunc
	done   chan struct{}

	run func(context.Context)
}

func (w *worker) Start(ctx context.Context) {
	ctx, w.cancel = context.WithCancel(ctx)
	go w.run(ctx)
}

func (w *worker) Stop() {
	w.cancel()
	close(w.done)
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
