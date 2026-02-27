package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SimpleEvent is a simple event implementation for testing.
type SimpleEvent struct {
	EventType string
	Data      string
}

// SimpleAggregate is a simple aggregate for testing.
type SimpleAggregate struct {
	id      string
	version int64
	events  []*SimpleEvent
	Name    string
}

// NewSimpleAggregate creates a new simple aggregate.
func NewSimpleAggregate(id, name string) *SimpleAggregate {
	agg := &SimpleAggregate{
		id:      id,
		version: 0,
		events:  make([]*SimpleEvent, 0),
		Name:    name,
	}

	// Add initial event
	event := &SimpleEvent{
		EventType: "Created",
		Data:      name,
	}
	agg.AddEvent(event)

	return agg
}

// GetID returns the aggregate's unique identifier.
func (a *SimpleAggregate) GetID() string {
	return a.id
}

// GetVersion returns the current version of the aggregate.
func (a *SimpleAggregate) GetVersion() int64 {
	return a.version
}

// GetUncommittedEvents returns all uncommitted events.
func (a *SimpleAggregate) GetUncommittedEvents() []*SimpleEvent {
	return a.events
}

// MarkEventsAsCommitted marks all events as committed.
func (a *SimpleAggregate) MarkEventsAsCommitted() {
	a.events = make([]*SimpleEvent, 0)
}

// AddEvent adds an event to the uncommitted events list.
func (a *SimpleAggregate) AddEvent(event *SimpleEvent) {
	a.events = append(a.events, event)
	a.version++
}

// Update updates the aggregate.
func (a *SimpleAggregate) Update(name string) {
	a.Name = name
	event := &SimpleEvent{
		EventType: "Updated",
		Data:      name,
	}
	a.AddEvent(event)
}

// applyEvent applies events to the aggregate.
func (a *SimpleAggregate) applyEvent(event *SimpleEvent) error {
	switch event.EventType {
	case "Created":
		a.Name = event.Data
	case "Updated":
		a.Name = event.Data
	}
	return nil
}

// LoadFromHistory loads the aggregate from a series of events.
func (a *SimpleAggregate) LoadFromHistory(events []*SimpleEvent) error {
	for _, event := range events {
		if err := a.applyEvent(event); err != nil {
			return err
		}
		a.version++
	}
	return nil
}

// SimpleRepository is a simple in-memory repository.
type SimpleRepository struct {
	entities map[string]*SimpleAggregate
}

// NewSimpleRepository creates a new simple repository.
func NewSimpleRepository() *SimpleRepository {
	return &SimpleRepository{
		entities: make(map[string]*SimpleAggregate),
	}
}

// Save saves an entity.
func (r *SimpleRepository) Save(ctx context.Context, entity *SimpleAggregate) error {
	r.entities[entity.GetID()] = entity
	return nil
}

// GetByID retrieves an entity by ID.
func (r *SimpleRepository) GetByID(ctx context.Context, id string) (*SimpleAggregate, error) {
	if entity, exists := r.entities[id]; exists {
		return entity, nil
	}
	return nil, ErrNotFound
}

// Exists checks if an entity exists.
func (r *SimpleRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, exists := r.entities[id]
	return exists, nil
}

// List retrieves all entities.
func (r *SimpleRepository) List(ctx context.Context, offset, limit int) ([]*SimpleAggregate, error) {
	entities := make([]*SimpleAggregate, 0, len(r.entities))
	for _, entity := range r.entities {
		entities = append(entities, entity)
	}
	return entities, nil
}

// Count returns the number of entities.
func (r *SimpleRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(r.entities)), nil
}

// Delete removes an entity.
func (r *SimpleRepository) Delete(ctx context.Context, id string) error {
	if _, exists := r.entities[id]; !exists {
		return ErrNotFound
	}
	delete(r.entities, id)
	return nil
}

func TestSimpleAggregate(t *testing.T) {
	// Create a new aggregate
	agg := NewSimpleAggregate("test-1", "John")

	// Verify initial state
	assert.Equal(t, "test-1", agg.GetID())
	assert.Equal(t, "John", agg.Name)
	assert.Equal(t, int64(1), agg.GetVersion())

	// Check events
	events := agg.GetUncommittedEvents()
	assert.Len(t, events, 1)
	assert.Equal(t, "Created", events[0].EventType)

	// Update the aggregate
	agg.Update("Jane")

	// Verify updated state
	assert.Equal(t, "Jane", agg.Name)
	assert.Equal(t, int64(2), agg.GetVersion())

	// Check events
	events = agg.GetUncommittedEvents()
	assert.Len(t, events, 2)
	assert.Equal(t, "Updated", events[1].EventType)

	// Mark events as committed
	agg.MarkEventsAsCommitted()
	assert.Len(t, agg.GetUncommittedEvents(), 0)
}

func TestSimpleRepository(t *testing.T) {
	// Create repository
	repo := NewSimpleRepository()

	// Create aggregate
	agg := NewSimpleAggregate("test-1", "John")

	// Save aggregate
	err := repo.Save(context.Background(), agg)
	require.NoError(t, err)

	// Check if exists
	exists, err := repo.Exists(context.Background(), "test-1")
	require.NoError(t, err)
	assert.True(t, exists)

	// Get aggregate
	retrieved, err := repo.GetByID(context.Background(), "test-1")
	require.NoError(t, err)
	assert.Equal(t, agg, retrieved)

	// Count aggregates
	count, err := repo.Count(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// Delete aggregate
	err = repo.Delete(context.Background(), "test-1")
	require.NoError(t, err)

	// Verify deleted
	exists, err = repo.Exists(context.Background(), "test-1")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestDomainErrors(t *testing.T) {
	// Test domain error creation
	err := NewError("TEST_ERROR", "test error message")
	assert.Equal(t, "TEST_ERROR: test error message", err.Error())

	// Test domain error with cause
	cause := NewError("CAUSE_ERROR", "cause error")
	errWithCause := NewErrorWithCause("TEST_ERROR", "test error message", cause)
	assert.Contains(t, errWithCause.Error(), "caused by")

	// Test error checking functions
	assert.True(t, IsNotFound(ErrNotFound))
	assert.False(t, IsNotFound(ErrInternal))
}
