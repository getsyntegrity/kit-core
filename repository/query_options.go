package repository

import (
	"fmt"
	"strings"
)

const (
	SortDirectionAsc  = "asc"
	SortDirectionDesc = "desc"
)

// QueryOptions represents query parameters for repository operations.
type QueryOptions struct {
	Offset         int
	Limit          int
	SortBy         string
	SortDir        string
	Filters        map[string]interface{}
	SearchTerm     string
	SearchFields   []string
	IncludeDeleted bool
	SelectFields   []string
	ExcludeFields  []string
}

// NewQueryOptions creates a new query options instance with defaults.
func NewQueryOptions() *QueryOptions {
	return &QueryOptions{
		Offset:        0,
		Limit:         10,
		SortBy:        "created_at",
		SortDir:       "desc",
		Filters:       make(map[string]interface{}),
		SearchFields:  []string{},
		SelectFields:  []string{},
		ExcludeFields: []string{},
	}
}

// WithPagination sets pagination parameters.
func (q *QueryOptions) WithPagination(offset, limit int) *QueryOptions {
	q.Offset = offset
	q.Limit = limit
	return q
}

// WithSorting sets sorting parameters.
func (q *QueryOptions) WithSorting(sortBy, sortDir string) *QueryOptions {
	q.SortBy = sortBy
	q.SortDir = strings.ToLower(sortDir)
	return q
}

// WithFilter adds a filter condition.
func (q *QueryOptions) WithFilter(field string, value interface{}) *QueryOptions {
	q.Filters[field] = value
	return q
}

// WithSearch sets search parameters.
func (q *QueryOptions) WithSearch(term string, fields []string) *QueryOptions {
	q.SearchTerm = term
	q.SearchFields = fields
	return q
}

// Validate validates the query options.
func (q *QueryOptions) Validate() error {
	if q.Offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}
	if q.Limit <= 0 {
		return fmt.Errorf("limit must be positive")
	}
	if q.Limit > 1000 {
		return fmt.Errorf("limit cannot exceed 1000")
	}
	if q.SortDir != "" && q.SortDir != SortDirectionAsc && q.SortDir != SortDirectionDesc {
		return fmt.Errorf("sort direction must be 'asc' or 'desc'")
	}
	return nil
}

// GetSortDirection returns the sort direction as a string.
func (q *QueryOptions) GetSortDirection() string {
	if q.SortDir == "" {
		return "desc"
	}
	return q.SortDir
}

// HasFilters checks if there are any filters applied.
func (q *QueryOptions) HasFilters() bool {
	return len(q.Filters) > 0
}

// HasSearch checks if there is a search term.
func (q *QueryOptions) HasSearch() bool {
	return q.SearchTerm != "" && len(q.SearchFields) > 0
}
