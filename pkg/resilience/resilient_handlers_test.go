package resilience

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Use MockEventHandler from mocks.go - no need to redefine here

func TestResilientEventHandler_NoResilience(t *testing.T) {
	// Test handler without any resilience patterns
	mockHandler := NewMockEventHandler()
	mockHandler.GetEventTypeFunc = func() string {
		return "test.event"
	}
	mockHandler.HandleFunc = func(ctx context.Context, event interface{}) error {
		return nil
	}

	config := Config{
		Retry: RetryConfig{
			Enabled: false,
		},
		CircuitBreaker: CircuitBreakerConfig{
			Enabled: false,
		},
		Timeout: TimeoutConfig{
			Enabled: false,
		},
	}

	resilientHandler := NewResilientEventHandler(mockHandler, nil, config, FixedClock{Time: time.Now()})

	ctx := context.Background()
	err := resilientHandler.Handle(ctx, "test event")

	require.NoError(t, err)
	assert.Equal(t, 1, mockHandler.CallCount)
}

func TestResilientEventHandler_WithRetry(t *testing.T) {
	// Test handler with retry pattern
	callCount := 0
	mockHandler := NewMockEventHandler()
	mockHandler.GetEventTypeFunc = func() string {
		return "test.event"
	}
	mockHandler.HandleFunc = func(ctx context.Context, event interface{}) error {
		callCount++
		if callCount < 3 {
			return errors.New("temporary error")
		}
		return nil
	}

	config := Config{
		Retry: RetryConfig{
			Enabled:     true,
			MaxAttempts: 3,
			Delay:       10 * time.Millisecond,
		},
		CircuitBreaker: CircuitBreakerConfig{
			Enabled: false,
		},
		Timeout: TimeoutConfig{
			Enabled: false,
		},
	}

	resilientHandler := NewResilientEventHandler(mockHandler, nil, config, FixedClock{Time: time.Now()})

	ctx := context.Background()
	err := resilientHandler.Handle(ctx, "test event")

	require.NoError(t, err)
	assert.Equal(t, 3, callCount)
}

func TestResilientEventHandler_WithTimeout(t *testing.T) {
	// Test handler with timeout pattern
	mockHandler := NewMockEventHandler()
	mockHandler.GetEventTypeFunc = func() string {
		return "test.event"
	}
	mockHandler.HandleFunc = func(ctx context.Context, event interface{}) error {
		// Simulate a long-running operation
		time.Sleep(200 * time.Millisecond)
		return nil
	}

	config := Config{
		Retry: RetryConfig{
			Enabled: false,
		},
		CircuitBreaker: CircuitBreakerConfig{
			Enabled: false,
		},
		Timeout: TimeoutConfig{
			Enabled:  true,
			Duration: 100 * time.Millisecond,
		},
	}

	resilientHandler := NewResilientEventHandler(mockHandler, nil, config, FixedClock{Time: time.Now()})

	ctx := context.Background()
	err := resilientHandler.Handle(ctx, "test event")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestResilientEventHandler_RetryMaxAttempts(t *testing.T) {
	// Test handler with retry that exceeds max attempts
	mockHandler := NewMockEventHandler()
	mockHandler.GetEventTypeFunc = func() string {
		return "test.event"
	}
	mockHandler.HandleFunc = func(ctx context.Context, event interface{}) error {
		return errors.New("persistent error")
	}

	config := Config{
		Retry: RetryConfig{
			Enabled:     true,
			MaxAttempts: 2,
			Delay:       10 * time.Millisecond,
		},
		CircuitBreaker: CircuitBreakerConfig{
			Enabled: false,
		},
		Timeout: TimeoutConfig{
			Enabled: false,
		},
	}

	resilientHandler := NewResilientEventHandler(mockHandler, nil, config, FixedClock{Time: time.Now()})

	ctx := context.Background()
	err := resilientHandler.Handle(ctx, "test event")

	require.Error(t, err)
	assert.Equal(t, "persistent error", err.Error())
	assert.Equal(t, 3, mockHandler.CallCount) // Initial attempt + 2 retries
}

func TestResilientEventHandler_ContextCancellation(t *testing.T) {
	// Test handler with context cancellation
	mockHandler := NewMockEventHandler()
	mockHandler.GetEventTypeFunc = func() string {
		return "test.event"
	}
	mockHandler.HandleFunc = func(ctx context.Context, event interface{}) error {
		// Simulate a long-running operation that respects context
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(200 * time.Millisecond):
			return nil
		}
	}

	config := Config{
		Retry: RetryConfig{
			Enabled: false, // Disable retry to make test more predictable
		},
		CircuitBreaker: CircuitBreakerConfig{
			Enabled: false,
		},
		Timeout: TimeoutConfig{
			Enabled: false,
		},
	}

	resilientHandler := NewResilientEventHandler(mockHandler, nil, config, FixedClock{Time: time.Now()})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := resilientHandler.Handle(ctx, "test event")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestResilientEventHandler_GetEventType(t *testing.T) {
	// Test that GetEventType is properly delegated
	mockHandler := NewMockEventHandler()
	mockHandler.GetEventTypeFunc = func() string {
		return "test.event"
	}
	mockHandler.HandleFunc = func(ctx context.Context, event interface{}) error {
		return nil
	}

	config := Config{}
	resilientHandler := NewResilientEventHandler(mockHandler, nil, config, FixedClock{Time: time.Now()})

	assert.Equal(t, "test.event", resilientHandler.GetEventType())
}
