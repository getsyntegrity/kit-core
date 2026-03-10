package validation

import (
	"testing"

	"github.com/pablogore/go-specs/specs"
	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	specs.Describe(t, "validator", func(s *specs.Spec) {
		s.When("NewValidator", func(s *specs.Spec) {
			s.It("returns non-nil validator with validate", func(ctx *specs.Context) {
				v := NewValidator()
				assert.NotNil(ctx.T, v)
				assert.NotNil(ctx.T, v.validate)
			})
			s.It("invokes tag name func for struct validation including json hyphen fields", func(ctx *specs.Context) {
				v := NewValidator()
				type S struct {
					Visible string `json:"visible" validate:"required"`
					Hidden  string `json:"-" validate:"required"`
				}
				err := v.Validate(&S{Visible: "x"})
				assert.Error(ctx.T, err)
				var verr Errors
				assert.ErrorAs(ctx.T, err, &verr)
				assert.Len(ctx.T, verr.Errors, 1)
				assert.Contains(ctx.T, verr.Errors[0].Message, "required")
			})
		})

		s.When("Validator_Validate", func(s *specs.Spec) {
			s.It("valid struct passes", func(ctx *specs.Context) {
				v := NewValidator()
				type S struct {
					Name string `json:"name" validate:"required"`
				}
				err := v.Validate(&S{Name: "ok"})
				assert.NoError(ctx.T, err)
			})
			s.It("invalid struct returns error with field", func(ctx *specs.Context) {
				v := NewValidator()
				type S struct {
					Name string `json:"name" validate:"required"`
				}
				err := v.Validate(&S{})
				assert.Error(ctx.T, err)
				var verr Errors
				assert.ErrorAs(ctx.T, err, &verr)
				assert.Len(ctx.T, verr.Errors, 1)
				assert.Equal(ctx.T, "name", verr.Errors[0].Field)
			})
		})

		s.When("Validator_ValidateVar", func(s *specs.Spec) {
			s.It("required: empty fails, non-empty passes", func(ctx *specs.Context) {
				v := NewValidator()
				err := v.ValidateVar("", "required")
				assert.Error(ctx.T, err)
				err = v.ValidateVar("x", "required")
				assert.NoError(ctx.T, err)
			})
			s.It("email: valid passes, invalid fails", func(ctx *specs.Context) {
				v := NewValidator()
				assert.NoError(ctx.T, v.ValidateVar("a@b.co", "email"))
				assert.Error(ctx.T, v.ValidateVar("invalid", "email"))
			})
			s.It("domain: empty passes, valid passes, invalid fails", func(ctx *specs.Context) {
				v := NewValidator()
				assert.NoError(ctx.T, v.ValidateVar("", "domain"))
				assert.NoError(ctx.T, v.ValidateVar("example.com", "domain"))
				assert.NoError(ctx.T, v.ValidateVar("a.co", "domain"))
				assert.NoError(ctx.T, v.ValidateVar("sub.example.com", "domain"))
				assert.Error(ctx.T, v.ValidateVar("invalid!.com", "domain"))
				assert.Error(ctx.T, v.ValidateVar("-x.com", "domain"))
			})
			s.It("ulid: empty passes, valid passes, invalid fails", func(ctx *specs.Context) {
				v := NewValidator()
				assert.NoError(ctx.T, v.ValidateVar("", "ulid"))
				assert.NoError(ctx.T, v.ValidateVar("01ARZ3NDEKTSV4RRFFQ69G5FAV", "ulid"))
				assert.Error(ctx.T, v.ValidateVar("not-a-ulid", "ulid"))
				assert.Error(ctx.T, v.ValidateVar("01ARZ3NDEKTSV4RRFFQ69G5FA", "ulid"))
			})
			s.It("positive: numeric positive passes, zero and negative fail", func(ctx *specs.Context) {
				v := NewValidator()
				assert.NoError(ctx.T, v.ValidateVar(1, "positive"))
				assert.NoError(ctx.T, v.ValidateVar(int64(1), "positive"))
				assert.NoError(ctx.T, v.ValidateVar(uint(1), "positive"))
				assert.NoError(ctx.T, v.ValidateVar(1.5, "positive"))
				assert.Error(ctx.T, v.ValidateVar(0, "positive"))
				assert.Error(ctx.T, v.ValidateVar(-1, "positive"))
				assert.Error(ctx.T, v.ValidateVar("x", "positive"))
			})
			s.It("min/max/len: respects constraints", func(ctx *specs.Context) {
				v := NewValidator()
				assert.Error(ctx.T, v.ValidateVar("ab", "min=3"))
				assert.NoError(ctx.T, v.ValidateVar("abc", "min=3"))
				assert.Error(ctx.T, v.ValidateVar("abcd", "max=3"))
				assert.NoError(ctx.T, v.ValidateVar("abc", "len=3"))
				assert.Error(ctx.T, v.ValidateVar("ab", "len=3"))
			})
			s.It("oneof: allowed value passes, other fails", func(ctx *specs.Context) {
				v := NewValidator()
				assert.NoError(ctx.T, v.ValidateVar("a", "oneof=a b c"))
				assert.Error(ctx.T, v.ValidateVar("x", "oneof=a b c"))
			})
		})

		s.When("Errors_Error", func(s *specs.Spec) {
			s.It("returns no validation errors when slice empty", func(ctx *specs.Context) {
				e := Errors{Errors: nil}
				assert.Equal(ctx.T, "no validation errors", e.Error())
			})
			s.It("returns no validation errors when slice has zero length", func(ctx *specs.Context) {
				e := Errors{Errors: []Error{}}
				assert.Equal(ctx.T, "no validation errors", e.Error())
			})
			s.It("formats multiple errors", func(ctx *specs.Context) {
				e := Errors{
					Errors: []Error{
						{Field: "A", Message: "msg1"},
						{Field: "B", Message: "msg2"},
					},
				}
				assert.Equal(ctx.T, "A: msg1; B: msg2", e.Error())
			})
		})

		s.When("Validator_ValidateVar message branches", func(s *specs.Spec) {
			s.It("min tag produces expected message", func(ctx *specs.Context) {
				v := NewValidator()
				err := v.ValidateVar("ab", "min=3")
				assert.Error(ctx.T, err)
				assert.Contains(ctx.T, err.Error(), "at least 3")
			})
			s.It("max tag produces expected message", func(ctx *specs.Context) {
				v := NewValidator()
				err := v.ValidateVar("abcd", "max=3")
				assert.Error(ctx.T, err)
				assert.Contains(ctx.T, err.Error(), "at most 3")
			})
			s.It("len tag produces expected message", func(ctx *specs.Context) {
				v := NewValidator()
				err := v.ValidateVar("ab", "len=3")
				assert.Error(ctx.T, err)
				assert.Contains(ctx.T, err.Error(), "exactly 3")
			})
			s.It("oneof tag produces expected message", func(ctx *specs.Context) {
				v := NewValidator()
				err := v.ValidateVar("x", "oneof=a b c")
				assert.Error(ctx.T, err)
				assert.Contains(ctx.T, err.Error(), "one of")
			})
			s.It("domain tag produces expected message", func(ctx *specs.Context) {
				v := NewValidator()
				err := v.ValidateVar("invalid!.com", "domain")
				assert.Error(ctx.T, err)
				assert.Contains(ctx.T, err.Error(), "valid domain")
			})
			s.It("ulid tag produces expected message", func(ctx *specs.Context) {
				v := NewValidator()
				err := v.ValidateVar("not-ulid", "ulid")
				assert.Error(ctx.T, err)
				assert.Contains(ctx.T, err.Error(), "valid ULID")
			})
			s.It("positive tag produces expected message", func(ctx *specs.Context) {
				v := NewValidator()
				err := v.ValidateVar(-1, "positive")
				assert.Error(ctx.T, err)
				assert.Contains(ctx.T, err.Error(), "positive")
			})
			s.It("unknown tag uses default message format", func(ctx *specs.Context) {
				v := NewValidator()
				// "alpha" requires only letters; "123" fails
				err := v.ValidateVar("123", "alpha")
				assert.Error(ctx.T, err)
				assert.Contains(ctx.T, err.Error(), "failed on the alpha tag")
			})
		})
	})
}
