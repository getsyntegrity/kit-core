package resilience

import (
	"time"
)

// Config represents the simplified resilience configuration.
type Config struct {
	Retry          RetryConfig          `yaml:"retry" json:"retry"`
	Timeout        TimeoutConfig        `yaml:"timeout" json:"timeout"`
	Fallback       FallbackConfig       `yaml:"fallback" json:"fallback"`
	CircuitBreaker CircuitBreakerConfig `yaml:"circuit_breaker" json:"circuit_breaker"`
}

// RetryConfig contains retry configuration.
type RetryConfig struct {
	Enabled     bool          `yaml:"enabled" json:"enabled"`
	MaxAttempts int           `yaml:"max_attempts" json:"max_attempts"`
	Delay       time.Duration `yaml:"delay" json:"delay"`
}

// HTTPRetryConfig contains HTTP-specific retry configuration with exponential backoff.
// This is used for HTTP client retry logic with configurable attempts, initial backoff, and max backoff.
type HTTPRetryConfig struct {
	// Attempts is the maximum number of retry attempts (0 disables retries)
	Attempts int `yaml:"attempts" json:"attempts"`
	// InitialBackoff is the initial backoff duration for exponential backoff
	InitialBackoff time.Duration `yaml:"initial_backoff" json:"initial_backoff"`
	// MaxBackoff is the maximum backoff duration (caps exponential growth)
	MaxBackoff time.Duration `yaml:"max_backoff" json:"max_backoff"`
}

// TimeoutConfig contains timeout configuration.
type TimeoutConfig struct {
	Enabled  bool          `yaml:"enabled" json:"enabled"`
	Duration time.Duration `yaml:"duration" json:"duration"`
}

// FallbackConfig contains fallback configuration.
type FallbackConfig struct {
	Enabled bool `yaml:"enabled" json:"enabled"`
}

// CircuitBreakerConfig contains circuit breaker configuration - KEEP IT SIMPLE.
type CircuitBreakerConfig struct {
	Enabled          bool          `yaml:"enabled" json:"enabled"`
	FailureThreshold int           `yaml:"failure_threshold" json:"failure_threshold"`
	Timeout          time.Duration `yaml:"timeout" json:"timeout"`
}

// Option is a function that configures a Config - KEEP IT SIMPLE.
type Option func(*Config)

// WithRetryConfig enables retry with simple parameters.
func WithRetryConfig(maxAttempts int, delay time.Duration) Option {
	return func(c *Config) {
		c.Retry = RetryConfig{
			Enabled:     true,
			MaxAttempts: maxAttempts,
			Delay:       delay,
		}
	}
}

// WithTimeoutConfig enables timeout with simple parameters.
func WithTimeoutConfig(duration time.Duration) Option {
	return func(c *Config) {
		c.Timeout = TimeoutConfig{
			Enabled:  true,
			Duration: duration,
		}
	}
}

// WithFallbackConfig enables fallback.
func WithFallbackConfig() Option {
	return func(c *Config) {
		c.Fallback = FallbackConfig{
			Enabled: true,
		}
	}
}

// WithCircuitBreakerConfig enables circuit breaker with simple parameters.
func WithCircuitBreakerConfig(failureThreshold int, timeout time.Duration) Option {
	return func(c *Config) {
		c.CircuitBreaker = CircuitBreakerConfig{
			Enabled:          true,
			FailureThreshold: failureThreshold,
			Timeout:          timeout,
		}
	}
}

// NewConfig creates a new configuration with the given options.
func NewConfig(options ...Option) Config {
	config := Config{
		Retry: RetryConfig{
			Enabled:     false,
			MaxAttempts: 3,
			Delay:       100 * time.Millisecond,
		},
		Timeout: TimeoutConfig{
			Enabled:  false,
			Duration: 30 * time.Second,
		},
		Fallback: FallbackConfig{
			Enabled: false,
		},
		CircuitBreaker: CircuitBreakerConfig{
			Enabled:          false,
			FailureThreshold: 5,
			Timeout:          60 * time.Second,
		},
	}

	for _, option := range options {
		option(&config)
	}

	return config
}

// DefaultConfig returns a default resilience configuration.
func DefaultConfig() Config {
	return NewConfig()
}

// IsEnabled returns true if any resilience pattern is enabled.
func (c Config) IsEnabled() bool {
	return c.Retry.Enabled || c.Timeout.Enabled || c.Fallback.Enabled || c.CircuitBreaker.Enabled
}
