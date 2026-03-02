package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQueryOptions(t *testing.T) {
	q := NewQueryOptions()
	assert.Equal(t, 0, q.Offset)
	assert.Equal(t, 10, q.Limit)
	assert.Equal(t, "created_at", q.SortBy)
	assert.Equal(t, "desc", q.SortDir)
	assert.NotNil(t, q.Filters)
	assert.NotNil(t, q.SearchFields)
	assert.NotNil(t, q.SelectFields)
	assert.NotNil(t, q.ExcludeFields)
}

func TestQueryOptions_WithPagination(t *testing.T) {
	q := NewQueryOptions().WithPagination(5, 20)
	assert.Equal(t, 5, q.Offset)
	assert.Equal(t, 20, q.Limit)
}

func TestQueryOptions_WithSorting(t *testing.T) {
	q := NewQueryOptions().WithSorting("name", "ASC")
	assert.Equal(t, "name", q.SortBy)
	assert.Equal(t, "asc", q.SortDir)
}

func TestQueryOptions_WithFilter(t *testing.T) {
	q := NewQueryOptions().WithFilter("status", "active")
	assert.Equal(t, "active", q.Filters["status"])
}

func TestQueryOptions_WithSearch(t *testing.T) {
	q := NewQueryOptions().WithSearch("foo", []string{"title", "body"})
	assert.Equal(t, "foo", q.SearchTerm)
	assert.Equal(t, []string{"title", "body"}, q.SearchFields)
}

func TestQueryOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		q       *QueryOptions
		wantErr bool
	}{
		{"valid default", NewQueryOptions(), false},
		{"valid custom", NewQueryOptions().WithPagination(0, 100), false},
		{"negative offset", &QueryOptions{Offset: -1, Limit: 10}, true},
		{"zero limit", &QueryOptions{Offset: 0, Limit: 0}, true},
		{"limit over 1000", &QueryOptions{Offset: 0, Limit: 1001}, true},
		{"invalid sort dir", &QueryOptions{Offset: 0, Limit: 10, SortDir: "invalid"}, true},
		{"asc ok", &QueryOptions{Offset: 0, Limit: 10, SortDir: "asc"}, false},
		{"desc ok", &QueryOptions{Offset: 0, Limit: 10, SortDir: "desc"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.q.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestQueryOptions_GetSortDirection(t *testing.T) {
	assert.Equal(t, "desc", NewQueryOptions().GetSortDirection())
	q := NewQueryOptions().WithSorting("x", "asc")
	assert.Equal(t, "asc", q.GetSortDirection())
}

func TestQueryOptions_HasFilters(t *testing.T) {
	assert.False(t, NewQueryOptions().HasFilters())
	assert.True(t, NewQueryOptions().WithFilter("k", "v").HasFilters())
}

func TestQueryOptions_HasSearch(t *testing.T) {
	assert.False(t, NewQueryOptions().HasSearch())
	assert.False(t, NewQueryOptions().WithSearch("term", nil).HasSearch())
	assert.False(t, NewQueryOptions().WithSearch("", []string{"a"}).HasSearch())
	assert.True(t, NewQueryOptions().WithSearch("term", []string{"a"}).HasSearch())
}
