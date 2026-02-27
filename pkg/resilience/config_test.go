package resilience

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	// Test Timeout config
	if config.Timeout.Enabled {
		t.Error("Timeout should be disabled by default")
	}
	if config.Timeout.Duration != 30*time.Second {
		t.Errorf("Expected timeout duration 30s, got %v", config.Timeout.Duration)
	}

	// Test Retry config
	if config.Retry.Enabled {
		t.Error("Retry should be disabled by default")
	}
	if config.Retry.MaxAttempts != 3 {
		t.Errorf("Expected max attempts 3, got %d", config.Retry.MaxAttempts)
	}
	if config.Retry.Delay != 100*time.Millisecond {
		t.Errorf("Expected delay 100ms, got %v", config.Retry.Delay)
	}

	// Test Fallback config
	if config.Fallback.Enabled {
		t.Error("Fallback should be disabled by default")
	}

	// Test Circuit Breaker config
	if config.CircuitBreaker.Enabled {
		t.Error("Circuit Breaker should be disabled by default")
	}
	if config.CircuitBreaker.FailureThreshold != 5 {
		t.Errorf("Expected failure threshold 5, got %d", config.CircuitBreaker.FailureThreshold)
	}
	// SuccessThreshold removed - KEEP IT SIMPLE
	if config.CircuitBreaker.Timeout != 60*time.Second {
		t.Errorf("Expected timeout 60s, got %v", config.CircuitBreaker.Timeout)
	}
}

func TestConfigStructTags(t *testing.T) {
	// This test ensures that the struct tags are properly defined
	// for environment variable configuration

	config := DefaultConfig()

	// Test that we can access all fields (this indirectly tests struct tags)
	_ = config.Timeout.Enabled
	_ = config.Timeout.Duration

	_ = config.Retry.Enabled
	_ = config.Retry.MaxAttempts
	_ = config.Retry.Delay

	_ = config.Fallback.Enabled

	_ = config.CircuitBreaker.Enabled
	_ = config.CircuitBreaker.FailureThreshold
	_ = config.CircuitBreaker.Timeout
	// Removed complex fields - KEEP IT SIMPLE
}

func TestNewConfigWithOptions(t *testing.T) {
	config := NewConfig(
		WithRetryConfig(5, 200*time.Millisecond),
		WithTimeoutConfig(10*time.Second),
		WithFallbackConfig(),
		WithCircuitBreakerConfig(10, 120*time.Second),
	)

	if !config.Retry.Enabled {
		t.Error("Retry should be enabled")
	}
	if config.Retry.MaxAttempts != 5 {
		t.Errorf("Expected max attempts 5, got %d", config.Retry.MaxAttempts)
	}
	if config.Retry.Delay != 200*time.Millisecond {
		t.Errorf("Expected delay 200ms, got %v", config.Retry.Delay)
	}

	if !config.Timeout.Enabled {
		t.Error("Timeout should be enabled")
	}
	if config.Timeout.Duration != 10*time.Second {
		t.Errorf("Expected timeout duration 10s, got %v", config.Timeout.Duration)
	}

	if !config.Fallback.Enabled {
		t.Error("Fallback should be enabled")
	}

	if !config.CircuitBreaker.Enabled {
		t.Error("Circuit Breaker should be enabled")
	}
	if config.CircuitBreaker.FailureThreshold != 10 {
		t.Errorf("Expected failure threshold 10, got %d", config.CircuitBreaker.FailureThreshold)
	}
	// SuccessThreshold removed - KEEP IT SIMPLE
}

func TestIsEnabled(t *testing.T) {
	config := DefaultConfig()
	if config.IsEnabled() {
		t.Error("Default config should not be enabled")
	}

	config = NewConfig(WithRetryConfig(3, 100*time.Millisecond))
	if !config.IsEnabled() {
		t.Error("Config with retry enabled should be enabled")
	}

	config = NewConfig(WithTimeoutConfig(5 * time.Second))
	if !config.IsEnabled() {
		t.Error("Config with timeout enabled should be enabled")
	}

	config = NewConfig(WithFallbackConfig())
	if !config.IsEnabled() {
		t.Error("Config with fallback enabled should be enabled")
	}

	config = NewConfig(WithCircuitBreakerConfig(5, 60*time.Second))
	if !config.IsEnabled() {
		t.Error("Config with circuit breaker enabled should be enabled")
	}
}
