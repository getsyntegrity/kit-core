package resilience

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewTimeout(t *testing.T) {
	config := TimeoutConfig{
		Enabled:  true,
		Duration: 100 * time.Millisecond,
	}

	timeout := NewTimeout(config)
	if timeout == nil {
		t.Error("Expected timeout to be not nil")
	}
}

func TestTimeoutExecute(t *testing.T) {
	config := TimeoutConfig{
		Enabled:  true,
		Duration: 100 * time.Millisecond,
	}

	timeout := NewTimeout(config)

	err := timeout.Execute(func() error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestTimeoutExecuteWithContext(t *testing.T) {
	config := TimeoutConfig{
		Enabled:  true,
		Duration: 100 * time.Millisecond,
	}

	timeout := NewTimeout(config)
	ctx := context.Background()

	err := timeout.ExecuteWithContext(ctx, func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestTimeoutExecuteWithTimeout(t *testing.T) {
	config := TimeoutConfig{
		Enabled:  true,
		Duration: 10 * time.Millisecond,
	}

	timeout := NewTimeout(config)

	err := timeout.Execute(func() error {
		time.Sleep(50 * time.Millisecond)
		return nil
	})

	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestTimeoutExecuteWithContextAndTimeout(t *testing.T) {
	config := TimeoutConfig{
		Enabled:  true,
		Duration: 10 * time.Millisecond,
	}

	timeout := NewTimeout(config)

	err := timeout.ExecuteWithContext(context.Background(), func(ctx context.Context) error {
		time.Sleep(50 * time.Millisecond)
		return nil
	})

	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestTimeoutDisabled(t *testing.T) {
	config := TimeoutConfig{
		Enabled:  false,
		Duration: 10 * time.Millisecond,
	}

	timeout := NewTimeout(config)

	err := timeout.Execute(func() error {
		time.Sleep(50 * time.Millisecond)
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestTimeoutDisabledWithContext(t *testing.T) {
	config := TimeoutConfig{
		Enabled:  false,
		Duration: 10 * time.Millisecond,
	}

	timeout := NewTimeout(config)

	err := timeout.ExecuteWithContext(context.Background(), func(ctx context.Context) error {
		time.Sleep(50 * time.Millisecond)
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestTimeoutContextCancellation(t *testing.T) {
	config := TimeoutConfig{
		Enabled:  true,
		Duration: 100 * time.Millisecond,
	}

	timeout := NewTimeout(config)
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context immediately
	cancel()

	err := timeout.ExecuteWithContext(ctx, func(ctx context.Context) error {
		time.Sleep(50 * time.Millisecond)
		return nil
	})

	if err == nil {
		t.Error("Expected context cancellation error, got nil")
	}
}

func TestWithTimeout(t *testing.T) {
	config := TimeoutConfig{
		Enabled:  true,
		Duration: 10 * time.Millisecond,
	}

	err := WithTimeout(func() error {
		time.Sleep(50 * time.Millisecond)
		return nil
	}, config)

	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestWithTimeoutAndContext(t *testing.T) {
	config := TimeoutConfig{
		Enabled:  true,
		Duration: 10 * time.Millisecond,
	}

	err := WithTimeoutAndContext(context.Background(), func(ctx context.Context) error {
		time.Sleep(50 * time.Millisecond)
		return nil
	}, config)

	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestTimeoutWithError(t *testing.T) {
	config := TimeoutConfig{
		Enabled:  true,
		Duration: 100 * time.Millisecond,
	}

	timeout := NewTimeout(config)

	err := timeout.Execute(func() error {
		return errors.New("operation error")
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "operation error" {
		t.Errorf("Expected 'operation error', got %v", err)
	}
}

func TestTimeoutWithContextAndError(t *testing.T) {
	config := TimeoutConfig{
		Enabled:  true,
		Duration: 100 * time.Millisecond,
	}

	timeout := NewTimeout(config)

	err := timeout.ExecuteWithContext(context.Background(), func(ctx context.Context) error {
		return errors.New("operation error")
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "operation error" {
		t.Errorf("Expected 'operation error', got %v", err)
	}
}
