package domain

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

// Entity represents any domain entity.
type Entity interface {
	GetID() string
	GetType() string
}

// MockableEventMessage represents an event message that can be mocked (without proto.Message embedding).
type MockableEventMessage interface {
	GetEventType() string
	GetEventData() []byte
}

// MockableBaseRepository defines the interface for basic repository operations (mockable version).
type MockableBaseRepository[T any] interface {
	Save(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id string) (*T, error)
	Exists(ctx context.Context, id string) (bool, error)
	List(ctx context.Context, offset, limit int) ([]*T, error)
	Count(ctx context.Context) (int64, error)
	Delete(ctx context.Context, id string) error
}

// MockableReadRepository defines the interface for read-only repository operations (mockable version).
type MockableReadRepository[T any] interface {
	GetByID(ctx context.Context, id string) (*T, error)
	Exists(ctx context.Context, id string) (bool, error)
	List(ctx context.Context, offset, limit int) ([]*T, error)
	Count(ctx context.Context) (int64, error)
}

// MockableWriteRepository defines the interface for write-only repository operations (mockable version).
type MockableWriteRepository[T any] interface {
	Save(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error
}

// MockableRepositoryFactory creates repository instances (mockable version).
type MockableRepositoryFactory interface {
	CreateRepository(config RepositoryConfig) (interface{}, error)
}

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
	GetRepository() interface{}
	GetCommandBus() CommandBus
	GetSnapshotProcessor() interface{}
}
