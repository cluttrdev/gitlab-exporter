package recorder

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

	"go.cluttr.dev/gitlab-exporter-clickhouse-recorder/internal/retry"
)

func (r *ClickHouseRecorder) CheckReadiness(ctx context.Context) error {
	if err := r.client.Ping(ctx); err != nil {
		return fmt.Errorf("error pinging clickhouse: %w", err)
	}

	return nil
}

func (r *ClickHouseRecorder) WaitForReady(ctx context.Context) error {
	seconds := func(d time.Duration) time.Duration {
		s := math.Ceil(d.Seconds())
		return time.Duration(s) * time.Second
	}

	return retry.Do(
		func(ctx context.Context) error {
			err := r.CheckReadiness(ctx)
			if err != nil {
				args := []any{
					"error", err,
				}

				v, ok := ctx.Value(retry.ContextValuesKey("retry")).(retry.ContextValues)
				if ok {
					args = append(args,
						"retry.attempt", v.Attempt,
						"retry.delay", fmt.Sprint(seconds(v.Delay)),
					)
				}

				slog.Error("Readiness check failed", args...)
			}
			return err
		},
		// as long as context is not done
		retry.WithContext(ctx),
		// with unlimited attempts
		retry.MaxAttempts(0),
		// with exponential backoff
		retry.WithBackoff(retry.Backoff{
			InitialDelay: 1 * time.Second,
			MaxDelay:     5 * time.Minute,
			Factor:       2.0,
			Jitter:       0.1, // +/- 10%
		}),
	)
}

func (r *ClickHouseRecorder) GetReady(ctx context.Context) error {
	slog.Debug("Checking readiness...")
	if err := r.WaitForReady(ctx); err != nil {
		return err
	}
	slog.Debug("Checking readiness... done")

	return nil
}

func (r *ClickHouseRecorder) WatchReadiness(ctx context.Context) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		var err error
		for {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
			default:
				err = r.CheckReadiness(ctx)
				errChan <- err
				if err == nil { // everything okay
					break /* select */
				}

				// readiness check failed, waiting for it to succeed again
				err = r.WaitForReady(ctx)
				errChan <- err
				if err != nil {
					// failed to get ready again, aborting
					return
				}
			}

			time.Sleep(3 * time.Second)
		}
	}()

	return errChan
}
