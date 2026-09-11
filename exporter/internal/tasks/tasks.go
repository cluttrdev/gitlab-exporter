package tasks

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab"
)

// outcome pairs a single unit of work's return value with its error so both can
// travel together over the results channel.
type outcome[R any] struct {
	value R
	err   error
}

type fanOutOptions struct {
	CollectOnError bool // append value even when work returns an error
}

// fanOut runs work for each item concurrently, bounded by glab's semaphore.
// Successful results are collected (order not guaranteed); errors are joined.
//
// T is the input item type, R is the per-item result type. Callers that produce
// several slices per item bundle them into a struct R and reduce afterwards.
func fanOut[T, R any](
	ctx context.Context,
	glab *gitlab.Client,
	items []T,
	work func(ctx context.Context, glab *gitlab.Client, item T) (R, error),
	opts fanOutOptions,
) ([]R, error) {
	var (
		wg sync.WaitGroup
		// Unbuffered on purpose: the worker's send blocks until the collector
		// receives, so wg.Done() (deferred, runs after the send) only completes
		// once the value is safely collected. This makes it impossible for `done`
		// to close while results remain unprocessed.
		results = make(chan outcome[R])
	)

	// Producer goroutine: acquire a semaphore slot per item, then spawn a worker.
	// It is itself tracked by the WaitGroup so `done` only closes after both
	// the producer and every worker have finished.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, item := range items {
			// Acquire blocks until a slot frees up; it errors only when ctx is
			// cancelled. Since that condition persists for all remaining items,
			// we stop the loop entirely rather than retry each one.
			if err := glab.Acquire(ctx, 1); err != nil {
				slog.Error("failed to acquire gitlab client", "error", err)
				break
			}
			wg.Add(1)
			go func() {
				defer glab.Release(1)
				defer wg.Done()
				v, err := work(ctx, glab, item)
				results <- outcome[R]{value: v, err: err}
			}()
		}
	}()

	// Closer goroutine: closes `done` once the producer and all workers finish.
	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	var (
		values []R
		errs   error
	)
loop:
	for {
		select {
		case <-done:
			// All producers/workers finished AND results is drained enough that
			// no more sends will arrive; safe to stop.
			break loop
		case r := <-results:
			if r.err != nil {
				errs = errors.Join(errs, r.err)
				// Success-only sites skip the value on error; collectOnError sites
				// fall through and append it.
				if !opts.CollectOnError {
					continue
				}
			}
			values = append(values, r.value)
		}
	}
	return values, errs
}
