package retry

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"
)

type ContextValuesKey string

type ContextValues struct {
	Attempt int
	Delay   time.Duration
}

func Do(fn func(ctx context.Context) error, opts ...Option) error {
	fn_ := func(ctx context.Context) (any, error) {
		return nil, fn(ctx)
	}

	_, err := DoWithData(fn_, opts...)
	return err
}

func DoWithData[T any](fn func(ctx context.Context) (T, error), opts ...Option) (T, error) {
	var t T

	cfg := defaultConfig()
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return t, err
		}
	}

	ctx := context.WithValue(
		cfg.context,
		ContextValuesKey("retry"),
		ContextValues{
			Attempt: 0,
			Delay:   cfg.backoff.InitialDelay,
		},
	)

	var (
		err     error
		attempt int           = 0
		delay   time.Duration = cfg.backoff.InitialDelay
	)

	ticker := time.NewTicker(delay)
	defer ticker.Stop()
loop:
	for {
		t, err = fn(ctx)
		if err == nil {
			return t, nil
		}

		if !cfg.retryIf(err) {
			break loop
		}

		// don't wait if this was the last attempt
		if cfg.maxAttempts > 0 && attempt == cfg.maxAttempts-1 {
			break loop
		}

		ticker.Reset(delay)
		select {
		case <-ticker.C:
			attempt++
			delay = cfg.backoff.Delay(attempt)

			ctx = context.WithValue(
				ctx,
				ContextValuesKey("retry"),
				ContextValues{
					Attempt: attempt,
					Delay:   delay,
				},
			)
			continue loop
		case <-cfg.context.Done():
			err = cfg.context.Err()
			break loop
		}
	}

	return t, err
}

type Config struct {
	maxAttempts int
	retryIf     func(error) bool
	backoff     Backoff
	context     context.Context
}

func defaultConfig() Config {
	return Config{
		maxAttempts: 5,
		retryIf:     func(error) bool { return true },
		backoff: Backoff{
			InitialDelay: 1 * time.Second,
			Factor:       2.0,
			Jitter:       0.1,
			MaxDelay:     120 * time.Second,
		},
		context: context.Background(),
	}
}

type Option func(*Config) error

func MaxAttempts(n int) Option {
	return func(c *Config) error {
		if n < 0 {
			return errors.New("Maximum retries must be non-negative")
		}
		c.maxAttempts = n
		return nil
	}
}

func RetryIf(fn func(error) bool) Option {
	return func(c *Config) error {
		c.retryIf = fn
		return nil
	}
}

func WithContext(ctx context.Context) Option {
	return func(c *Config) error {
		if ctx == nil {
			return errors.New("Context must not be nil")
		}
		c.context = ctx
		return nil
	}
}

func WithBackoff(b Backoff) Option {
	return func(c *Config) error {
		c.backoff = b
		return nil
	}
}

type Backoff struct {
	// How long to wait before first retry
	InitialDelay time.Duration
	// Upper bound on backoff
	MaxDelay time.Duration
	// Factor with which to multiply backoff after a failed retry
	Factor float64
	// By how much to randomize backoff
	Jitter float64
}

func (b *Backoff) Delay(attempt int) time.Duration {
	if b == nil {
		return time.Duration(0)
	}

	if attempt < 0 {
		return b.InitialDelay
	}

	delay := math.Pow(b.Factor, float64(attempt)) * float64(b.InitialDelay)
	if delay > float64(b.MaxDelay) {
		delay = float64(b.MaxDelay)
	}

	// if jitter is 0.1, then multiply delay by random value in [0.9, 1.1)
	r := -1 + 2*rand.Float64() // pseudo-random number in [-1, 1)
	delay = delay * (1 + b.Jitter*r)

	return time.Duration(delay)
}
