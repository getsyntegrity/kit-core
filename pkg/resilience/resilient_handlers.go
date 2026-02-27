package resilience

import (
	"context"
	"time"
)

// EventHandler is a simple interface for handling events.
type EventHandler interface {
	GetEventType() string
	Handle(ctx context.Context, event interface{}) error
}

// ResilientEventHandler wraps a handler with resilience patterns.
type ResilientEventHandler struct {
	handler          EventHandler
	fallbackFn       func(ctx context.Context, event interface{}) error
	resilienceConfig Config
	circuitBreaker   *CircuitBreaker
}

// NewResilientEventHandler creates a new resilient event handler.
// Callers must inject a Clock (e.g. from kit-runtime for real time).
func NewResilientEventHandler(handler EventHandler, fallbackFn func(ctx context.Context, event interface{}) error, config Config, clock Clock) *ResilientEventHandler {
	return &ResilientEventHandler{
		handler:          handler,
		fallbackFn:       fallbackFn,
		resilienceConfig: config,
		circuitBreaker:   NewCircuitBreaker(config.CircuitBreaker, clock),
	}
}

// GetEventType returns the event type.
func (h *ResilientEventHandler) GetEventType() string {
	return h.handler.GetEventType()
}

// Handle handles the event with resilience patterns.
func (h *ResilientEventHandler) Handle(ctx context.Context, event interface{}) error {
	// Create primary operation
	primaryOperation := func(ctx context.Context) error {
		return h.handler.Handle(ctx, event)
	}

	// Apply simple resilience patterns
	return h.executeWithResilience(ctx, primaryOperation)
}

// executeWithResilience applies simple resilience patterns.
func (h *ResilientEventHandler) executeWithResilience(ctx context.Context, operation func(ctx context.Context) error) error {
	// Simple retry pattern
	if h.resilienceConfig.Retry.Enabled {
		return h.executeWithRetry(ctx, operation)
	}

	// Simple circuit breaker pattern
	if h.resilienceConfig.CircuitBreaker.Enabled {
		return h.executeWithCircuitBreaker(ctx, operation)
	}

	// Simple timeout pattern
	if h.resilienceConfig.Timeout.Enabled {
		return h.executeWithTimeout(ctx, operation)
	}

	// No resilience patterns enabled, execute directly
	return operation(ctx)
}

// executeWithRetry applies simple retry logic.
func (h *ResilientEventHandler) executeWithRetry(ctx context.Context, operation func(ctx context.Context) error) error {
	var lastErr error

	for attempt := 0; attempt <= h.resilienceConfig.Retry.MaxAttempts; attempt++ {
		if attempt > 0 {
			// Simple exponential backoff
			backoff := h.resilienceConfig.Retry.Delay * time.Duration(1<<(attempt-1))
			if backoff > 5*time.Second {
				backoff = 5 * time.Second // Cap at 5 seconds
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		err := operation(ctx)
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if !h.isRetryableError(err) {
			return err
		}
	}

	return lastErr
}

// executeWithCircuitBreaker applies circuit breaker logic with proper state management.
func (h *ResilientEventHandler) executeWithCircuitBreaker(ctx context.Context, operation func(ctx context.Context) error) error {
	// Execute operation through circuit breaker using the existing implementation
	err := h.circuitBreaker.ExecuteWithContext(ctx, operation)
	// If circuit breaker is open, try fallback if available
	if err != nil {
		if err.Error() == "circuit breaker is open" {
			if h.fallbackFn != nil && h.resilienceConfig.Fallback.Enabled {
				return h.fallbackFn(ctx, nil)
			}
		}
	}

	return err
}

// executeWithTimeout applies simple timeout logic.
func (h *ResilientEventHandler) executeWithTimeout(ctx context.Context, operation func(ctx context.Context) error) error {
	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, h.resilienceConfig.Timeout.Duration)
	defer cancel()

	// Execute operation with timeout
	done := make(chan error, 1)
	go func() {
		done <- operation(timeoutCtx)
	}()

	select {
	case err := <-done:
		return err
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	}
}

// isRetryableError checks if an error is retryable.
func (h *ResilientEventHandler) isRetryableError(err error) bool {
	// Simple retryable error check
	// In a real implementation, you'd check error types, status codes, etc.
	return err != nil
}

// GetCircuitBreakerState returns the current state of the circuit breaker.
func (h *ResilientEventHandler) GetCircuitBreakerState() string {
	return h.circuitBreaker.GetState()
}

// GetCircuitBreakerStats returns circuit breaker statistics.
func (h *ResilientEventHandler) GetCircuitBreakerStats() map[string]interface{} {
	return h.circuitBreaker.GetStats()
}

// ResetCircuitBreaker resets the circuit breaker to CLOSED state.
func (h *ResilientEventHandler) ResetCircuitBreaker() {
	h.circuitBreaker.Reset()
}
