package domain

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

// =============================================================================
// DOMAIN ENTITY INTERFACES
// =============================================================================

// Entity represents any domain entity.
type Entity interface {
	GetID() string
	GetType() string
}

// MockableEventMessage represents an event message that can be mocked (without proto.Message embedding).
type MockableEventMessage interface {
	// Add any methods that need to be mocked
	GetEventType() string
	GetEventData() []byte
}

// MockableBaseRepository defines the interface for basic repository operations (mockable version).
type MockableBaseRepository[T any] interface {
	// Save saves an entity to the repository
	Save(ctx context.Context, entity *T) error

	// GetByID retrieves an entity by its ID
	GetByID(ctx context.Context, id string) (*T, error)

	// Exists checks if an entity exists with the given ID
	Exists(ctx context.Context, id string) (bool, error)

	// List retrieves a list of entities with pagination
	List(ctx context.Context, offset, limit int) ([]*T, error)

	// Count returns the total number of entities
	Count(ctx context.Context) (int64, error)

	// Delete removes an entity by its ID
	Delete(ctx context.Context, id string) error
}

// MockableReadRepository defines the interface for read-only repository operations (mockable version).
type MockableReadRepository[T any] interface {
	// GetByID retrieves an entity by its ID
	GetByID(ctx context.Context, id string) (*T, error)

	// Exists checks if an entity exists with the given ID
	Exists(ctx context.Context, id string) (bool, error)

	// List retrieves a list of entities with pagination
	List(ctx context.Context, offset, limit int) ([]*T, error)

	// Count returns the total number of entities
	Count(ctx context.Context) (int64, error)
}

// MockableWriteRepository defines the interface for write-only repository operations (mockable version).
type MockableWriteRepository[T any] interface {
	// Save saves an entity to the repository
	Save(ctx context.Context, entity *T) error

	// Delete removes an entity by its ID
	Delete(ctx context.Context, id string) error
}

// MockableRepositoryFactory creates repository instances (mockable version).
type MockableRepositoryFactory interface {
	// CreateRepository creates a new repository instance
	CreateRepository(config RepositoryConfig) (interface{}, error)
}

// =============================================================================
// COMMAND INTERFACES
// =============================================================================

// CommandHandler defines the interface for command handlers with generics.
type CommandHandler[T Entity] interface {
	Handle(ctx context.Context, cmd Command[T]) error
}

// Command defines the interface for commands with generics.
type Command[T Entity] interface {
	CommandType() string
	AggregateID() string
	GetEntity() T
}

// =============================================================================
// INFRASTRUCTURE INTERFACES (from interfaces package)
// =============================================================================

// CommandBus defines the interface for command bus operations.
type CommandBus interface {
	Send(ctx context.Context, cmd interface{}) error
}

// EventBus defines the interface for event bus operations.
type EventBus interface {
	Publish(topic string, messages ...*message.Message) error
	Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
	Close() error
}

// ApplicationContainer defines the interface for application dependencies with generics.
type ApplicationContainer[T Entity] interface {
	GetEventStore() EventStore
	GetEventBus() EventBus
	GetRepository() interface{} // Generic repository interface

	GetCommandBus() CommandBus
	GetSnapshotProcessor() interface{}
}
