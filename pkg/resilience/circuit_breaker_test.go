package resilience

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewCircuitBreaker(t *testing.T) {
	config := CircuitBreakerConfig{
		Enabled:          true,
		FailureThreshold: 5,
		Timeout:          60 * time.Second,
	}

	cb := NewCircuitBreaker(config, FixedClock{Time: time.Now()})
	if cb == nil {
		t.Error("Expected circuit breaker to be not nil")
	}
	if cb.GetState() != "CLOSED" {
		t.Errorf("Expected initial state to be CLOSED, got %s", cb.GetState())
	}
}

func TestCircuitBreakerExecute(t *testing.T) {
	config := CircuitBreakerConfig{
		Enabled:          true,
		FailureThreshold: 5,
		Timeout:          60 * time.Second,
	}

	cb := NewCircuitBreaker(config, FixedClock{Time: time.Now()})

	err := cb.Execute(func() error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if cb.GetState() != "CLOSED" {
		t.Errorf("Expected state to be CLOSED, got %s", cb.GetState())
	}
}

func TestCircuitBreakerExecuteWithContext(t *testing.T) {
	config := CircuitBreakerConfig{
		Enabled:          true,
		FailureThreshold: 5,
		Timeout:          60 * time.Second,
	}

	cb := NewCircuitBreaker(config, FixedClock{Time: time.Now()})

	err := cb.ExecuteWithContext(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if cb.GetState() != "CLOSED" {
		t.Errorf("Expected state to be CLOSED, got %s", cb.GetState())
	}
}

func TestCircuitBreakerOpenOnFailure(t *testing.T) {
	config := CircuitBreakerConfig{
		Enabled:          true,
		FailureThreshold: 3, // Lower threshold for faster test
		Timeout:          60 * time.Second,
	}

	cb := NewCircuitBreaker(config, FixedClock{Time: time.Now()})

	// Execute multiple failures to trigger circuit breaker
	for i := 0; i < 3; i++ {
		err := cb.Execute(func() error {
			return errors.New("operation failed")
		})
		if err == nil {
			t.Error("Expected error, got nil")
		}
	}

	// Check if circuit breaker is open
	if cb.GetState() != "OPEN" {
		t.Errorf("Expected circuit breaker to be OPEN, got %s", cb.GetState())
	}
}

func TestCircuitBreakerDisabled(t *testing.T) {
	config := CircuitBreakerConfig{
		Enabled:          false,
		FailureThreshold: 5,
		Timeout:          60 * time.Second,
	}

	cb := NewCircuitBreaker(config, FixedClock{Time: time.Now()})

	// Even with failures, circuit breaker should not open when disabled
	for i := 0; i < 10; i++ {
		err := cb.Execute(func() error {
			return errors.New("operation failed")
		})
		if err == nil {
			t.Error("Expected error, got nil")
		}
	}

	// Should still be closed when disabled
	if cb.GetState() != "CLOSED" {
		t.Errorf("Expected circuit breaker to be CLOSED when disabled, got %s", cb.GetState())
	}
}

func TestCircuitBreakerGetStats(t *testing.T) {
	config := CircuitBreakerConfig{
		Enabled:          true,
		FailureThreshold: 5,
		Timeout:          60 * time.Second,
	}

	cb := NewCircuitBreaker(config, FixedClock{Time: time.Now()})

	// Execute a failure
	err := cb.Execute(func() error {
		return errors.New("operation failed")
	})
	if err == nil {
		t.Error("Expected error, got nil")
	}

	stats := cb.GetStats()
	if stats == nil {
		t.Error("Expected stats to be not nil")
	}

	// Check basic stats
	if stats["state"] != "CLOSED" {
		t.Errorf("Expected state to be CLOSED, got %s", stats["state"])
	}

	if failureCount, ok := stats["failure_count"].(int64); !ok || failureCount != 1 {
		t.Errorf("Expected failure count to be 1, got %v", failureCount)
	}
}

func TestCircuitBreakerReset(t *testing.T) {
	config := CircuitBreakerConfig{
		Enabled:          true,
		FailureThreshold: 2,
		Timeout:          60 * time.Second,
	}

	cb := NewCircuitBreaker(config, FixedClock{Time: time.Now()})

	// Open the circuit breaker
	for i := 0; i < 2; i++ {
		cb.Execute(func() error {
			return errors.New("operation failed")
		})
	}

	if cb.GetState() != "OPEN" {
		t.Errorf("Expected circuit breaker to be OPEN, got %s", cb.GetState())
	}

	// Reset the circuit breaker
	cb.Reset()

	if cb.GetState() != "CLOSED" {
		t.Errorf("Expected circuit breaker to be CLOSED after reset, got %s", cb.GetState())
	}

	stats := cb.GetStats()
	if failureCount, ok := stats["failure_count"].(int64); !ok || failureCount != 0 {
		t.Errorf("Expected failure count to be 0 after reset, got %v", failureCount)
	}
}
