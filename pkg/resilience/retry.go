package resilience

import (
	"context"
	"time"
)

// Retry provides simple retry functionality.
type Retry struct {
	config RetryConfig
}

// NewRetry creates a new retry with the given configuration.
func NewRetry(config RetryConfig) *Retry {
	return &Retry{
		config: config,
	}
}

// Execute executes an operation with retry logic.
func (r *Retry) Execute(operation func() error) error {
	return r.ExecuteWithContext(context.Background(), func(_ context.Context) error {
		return operation()
	})
}

// ExecuteWithContext executes an operation with retry logic and context.
func (r *Retry) ExecuteWithContext(ctx context.Context, operation func(ctx context.Context) error) error {
	if !r.config.Enabled {
		return operation(ctx)
	}

	var lastErr error
	for attempt := 0; attempt <= r.config.MaxAttempts; attempt++ {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Execute operation
		err := operation(ctx)
		if err == nil {
			return nil // Success
		}

		lastErr = err

		// Don't retry on last attempt
		if attempt == r.config.MaxAttempts {
			break
		}

		// Wait before retry
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(r.config.Delay):
		}
	}

	return lastErr
}

// WithRetry is a convenience function for simple retry operations.
func WithRetry(operation func() error, config RetryConfig) error {
	retry := NewRetry(config)
	return retry.Execute(operation)
}

// WithRetryAndContext is a convenience function for retry operations with context.
func WithRetryAndContext(ctx context.Context, operation func(ctx context.Context) error, config RetryConfig) error {
	retry := NewRetry(config)
	return retry.ExecuteWithContext(ctx, operation)
}
