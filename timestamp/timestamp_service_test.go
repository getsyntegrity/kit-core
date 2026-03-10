package timestamp

import (
	"testing"
	"time"

	"github.com/pablogore/go-specs/specs"
	"github.com/stretchr/testify/assert"
)

func TestTimestamp(t *testing.T) {
	specs.Describe(t, "timestamp", func(s *specs.Spec) {
		s.When("NewTimestampService", func(s *specs.Spec) {
			s.It("returns non-nil service with time provider", func(ctx *specs.Context) {
				svc := NewTimestampService()
				assert.NotNil(ctx.T, svc)
				assert.NotNil(ctx.T, svc.timeProvider)
			})
		})

		s.When("NewTimestampServiceWithProvider", func(s *specs.Spec) {
			s.It("uses provided provider and returns fixed time from GetCurrentTimestamp", func(ctx *specs.Context) {
				fixed := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
				p := NewMockTimeProvider(fixed)
				svc := NewTimestampServiceWithProvider(p)
				assert.NotNil(ctx.T, svc)
				assert.Equal(ctx.T, fixed, svc.GetCurrentTimestamp())
			})
		})

		s.When("TimestampService_GetCurrentTimestamp", func(s *specs.Spec) {
			s.It("returns provider Now", func(ctx *specs.Context) {
				fixed := time.Date(2025, 2, 1, 12, 0, 0, 0, time.UTC)
				svc := NewTimestampServiceWithProvider(NewMockTimeProvider(fixed))
				assert.Equal(ctx.T, fixed, svc.GetCurrentTimestamp())
			})
		})

		s.When("MockTimeProvider", func(s *specs.Spec) {
			s.It("SetFixedTime updates returned time", func(ctx *specs.Context) {
				p := NewMockTimeProvider(time.Now())
				newTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
				p.SetFixedTime(newTime)
				assert.Equal(ctx.T, newTime, p.Now())
			})
			s.It("AdvanceTime advances fixed time", func(ctx *specs.Context) {
				t0 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
				p := NewMockTimeProvider(t0)
				p.AdvanceTime(24 * time.Hour)
				assert.Equal(ctx.T, t0.Add(24*time.Hour), p.Now())
			})
		})

		s.When("DefaultTimeProvider_Now", func(s *specs.Spec) {
			s.It("returns non-zero time", func(ctx *specs.Context) {
				p := &DefaultTimeProvider{}
				t1 := p.Now()
				assert.False(ctx.T, t1.IsZero())
			})
		})
	})
}
