// Package readside provides read-side projection infrastructure.
// This package provides offset store implementations for checkpoint persistence.
package readside

import (
	"context"
	"errors"
	"fmt"
)

// OffsetStore defines the interface for checkpoint persistence.
type OffsetStore interface {
	// LoadLastVersion loads the last processed version for a projection.
	// Returns (version, nil) if projection exists, or (0, ErrNoCheckpoint) if projection has never been processed.
	LoadLastVersion(ctx context.Context, projectionName string) (int64, error)

	// SaveLastVersion saves the last processed version checkpoint for a projection.
	// Returns (wasSaved, error) where wasSaved indicates if a new version was actually saved.
	SaveLastVersion(ctx context.Context, projectionName string, version int64) (bool, error)
}

// Querier defines the interface for database queries used by offset stores.
type Querier interface {
	// LoadLastVersion loads the last processed version for a projection.
	// Returns (version, nil) if projection exists, or (0, ErrNoCheckpoint) if projection has never been processed.
	LoadLastVersion(ctx context.Context, projectionName string) (int64, error)

	// SaveLastVersionIfContiguous saves a version checkpoint if it's contiguous.
	// Returns the saved version or an error if the version is not contiguous.
	SaveLastVersionIfContiguous(ctx context.Context, params SaveLastVersionIfContiguousParams) (int64, error)
}

// SaveLastVersionIfContiguousParams holds parameters for saving a contiguous version.
type SaveLastVersionIfContiguousParams struct {
	// ProjectionName is the name of the projection.
	ProjectionName string

	// LastVersion is the version to save (must be contiguous with existing version).
	LastVersion int64
}

// PostgresOffsetStore implements OffsetStore using Postgres and a Querier.
type PostgresOffsetStore struct {
	querier Querier
}

// NewPostgresOffsetStore creates a new Postgres-based offset store.
func NewPostgresOffsetStore(querier Querier) OffsetStore {
	if querier == nil {
		panic("querier cannot be nil")
	}
	return &PostgresOffsetStore{
		querier: querier,
	}
}

// LoadLastVersion loads the last processed version for a projection.
func (s *PostgresOffsetStore) LoadLastVersion(ctx context.Context, projectionName string) (int64, error) {
	return s.querier.LoadLastVersion(ctx, projectionName)
}

// SaveLastVersion saves the last processed version checkpoint for a projection.
func (s *PostgresOffsetStore) SaveLastVersion(ctx context.Context, projectionName string, version int64) (bool, error) {
	// Try to save as contiguous
	savedVersion, err := s.querier.SaveLastVersionIfContiguous(ctx, SaveLastVersionIfContiguousParams{
		ProjectionName: projectionName,
		LastVersion:    version,
	})
	if err != nil {
		if errors.Is(err, ErrVersionNotContiguous) {
			// Version not contiguous (out-of-order event) - not an error, just return false
			return false, nil
		}
		// Domain errors should be propagated without wrapping to prevent raw DB errors from escaping
		if errors.Is(err, ErrInvariantViolation) || errors.Is(err, ErrInsertFailedForNewProjection) {
			return false, err
		}
		// Other errors (DB failures, etc.)
		return false, fmt.Errorf("save last version: %w", err)
	}

	// Check if version was actually saved (savedVersion == version means it was saved)
	wasSaved := savedVersion == version
	return wasSaved, nil
}

// Errors
var (
	// ErrVersionNotContiguous indicates that the version is not contiguous (out-of-order event).
	ErrVersionNotContiguous = errors.New("version not contiguous")

	// ErrInsertFailedForNewProjection indicates that inserting a new projection failed.
	ErrInsertFailedForNewProjection = errors.New("insert failed for new projection")

	// ErrContiguousSaveFailed indicates that saving a contiguous version failed.
	ErrContiguousSaveFailed = errors.New("contiguous save failed")

	// ErrInvariantViolation indicates a deterministic invariant violation.
	ErrInvariantViolation = errors.New("invariant violation")

	// ErrNoCheckpoint indicates that no checkpoint exists for the projection (new projection).
	ErrNoCheckpoint = errors.New("no checkpoint found")
)
