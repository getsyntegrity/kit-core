package repository

import (
	"testing"

	"github.com/pablogore/go-specs/specs"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	specs.Describe(t, "repository", func(s *specs.Spec) {
		s.When("NewBaseRepository", func(s *specs.Spec) {
			s.It("returns non-nil repository", func(ctx *specs.Context) {
				r := NewBaseRepository[string]()
				assert.NotNil(ctx.T, r)
			})
		})

		s.When("BaseRepository_ValidatePagination", func(s *specs.Spec) {
			s.It("covers all pagination validation cases", func(ctx *specs.Context) {
				r := NewBaseRepository[string]()
				tests := []struct {
					name    string
					offset  int
					limit   int
					wantErr bool
				}{
					{"valid", 0, 10, false},
					{"negative offset", -1, 10, true},
					{"zero limit", 0, 0, true},
					{"negative limit", 0, -1, true},
					{"limit over 1000", 0, 1001, true},
				}
				for _, tt := range tests {
					err := r.ValidatePagination(tt.offset, tt.limit)
					if tt.wantErr {
						assert.Error(ctx.T, err)
					} else {
						assert.NoError(ctx.T, err)
					}
				}
			})
		})

		s.When("BaseRepository_BuildPagination", func(s *specs.Spec) {
			s.It("covers all build pagination cases", func(ctx *specs.Context) {
				r := NewBaseRepository[string]()
				tests := []struct {
					name       string
					offset     int
					limit      int
					wantOffset int
					wantLimit  int
				}{
					{"valid", 0, 10, 0, 10},
					{"negative offset clamped", -5, 10, 0, 10},
					{"zero limit default", 0, 0, 0, 10},
					{"negative limit default", 0, -1, 0, 10},
					{"over 1000 capped", 0, 2000, 0, 1000},
				}
				for _, tt := range tests {
					p := r.BuildPagination(tt.offset, tt.limit)
					assert.Equal(ctx.T, tt.wantOffset, p.Offset)
					assert.Equal(ctx.T, tt.wantLimit, p.Limit)
				}
			})
		})
	})
}
