// Transaction: driver-agnostic unit-of-work contract.
// Implementations (e.g. PostgreSQL) live in kit-runtime; no driver types here.
package repository

import (
	"context"
	"time"
)

// TransactionOptions holds driver-agnostic transaction options.
// No database-specific enums or types.
type TransactionOptions struct {
	ReadOnly bool
	Timeout  time.Duration
}

// UnitOfWork runs a function in a single transactional unit.
// Commit on success, rollback on error or panic.
// Implementations use real drivers (e.g. pgx) in kit-runtime only.
type UnitOfWork interface {
	Do(ctx context.Context, opts TransactionOptions, fn func(ctx context.Context) error) error
}
