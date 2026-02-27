package resilience

import "time"

// Clock provides the current time. Callers must inject an implementation; kit-core does not use wall clock.
// A real implementation (e.g. time.Now()) lives in kit-runtime or the application.
type Clock interface {
	Now() time.Time
}

// FixedClock returns a fixed time. Use in tests or when deterministic time is required.
type FixedClock struct {
	Time time.Time
}

// Now returns the fixed time.
func (c FixedClock) Now() time.Time {
	return c.Time
}
