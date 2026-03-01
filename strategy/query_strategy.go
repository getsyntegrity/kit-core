// Package strategy provides query strategy contracts.
package strategy

import (
	"context"

	"github.com/getsyntegrity/kit-core/repository"
)

// QueryStrategy defines the interface for different query strategies.
type QueryStrategy[T any] interface {
	Execute(ctx context.Context, tenantID string, repo repository.QueryableRepository[T], options *repository.QueryOptions) ([]T, error)
	GetName() string
}

// SimpleQueryStrategy implements a simple query strategy.
type SimpleQueryStrategy[T any] struct{}

// Execute executes a simple query.
func (s *SimpleQueryStrategy[T]) Execute(ctx context.Context, tenantID string, repo repository.QueryableRepository[T], options *repository.QueryOptions) ([]T, error) {
	return repo.Query(ctx, tenantID, options)
}

// GetName returns the strategy name.
func (s *SimpleQueryStrategy[T]) GetName() string {
	return "simple"
}

// PaginatedQueryStrategy implements a paginated query strategy.
type PaginatedQueryStrategy[T any] struct{}

// Execute executes a paginated query.
func (s *PaginatedQueryStrategy[T]) Execute(ctx context.Context, tenantID string, repo repository.QueryableRepository[T], options *repository.QueryOptions) ([]T, error) {
	if options.Limit == 0 {
		options.Limit = 10
	}
	return repo.Query(ctx, tenantID, options)
}

// GetName returns the strategy name.
func (s *PaginatedQueryStrategy[T]) GetName() string {
	return "paginated"
}

// SearchQueryStrategy implements a search query strategy.
type SearchQueryStrategy[T any] struct{}

// Execute executes a search query.
func (s *SearchQueryStrategy[T]) Execute(ctx context.Context, tenantID string, repo repository.QueryableRepository[T], options *repository.QueryOptions) ([]T, error) {
	if !options.HasSearch() {
		return repo.Query(ctx, tenantID, options)
	}
	return repo.Search(ctx, tenantID, options.SearchTerm, options.SearchFields, options.Limit)
}

// GetName returns the strategy name.
func (s *SearchQueryStrategy[T]) GetName() string {
	return "search"
}

// FilteredQueryStrategy implements a filtered query strategy.
type FilteredQueryStrategy[T any] struct{}

// Execute executes a filtered query.
func (s *FilteredQueryStrategy[T]) Execute(ctx context.Context, tenantID string, repo repository.QueryableRepository[T], options *repository.QueryOptions) ([]T, error) {
	if !options.HasFilters() {
		return repo.Query(ctx, tenantID, options)
	}
	return repo.FindByFields(ctx, tenantID, options.Filters)
}

// GetName returns the strategy name.
func (s *FilteredQueryStrategy[T]) GetName() string {
	return "filtered"
}
