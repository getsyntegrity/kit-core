package repository

import (
	"fmt"
	"strings"
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
		SortDir:       SortDirectionDesc,
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

// WithFilters sets multiple filter conditions.
func (q *QueryOptions) WithFilters(filters map[string]interface{}) *QueryOptions {
	for field, value := range filters {
		q.Filters[field] = value
	}
	return q
}

// WithSearch sets search parameters.
func (q *QueryOptions) WithSearch(term string, fields []string) *QueryOptions {
	q.SearchTerm = term
	q.SearchFields = fields
	return q
}

// WithIncludeDeleted includes deleted records in the query.
func (q *QueryOptions) WithIncludeDeleted(include bool) *QueryOptions {
	q.IncludeDeleted = include
	return q
}

// WithSelectFields specifies which fields to select.
func (q *QueryOptions) WithSelectFields(fields []string) *QueryOptions {
	q.SelectFields = fields
	return q
}

// WithExcludeFields specifies which fields to exclude.
func (q *QueryOptions) WithExcludeFields(fields []string) *QueryOptions {
	q.ExcludeFields = fields
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
		return SortDirectionDesc
	}
	return q.SortDir
}

// IsAscending checks if the sort direction is ascending.
func (q *QueryOptions) IsAscending() bool {
	return q.GetSortDirection() == SortDirectionAsc
}

// IsDescending checks if the sort direction is descending.
func (q *QueryOptions) IsDescending() bool {
	return q.GetSortDirection() == SortDirectionDesc
}

// HasFilters checks if there are any filters applied.
func (q *QueryOptions) HasFilters() bool {
	return len(q.Filters) > 0
}

// HasSearch checks if there is a search term.
func (q *QueryOptions) HasSearch() bool {
	return q.SearchTerm != "" && len(q.SearchFields) > 0
}

// GetFilterValue gets a filter value by field name.
func (q *QueryOptions) GetFilterValue(field string) (interface{}, bool) {
	value, exists := q.Filters[field]
	return value, exists
}

// RemoveFilter removes a filter by field name.
func (q *QueryOptions) RemoveFilter(field string) {
	delete(q.Filters, field)
}

// ClearFilters removes all filters.
func (q *QueryOptions) ClearFilters() {
	q.Filters = make(map[string]interface{})
}

// BuildWhereClause builds a WHERE clause from filters.
func (q *QueryOptions) BuildWhereClause() (string, []interface{}, error) {
	if !q.HasFilters() {
		return "", []interface{}{}, nil
	}
	var conditions []string
	var args []interface{}
	argIndex := 1
	for field, value := range q.Filters {
		if value == nil {
			conditions = append(conditions, fmt.Sprintf("%s IS NULL", field))
		} else {
			conditions = append(conditions, fmt.Sprintf("%s = $%d", field, argIndex))
			args = append(args, value)
			argIndex++
		}
	}
	return strings.Join(conditions, " AND "), args, nil
}

// BuildOrderByClause builds an ORDER BY clause.
func (q *QueryOptions) BuildOrderByClause() string {
	if q.SortBy == "" {
		return ""
	}
	return fmt.Sprintf("%s %s", q.SortBy, q.GetSortDirection())
}

// BuildLimitOffsetClause builds LIMIT and OFFSET clauses.
func (q *QueryOptions) BuildLimitOffsetClause() string {
	return fmt.Sprintf("LIMIT %d OFFSET %d", q.Limit, q.Offset)
}

// Clone creates a copy of the query options.
func (q *QueryOptions) Clone() *QueryOptions {
	clone := &QueryOptions{
		Offset:         q.Offset,
		Limit:          q.Limit,
		SortBy:         q.SortBy,
		SortDir:        q.SortDir,
		SearchTerm:     q.SearchTerm,
		IncludeDeleted: q.IncludeDeleted,
		Filters:        make(map[string]interface{}),
		SearchFields:   make([]string, len(q.SearchFields)),
		SelectFields:   make([]string, len(q.SelectFields)),
		ExcludeFields:  make([]string, len(q.ExcludeFields)),
	}
	for k, v := range q.Filters {
		clone.Filters[k] = v
	}
	copy(clone.SearchFields, q.SearchFields)
	copy(clone.SelectFields, q.SelectFields)
	copy(clone.ExcludeFields, q.ExcludeFields)
	return clone
}
