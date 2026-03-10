// Package fake provides deterministic in-memory implementations of interfaces for testing.
package fake

import (
	"sync"
	"time"

	"github.com/getsyntegrity/kit-core/clock"
)

// Clock is a deterministic clock for tests. Now() returns the current fixed time;
// Advance() advances it. After(d) returns a channel that fires after d from the
// current fixed time (simulated by advancing and closing).
// Thread-safe.
type Clock struct {
	mu  sync.Mutex
	now time.Time
}

// NewClock returns a fake clock with the given initial time.
func NewClock(initial time.Time) *Clock {
	return &Clock{now: initial}
}

// Now returns the current fixed time.
func (c *Clock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

// Advance advances the clock by d.
func (c *Clock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

// Set sets the clock to t.
func (c *Clock) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = t
}

// After returns a channel that receives the current time after duration d.
// Deterministic: clock is advanced by d, then the new time is sent on a buffered channel and closed.
func (c *Clock) After(d time.Duration) <-chan time.Time {
	ch := make(chan time.Time, 1)
	c.mu.Lock()
	c.now = c.now.Add(d)
	t := c.now
	c.mu.Unlock()
	ch <- t
	close(ch)
	return ch
}

// Ensure Clock implements clock.Clock.
var _ clock.Clock = (*Clock)(nil)
