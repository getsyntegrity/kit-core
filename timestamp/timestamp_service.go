package timestamp

import "time"

// TimestampService handles timestamp operations in the application layer.
//
//nolint:revive // exported: timestamp-prefixed name is intentional for clarity at call sites
type TimestampService struct {
	timeProvider TimeProvider
}

// DefaultTimeProvider provides the default time implementation.
type DefaultTimeProvider struct{}

// Now returns the current time.
func (p *DefaultTimeProvider) Now() time.Time {
	return time.Now()
}

// NewTimestampService creates a new timestamp service with default time provider.
func NewTimestampService() *TimestampService {
	return &TimestampService{
		timeProvider: &DefaultTimeProvider{},
	}
}

// NewTimestampServiceWithProvider creates a new timestamp service with custom time provider.
func NewTimestampServiceWithProvider(provider TimeProvider) *TimestampService {
	return &TimestampService{timeProvider: provider}
}

// GetCurrentTimestamp returns the current timestamp.
func (s *TimestampService) GetCurrentTimestamp() time.Time {
	return s.timeProvider.Now()
}
