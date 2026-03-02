package setutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSet(t *testing.T) {
	s := NewSet()
	assert.NotNil(t, s)
	assert.Empty(t, s.Items())
}

func TestSet_Add(t *testing.T) {
	s := NewSet()
	s.Add("a")
	s.Add("b")
	s.Add("a")
	assert.True(t, s.Contains("a"))
	assert.True(t, s.Contains("b"))
	assert.Len(t, s.Items(), 2)
}

func TestSet_Contains(t *testing.T) {
	s := NewSet()
	assert.False(t, s.Contains("x"))
	s.Add("x")
	assert.True(t, s.Contains("x"))
}

func TestSet_Items(t *testing.T) {
	s := NewSet()
	assert.Empty(t, s.Items())
	s.Add("b")
	s.Add("a")
	items := s.Items()
	assert.Len(t, items, 2)
	assert.Contains(t, items, "a")
	assert.Contains(t, items, "b")
}
