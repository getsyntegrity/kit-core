package setutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOptionalVal(t *testing.T) {
	o := NewOptionalVal(42)
	assert.True(t, o.IsValid())
	assert.Equal(t, 42, o.Value)
}

func TestEmptyOptionalVal(t *testing.T) {
	o := EmptyOptionalVal[string]()
	assert.False(t, o.IsValid())
	assert.Equal(t, "", o.Value)
}

func TestOptionalVal_Get(t *testing.T) {
	o := NewOptionalVal("ok")
	v, err := o.Get()
	assert.NoError(t, err)
	assert.Equal(t, "ok", v)
}

func TestOptionalVal_Get_Invalid(t *testing.T) {
	o := EmptyOptionalVal[int]()
	_, err := o.Get()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no valid value")
}

func TestOptionalVal_OrElse(t *testing.T) {
	o := NewOptionalVal(10)
	assert.Equal(t, 10, o.OrElse(99))
}

func TestOptionalVal_OrElse_Invalid(t *testing.T) {
	o := EmptyOptionalVal[int]()
	assert.Equal(t, 99, o.OrElse(99))
}
