package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericValidationError_Error(t *testing.T) {
	e := GenericValidationError{Message: "field is required"}
	assert.Equal(t, "field is required", e.Error())
}

func TestGenericValidationErrors(t *testing.T) {
	t.Run("empty errors slice", func(t *testing.T) {
		v := GenericValidationErrors{Errors: nil}
		assert.Nil(t, v.Errors)
	})
	t.Run("single error in slice", func(t *testing.T) {
		v := GenericValidationErrors{
			Errors: []GenericValidationError{
				{Field: "email", Message: "invalid"},
			},
		}
		assert.Len(t, v.Errors, 1)
		assert.Equal(t, "invalid", v.Errors[0].Error())
	})
	t.Run("multiple errors in slice", func(t *testing.T) {
		v := GenericValidationErrors{
			Errors: []GenericValidationError{
				{Field: "email", Message: "invalid"},
				{Field: "name", Message: "required"},
			},
		}
		assert.Len(t, v.Errors, 2)
	})
}
