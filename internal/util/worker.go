package util

import (
	"fmt"
	"sync"
)

type Task interface {
	// Execute performs the work
	Execute() error

	// OnError handles any error returned from Execute()
	OnError(error)
}

type Pool struct {
	numWorkers int
	tasks      chan Task

	// ensure the pool can only be started once
	start sync.Once
	// ensure the pool can only be stopped once
	stop sync.Once

	// close to signal the workers to stop working
	quit chan struct{}
}

var ErrNoWorkers = fmt.Errorf("attempting to create worker pool with less than 1 worker")
var ErrNegativeChannelSize = fmt.Errorf("attempting to create worker pool with a negative channel size")

func NewPool(numWorkers int, channelSize int) (*Pool, error) {
	if numWorkers <= 0 {
		return nil, ErrNoWorkers
	}
	if channelSize < 0 {
		return nil, ErrNegativeChannelSize
	}

	tasks := make(chan Task, channelSize)

	return &Pool{
		numWorkers: numWorkers,
		tasks:      tasks,

		start: sync.Once{},
		stop:  sync.Once{},

		quit: make(chan struct{}),
	}, nil
}

// Start gets the workerpool ready to process jobs, and should only be called once
func (p *Pool) Start() {
	p.start.Do(func() {
		p.startWorkers()
	})
}

// Stop stops the workerpool, tears down any required resources,
// and should only be called once
func (p *Pool) Stop() {
	p.stop.Do(func() {
		close(p.quit)
	})
}

// AddWork adds a task for the worker pool to process. It is only valid after
// Start() has been called and before Stop() has been called.
// If the channel buffer is full (or 0) and all workers are occupied, this will
// hang until work is consumed or Stop() is called.
func (p *Pool) AddWork(t Task) {
	select {
	case p.tasks <- t:
	case <-p.quit:
	}
}

// AddWorkNonBlocking adds work to the Pool and returns immediately
func (p *Pool) AddWorkNonBlocking(t Task) {
	go p.AddWork(t)
}

func (p *Pool) startWorkers() {
	for i := 0; i < p.numWorkers; i++ {
		go func() {
			for {
				select {
				case <-p.quit:
					return
				case task, ok := <-p.tasks:
					if !ok {
						return
					}

					if err := task.Execute(); err != nil {
						task.OnError(err)
					}
				}
			}
		}()
	}
}

type SimpleTask struct {
	execFunc func() error

	Error error
}

func NewTask(f func() error) *SimpleTask {
	return &SimpleTask{
		execFunc: f,
	}
}

func (t *SimpleTask) Execute() error {
	return t.execFunc()
}

func (t *SimpleTask) OnError(err error) {
	t.Error = err
}
