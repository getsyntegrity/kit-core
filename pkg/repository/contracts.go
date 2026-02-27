// Contracts: repository interfaces and base types.
// Repository, QueryableRepository (see queryable.go), and domain types.
package repository

import (
	"context"
	"fmt"
)

// Repository defines the interface for generic repository operations.
// CRITICAL: All methods require tenant_id - cross-tenant access is impossible by design.
type Repository[T any] interface {
	Save(ctx context.Context, entity T) error
	GetByID(ctx context.Context, tenantID, id string) (T, error)
	Exists(ctx context.Context, tenantID, id string) (bool, error)
	List(ctx context.Context, tenantID string, offset, limit int) ([]T, error)
	Count(ctx context.Context, tenantID string) (int64, error)
	Delete(ctx context.Context, tenantID, id string) error
}

// Pagination represents pagination parameters.
type Pagination struct {
	Offset int
	Limit  int
}

// Result represents a paginated result from repository.
type Result[T any] struct {
	Data  []T
	Total int64
	Page  Pagination
}

// BaseRepository provides a base implementation for repositories.
type BaseRepository[T any] struct{}

// NewBaseRepository creates a new base repository.
func NewBaseRepository[T any]() *BaseRepository[T] {
	return &BaseRepository[T]{}
}

// ValidatePagination validates pagination parameters.
func (r *BaseRepository[T]) ValidatePagination(offset, limit int) error {
	if offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}
	if limit <= 0 {
		return fmt.Errorf("limit must be positive")
	}
	if limit > 1000 {
		return fmt.Errorf("limit cannot exceed 1000")
	}
	return nil
}

// BuildPagination builds pagination parameters with defaults.
func (r *BaseRepository[T]) BuildPagination(offset, limit int) Pagination {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 1000 {
		limit = 1000
	}
	return Pagination{
		Offset: offset,
		Limit:  limit,
	}
}

// DomainEntity represents any domain entity.
type DomainEntity interface {
	GetID() string
	GetType() string
}

// DomainRepository defines the interface for read model operations with domain entities.
// CRITICAL: All methods require tenant_id - cross-tenant access is impossible by design.
type DomainRepository[T DomainEntity] interface {
	Save(ctx context.Context, entity T) error
	GetByID(ctx context.Context, tenantID, id string) (T, error)
	Exists(ctx context.Context, tenantID, id string) (bool, error)
	List(ctx context.Context, tenantID string, offset, limit int) ([]T, error)
	Count(ctx context.Context, tenantID string) (int64, error)
	Delete(ctx context.Context, tenantID, id string) error
}
