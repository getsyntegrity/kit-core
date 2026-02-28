package listeners

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSaveCheckpointWithRetry_NonRetryableErrors(t *testing.T) {
	t.Run("version not contiguous is not retried", func(t *testing.T) {
		attempts := 0
		l := &GenericListener{
			retryPolicy: RetryPolicy{
				MaxAttempts:    5,
				InitialBackoff: time.Millisecond,
				MaxBackoff:     time.Millisecond,
			},
			saveCheckpointFunc: func(ctx context.Context, version int64) (bool, error) {
				attempts++
				return false, ErrVersionNotContiguous
			},
		}

		_, err := l.saveCheckpointWithRetry(context.Background(), 10)
		if !errors.Is(err, ErrVersionNotContiguous) {
			t.Fatalf("expected ErrVersionNotContiguous, got %v", err)
		}
		if attempts != 1 {
			t.Fatalf("expected 1 attempt, got %d", attempts)
		}
	})

	t.Run("invariant violation is not retried", func(t *testing.T) {
		attempts := 0
		l := &GenericListener{
			retryPolicy: RetryPolicy{
				MaxAttempts:    5,
				InitialBackoff: time.Millisecond,
				MaxBackoff:     time.Millisecond,
			},
			saveCheckpointFunc: func(ctx context.Context, version int64) (bool, error) {
				attempts++
				return false, ErrInvariantViolation
			},
		}

		_, err := l.saveCheckpointWithRetry(context.Background(), 10)
		if !errors.Is(err, ErrInvariantViolation) {
			t.Fatalf("expected ErrInvariantViolation, got %v", err)
		}
		if attempts != 1 {
			t.Fatalf("expected 1 attempt, got %d", attempts)
		}
	})
}

func TestSaveCheckpointWithRetry_OnlyTransientIsRetried(t *testing.T) {
	t.Run("transient infra error is retried", func(t *testing.T) {
		attempts := 0
		l := &GenericListener{
			retryPolicy: RetryPolicy{
				MaxAttempts:    3,
				InitialBackoff: time.Millisecond,
				MaxBackoff:     time.Millisecond,
			},
			saveCheckpointFunc: func(ctx context.Context, version int64) (bool, error) {
				attempts++
				if attempts < 3 {
					return false, ErrTransientInfrastructure
				}
				return true, nil
			},
		}

		wasSaved, err := l.saveCheckpointWithRetry(context.Background(), 10)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !wasSaved {
			t.Fatalf("expected checkpoint saved")
		}
		if attempts != 3 {
			t.Fatalf("expected 3 attempts, got %d", attempts)
		}
	})

	t.Run("unknown error is not retried", func(t *testing.T) {
		attempts := 0
		unknownErr := errors.New("unknown")
		l := &GenericListener{
			retryPolicy: RetryPolicy{
				MaxAttempts:    5,
				InitialBackoff: time.Millisecond,
				MaxBackoff:     time.Millisecond,
			},
			saveCheckpointFunc: func(ctx context.Context, version int64) (bool, error) {
				attempts++
				return false, unknownErr
			},
		}

		_, err := l.saveCheckpointWithRetry(context.Background(), 10)
		if !errors.Is(err, unknownErr) {
			t.Fatalf("expected unknown error, got %v", err)
		}
		if attempts != 1 {
			t.Fatalf("expected 1 attempt, got %d", attempts)
		}
	})
}
