package resilience

import (
	"context"
)

// Fallback provides simple fallback functionality.
type Fallback struct {
	config FallbackConfig
}

// NewFallback creates a new fallback with the given configuration.
func NewFallback(config FallbackConfig) *Fallback {
	return &Fallback{
		config: config,
	}
}

// Execute executes an operation with fallback logic.
func (f *Fallback) Execute(operation func() error, fallbackFn func(error) error) error {
	return f.ExecuteWithContext(context.Background(), func(_ context.Context) error {
		return operation()
	}, func(_ context.Context, err error) error {
		return fallbackFn(err)
	})
}

// ExecuteWithContext executes an operation with fallback logic and context.
func (f *Fallback) ExecuteWithContext(ctx context.Context, operation func(ctx context.Context) error, fallbackFn func(ctx context.Context, err error) error) error {
	if !f.config.Enabled {
		return operation(ctx)
	}

	// Try primary operation
	err := operation(ctx)
	if err == nil {
		return nil // Success
	}

	// Execute fallback if primary operation failed
	return fallbackFn(ctx, err)
}

// WithFallbackOperation is a convenience function for simple fallback operations.
func WithFallbackOperation(operation func() error, fallbackFn func(error) error, config FallbackConfig) error {
	fallback := NewFallback(config)
	return fallback.Execute(operation, fallbackFn)
}

// WithFallbackAndContext is a convenience function for fallback operations with context.
func WithFallbackAndContext(ctx context.Context, operation func(ctx context.Context) error, fallbackFn func(ctx context.Context, err error) error, config FallbackConfig) error {
	fallback := NewFallback(config)
	return fallback.ExecuteWithContext(ctx, operation, fallbackFn)
}
