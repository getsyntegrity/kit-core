// Package listeners provides listener contracts for event processing.
package listeners

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Event represents a raw event from an event stream.
type Event interface {
	Version() int64
	Type() string
}

// RetryPolicy configures retry behavior for checkpoint persistence.
type RetryPolicy struct {
	MaxAttempts    int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

// DefaultRetryPolicy returns the default retry policy for checkpoint persistence.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxAttempts:    5,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     2 * time.Second,
	}
}

// ListenerStats holds statistics for the generic listener.
type ListenerStats struct {
	EventsProcessed  int64
	EventsFailed     int64
	CheckpointsSaved int64
	LastEventTime    time.Time
}

// ListenerOption configures a GenericListener (implementation lives in kit-runtime).
type ListenerOption func(interface{}) error

// WithRecv sets the function that receives events from the stream.
func WithRecv(recvFunc func(ctx context.Context) (Event, error)) ListenerOption {
	return func(l interface{}) error {
		if recvFunc == nil {
			return fmt.Errorf("recv function cannot be nil")
		}
		return nil
	}
}

// WithHandler sets the function that handles events.
func WithHandler(handlerFunc func(ctx context.Context, event Event) error) ListenerOption {
	return func(l interface{}) error {
		if handlerFunc == nil {
			return fmt.Errorf("handler function cannot be nil")
		}
		return nil
	}
}

// WithCheckpointStore sets the functions that load and save checkpoints.
func WithCheckpointStore(loadFunc func(ctx context.Context) (int64, error), saveFunc func(ctx context.Context, version int64) (bool, error)) ListenerOption {
	return func(l interface{}) error {
		if loadFunc == nil || saveFunc == nil {
			return fmt.Errorf("load and save checkpoint functions cannot be nil")
		}
		return nil
	}
}

// WithRetryPolicy sets the retry policy for checkpoint persistence.
func WithRetryPolicy(policy RetryPolicy) ListenerOption {
	return func(l interface{}) error {
		return nil
	}
}

// Errors
var (
	ErrVersionNotContiguous    = errors.New("version not contiguous")
	ErrInvariantViolation      = errors.New("invariant violation")
	ErrTransientInfrastructure = errors.New("transient infrastructure error")
	ErrNoCheckpoint            = errors.New("no checkpoint found")
)
