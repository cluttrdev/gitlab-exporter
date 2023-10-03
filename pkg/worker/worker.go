package worker

import (
	"context"
)

type Worker interface {
	Start()
	Stop()
	Done() <-chan struct{}
}

type worker struct {
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}

	run func(*worker)
}

func (w *worker) Start() {
	w.ctx, w.cancel = context.WithCancel(context.Background())
	go w.run(w)
}

func (w *worker) Stop() {
	w.cancel()
	close(w.done)
}

func (w *worker) Done() <-chan struct{} {
	return w.done
}

func newWorker(run func(*worker)) Worker {
	return &worker{
		done: make(chan struct{}),
		run:  run,
	}
}
