package resilience

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewRetry(t *testing.T) {
	config := RetryConfig{
		Enabled:     true,
		MaxAttempts: 3,
		Delay:       100 * time.Millisecond,
	}

	retry := NewRetry(config)
	if retry == nil {
		t.Error("Expected retry to be not nil")
	}
}

func TestRetryExecute(t *testing.T) {
	config := RetryConfig{
		Enabled:     true,
		MaxAttempts: 3,
		Delay:       10 * time.Millisecond,
	}

	retry := NewRetry(config)

	err := retry.Execute(func() error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRetryExecuteWithContext(t *testing.T) {
	config := RetryConfig{
		Enabled:     true,
		MaxAttempts: 3,
		Delay:       10 * time.Millisecond,
	}

	retry := NewRetry(config)
	ctx := context.Background()

	err := retry.ExecuteWithContext(ctx, func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRetryExecuteWithFailure(t *testing.T) {
	config := RetryConfig{
		Enabled:     true,
		MaxAttempts: 2,
		Delay:       10 * time.Millisecond,
	}

	retry := NewRetry(config)
	attempts := 0

	err := retry.Execute(func() error {
		attempts++
		return errors.New("operation failed")
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if attempts != 3 { // MaxAttempts + 1
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetryExecuteWithContextAndFailure(t *testing.T) {
	config := RetryConfig{
		Enabled:     true,
		MaxAttempts: 2,
		Delay:       10 * time.Millisecond,
	}

	retry := NewRetry(config)
	attempts := 0

	err := retry.ExecuteWithContext(context.Background(), func(ctx context.Context) error {
		attempts++
		return errors.New("operation failed")
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if attempts != 3 { // MaxAttempts + 1
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetryDisabled(t *testing.T) {
	config := RetryConfig{
		Enabled:     false,
		MaxAttempts: 3,
		Delay:       10 * time.Millisecond,
	}

	retry := NewRetry(config)
	attempts := 0

	err := retry.Execute(func() error {
		attempts++
		return errors.New("operation failed")
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

func TestRetryDisabledWithContext(t *testing.T) {
	config := RetryConfig{
		Enabled:     false,
		MaxAttempts: 3,
		Delay:       10 * time.Millisecond,
	}

	retry := NewRetry(config)
	attempts := 0

	err := retry.ExecuteWithContext(context.Background(), func(ctx context.Context) error {
		attempts++
		return errors.New("operation failed")
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

func TestRetryContextCancellation(t *testing.T) {
	config := RetryConfig{
		Enabled:     true,
		MaxAttempts: 5,
		Delay:       100 * time.Millisecond,
	}

	retry := NewRetry(config)
	ctx, cancel := context.WithCancel(context.Background())
	attempts := 0

	// Cancel context after a short delay
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := retry.ExecuteWithContext(ctx, func(ctx context.Context) error {
		attempts++
		return errors.New("operation failed")
	})

	if err == nil {
		t.Error("Expected context cancellation error, got nil")
	}
	if attempts < 1 {
		t.Errorf("Expected at least 1 attempt, got %d", attempts)
	}
}

func TestWithRetry(t *testing.T) {
	config := RetryConfig{
		Enabled:     true,
		MaxAttempts: 2,
		Delay:       10 * time.Millisecond,
	}

	attempts := 0
	err := WithRetry(func() error {
		attempts++
		return errors.New("operation failed")
	}, config)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if attempts != 3 { // MaxAttempts + 1
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestWithRetryAndContext(t *testing.T) {
	config := RetryConfig{
		Enabled:     true,
		MaxAttempts: 2,
		Delay:       10 * time.Millisecond,
	}

	attempts := 0
	err := WithRetryAndContext(context.Background(), func(ctx context.Context) error {
		attempts++
		return errors.New("operation failed")
	}, config)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if attempts != 3 { // MaxAttempts + 1
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetrySuccessAfterFailure(t *testing.T) {
	config := RetryConfig{
		Enabled:     true,
		MaxAttempts: 3,
		Delay:       10 * time.Millisecond,
	}

	retry := NewRetry(config)
	attempts := 0

	err := retry.Execute(func() error {
		attempts++
		if attempts < 3 {
			return errors.New("operation failed")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}
