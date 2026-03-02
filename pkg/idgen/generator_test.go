package idgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidULID(t *testing.T) {
	assert.False(t, IsValidULID(""))
	assert.True(t, IsValidULID("01ARZ3NDEKTSV4RRFFQ69G5FAV"))
	assert.False(t, IsValidULID("invalid"))
	assert.False(t, IsValidULID("01ARZ3NDEKTSV4RRFFQ69G5FA")) // too short
}

func TestParseULID(t *testing.T) {
	u, err := ParseULID("01ARZ3NDEKTSV4RRFFQ69G5FAV")
	assert.NoError(t, err)
	assert.Equal(t, "01ARZ3NDEKTSV4RRFFQ69G5FAV", u.String())
	_, err = ParseULID("bad")
	assert.Error(t, err)
}

func TestMustParseULID(t *testing.T) {
	u := MustParseULID("01ARZ3NDEKTSV4RRFFQ69G5FAV")
	assert.Equal(t, "01ARZ3NDEKTSV4RRFFQ69G5FAV", u.String())
}

func TestMustParseULID_Panic(t *testing.T) {
	assert.Panics(t, func() { MustParseULID("invalid") })
}
