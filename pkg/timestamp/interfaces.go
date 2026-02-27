package timestamp

import "time"

// TimestampServiceInterface represents the interface for timestamp operations
type TimestampServiceInterface interface {
	// GetCurrentTimestamp returns the current timestamp
	GetCurrentTimestamp() time.Time

	// UpdateTimestamps updates created_at and updated_at for new records
	UpdateTimestamps(createdAt, updatedAt *time.Time)

	// UpdateUpdatedAt updates only the updated_at timestamp for existing records
	UpdateUpdatedAt(updatedAt *time.Time)

	// UpdateCreatedAt updates only the created_at timestamp for new records
	UpdateCreatedAt(createdAt *time.Time)

	// IsAfter checks if the first timestamp is after the second
	IsAfter(t1, t2 time.Time) bool

	// IsBefore checks if the first timestamp is before the second
	IsBefore(t1, t2 time.Time) bool

	// IsEqual checks if two timestamps are equal (within a small tolerance)
	IsEqual(t1, t2 time.Time) bool

	// DurationSince returns the duration since the given timestamp
	DurationSince(t time.Time) time.Duration

	// DurationUntil returns the duration until the given timestamp
	DurationUntil(t time.Time) time.Duration

	// AddDuration adds a duration to a timestamp
	AddDuration(t time.Time, d time.Duration) time.Time

	// SubtractDuration subtracts a duration from a timestamp
	SubtractDuration(t time.Time, d time.Duration) time.Time

	// FormatTimestamp formats a timestamp using the given layout
	FormatTimestamp(t time.Time, layout string) string

	// ParseTimestamp parses a timestamp string using the given layout
	ParseTimestamp(layout, value string) (time.Time, error)

	// GetTimestampInTimezone returns a timestamp in the specified timezone
	GetTimestampInTimezone(t time.Time, timezone string) (time.Time, error)

	// GetCurrentTimestampInTimezone returns the current timestamp in the specified timezone
	GetCurrentTimestampInTimezone(timezone string) (time.Time, error)

	// IsTimestampValid checks if a timestamp is valid (not zero)
	IsTimestampValid(t time.Time) bool

	// GetZeroTimestamp returns a zero timestamp
	GetZeroTimestamp() time.Time

	// GetUnixTimestamp returns the Unix timestamp (seconds since epoch)
	GetUnixTimestamp(t time.Time) int64

	// GetUnixNanoTimestamp returns the Unix timestamp in nanoseconds
	GetUnixNanoTimestamp(t time.Time) int64

	// FromUnixTimestamp creates a timestamp from Unix timestamp (seconds since epoch)
	FromUnixTimestamp(sec int64) time.Time

	// FromUnixNanoTimestamp creates a timestamp from Unix timestamp in nanoseconds
	FromUnixNanoTimestamp(nsec int64) time.Time

	// GetCurrentUnixTimestamp returns the current Unix timestamp
	GetCurrentUnixTimestamp() int64

	// GetCurrentUnixNanoTimestamp returns the current Unix timestamp in nanoseconds
	GetCurrentUnixNanoTimestamp() int64
}

// TimeProvider represents a time provider interface for testing
type TimeProvider interface {
	// Now returns the current time
	Now() time.Time
}

// MockTimeProvider provides a mock time implementation for testing
type MockTimeProvider struct {
	FixedTime time.Time
}

// Now returns the fixed time
func (p *MockTimeProvider) Now() time.Time {
	return p.FixedTime
}

// NewMockTimeProvider creates a new mock time provider with the given fixed time
func NewMockTimeProvider(fixedTime time.Time) *MockTimeProvider {
	return &MockTimeProvider{
		FixedTime: fixedTime,
	}
}

// SetFixedTime sets the fixed time for the mock provider
func (p *MockTimeProvider) SetFixedTime(t time.Time) {
	p.FixedTime = t
}

// AdvanceTime advances the fixed time by the given duration
func (p *MockTimeProvider) AdvanceTime(d time.Duration) {
	p.FixedTime = p.FixedTime.Add(d)
}

// TimestampFormatter represents a timestamp formatter interface
type TimestampFormatter interface {
	// Format formats a timestamp using the default layout
	Format(t time.Time) string

	// FormatWithLayout formats a timestamp using the given layout
	FormatWithLayout(t time.Time, layout string) string

	// Parse parses a timestamp string using the default layout
	Parse(value string) (time.Time, error)

	// ParseWithLayout parses a timestamp string using the given layout
	ParseWithLayout(layout, value string) (time.Time, error)
}

// DefaultTimestampFormatter provides default timestamp formatting
type DefaultTimestampFormatter struct {
	DefaultLayout string
}

// NewDefaultTimestampFormatter creates a new default timestamp formatter
func NewDefaultTimestampFormatter() *DefaultTimestampFormatter {
	return &DefaultTimestampFormatter{
		DefaultLayout: time.RFC3339,
	}
}

// NewDefaultTimestampFormatterWithLayout creates a new default timestamp formatter with custom layout
func NewDefaultTimestampFormatterWithLayout(layout string) *DefaultTimestampFormatter {
	return &DefaultTimestampFormatter{
		DefaultLayout: layout,
	}
}

// Format formats a timestamp using the default layout
func (f *DefaultTimestampFormatter) Format(t time.Time) string {
	return t.Format(f.DefaultLayout)
}

// FormatWithLayout formats a timestamp using the given layout
func (f *DefaultTimestampFormatter) FormatWithLayout(t time.Time, layout string) string {
	return t.Format(layout)
}

// Parse parses a timestamp string using the default layout
func (f *DefaultTimestampFormatter) Parse(value string) (time.Time, error) {
	return time.Parse(f.DefaultLayout, value)
}

// ParseWithLayout parses a timestamp string using the given layout
func (f *DefaultTimestampFormatter) ParseWithLayout(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

// SetDefaultLayout sets the default layout for formatting
func (f *DefaultTimestampFormatter) SetDefaultLayout(layout string) {
	f.DefaultLayout = layout
}

// GetDefaultLayout returns the default layout
func (f *DefaultTimestampFormatter) GetDefaultLayout() string {
	return f.DefaultLayout
}
