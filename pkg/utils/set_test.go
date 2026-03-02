package setutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSet_Utils(t *testing.T) {
	s := NewSet()
	assert.NotNil(t, s)
	assert.Empty(t, s.Items())
}

func TestSet_Add_Utils(t *testing.T) {
	s := NewSet()
	s.Add("a")
	s.Add("b")
	s.Add("a")
	assert.True(t, s.Contains("a"))
	assert.True(t, s.Contains("b"))
	assert.Len(t, s.Items(), 2)
}
