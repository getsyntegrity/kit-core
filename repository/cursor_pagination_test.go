package repository

import (
	"testing"
	"time"

	"github.com/pablogore/go-specs/specs"
	"github.com/stretchr/testify/assert"
)

func TestCursorPagination(t *testing.T) {
	specs.Describe(t, "cursor pagination", func(s *specs.Spec) {
		s.When("NewCursorPagination", func(s *specs.Spec) {
			s.It("sets limit and default sort field", func(ctx *specs.Context) {
				cp := NewCursorPagination(10)
				assert.Equal(ctx.T, 10, cp.Limit)
				assert.Len(ctx.T, cp.SortFields, 1)
				assert.Equal(ctx.T, "created_at", cp.SortFields[0].Field)
				assert.Equal(ctx.T, "desc", cp.SortFields[0].Direction)
			})
			s.It("with sort fields uses given sort", func(ctx *specs.Context) {
				cp := NewCursorPagination(20, SortField{Field: "name", Direction: "asc"})
				assert.Equal(ctx.T, 20, cp.Limit)
				assert.Len(ctx.T, cp.SortFields, 1)
				assert.Equal(ctx.T, "name", cp.SortFields[0].Field)
				assert.Equal(ctx.T, "asc", cp.SortFields[0].Direction)
			})
		})

		s.When("CursorPagination_WithCursor", func(s *specs.Spec) {
			s.It("sets cursor", func(ctx *specs.Context) {
				cp := NewCursorPagination(10)
				c := &Cursor{Timestamp: time.Now()}
				cp.WithCursor(c)
				assert.Equal(ctx.T, c, cp.Cursor)
			})
		})

		s.When("CursorPagination_Validate", func(s *specs.Spec) {
			s.It("covers all validation cases", func(ctx *specs.Context) {
				tests := []struct {
					name    string
					cp      *CursorPagination
					wantErr bool
				}{
					{"valid default", NewCursorPagination(10), false},
					{"valid with sort", NewCursorPagination(5, SortField{Field: "id", Direction: "asc"}), false},
					{"limit zero", &CursorPagination{Limit: 0, SortFields: []SortField{{Field: "x", Direction: "asc"}}}, true},
					{"limit negative", &CursorPagination{Limit: -1, SortFields: []SortField{{Field: "x", Direction: "asc"}}}, true},
					{"limit over 1000", &CursorPagination{Limit: 1001, SortFields: []SortField{{Field: "x", Direction: "asc"}}}, true},
					{"no sort fields", &CursorPagination{Limit: 10, SortFields: nil}, true},
					{"empty sort field name", &CursorPagination{Limit: 10, SortFields: []SortField{{Field: "", Direction: "asc"}}}, true},
					{"invalid sort direction", &CursorPagination{Limit: 10, SortFields: []SortField{{Field: "x", Direction: "invalid"}}}, true},
				}
				for _, tt := range tests {
					err := tt.cp.Validate()
					if tt.wantErr {
						assert.Error(ctx.T, err)
					} else {
						assert.NoError(ctx.T, err)
					}
				}
			})
		})

		s.When("EncodeCursor", func(s *specs.Spec) {
			s.It("nil cursor returns empty string", func(ctx *specs.Context) {
				encoded, err := EncodeCursor(nil)
				assert.NoError(ctx.T, err)
				assert.Empty(ctx.T, encoded)
			})
			s.It("valid cursor with Values returns non-empty string", func(ctx *specs.Context) {
				c := &Cursor{
					Timestamp: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					Values:    map[string]interface{}{"id": "123"},
				}
				encoded, err := EncodeCursor(c)
				assert.NoError(ctx.T, err)
				assert.NotEmpty(ctx.T, encoded)
			})
			s.It("cursor with Direction marshals and encodes", func(ctx *specs.Context) {
				c := &Cursor{
					Timestamp: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					Direction: map[string]string{"created_at": "desc"},
					Values:    map[string]interface{}{"id": "x"},
				}
				encoded, err := EncodeCursor(c)
				assert.NoError(ctx.T, err)
				assert.NotEmpty(ctx.T, encoded)
				dec, err := DecodeCursor(encoded)
				assert.NoError(ctx.T, err)
				assert.NotNil(ctx.T, dec)
				assert.Equal(ctx.T, "x", dec.Values["id"])
			})
			s.It("returns error when cursor cannot be marshalled", func(ctx *specs.Context) {
				// Values containing a channel cannot be JSON-marshalled
				c := &Cursor{
					Timestamp: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					Values:    map[string]interface{}{"ch": make(chan int)},
				}
				encoded, err := EncodeCursor(c)
				assert.Error(ctx.T, err)
				assert.Empty(ctx.T, encoded)
				assert.Contains(ctx.T, err.Error(), "marshal")
			})
		})

		s.When("DecodeCursor", func(s *specs.Spec) {
			s.It("empty string returns nil cursor", func(ctx *specs.Context) {
				c, err := DecodeCursor("")
				assert.NoError(ctx.T, err)
				assert.Nil(ctx.T, c)
			})
			s.It("invalid base64 returns error", func(ctx *specs.Context) {
				_, err := DecodeCursor("!!!")
				assert.Error(ctx.T, err)
			})
			s.It("invalid json returns error", func(ctx *specs.Context) {
				_, err := DecodeCursor("bm90LWpzb24=")
				assert.Error(ctx.T, err)
			})
			s.It("roundtrip preserves timestamp", func(ctx *specs.Context) {
				orig := &Cursor{
					Timestamp: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					Values:    map[string]interface{}{"id": "x"},
				}
				enc, err := EncodeCursor(orig)
				assert.NoError(ctx.T, err)
				dec, err := DecodeCursor(enc)
				assert.NoError(ctx.T, err)
				assert.NotNil(ctx.T, dec)
				assert.Equal(ctx.T, orig.Timestamp.Unix(), dec.Timestamp.Unix())
			})
		})

		s.When("CursorPagination_BuildOrderByClause", func(s *specs.Spec) {
			s.It("no sort fields returns empty", func(ctx *specs.Context) {
				cp := &CursorPagination{SortFields: nil}
				assert.Empty(ctx.T, cp.BuildOrderByClause())
			})
			s.It("single field returns one clause", func(ctx *specs.Context) {
				cp := &CursorPagination{SortFields: []SortField{{Field: "name", Direction: "asc"}}}
				assert.Equal(ctx.T, "name asc", cp.BuildOrderByClause())
			})
			s.It("multiple fields returns comma-separated", func(ctx *specs.Context) {
				cp := &CursorPagination{
					SortFields: []SortField{
						{Field: "created_at", Direction: "desc"},
						{Field: "id", Direction: "asc"},
					},
				}
				assert.Equal(ctx.T, "created_at desc, id asc", cp.BuildOrderByClause())
			})
		})
	})
}
