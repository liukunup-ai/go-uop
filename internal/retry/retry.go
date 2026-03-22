package retry

import (
	"context"
	"fmt"
	"time"
)

type Config struct {
	MaxAttempts int
	Delay       time.Duration
	MaxDelay    time.Duration
	Backoff     float64
}

func DefaultConfig() Config {
	return Config{
		MaxAttempts: 3,
		Delay:       100 * time.Millisecond,
		MaxDelay:    30 * time.Second,
		Backoff:     2.0,
	}
}

type RetryableFunc func() error

func Do(fn RetryableFunc, config Config) error {
	var err error
	delay := config.Delay

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}

		if attempt >= config.MaxAttempts {
			break
		}

		time.Sleep(delay)

		delay = time.Duration(float64(delay) * config.Backoff)
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}
	}

	return fmt.Errorf("after %d attempts: %w", config.MaxAttempts, err)
}

func DoWithContext(ctx context.Context, fn func(context.Context) error, config Config) error {
	var err error
	delay := config.Delay

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		err = fn(ctx)
		if err == nil {
			return nil
		}

		if attempt >= config.MaxAttempts {
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}

		delay = time.Duration(float64(delay) * config.Backoff)
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}
	}

	return fmt.Errorf("after %d attempts: %w", config.MaxAttempts, err)
}

func DoUntilSuccess(fn RetryableFunc, maxDuration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), maxDuration)
	defer cancel()

	var err error
	delay := 100 * time.Millisecond
	maxDelay := maxDuration / 2
	if maxDelay < 100*time.Millisecond {
		maxDelay = 100 * time.Millisecond
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout after %s: %w", maxDuration, err)
		default:
		}

		err = fn()
		if err == nil {
			return nil
		}

		time.Sleep(delay)
		if delay < maxDelay {
			delay *= 2
		}
	}
}

type Result[T any] struct {
	Value T
	Err   error
}

func DoWithResult[T any](fn func() (T, error), config Config) (T, error) {
	var zero T
	var err error
	delay := config.Delay

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		var val T
		val, err = fn()
		if err == nil {
			return val, nil
		}

		if attempt >= config.MaxAttempts {
			break
		}

		time.Sleep(delay)
		delay = time.Duration(float64(delay) * config.Backoff)
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}
	}

	return zero, fmt.Errorf("after %d attempts: %w", config.MaxAttempts, err)
}
