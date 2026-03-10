package repository

import (
	"testing"

	"github.com/pablogore/go-specs/specs"
	"github.com/stretchr/testify/assert"
)

func TestQueryOptions(t *testing.T) {
	specs.Describe(t, "query options", func(s *specs.Spec) {
		s.When("NewQueryOptions", func(s *specs.Spec) {
			s.It("returns default offset, limit, sort and non-nil maps", func(ctx *specs.Context) {
				q := NewQueryOptions()
				assert.Equal(ctx.T, 0, q.Offset)
				assert.Equal(ctx.T, 10, q.Limit)
				assert.Equal(ctx.T, "created_at", q.SortBy)
				assert.Equal(ctx.T, "desc", q.SortDir)
				assert.NotNil(ctx.T, q.Filters)
				assert.NotNil(ctx.T, q.SearchFields)
				assert.NotNil(ctx.T, q.SelectFields)
				assert.NotNil(ctx.T, q.ExcludeFields)
			})
		})

		s.When("QueryOptions_WithPagination", func(s *specs.Spec) {
			s.It("sets offset and limit", func(ctx *specs.Context) {
				q := NewQueryOptions().WithPagination(5, 20)
				assert.Equal(ctx.T, 5, q.Offset)
				assert.Equal(ctx.T, 20, q.Limit)
			})
		})

		s.When("QueryOptions_WithSorting", func(s *specs.Spec) {
			s.It("sets sort by and direction", func(ctx *specs.Context) {
				q := NewQueryOptions().WithSorting("name", "ASC")
				assert.Equal(ctx.T, "name", q.SortBy)
				assert.Equal(ctx.T, "asc", q.SortDir)
			})
		})

		s.When("QueryOptions_WithFilter", func(s *specs.Spec) {
			s.It("sets filter value", func(ctx *specs.Context) {
				q := NewQueryOptions().WithFilter("status", "active")
				assert.Equal(ctx.T, "active", q.Filters["status"])
			})
		})

		s.When("QueryOptions_WithSearch", func(s *specs.Spec) {
			s.It("sets search term and fields", func(ctx *specs.Context) {
				q := NewQueryOptions().WithSearch("foo", []string{"title", "body"})
				assert.Equal(ctx.T, "foo", q.SearchTerm)
				assert.Equal(ctx.T, []string{"title", "body"}, q.SearchFields)
			})
		})

		s.When("QueryOptions_Validate", func(s *specs.Spec) {
			s.It("covers all validation cases", func(ctx *specs.Context) {
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
					err := tt.q.Validate()
					if tt.wantErr {
						assert.Error(ctx.T, err)
					} else {
						assert.NoError(ctx.T, err)
					}
				}
			})
		})

		s.When("QueryOptions_GetSortDirection", func(s *specs.Spec) {
			s.It("returns desc when SortDir empty", func(ctx *specs.Context) {
				q := NewQueryOptions()
				q.SortDir = ""
				assert.Equal(ctx.T, "desc", q.GetSortDirection())
			})
			s.It("returns asc when set", func(ctx *specs.Context) {
				q := NewQueryOptions().WithSorting("x", "asc")
				assert.Equal(ctx.T, "asc", q.GetSortDirection())
			})
			s.It("returns desc when set", func(ctx *specs.Context) {
				q := NewQueryOptions().WithSorting("x", "DESC")
				assert.Equal(ctx.T, "desc", q.GetSortDirection())
			})
		})

		s.When("QueryOptions_HasFilters", func(s *specs.Spec) {
			s.It("returns false when empty, true when filter set", func(ctx *specs.Context) {
				assert.False(ctx.T, NewQueryOptions().HasFilters())
				assert.True(ctx.T, NewQueryOptions().WithFilter("k", "v").HasFilters())
			})
		})

		s.When("QueryOptions_HasSearch", func(s *specs.Spec) {
			s.It("returns true only when term and fields set", func(ctx *specs.Context) {
				assert.False(ctx.T, NewQueryOptions().HasSearch())
				assert.False(ctx.T, NewQueryOptions().WithSearch("term", nil).HasSearch())
				assert.False(ctx.T, NewQueryOptions().WithSearch("", []string{"a"}).HasSearch())
				assert.True(ctx.T, NewQueryOptions().WithSearch("term", []string{"a"}).HasSearch())
			})
		})
	})
}
