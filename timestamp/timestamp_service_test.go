package timestamp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTimestampService(t *testing.T) {
	s := NewTimestampService()
	assert.NotNil(t, s)
	assert.NotNil(t, s.timeProvider)
}

func TestNewTimestampServiceWithProvider(t *testing.T) {
	fixed := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	p := NewMockTimeProvider(fixed)
	s := NewTimestampServiceWithProvider(p)
	assert.NotNil(t, s)
	assert.Equal(t, fixed, s.GetCurrentTimestamp())
}

func TestTimestampService_GetCurrentTimestamp(t *testing.T) {
	fixed := time.Date(2025, 2, 1, 12, 0, 0, 0, time.UTC)
	s := NewTimestampServiceWithProvider(NewMockTimeProvider(fixed))
	assert.Equal(t, fixed, s.GetCurrentTimestamp())
}

func TestMockTimeProvider_SetFixedTime(t *testing.T) {
	p := NewMockTimeProvider(time.Now())
	newTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	p.SetFixedTime(newTime)
	assert.Equal(t, newTime, p.Now())
}

func TestMockTimeProvider_AdvanceTime(t *testing.T) {
	t0 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	p := NewMockTimeProvider(t0)
	p.AdvanceTime(24 * time.Hour)
	assert.Equal(t, t0.Add(24*time.Hour), p.Now())
}

func TestDefaultTimeProvider_Now(t *testing.T) {
	p := &DefaultTimeProvider{}
	t1 := p.Now()
	assert.False(t, t1.IsZero())
}
