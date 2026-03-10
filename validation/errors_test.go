package validation

import (
	"testing"

	"github.com/pablogore/go-specs/specs"
	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	specs.Describe(t, "validation errors", func(s *specs.Spec) {
		s.When("GenericValidationError_Error", func(s *specs.Spec) {
			s.It("returns the message", func(ctx *specs.Context) {
				e := GenericValidationError{Message: "field is required"}
				assert.Equal(ctx.T, "field is required", e.Error())
			})
		})

		s.When("GenericValidationErrors", func(s *specs.Spec) {
			s.It("empty errors slice", func(ctx *specs.Context) {
				v := GenericValidationErrors{Errors: nil}
				assert.Nil(ctx.T, v.Errors)
			})
			s.It("single error in slice", func(ctx *specs.Context) {
				v := GenericValidationErrors{
					Errors: []GenericValidationError{
						{Field: "email", Message: "invalid"},
					},
				}
				assert.Len(ctx.T, v.Errors, 1)
				assert.Equal(ctx.T, "invalid", v.Errors[0].Error())
			})
			s.It("multiple errors in slice", func(ctx *specs.Context) {
				v := GenericValidationErrors{
					Errors: []GenericValidationError{
						{Field: "email", Message: "invalid"},
						{Field: "name", Message: "required"},
					},
				}
				assert.Len(ctx.T, v.Errors, 2)
			})
		})
	})
}
