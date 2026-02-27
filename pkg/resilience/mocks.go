// Package resilience provides resilience patterns for distributed systems.
package resilience

import (
	"context"
)

// MockEventHandler implements EventHandler for testing.
type MockEventHandler struct {
	GetEventTypeFunc func() string
	HandleFunc       func(ctx context.Context, event interface{}) error
	CallCount        int // Track number of Handle calls for testing
}

// GetEventType returns the event type this handler processes.
func (m *MockEventHandler) GetEventType() string {
	if m.GetEventTypeFunc != nil {
		return m.GetEventTypeFunc()
	}
	return "test-event" // Default behavior
}

// Handle handles an event.
func (m *MockEventHandler) Handle(ctx context.Context, event interface{}) error {
	m.CallCount++
	if m.HandleFunc != nil {
		return m.HandleFunc(ctx, event)
	}
	return nil // Default behavior
}

// NewMockEventHandler creates a new mock EventHandler.
func NewMockEventHandler() *MockEventHandler {
	return &MockEventHandler{}
}

// NewMockEventHandlerWithFunc creates a new mock EventHandler with custom functions.
func NewMockEventHandlerWithFunc(
	getEventTypeFunc func() string,
	handleFunc func(ctx context.Context, event interface{}) error,
) *MockEventHandler {
	return &MockEventHandler{
		GetEventTypeFunc: getEventTypeFunc,
		HandleFunc:       handleFunc,
	}
}

// MockCircuitBreaker implements CircuitBreakerInterface for testing.
type MockCircuitBreaker struct {
	ExecuteFunc            func(operation func() error) error
	ExecuteWithContextFunc func(ctx context.Context, operation func(ctx context.Context) error) error
	GetStateFunc           func() string
	GetStatsFunc           func() map[string]interface{}
	ResetFunc              func()
}

// Execute executes an operation with circuit breaker logic.
func (m *MockCircuitBreaker) Execute(operation func() error) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(operation)
	}
	return operation() // Default behavior: execute operation
}

// ExecuteWithContext executes an operation with circuit breaker logic and context.
func (m *MockCircuitBreaker) ExecuteWithContext(ctx context.Context, operation func(ctx context.Context) error) error {
	if m.ExecuteWithContextFunc != nil {
		return m.ExecuteWithContextFunc(ctx, operation)
	}
	return operation(ctx) // Default behavior: execute operation
}

// GetState returns the current state as a string.
func (m *MockCircuitBreaker) GetState() string {
	if m.GetStateFunc != nil {
		return m.GetStateFunc()
	}
	return "CLOSED" // Default behavior: closed state
}

// GetStats returns statistics.
func (m *MockCircuitBreaker) GetStats() map[string]interface{} {
	if m.GetStatsFunc != nil {
		return m.GetStatsFunc()
	}
	return map[string]interface{}{
		"state": "CLOSED",
	} // Default behavior
}

// Reset resets the circuit breaker to closed state.
func (m *MockCircuitBreaker) Reset() {
	if m.ResetFunc != nil {
		m.ResetFunc()
	}
	// Default behavior: no-op
}

// NewMockCircuitBreaker creates a new mock CircuitBreaker.
func NewMockCircuitBreaker() *MockCircuitBreaker {
	return &MockCircuitBreaker{}
}

// SetupDefaultMockCircuitBreaker configures a mock with default behavior.
func SetupDefaultMockCircuitBreaker(m *MockCircuitBreaker) {
	// Default behavior is already set in the struct methods
	// This function exists for compatibility and can be extended if needed
}
