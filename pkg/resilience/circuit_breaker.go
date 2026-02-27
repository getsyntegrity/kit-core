// Package resilience provides circuit breaker functionality.
package resilience

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

// safeInt32 safely converts int to int32 with bounds checking.
func safeInt32(value int) int32 {
	const maxInt32 = 2147483647
	const minInt32 = -2147483648

	if value > maxInt32 {
		return maxInt32
	}
	if value < minInt32 {
		return minInt32
	}
	return int32(value)
}

// CircuitBreakerState represents the state of the circuit breaker.
type CircuitBreakerState int

const (
	// CLOSED represents the closed state of the circuit breaker.
	CLOSED CircuitBreakerState = iota
	// OPEN represents the open state of the circuit breaker.
	OPEN
	// HalfOpen represents the half-open state of the circuit breaker.
	HalfOpen
)

// String returns the string representation of the circuit breaker state.
func (s CircuitBreakerState) String() string {
	switch s {
	case CLOSED:
		return "CLOSED"
	case OPEN:
		return "OPEN"
	case HalfOpen:
		return "HalfOpen"
	default:
		return "UNKNOWN"
	}
}

// CircuitBreaker provides circuit breaker functionality - KEEP IT SIMPLE.
type CircuitBreaker struct {
	config CircuitBreakerConfig
	clock  Clock
	state  int32 // atomic state

	// Simple counters - KEEP IT SIMPLE
	failureCount    int64
	lastFailureTime time.Time
}

// NewCircuitBreaker creates a new circuit breaker with the given configuration and clock. KEEP IT SIMPLE.
// Callers must inject a Clock (e.g. from kit-runtime for real time).
func NewCircuitBreaker(config CircuitBreakerConfig, clock Clock) *CircuitBreaker {
	cb := &CircuitBreaker{
		config: config,
		clock:  clock,
		state:  int32(CLOSED),
	}
	return cb
}

// Execute executes an operation with circuit breaker logic.
func (cb *CircuitBreaker) Execute(operation func() error) error {
	return cb.ExecuteWithContext(context.Background(), func(_ context.Context) error {
		return operation()
	})
}

// ExecuteWithContext executes an operation with circuit breaker logic and context - KEEP IT SIMPLE.
func (cb *CircuitBreaker) ExecuteWithContext(ctx context.Context, operation func(ctx context.Context) error) error {
	if !cb.config.Enabled {
		return operation(ctx)
	}

	// Check if circuit breaker is open
	if cb.getState() == OPEN {
		if cb.clock.Now().Sub(cb.lastFailureTime) < cb.config.Timeout {
			return errors.New("circuit breaker is open")
		}
		// Try to transition to half-open
		cb.transitionToHalfOpen()
	}

	// Execute operation
	err := operation(ctx)

	// Update state based on result - KEEP IT SIMPLE
	if err != nil {
		cb.onFailure()
	} else {
		cb.onSuccess()
	}

	return err
}

// getState returns the current state of the circuit breaker.
func (cb *CircuitBreaker) getState() CircuitBreakerState {
	return CircuitBreakerState(atomic.LoadInt32(&cb.state))
}

// setState sets the state of the circuit breaker.
func (cb *CircuitBreaker) setState(state CircuitBreakerState) {
	atomic.StoreInt32(&cb.state, safeInt32(int(state)))
}

// transitionToHalfOpen attempts to transition to half-open state - KEEP IT SIMPLE.
func (cb *CircuitBreaker) transitionToHalfOpen() {
	if cb.getState() == OPEN && cb.clock.Now().Sub(cb.lastFailureTime) >= cb.config.Timeout {
		cb.setState(HalfOpen)
		// Reset failure count for half-open state
		atomic.StoreInt64(&cb.failureCount, 0)
	}
}

// onFailure handles a failed call - KEEP IT SIMPLE.
func (cb *CircuitBreaker) onFailure() {
	cb.lastFailureTime = cb.clock.Now()
	atomic.AddInt64(&cb.failureCount, 1)

	// Check if we should open the circuit
	if atomic.LoadInt64(&cb.failureCount) >= int64(cb.config.FailureThreshold) {
		cb.setState(OPEN)
	}
}

// onSuccess handles a successful call - KEEP IT SIMPLE.
func (cb *CircuitBreaker) onSuccess() {
	// Reset failure count on success
	atomic.StoreInt64(&cb.failureCount, 0)

	// If we were in half-open, go back to closed
	if cb.getState() == HalfOpen {
		cb.setState(CLOSED)
	}
}

// GetState returns the current state as a string.
func (cb *CircuitBreaker) GetState() string {
	return cb.getState().String()
}

// GetStats returns simple statistics - KEEP IT SIMPLE.
func (cb *CircuitBreaker) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"state":         cb.GetState(),
		"failure_count": atomic.LoadInt64(&cb.failureCount),
		"last_failure":  cb.lastFailureTime,
		"config":        cb.config,
	}
}

// CircuitBreakerInterface defines the interface for circuit breaker operations.
type CircuitBreakerInterface interface {
	Execute(operation func() error) error
	ExecuteWithContext(ctx context.Context, operation func(ctx context.Context) error) error
	GetState() string
	GetStats() map[string]interface{}
	Reset()
}

// Reset resets the circuit breaker to closed state.
func (cb *CircuitBreaker) Reset() {
	cb.setState(CLOSED)
	atomic.StoreInt64(&cb.failureCount, 0)
}
