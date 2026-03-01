package timestamp

import "time"

// TimeProvider represents a time provider interface for testing.
type TimeProvider interface {
	Now() time.Time
}

// MockTimeProvider provides a mock time implementation for testing.
type MockTimeProvider struct {
	FixedTime time.Time
}

// Now returns the fixed time.
func (p *MockTimeProvider) Now() time.Time {
	return p.FixedTime
}

// NewMockTimeProvider creates a new mock time provider with the given fixed time.
func NewMockTimeProvider(fixedTime time.Time) *MockTimeProvider {
	return &MockTimeProvider{FixedTime: fixedTime}
}

// SetFixedTime sets the fixed time for the mock provider.
func (p *MockTimeProvider) SetFixedTime(t time.Time) {
	p.FixedTime = t
}

// AdvanceTime advances the fixed time by the given duration.
func (p *MockTimeProvider) AdvanceTime(d time.Duration) {
	p.FixedTime = p.FixedTime.Add(d)
}
