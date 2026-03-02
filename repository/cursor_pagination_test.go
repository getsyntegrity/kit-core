package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCursorPagination(t *testing.T) {
	cp := NewCursorPagination(10)
	assert.Equal(t, 10, cp.Limit)
	assert.Len(t, cp.SortFields, 1)
	assert.Equal(t, "created_at", cp.SortFields[0].Field)
	assert.Equal(t, "desc", cp.SortFields[0].Direction)
}

func TestNewCursorPagination_WithSortFields(t *testing.T) {
	cp := NewCursorPagination(20, SortField{Field: "name", Direction: "asc"})
	assert.Equal(t, 20, cp.Limit)
	assert.Len(t, cp.SortFields, 1)
	assert.Equal(t, "name", cp.SortFields[0].Field)
	assert.Equal(t, "asc", cp.SortFields[0].Direction)
}

func TestCursorPagination_WithCursor(t *testing.T) {
	cp := NewCursorPagination(10)
	c := &Cursor{Timestamp: time.Now()}
	cp.WithCursor(c)
	assert.Equal(t, c, cp.Cursor)
}

func TestCursorPagination_Validate(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cp.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEncodeCursor(t *testing.T) {
	t.Run("nil cursor", func(t *testing.T) {
		s, err := EncodeCursor(nil)
		assert.NoError(t, err)
		assert.Empty(t, s)
	})
	t.Run("valid cursor", func(t *testing.T) {
		c := &Cursor{
			Timestamp: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Values:    map[string]interface{}{"id": "123"},
		}
		s, err := EncodeCursor(c)
		assert.NoError(t, err)
		assert.NotEmpty(t, s)
	})
}

func TestDecodeCursor(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		c, err := DecodeCursor("")
		assert.NoError(t, err)
		assert.Nil(t, c)
	})
	t.Run("invalid base64", func(t *testing.T) {
		_, err := DecodeCursor("!!!")
		assert.Error(t, err)
	})
	t.Run("invalid json", func(t *testing.T) {
		// valid base64 but not valid cursor JSON
		_, err := DecodeCursor("bm90LWpzb24=")
		assert.Error(t, err)
	})
	t.Run("roundtrip", func(t *testing.T) {
		orig := &Cursor{
			Timestamp: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Values:    map[string]interface{}{"id": "x"},
		}
		enc, err := EncodeCursor(orig)
		assert.NoError(t, err)
		dec, err := DecodeCursor(enc)
		assert.NoError(t, err)
		assert.NotNil(t, dec)
		assert.Equal(t, orig.Timestamp.Unix(), dec.Timestamp.Unix())
	})
}

func TestCursorPagination_BuildOrderByClause(t *testing.T) {
	t.Run("no sort fields", func(t *testing.T) {
		cp := &CursorPagination{SortFields: nil}
		assert.Empty(t, cp.BuildOrderByClause())
	})
	t.Run("single field", func(t *testing.T) {
		cp := &CursorPagination{SortFields: []SortField{{Field: "name", Direction: "asc"}}}
		assert.Equal(t, "name asc", cp.BuildOrderByClause())
	})
	t.Run("multiple fields", func(t *testing.T) {
		cp := &CursorPagination{
			SortFields: []SortField{
				{Field: "created_at", Direction: "desc"},
				{Field: "id", Direction: "asc"},
			},
		}
		assert.Equal(t, "created_at desc, id asc", cp.BuildOrderByClause())
	})
}
