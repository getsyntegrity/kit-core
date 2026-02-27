// Package memory provides an in-memory repository implementation for testing and development.
package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/getsyntegrity/kit-core/pkg/repository"
)

// MemoryRepository implements an in-memory repository for testing and development.
type MemoryRepository[T any] struct {
	entities map[string]T
	mu       sync.RWMutex
}

// NewMemoryRepository creates a new in-memory repository.
func NewMemoryRepository[T any]() *MemoryRepository[T] {
	return &MemoryRepository[T]{
		entities: make(map[string]T),
	}
}

// Save saves an entity to the repository.
func (r *MemoryRepository[T]) Save(_ context.Context, entity T) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("%v", entity)
	r.entities[key] = entity
	return nil
}

// SaveWithID saves an entity with a specific ID.
func (r *MemoryRepository[T]) SaveWithID(_ context.Context, id string, entity T) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entities[id] = entity
	return nil
}

// GetByID retrieves an entity by its ID.
func (r *MemoryRepository[T]) GetByID(_ context.Context, _ string, entityID string) (T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var zero T
	entity, exists := r.entities[entityID]
	if !exists {
		return zero, repository.NotFoundError("entity", entityID)
	}
	return entity, nil
}

// Exists checks if an entity exists with the given ID.
func (r *MemoryRepository[T]) Exists(_ context.Context, _ string, id string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.entities[id]
	return exists, nil
}

// List retrieves a list of entities with pagination.
func (r *MemoryRepository[T]) List(_ context.Context, _ string, offset, limit int) ([]T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entities := make([]T, 0, len(r.entities))
	for _, entity := range r.entities {
		entities = append(entities, entity)
	}
	start := offset
	end := offset + limit
	if start >= len(entities) {
		return []T{}, nil
	}
	if end > len(entities) {
		end = len(entities)
	}
	return entities[start:end], nil
}

// Count returns the total number of entities.
func (r *MemoryRepository[T]) Count(_ context.Context, _ string) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return int64(len(r.entities)), nil
}

// Delete removes an entity from the repository.
func (r *MemoryRepository[T]) Delete(_ context.Context, _ string, entityID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.entities[entityID]; !exists {
		return repository.NotFoundError("entity", entityID)
	}
	delete(r.entities, entityID)
	return nil
}

// Clear removes all entities (useful for testing).
func (r *MemoryRepository[T]) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entities = make(map[string]T)
}

// GetAll returns all entities (useful for testing).
func (r *MemoryRepository[T]) GetAll() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entities := make([]T, 0, len(r.entities))
	for _, entity := range r.entities {
		entities = append(entities, entity)
	}
	return entities
}
