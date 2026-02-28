// Package listeners provides generic listener infrastructure for event processing.
// This package provides the core listener implementation that projection listeners use.
package listeners

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Event represents a raw event from an event stream.
// This interface is implemented by infrastructure event types (e.g., eventstreamv2.EventMessage).
type Event interface {
	// Version returns the event version.
	Version() int64

	// Type returns the event type.
	Type() string
}

// RetryPolicy configures retry behavior for checkpoint persistence.
type RetryPolicy struct {
	// MaxAttempts is the maximum number of retry attempts (0 disables retries).
	MaxAttempts int

	// InitialBackoff is the initial backoff duration for exponential backoff.
	InitialBackoff time.Duration

	// MaxBackoff is the maximum backoff duration (caps exponential growth).
	MaxBackoff time.Duration
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

// GenericListener is the core listener implementation that handles event processing with retry logic.
type GenericListener struct {
	recvFunc              func(ctx context.Context) (Event, error)
	handlerFunc           func(ctx context.Context, event Event) error
	loadCheckpointFunc    func(ctx context.Context) (int64, error)
	saveCheckpointFunc    func(ctx context.Context, version int64) (bool, error)
	retryPolicy           RetryPolicy
	onFatalErrorFunc      func(ctx context.Context, err error)
	onCheckpointSavedFunc func(ctx context.Context, version int64, wasSaved bool)

	// Runtime state
	ctx    context.Context
	cancel context.CancelFunc

	// Statistics
	eventsProcessed  int64
	eventsFailed     int64
	checkpointsSaved int64
	lastEventTime    time.Time
}

// ListenerOption configures a GenericListener.
type ListenerOption func(*GenericListener) error

// WithRecv sets the function that receives events from the stream.
func WithRecv(recvFunc func(ctx context.Context) (Event, error)) ListenerOption {
	return func(l *GenericListener) error {
		if recvFunc == nil {
			return fmt.Errorf("recv function cannot be nil")
		}
		l.recvFunc = recvFunc
		return nil
	}
}

// WithHandler sets the function that handles events.
func WithHandler(handlerFunc func(ctx context.Context, event Event) error) ListenerOption {
	return func(l *GenericListener) error {
		if handlerFunc == nil {
			return fmt.Errorf("handler function cannot be nil")
		}
		l.handlerFunc = handlerFunc
		return nil
	}
}

// WithCheckpointStore sets the functions that load and save checkpoints.
func WithCheckpointStore(loadFunc func(ctx context.Context) (int64, error), saveFunc func(ctx context.Context, version int64) (bool, error)) ListenerOption {
	return func(l *GenericListener) error {
		if loadFunc == nil {
			return fmt.Errorf("load checkpoint function cannot be nil")
		}
		if saveFunc == nil {
			return fmt.Errorf("save checkpoint function cannot be nil")
		}
		l.loadCheckpointFunc = loadFunc
		l.saveCheckpointFunc = saveFunc
		return nil
	}
}

// WithRetryPolicy sets the retry policy for checkpoint persistence.
func WithRetryPolicy(policy RetryPolicy) ListenerOption {
	return func(l *GenericListener) error {
		l.retryPolicy = policy
		return nil
	}
}

// WithOnFatalError sets the callback for fatal errors.
func WithOnFatalError(onFatalErrorFunc func(ctx context.Context, err error)) ListenerOption {
	return func(l *GenericListener) error {
		l.onFatalErrorFunc = onFatalErrorFunc
		return nil
	}
}

// WithOnCheckpointSaved sets the callback for successful checkpoint saves.
func WithOnCheckpointSaved(onCheckpointSavedFunc func(ctx context.Context, version int64, wasSaved bool)) ListenerOption {
	return func(l *GenericListener) error {
		l.onCheckpointSavedFunc = onCheckpointSavedFunc
		return nil
	}
}

// New creates a new GenericListener with the provided options.
func New(opts ...ListenerOption) (*GenericListener, error) {
	l := &GenericListener{
		retryPolicy: DefaultRetryPolicy(),
	}

	for _, opt := range opts {
		if err := opt(l); err != nil {
			return nil, fmt.Errorf("invalid listener option: %w", err)
		}
	}

	// Validate required functions
	if l.recvFunc == nil {
		return nil, fmt.Errorf("recv function is required")
	}
	if l.handlerFunc == nil {
		return nil, fmt.Errorf("handler function is required")
	}
	if l.loadCheckpointFunc == nil {
		return nil, fmt.Errorf("load checkpoint function is required")
	}
	if l.saveCheckpointFunc == nil {
		return nil, fmt.Errorf("save checkpoint function is required")
	}

	return l, nil
}

// Run runs the listener, processing events until context is cancelled or a fatal error occurs.
func (l *GenericListener) Run(ctx context.Context) error {
	// Create listener context (owned by listener for cancellation)
	l.ctx, l.cancel = context.WithCancel(ctx)
	defer l.cancel()

	// Load initial checkpoint
	_, err := l.loadCheckpointFunc(l.ctx)
	if err != nil {
		if errors.Is(err, ErrNoCheckpoint) {
			// New projection: start from version 0
			// Checkpoint will be loaded on first event
		} else {
			return fmt.Errorf("load checkpoint: %w", err)
		}
	}

	// Process events
	for {
		// Check context cancellation
		select {
		case <-l.ctx.Done():
			return l.ctx.Err()
		default:
		}

		// Receive event
		event, err := l.recvFunc(l.ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return err
			}
			return fmt.Errorf("receive event: %w", err)
		}

		// Handle event
		if err := l.handlerFunc(l.ctx, event); err != nil {
			l.eventsFailed++
			return fmt.Errorf("handle event: %w", err)
		}
		l.eventsProcessed++
		l.lastEventTime = time.Now()

		// Save checkpoint with retry
		eventVersion := event.Version()
		wasSaved, err := l.saveCheckpointWithRetry(l.ctx, eventVersion)
		if err != nil {
			// Check if error is context cancellation (graceful shutdown)
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return err
			}
			// ErrVersionNotContiguous is expected for out-of-order events - continue processing
			if errors.Is(err, ErrVersionNotContiguous) {
				// Notify checkpoint saved (even though it wasn't saved due to non-contiguous version)
				if l.onCheckpointSavedFunc != nil {
					l.onCheckpointSavedFunc(l.ctx, eventVersion, false)
				}
				// Continue processing next event
				continue
			}
			// Fatal error: checkpoint save failed after all retries (only for ErrInvariantViolation or other errors)
			if l.onFatalErrorFunc != nil {
				l.onFatalErrorFunc(l.ctx, err)
			}
			return fmt.Errorf("persist checkpoint after retries: %w", err)
		}

		// Notify checkpoint saved
		if wasSaved {
			l.checkpointsSaved++
		}
		if l.onCheckpointSavedFunc != nil {
			l.onCheckpointSavedFunc(l.ctx, eventVersion, wasSaved)
		}
	}
}

