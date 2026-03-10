package errorschain

import (
	"errors"
	"testing"

	"github.com/pablogore/go-specs/specs"
	"github.com/stretchr/testify/assert"
)

func TestErrorsChain(t *testing.T) {
	specs.Describe(t, "errorschain", func(s *specs.Spec) {
		s.When("New", func(s *specs.Spec) {
			s.It("returns non-nil chain with no error", func(ctx *specs.Context) {
				c := New()
				assert.NotNil(ctx.T, c)
				assert.Nil(ctx.T, c.Error())
			})
			s.It("accepts options such as ReturnFirst", func(ctx *specs.Context) {
				c := New(ReturnFirst())
				assert.NotNil(ctx.T, c)
			})
		})

		s.When("Chain_AddError", func(s *specs.Spec) {
			s.It("collects multiple errors", func(ctx *specs.Context) {
				err1 := errors.New("err1")
				err2 := errors.New("err2")
				c := New()
				c.AddError(err1).AddError(err2)
				err := c.Error()
				assert.Error(ctx.T, err)
				assert.Contains(ctx.T, err.Error(), "err1")
				assert.Contains(ctx.T, err.Error(), "err2")
			})
			s.It("ignores nil errors", func(ctx *specs.Context) {
				c := New()
				c.AddError(nil)
				assert.Nil(ctx.T, c.Error())
			})
		})

		s.When("Chain_AddErrorFn", func(s *specs.Spec) {
			s.It("evaluates function and adds error", func(ctx *specs.Context) {
				c := New()
				c.AddErrorFn(func() error { return errors.New("fn error") })
				err := c.Error()
				assert.Error(ctx.T, err)
				assert.Contains(ctx.T, err.Error(), "fn error")
			})
			s.It("ignores when function returns nil", func(ctx *specs.Context) {
				c := New()
				c.AddErrorFn(func() error { return nil })
				assert.Nil(ctx.T, c.Error())
			})
		})

		s.When("Chain_ReturnFirst", func(s *specs.Spec) {
			s.It("returns first error only", func(ctx *specs.Context) {
				c := New(ReturnFirst())
				c.AddError(errors.New("first")).AddError(errors.New("second"))
				err := c.Error()
				assert.Error(ctx.T, err)
				assert.Equal(ctx.T, "first", err.Error())
			})
			s.It("returns first when only first is added then second", func(ctx *specs.Context) {
				c := New(ReturnFirst())
				c.AddError(errors.New("first"))
				c.AddError(errors.New("second"))
				err := c.Error()
				assert.Equal(ctx.T, "first", err.Error())
			})
			s.It("works with AddErrorFn", func(ctx *specs.Context) {
				c := New(ReturnFirst())
				c.AddErrorFn(func() error { return errors.New("first") })
				c.AddErrorFn(func() error { return errors.New("second") })
				err := c.Error()
				assert.Equal(ctx.T, "first", err.Error())
			})
		})

		s.When("Chain_ReturnAll", func(s *specs.Spec) {
			s.It("returns combined error", func(ctx *specs.Context) {
				c := New(ReturnAll())
				c.AddError(errors.New("a")).AddError(errors.New("b"))
				err := c.Error()
				assert.Error(ctx.T, err)
			})
		})

		s.When("Chain_Error", func(s *specs.Spec) {
			s.It("ReturnFirst with no errors returns nil", func(ctx *specs.Context) {
				c := New(ReturnFirst())
				err := c.Error()
				assert.Nil(ctx.T, err)
			})
			s.It("ReturnAll with multiple errors combines them", func(ctx *specs.Context) {
				c := New(ReturnAll())
				c.AddError(errors.New("err1")).AddError(errors.New("err2"))
				err := c.Error()
				assert.Error(ctx.T, err)
				assert.Contains(ctx.T, err.Error(), "err1")
				assert.Contains(ctx.T, err.Error(), "err2")
			})
		})
	})
}
