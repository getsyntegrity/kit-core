package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValidator(t *testing.T) {
	v := NewValidator()
	assert.NotNil(t, v)
	assert.NotNil(t, v.validate)
}

func TestValidator_Validate_ValidStruct(t *testing.T) {
	v := NewValidator()
	type S struct {
		Name string `json:"name" validate:"required"`
	}
	err := v.Validate(&S{Name: "ok"})
	assert.NoError(t, err)
}

func TestValidator_Validate_InvalidStruct(t *testing.T) {
	v := NewValidator()
	type S struct {
		Name string `json:"name" validate:"required"`
	}
	err := v.Validate(&S{})
	assert.Error(t, err)
	var verr Errors
	assert.ErrorAs(t, err, &verr)
	assert.Len(t, verr.Errors, 1)
	assert.Equal(t, "name", verr.Errors[0].Field)
}

func TestValidator_ValidateVar(t *testing.T) {
	v := NewValidator()
	err := v.ValidateVar("", "required")
	assert.Error(t, err)
	err = v.ValidateVar("x", "required")
	assert.NoError(t, err)
}

func TestValidator_ValidateVar_Email(t *testing.T) {
	v := NewValidator()
	assert.NoError(t, v.ValidateVar("a@b.co", "email"))
	assert.Error(t, v.ValidateVar("invalid", "email"))
}

func TestValidator_ValidateVar_Domain(t *testing.T) {
	v := NewValidator()
	assert.NoError(t, v.ValidateVar("example.com", "domain"))
	assert.Error(t, v.ValidateVar("invalid!.com", "domain"))
}

func TestValidator_ValidateVar_ULID(t *testing.T) {
	v := NewValidator()
	assert.NoError(t, v.ValidateVar("01ARZ3NDEKTSV4RRFFQ69G5FAV", "ulid"))
	assert.Error(t, v.ValidateVar("not-a-ulid", "ulid"))
}

func TestValidator_ValidateVar_Positive(t *testing.T) {
	v := NewValidator()
	assert.NoError(t, v.ValidateVar(1, "positive"))
	assert.NoError(t, v.ValidateVar(int64(1), "positive"))
	assert.NoError(t, v.ValidateVar(uint(1), "positive"))
	assert.NoError(t, v.ValidateVar(1.5, "positive"))
	assert.Error(t, v.ValidateVar(0, "positive"))
	assert.Error(t, v.ValidateVar(-1, "positive"))
	assert.Error(t, v.ValidateVar("x", "positive")) // invalid type
}

func TestValidator_ValidateVar_MinMaxLen(t *testing.T) {
	v := NewValidator()
	assert.Error(t, v.ValidateVar("ab", "min=3"))
	assert.NoError(t, v.ValidateVar("abc", "min=3"))
	assert.Error(t, v.ValidateVar("abcd", "max=3"))
	assert.NoError(t, v.ValidateVar("abc", "len=3"))
	assert.Error(t, v.ValidateVar("ab", "len=3"))
}

func TestValidator_ValidateVar_Oneof(t *testing.T) {
	v := NewValidator()
	assert.NoError(t, v.ValidateVar("a", "oneof=a b c"))
	assert.Error(t, v.ValidateVar("x", "oneof=a b c"))
}

func TestErrors_Error(t *testing.T) {
	e := Errors{
		Errors: []Error{
			{Field: "A", Message: "msg1"},
			{Field: "B", Message: "msg2"},
		},
	}
	assert.Equal(t, "A: msg1; B: msg2", e.Error())
}
