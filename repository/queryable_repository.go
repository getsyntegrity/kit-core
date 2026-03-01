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