// Stop stops the listener by cancelling its context.
// This returns immediately after cancelling context (does not wait for Run() to return).
func (l *GenericListener) Stop(ctx context.Context) error {
	if l.cancel != nil {
		l.cancel()
	}
	return nil
}

// GetStats returns listener statistics.
func (l *GenericListener) GetStats() ListenerStats {
	return ListenerStats{
		EventsProcessed:  l.eventsProcessed,
		EventsFailed:     l.eventsFailed,
		CheckpointsSaved: l.checkpointsSaved,
		LastEventTime:    l.lastEventTime,
	}
}

// saveCheckpointWithRetry saves the checkpoint with retry logic.
// Retry policy:
//   - Do not retry ErrVersionNotContiguous (expected out-of-order case)
//   - Do not retry ErrInvariantViolation (deterministic failure)
//   - Retry only ErrTransientInfrastructure
//   - Any other error is treated as non-retryable
func (l *GenericListener) saveCheckpointWithRetry(ctx context.Context, version int64) (bool, error) {
	var lastErr error
	backoff := l.retryPolicy.InitialBackoff

	for attempt := 0; attempt < l.retryPolicy.MaxAttempts; attempt++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
		}

		wasSaved, err := l.saveCheckpointFunc(ctx, version)
		if err == nil {
			return wasSaved, nil
		}

		// Non-retryable classes: return immediately.
		if errors.Is(err, ErrVersionNotContiguous) || errors.Is(err, ErrInvariantViolation) {
			return false, err
		}

		// Retry only transient infrastructure failures.
		if !isRetryableCheckpointError(err) {
			return false, err
		}

		lastErr = err

		// Exponential backoff (don't sleep on last attempt)
		if attempt < l.retryPolicy.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return false, ctx.Err()
			case <-time.After(backoff):
			}

			// Increase backoff for next attempt
			backoff *= 2
			if backoff > l.retryPolicy.MaxBackoff {
				backoff = l.retryPolicy.MaxBackoff
			}
		}
	}

	return false, fmt.Errorf("failed after %d attempts: %w", l.retryPolicy.MaxAttempts, lastErr)
}

// isRetryableCheckpointError returns true only for transient infrastructure failures.
func isRetryableCheckpointError(err error) bool {
	return errors.Is(err, ErrTransientInfrastructure)
}

// Errors
var (
	// ErrVersionNotContiguous indicates that the version is not contiguous (out-of-order event).
	// This is a non-retryable error - the listener should continue processing.
	ErrVersionNotContiguous = errors.New("version not contiguous")

	// ErrInvariantViolation indicates a deterministic invariant violation.
	// This is a non-retryable error - the listener should fail-fast.
	ErrInvariantViolation = errors.New("invariant violation")

	// ErrTransientInfrastructure indicates a transient infrastructure failure.
	// This is retryable by checkpoint persistence retry loop.
	ErrTransientInfrastructure = errors.New("transient infrastructure error")

	// ErrNoCheckpoint indicates that no checkpoint exists for the projection (new projection).
	ErrNoCheckpoint = errors.New("no checkpoint found")
)
