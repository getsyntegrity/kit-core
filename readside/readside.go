// Package readside provides read-side projection contracts (OffsetStore, Querier).
// Implementations (e.g. PostgresOffsetStore) live in kit-runtime.
package readside

import (
	"context"
	"errors"
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
	LoadLastVersion(ctx context.Context, projectionName string) (int64, error)
	SaveLastVersionIfContiguous(ctx context.Context, params SaveLastVersionIfContiguousParams) (int64, error)
}

// SaveLastVersionIfContiguousParams holds parameters for saving a contiguous version.
type SaveLastVersionIfContiguousParams struct {
	ProjectionName string
	LastVersion    int64
}

// Errors
var (
	ErrVersionNotContiguous         = errors.New("version not contiguous")
	ErrInsertFailedForNewProjection = errors.New("insert failed for new projection")
	ErrContiguousSaveFailed         = errors.New("contiguous save failed")
	ErrInvariantViolation           = errors.New("invariant violation")
	ErrNoCheckpoint                 = errors.New("no checkpoint found")
)
