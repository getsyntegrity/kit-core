package domain

import (
	"context"

	"google.golang.org/protobuf/proto"
)

// EventStore defines the interface for event store operations.
type EventStore interface {
	// SaveEvents saves domain events to the event store
	SaveEvents(ctx context.Context, aggregateID string, aggregateType string, events []proto.Message) error

	// GetEvents retrieves events for an aggregate
	GetEvents(ctx context.Context, aggregateID string) ([]proto.Message, error)

	// GetEventsSince retrieves events for an aggregate since a specific version
	GetEventsSince(ctx context.Context, aggregateID string, version int64) ([]proto.Message, error)

	// GetLatestVersion gets the latest version for an aggregate
	GetLatestVersion(ctx context.Context, aggregateID string) (int64, error)
}

// EventRecord represents an event stored in the event store.
type EventRecord struct {
	ID            string
	AggregateID   string
	AggregateType string
	EventType     string
	EventData     []byte
	Version       int64
	Timestamp     int64
	Metadata      map[string]string
}

// EventSerializer defines how events are serialized/deserialized.
type EventSerializer interface {
	// SerializeEvent serializes an event to bytes
	SerializeEvent(event proto.Message) ([]byte, error)

	// DeserializeEvent deserializes bytes to an event
	DeserializeEvent(eventType string, data []byte) (proto.Message, error)

	// GetEventType returns the event type for a given event
	GetEventType(event proto.Message) string
}
