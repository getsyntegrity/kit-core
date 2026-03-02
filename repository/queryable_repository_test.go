package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQueryResult(t *testing.T) {
	data := []string{"a", "b"}
	r := NewQueryResult(data, 10, 1, 5)
	assert.Equal(t, data, r.Data)
	assert.Equal(t, int64(10), r.Total)
	assert.Equal(t, 1, r.Page)
	assert.Equal(t, 5, r.PageSize)
	assert.Equal(t, 2, r.TotalPages) // ceil(10/5)
	assert.True(t, r.HasNext)
	assert.False(t, r.HasPrev)
}

func TestNewQueryResult_LastPage(t *testing.T) {
	data := []string{"a"}
	r := NewQueryResult(data, 10, 2, 5)
	assert.Equal(t, 2, r.TotalPages)
	assert.False(t, r.HasNext)
	assert.True(t, r.HasPrev)
}

func TestNewQueryResult_SinglePage(t *testing.T) {
	data := []string{"a", "b"}
	r := NewQueryResult(data, 2, 1, 5)
	assert.Equal(t, 1, r.TotalPages)
	assert.False(t, r.HasNext)
	assert.False(t, r.HasPrev)
}
