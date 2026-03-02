package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBaseRepository(t *testing.T) {
	r := NewBaseRepository[string]()
	assert.NotNil(t, r)
}

func TestBaseRepository_ValidatePagination(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			err := r.ValidatePagination(tt.offset, tt.limit)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBaseRepository_BuildPagination(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			p := r.BuildPagination(tt.offset, tt.limit)
			assert.Equal(t, tt.wantOffset, p.Offset)
			assert.Equal(t, tt.wantLimit, p.Limit)
		})
	}
}
