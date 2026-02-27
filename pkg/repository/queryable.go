package repository

import "context"

// QueryableRepository extends the basic repository with advanced query capabilities.
// CRITICAL: All methods require tenant_id - cross-tenant access is impossible by design.
type QueryableRepository[T any] interface {
	Repository[T]

	Query(ctx context.Context, tenantID string, options *QueryOptions) ([]T, error)
	QueryWithCount(ctx context.Context, tenantID string, options *QueryOptions) ([]T, int64, error)
	QueryWithCursor(ctx context.Context, tenantID string, pagination *CursorPagination, filters map[string]interface{}) (*CursorResult[T], error)
	FindByField(ctx context.Context, tenantID, field string, value interface{}) ([]T, error)
	FindByFields(ctx context.Context, tenantID string, filters map[string]interface{}) ([]T, error)
	Search(ctx context.Context, tenantID, term string, fields []string, limit int) ([]T, error)
	GetByIDs(ctx context.Context, tenantID string, ids []string) ([]T, error)
	BulkSave(ctx context.Context, entities []T) error
	BulkDelete(ctx context.Context, tenantID string, ids []string) error
	SoftDelete(ctx context.Context, tenantID, id string) error
	Restore(ctx context.Context, tenantID, id string) error
}

// BaseQueryableRepository provides a base implementation for queryable repositories.
type BaseQueryableRepository[T any] struct {
	BaseRepository[T]
}

// NewBaseQueryableRepository creates a new base queryable repository.
func NewBaseQueryableRepository[T any]() *BaseQueryableRepository[T] {
	return &BaseQueryableRepository[T]{
		BaseRepository: *NewBaseRepository[T](),
	}
}

// QueryResult represents a paginated query result.
type QueryResult[T any] struct {
	Data       []T
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
	HasNext    bool
	HasPrev    bool
}

// NewQueryResult creates a new query result.
func NewQueryResult[T any](data []T, total int64, page, pageSize int) *QueryResult[T] {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	return &QueryResult[T]{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// QueryOptionsBuilder provides a fluent interface for building query options.
type QueryOptionsBuilder struct {
	options *QueryOptions
}

// NewQueryOptionsBuilder creates a new query options builder.
func NewQueryOptionsBuilder() *QueryOptionsBuilder {
	return &QueryOptionsBuilder{options: NewQueryOptions()}
}

// Pagination sets pagination parameters.
func (b *QueryOptionsBuilder) Pagination(offset, limit int) *QueryOptionsBuilder {
	b.options.WithPagination(offset, limit)
	return b
}

// Sort sets sorting parameters.
func (b *QueryOptionsBuilder) Sort(sortBy, sortDir string) *QueryOptionsBuilder {
	b.options.WithSorting(sortBy, sortDir)
	return b
}

// Filter adds a filter condition.
func (b *QueryOptionsBuilder) Filter(field string, value interface{}) *QueryOptionsBuilder {
	b.options.WithFilter(field, value)
	return b
}

// Filters sets multiple filter conditions.
func (b *QueryOptionsBuilder) Filters(filters map[string]interface{}) *QueryOptionsBuilder {
	b.options.WithFilters(filters)
	return b
}

// Search sets search parameters.
func (b *QueryOptionsBuilder) Search(term string, fields []string) *QueryOptionsBuilder {
	b.options.WithSearch(term, fields)
	return b
}

// IncludeDeleted includes deleted records.
func (b *QueryOptionsBuilder) IncludeDeleted(include bool) *QueryOptionsBuilder {
	b.options.WithIncludeDeleted(include)
	return b
}

// SelectFields specifies which fields to select.
func (b *QueryOptionsBuilder) SelectFields(fields []string) *QueryOptionsBuilder {
	b.options.WithSelectFields(fields)
	return b
}

// ExcludeFields specifies which fields to exclude.
func (b *QueryOptionsBuilder) ExcludeFields(fields []string) *QueryOptionsBuilder {
	b.options.WithExcludeFields(fields)
	return b
}

// Build returns the built query options.
func (b *QueryOptionsBuilder) Build() *QueryOptions {
	return b.options
}

// QueryBuilder provides a fluent interface for building repository queries.
type QueryBuilder[T any] struct {
	repository QueryableRepository[T]
	tenantID   string
	options    *QueryOptions
}

// NewQueryBuilder creates a new repository query builder.
func NewQueryBuilder[T any](repository QueryableRepository[T], tenantID string) *QueryBuilder[T] {
	return &QueryBuilder[T]{
		repository: repository,
		tenantID:   tenantID,
		options:    NewQueryOptions(),
	}
}

// Pagination sets pagination parameters.
func (b *QueryBuilder[T]) Pagination(offset, limit int) *QueryBuilder[T] {
	b.options.WithPagination(offset, limit)
	return b
}

// Sort sets sorting parameters.
func (b *QueryBuilder[T]) Sort(sortBy, sortDir string) *QueryBuilder[T] {
	b.options.WithSorting(sortBy, sortDir)
	return b
}

// Filter adds a filter condition.
func (b *QueryBuilder[T]) Filter(field string, value interface{}) *QueryBuilder[T] {
	b.options.WithFilter(field, value)
	return b
}

// Filters sets multiple filter conditions.
func (b *QueryBuilder[T]) Filters(filters map[string]interface{}) *QueryBuilder[T] {
	b.options.WithFilters(filters)
	return b
}

// Search sets search parameters.
func (b *QueryBuilder[T]) Search(term string, fields []string) *QueryBuilder[T] {
	b.options.WithSearch(term, fields)
	return b
}

// IncludeDeleted includes deleted records.
func (b *QueryBuilder[T]) IncludeDeleted(include bool) *QueryBuilder[T] {
	b.options.WithIncludeDeleted(include)
	return b
}

// SelectFields specifies which fields to select.
func (b *QueryBuilder[T]) SelectFields(fields []string) *QueryBuilder[T] {
	b.options.WithSelectFields(fields)
	return b
}

// ExcludeFields specifies which fields to exclude.
func (b *QueryBuilder[T]) ExcludeFields(fields []string) *QueryBuilder[T] {
	b.options.WithExcludeFields(fields)
	return b
}

// Execute executes the query and returns the results.
func (b *QueryBuilder[T]) Execute(ctx context.Context) ([]T, error) {
	return b.repository.Query(ctx, b.tenantID, b.options)
}

// ExecuteWithCount executes the query and returns the results with total count.
func (b *QueryBuilder[T]) ExecuteWithCount(ctx context.Context) ([]T, int64, error) {
	return b.repository.QueryWithCount(ctx, b.tenantID, b.options)
}

// ExecuteResult executes the query and returns a paginated result.
func (b *QueryBuilder[T]) ExecuteResult(ctx context.Context) (*QueryResult[T], error) {
	data, total, err := b.repository.QueryWithCount(ctx, b.tenantID, b.options)
	if err != nil {
		return nil, err
	}
	page := (b.options.Offset / b.options.Limit) + 1
	return NewQueryResult(data, total, page, b.options.Limit), nil
}
