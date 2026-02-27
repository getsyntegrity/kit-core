package domain

import (
	"context"
	"errors"

	"google.golang.org/protobuf/proto"
)

// Domain errors.
var (
	ErrEntityNotFound = errors.New("entity not found")
)

// EventMessage represents a protobuf event message.
type EventMessage interface {
	proto.Message
}

// BaseRepository defines the interface for basic repository operations.
// CRITICAL: All methods require tenant_id - cross-tenant access is impossible by design.
type BaseRepository[T any] interface {
	// Save saves an entity to the repository (tenant_id extracted from entity)
	Save(ctx context.Context, entity *T) error

	// GetByID retrieves an entity by its ID within a specific tenant
	// CRITICAL: tenant_id is REQUIRED
	GetByID(ctx context.Context, tenantID, id string) (*T, error)

	// Exists checks if an entity exists within a specific tenant
	// CRITICAL: tenant_id is REQUIRED
	Exists(ctx context.Context, tenantID, id string) (bool, error)

	// List retrieves a list of entities within a specific tenant
	// CRITICAL: tenant_id is REQUIRED
	List(ctx context.Context, tenantID string, offset, limit int) ([]*T, error)

	// Count returns the total number of entities within a specific tenant
	// CRITICAL: tenant_id is REQUIRED
	Count(ctx context.Context, tenantID string) (int64, error)

	// Delete removes an entity by its ID within a specific tenant
	// CRITICAL: tenant_id is REQUIRED
	Delete(ctx context.Context, tenantID, id string) error
}

// ReadRepository defines the interface for read-only repository operations.
// CRITICAL: All methods require tenant_id - cross-tenant access is impossible by design.
type ReadRepository[T any] interface {
	// GetByID retrieves an entity by its ID within a specific tenant
	// CRITICAL: tenant_id is REQUIRED
	GetByID(ctx context.Context, tenantID, id string) (*T, error)

	// Exists checks if an entity exists within a specific tenant
	// CRITICAL: tenant_id is REQUIRED
	Exists(ctx context.Context, tenantID, id string) (bool, error)

	// List retrieves a list of entities within a specific tenant
	// CRITICAL: tenant_id is REQUIRED
	List(ctx context.Context, tenantID string, offset, limit int) ([]*T, error)

	// Count returns the total number of entities within a specific tenant
	// CRITICAL: tenant_id is REQUIRED
	Count(ctx context.Context, tenantID string) (int64, error)
}

// WriteRepository defines the interface for write-only repository operations.
type WriteRepository[T any] interface {
	// Save saves an entity to the repository (tenant_id extracted from entity)
	Save(ctx context.Context, entity *T) error

	// Delete removes an entity by its ID within a specific tenant
	// CRITICAL: tenant_id is REQUIRED
	Delete(ctx context.Context, tenantID, id string) error
}

// RepositoryFactory creates repository instances.
type RepositoryFactory interface {
	// CreateRepository creates a new repository instance
	CreateRepository(config RepositoryConfig) (interface{}, error)
}

// RepositoryConfig holds configuration for repositories.
type RepositoryConfig struct {
	Type     string            // "postgres", "mysql", "inmemory"
	DSN      string            // Database connection string
	Table    string            // Table name
	Metadata map[string]string // Additional metadata
}

// =============================================================================
// TENANT-SCOPED REPOSITORY INTERFACES
// =============================================================================

// TenantScopedRepository defines the interface for tenant-scoped repository operations.
// All methods REQUIRE tenant_id - cross-tenant access is impossible by design.
type TenantScopedRepository[T any] interface {
	// GetByID retrieves an entity by its ID within a specific tenant
	GetByID(ctx context.Context, tenantID, id string) (*T, error)

	// Exists checks if an entity exists within a specific tenant
	Exists(ctx context.Context, tenantID, id string) (bool, error)

	// List retrieves a list of entities within a specific tenant
	List(ctx context.Context, tenantID string, offset, limit int) ([]*T, error)

	// Count returns the total number of entities within a specific tenant
	Count(ctx context.Context, tenantID string) (int64, error)

	// Save saves an entity to the repository (tenant_id must match)
	Save(ctx context.Context, tenantID string, entity *T) error

	// Delete removes an entity by its ID within a specific tenant
	Delete(ctx context.Context, tenantID, id string) error
}

// TenantScopedReadRepository defines the interface for tenant-scoped read operations.
type TenantScopedReadRepository[T any] interface {
	// GetByID retrieves an entity by its ID within a specific tenant
	GetByID(ctx context.Context, tenantID, id string) (*T, error)

	// Exists checks if an entity exists within a specific tenant
	Exists(ctx context.Context, tenantID, id string) (bool, error)

	// List retrieves a list of entities within a specific tenant
	List(ctx context.Context, tenantID string, offset, limit int) ([]*T, error)

	// Count returns the total number of entities within a specific tenant
	Count(ctx context.Context, tenantID string) (int64, error)
}

// TenantScopedWriteRepository defines the interface for tenant-scoped write operations.
type TenantScopedWriteRepository[T any] interface {
	// Save saves an entity to the repository (tenant_id must match)
	Save(ctx context.Context, tenantID string, entity *T) error

	// Delete removes an entity by its ID within a specific tenant
	Delete(ctx context.Context, tenantID, id string) error
}
