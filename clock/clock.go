package clock

import "time"

// Clock represents the time abstraction (hexagonal port).
// Implementations live outside kit-core (e.g. in runtime).
// Allows getting current time and scheduling after a duration.
type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}
