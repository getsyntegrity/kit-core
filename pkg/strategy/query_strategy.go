// Package strategy provides query strategy implementations.
package strategy

import (
	"context"

	"github.com/getsyntegrity/kit-core/pkg/repository"
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
	// Ensure pagination is set
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

	// Use the search functionality
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

	// Use the FindByFields functionality
	return repo.FindByFields(ctx, tenantID, options.Filters)
}

// GetName returns the strategy name.
func (s *FilteredQueryStrategy[T]) GetName() string {
	return "filtered"
}

// QueryStrategyFactory provides a factory for creating query strategies.
type QueryStrategyFactory[T any] struct{}

// NewQueryStrategyFactory creates a new query strategy factory.
func NewQueryStrategyFactory[T any]() *QueryStrategyFactory[T] {
	return &QueryStrategyFactory[T]{}
}

// CreateStrategy creates a query strategy based on the strategy type.
func (f *QueryStrategyFactory[T]) CreateStrategy(strategyType string) QueryStrategy[T] {
	switch strategyType {
	case "simple":
		return &SimpleQueryStrategy[T]{}
	case "paginated":
		return &PaginatedQueryStrategy[T]{}
	case "search":
		return &SearchQueryStrategy[T]{}
	case "filtered":
		return &FilteredQueryStrategy[T]{}
	default:
		return &SimpleQueryStrategy[T]{}
	}
}

// QueryContext holds the context for executing queries with strategies.
type QueryContext[T any] struct {
	repo     repository.QueryableRepository[T]
	strategy QueryStrategy[T]
}

// NewQueryContext creates a new query context.
func NewQueryContext[T any](repo repository.QueryableRepository[T], strategy QueryStrategy[T]) *QueryContext[T] {
	return &QueryContext[T]{
		repo:     repo,
		strategy: strategy,
	}
}

// ExecuteQuery executes a query using the current strategy.
func (qc *QueryContext[T]) ExecuteQuery(ctx context.Context, tenantID string, options *repository.QueryOptions) ([]T, error) {
	return qc.strategy.Execute(ctx, tenantID, qc.repo, options)
}

// SetStrategy changes the query strategy.
func (qc *QueryContext[T]) SetStrategy(strategy QueryStrategy[T]) {
	qc.strategy = strategy
}

// GetStrategyName returns the current strategy name.
func (qc *QueryContext[T]) GetStrategyName() string {
	return qc.strategy.GetName()
}
