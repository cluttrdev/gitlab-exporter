package worker

import (
	"context"
)

type Pool struct {
	workers []*worker

	queue chan func(context.Context)
}

func NewWorkerPool(size int) *Pool {
	pool := &Pool{
		workers: make([]*worker, 0, size),
		queue:   make(chan func(context.Context), size),
	}

	for i := 0; i < size; i++ {
		pool.workers = append(
			pool.workers,
			&worker{
				done: make(chan struct{}),
				run: func(ctx context.Context) {
					for {
						select {
						case <-ctx.Done():
							return
						case task, ok := <-pool.queue:
							if !ok {
								return
							}
							task(ctx)
						}
					}
				},
			},
		)
	}

	return pool
}

func (p *Pool) Start(ctx context.Context) {
	for _, w := range p.workers {
		go w.Start(ctx)
	}
}

func (p *Pool) Stop() {
	for _, w := range p.workers {
		w.Stop()
		<-w.Done()
	}
}

func (p *Pool) Submit(task func(context.Context)) {
	p.queue <- task
}
