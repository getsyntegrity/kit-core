package errorschain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c := New()
	assert.NotNil(t, c)
	assert.Nil(t, c.Error())
}

func TestNew_WithOptions(t *testing.T) {
	c := New(ReturnFirst())
	assert.NotNil(t, c)
}

func TestChain_AddError(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")
	c := New()
	c.AddError(err1).AddError(err2)
	err := c.Error()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "err1")
	assert.Contains(t, err.Error(), "err2")
}

func TestChain_AddError_NilIgnored(t *testing.T) {
	c := New()
	c.AddError(nil)
	assert.Nil(t, c.Error())
}

func TestChain_AddErrorFn(t *testing.T) {
	c := New()
	c.AddErrorFn(func() error { return errors.New("fn error") })
	err := c.Error()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fn error")
}

func TestChain_AddErrorFn_NilIgnored(t *testing.T) {
	c := New()
	c.AddErrorFn(func() error { return nil })
	assert.Nil(t, c.Error())
}

func TestChain_ReturnFirst(t *testing.T) {
	c := New(ReturnFirst())
	c.AddError(errors.New("first")).AddError(errors.New("second"))
	err := c.Error()
	assert.Error(t, err)
	assert.Equal(t, "first", err.Error())
}

func TestChain_ReturnFirst_OnlyFirstErrorAdded(t *testing.T) {
	c := New(ReturnFirst())
	c.AddError(errors.New("first"))
	c.AddError(errors.New("second")) // should be ignored when returnFirst
	err := c.Error()
	assert.Equal(t, "first", err.Error())
}

func TestChain_ReturnFirst_AddErrorFn(t *testing.T) {
	c := New(ReturnFirst())
	c.AddErrorFn(func() error { return errors.New("first") })
	c.AddErrorFn(func() error { return errors.New("second") })
	err := c.Error()
	assert.Equal(t, "first", err.Error())
}

func TestChain_ReturnAll(t *testing.T) {
	c := New(ReturnAll())
	c.AddError(errors.New("a")).AddError(errors.New("b"))
	err := c.Error()
	assert.Error(t, err)
}
