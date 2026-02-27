package resilience

import (
	"context"
)

// Timeout provides simple timeout functionality.
type Timeout struct {
	config TimeoutConfig
}

// NewTimeout creates a new timeout with the given configuration.
func NewTimeout(config TimeoutConfig) *Timeout {
	return &Timeout{
		config: config,
	}
}

// Execute executes an operation with timeout.
func (t *Timeout) Execute(operation func() error) error {
	return t.ExecuteWithContext(context.Background(), func(_ context.Context) error {
		return operation()
	})
}

// ExecuteWithContext executes an operation with timeout and context.
func (t *Timeout) ExecuteWithContext(ctx context.Context, operation func(ctx context.Context) error) error {
	if !t.config.Enabled {
		return operation(ctx)
	}

	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, t.config.Duration)
	defer cancel()

	// Execute operation with timeout
	err := operation(timeoutCtx)
	if err != nil {
		return err
	}

	// Check if timeout occurred
	select {
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	default:
		return nil
	}
}

// WithTimeout is a convenience function for simple timeout operations.
func WithTimeout(operation func() error, config TimeoutConfig) error {
	timeout := NewTimeout(config)
	return timeout.Execute(operation)
}

// WithTimeoutAndContext is a convenience function for timeout operations with context.
func WithTimeoutAndContext(ctx context.Context, operation func(ctx context.Context) error, config TimeoutConfig) error {
	timeout := NewTimeout(config)
	return timeout.ExecuteWithContext(ctx, operation)
}
