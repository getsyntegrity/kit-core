package repository

import (
	"testing"

	"github.com/pablogore/go-specs/specs"
	"github.com/stretchr/testify/assert"
)

func TestQueryableRepository(t *testing.T) {
	specs.Describe(t, "queryable repository", func(s *specs.Spec) {
		s.When("NewQueryResult", func(s *specs.Spec) {
			s.It("sets data, total, page, page size, total pages and nav flags", func(ctx *specs.Context) {
				data := []string{"a", "b"}
				r := NewQueryResult(data, 10, 1, 5)
				assert.Equal(ctx.T, data, r.Data)
				assert.Equal(ctx.T, int64(10), r.Total)
				assert.Equal(ctx.T, 1, r.Page)
				assert.Equal(ctx.T, 5, r.PageSize)
				assert.Equal(ctx.T, 2, r.TotalPages)
				assert.True(ctx.T, r.HasNext)
				assert.False(ctx.T, r.HasPrev)
			})
			s.It("last page has no next", func(ctx *specs.Context) {
				data := []string{"a"}
				r := NewQueryResult(data, 10, 2, 5)
				assert.Equal(ctx.T, 2, r.TotalPages)
				assert.False(ctx.T, r.HasNext)
				assert.True(ctx.T, r.HasPrev)
			})
			s.It("single page has no next or prev", func(ctx *specs.Context) {
				data := []string{"a", "b"}
				r := NewQueryResult(data, 2, 1, 5)
				assert.Equal(ctx.T, 1, r.TotalPages)
				assert.False(ctx.T, r.HasNext)
				assert.False(ctx.T, r.HasPrev)
			})
		})
	})
}
