package timestamp

import (
	"time"
)

// TimestampService handles timestamp operations in the application layer.
// Time is provided via an injected TimeProvider; no default real clock in core.
type TimestampService struct {
	timeProvider TimeProvider
}

// NewTimestampServiceWithProvider creates a new timestamp service with the given time provider.
// Callers must inject a TimeProvider (e.g. from kit-runtime for real time).
func NewTimestampServiceWithProvider(provider TimeProvider) *TimestampService {
	return &TimestampService{
		timeProvider: provider,
	}
}

// GetCurrentTimestamp returns the current timestamp
func (s *TimestampService) GetCurrentTimestamp() time.Time {
	return s.timeProvider.Now()
}

// UpdateTimestamps updates created_at and updated_at for new records
func (s *TimestampService) UpdateTimestamps(createdAt, updatedAt *time.Time) {
	now := s.GetCurrentTimestamp()

	if createdAt != nil {
		*createdAt = now
	}
	if updatedAt != nil {
		*updatedAt = now
	}
}

// UpdateUpdatedAt updates only the updated_at timestamp for existing records
func (s *TimestampService) UpdateUpdatedAt(updatedAt *time.Time) {
	if updatedAt != nil {
		*updatedAt = s.GetCurrentTimestamp()
	}
}

// UpdateCreatedAt updates only the created_at timestamp for new records
func (s *TimestampService) UpdateCreatedAt(createdAt *time.Time) {
	if createdAt != nil {
		*createdAt = s.GetCurrentTimestamp()
	}
}

// IsAfter checks if the first timestamp is after the second
func (s *TimestampService) IsAfter(t1, t2 time.Time) bool {
	return t1.After(t2)
}

// IsBefore checks if the first timestamp is before the second
func (s *TimestampService) IsBefore(t1, t2 time.Time) bool {
	return t1.Before(t2)
}

// IsEqual checks if two timestamps are equal (within a small tolerance)
func (s *TimestampService) IsEqual(t1, t2 time.Time) bool {
	return t1.Equal(t2)
}

// DurationSince returns the duration since the given timestamp
func (s *TimestampService) DurationSince(t time.Time) time.Duration {
	return s.GetCurrentTimestamp().Sub(t)
}

// DurationUntil returns the duration until the given timestamp
func (s *TimestampService) DurationUntil(t time.Time) time.Duration {
	return t.Sub(s.GetCurrentTimestamp())
}

// AddDuration adds a duration to a timestamp
func (s *TimestampService) AddDuration(t time.Time, d time.Duration) time.Time {
	return t.Add(d)
}

// SubtractDuration subtracts a duration from a timestamp
func (s *TimestampService) SubtractDuration(t time.Time, d time.Duration) time.Time {
	return t.Add(-d)
}

// FormatTimestamp formats a timestamp using the given layout
func (s *TimestampService) FormatTimestamp(t time.Time, layout string) string {
	return t.Format(layout)
}

// ParseTimestamp parses a timestamp string using the given layout
func (s *TimestampService) ParseTimestamp(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

// GetTimestampInTimezone returns a timestamp in the specified timezone
func (s *TimestampService) GetTimestampInTimezone(t time.Time, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}

// GetCurrentTimestampInTimezone returns the current timestamp in the specified timezone
func (s *TimestampService) GetCurrentTimestampInTimezone(timezone string) (time.Time, error) {
	now := s.GetCurrentTimestamp()
	return s.GetTimestampInTimezone(now, timezone)
}

// IsTimestampValid checks if a timestamp is valid (not zero)
func (s *TimestampService) IsTimestampValid(t time.Time) bool {
	return !t.IsZero()
}

// GetZeroTimestamp returns a zero timestamp
func (s *TimestampService) GetZeroTimestamp() time.Time {
	return time.Time{}
}

// GetUnixTimestamp returns the Unix timestamp (seconds since epoch)
func (s *TimestampService) GetUnixTimestamp(t time.Time) int64 {
	return t.Unix()
}

// GetUnixNanoTimestamp returns the Unix timestamp in nanoseconds
func (s *TimestampService) GetUnixNanoTimestamp(t time.Time) int64 {
	return t.UnixNano()
}

// FromUnixTimestamp creates a timestamp from Unix timestamp (seconds since epoch)
func (s *TimestampService) FromUnixTimestamp(sec int64) time.Time {
	return time.Unix(sec, 0)
}

// FromUnixNanoTimestamp creates a timestamp from Unix timestamp in nanoseconds
func (s *TimestampService) FromUnixNanoTimestamp(nsec int64) time.Time {
	return time.Unix(0, nsec)
}

// GetCurrentUnixTimestamp returns the current Unix timestamp
func (s *TimestampService) GetCurrentUnixTimestamp() int64 {
	return s.GetUnixTimestamp(s.GetCurrentTimestamp())
}

// GetCurrentUnixNanoTimestamp returns the current Unix timestamp in nanoseconds
func (s *TimestampService) GetCurrentUnixNanoTimestamp() int64 {
	return s.GetUnixNanoTimestamp(s.GetCurrentTimestamp())
}

// SetTimeProvider sets a custom time provider (useful for testing)
func (s *TimestampService) SetTimeProvider(provider TimeProvider) {
	s.timeProvider = provider
}

// GetTimeProvider returns the current time provider
func (s *TimestampService) GetTimeProvider() TimeProvider {
	return s.timeProvider
}
