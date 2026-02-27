package resilience

import (
	"context"
	"errors"
	"testing"
)

func TestNewFallback(t *testing.T) {
	config := FallbackConfig{
		Enabled: true,
	}

	fallback := NewFallback(config)
	if fallback == nil {
		t.Error("Expected fallback to be not nil")
	}
}

func TestFallbackExecute(t *testing.T) {
	config := FallbackConfig{
		Enabled: true,
	}

	fallback := NewFallback(config)

	// Test successful main operation
	err := fallback.Execute(func() error {
		return nil
	}, func(err error) error {
		return errors.New("fallback should not be called")
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestFallbackExecuteWithFallback(t *testing.T) {
	config := FallbackConfig{
		Enabled: true,
	}

	fallback := NewFallback(config)

	// Test main operation fails, fallback succeeds
	err := fallback.Execute(func() error {
		return errors.New("main failed")
	}, func(err error) error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestFallbackExecuteWithContext(t *testing.T) {
	config := FallbackConfig{
		Enabled: true,
	}

	fallback := NewFallback(config)

	// Test successful main operation with context
	err := fallback.ExecuteWithContext(context.Background(), func(ctx context.Context) error {
		return nil
	}, func(ctx context.Context, err error) error {
		return errors.New("fallback should not be called")
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestFallbackExecuteWithContextAndFallback(t *testing.T) {
	config := FallbackConfig{
		Enabled: true,
	}

	fallback := NewFallback(config)

	// Test main operation fails, fallback succeeds with context
	err := fallback.ExecuteWithContext(context.Background(), func(ctx context.Context) error {
		return errors.New("main failed")
	}, func(ctx context.Context, err error) error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestFallbackDisabled(t *testing.T) {
	config := FallbackConfig{
		Enabled: false,
	}

	fallback := NewFallback(config)

	// Test that fallback is not used when disabled
	err := fallback.Execute(func() error {
		return errors.New("main failed")
	}, func(err error) error {
		return errors.New("fallback should not be called")
	})

	if err == nil {
		t.Error("Expected error from main operation")
	}
	if err.Error() != "main failed" {
		t.Errorf("Expected 'main failed' error, got %v", err)
	}
}

func TestFallbackDisabledWithContext(t *testing.T) {
	config := FallbackConfig{
		Enabled: false,
	}

	fallback := NewFallback(config)

	// Test that fallback is not used when disabled with context
	err := fallback.ExecuteWithContext(context.Background(), func(ctx context.Context) error {
		return errors.New("main failed")
	}, func(ctx context.Context, err error) error {
		return errors.New("fallback should not be called")
	})

	if err == nil {
		t.Error("Expected error from main operation")
	}
	if err.Error() != "main failed" {
		t.Errorf("Expected 'main failed' error, got %v", err)
	}
}

func TestWithFallbackOperation(t *testing.T) {
	config := FallbackConfig{
		Enabled: true,
	}

	// Test successful main operation
	err := WithFallbackOperation(func() error {
		return nil
	}, func(err error) error {
		return errors.New("fallback should not be called")
	}, config)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestWithFallbackOperationWithFallback(t *testing.T) {
	config := FallbackConfig{
		Enabled: true,
	}

	// Test main operation fails, fallback succeeds
	err := WithFallbackOperation(func() error {
		return errors.New("main failed")
	}, func(err error) error {
		return nil
	}, config)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestWithFallbackAndContext(t *testing.T) {
	config := FallbackConfig{
		Enabled: true,
	}

	// Test successful main operation with context
	err := WithFallbackAndContext(context.Background(), func(ctx context.Context) error {
		return nil
	}, func(ctx context.Context, err error) error {
		return errors.New("fallback should not be called")
	}, config)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestWithFallbackAndContextWithFallback(t *testing.T) {
	config := FallbackConfig{
		Enabled: true,
	}

	// Test main operation fails, fallback succeeds with context
	err := WithFallbackAndContext(context.Background(), func(ctx context.Context) error {
		return errors.New("main failed")
	}, func(ctx context.Context, err error) error {
		return nil
	}, config)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
