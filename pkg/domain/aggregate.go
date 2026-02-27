// Package domain provides generic domain patterns and interfaces for event sourcing.
package domain

import (
	"google.golang.org/protobuf/proto"
)

// Aggregate represents a domain aggregate that can produce and consume events.
type Aggregate interface {
	// GetID returns the aggregate's unique identifier
	GetID() string

	// GetUncommittedEvents returns all uncommitted events
	GetUncommittedEvents() []proto.Message

	// MarkEventsAsCommitted marks all events as committed
	MarkEventsAsCommitted()

	// LoadFromHistory loads the aggregate from a series of events
	LoadFromHistory(events []proto.Message) error

	// GetVersion returns the current version of the aggregate
	GetVersion() int64
}

// BaseAggregate provides a base implementation for aggregates.
type BaseAggregate struct {
	id      string
	version int64
	events  []proto.Message
}

// NewBaseAggregate creates a new base aggregate.
func NewBaseAggregate(id string) *BaseAggregate {
	return &BaseAggregate{
		id:      id,
		version: 0,
		events:  make([]proto.Message, 0),
	}
}

// GetID returns the aggregate's unique identifier.
func (a *BaseAggregate) GetID() string {
	return a.id
}

// GetVersion returns the current version of the aggregate.
func (a *BaseAggregate) GetVersion() int64 {
	return a.version
}

// GetUncommittedEvents returns all uncommitted events.
func (a *BaseAggregate) GetUncommittedEvents() []proto.Message {
	return a.events
}

// MarkEventsAsCommitted marks all events as committed.
func (a *BaseAggregate) MarkEventsAsCommitted() {
	a.events = make([]proto.Message, 0)
}

// AddEvent adds an event to the uncommitted events list.
func (a *BaseAggregate) AddEvent(event proto.Message) {
	a.events = append(a.events, event)
	a.version++
}

// LoadFromHistory loads the aggregate from a series of events.
func (a *BaseAggregate) LoadFromHistory(events []proto.Message) error {
	for _, event := range events {
		if err := a.applyEvent(event); err != nil {
			return err
		}
		a.version++
	}
	return nil
}

// applyEvent applies a single event to the aggregate
// This method should be overridden by specific aggregates.
func (a *BaseAggregate) applyEvent(_ proto.Message) error {
	// Base implementation does nothing
	// Specific aggregates should override this method
	return nil
}

// SetID sets the aggregate's unique identifier.
func (a *BaseAggregate) SetID(id string) {
	a.id = id
}

// SetVersion sets the aggregate's version.
func (a *BaseAggregate) SetVersion(version int64) {
	a.version = version
}
