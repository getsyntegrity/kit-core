package domain

import (
	"errors"
	"testing"

	"github.com/pablogore/go-specs/specs"
	"github.com/stretchr/testify/assert"
)

func TestDomainErrors(t *testing.T) {
	specs.Describe(t, "domain errors", func(s *specs.Spec) {
		s.When("Error", func(s *specs.Spec) {
			s.It("formats without cause", func(ctx *specs.Context) {
				err := &Error{Code: "TEST_CODE", Message: "test message"}
				assert.Equal(ctx.T, "TEST_CODE: test message", err.Error())
			})
			s.It("formats with cause", func(ctx *specs.Context) {
				err := &Error{
					Code:    "TEST_CODE",
					Message: "test message",
					Cause:   errors.New("original error"),
				}
				assert.Equal(ctx.T, "TEST_CODE: test message (caused by: original error)", err.Error())
			})
			s.It("Unwrap returns cause", func(ctx *specs.Context) {
				originalErr := errors.New("original error")
				err := &Error{Code: "TEST_CODE", Message: "test message", Cause: originalErr}
				assert.Equal(ctx.T, originalErr, err.Unwrap())
			})
			s.It("Unwrap returns nil when no cause", func(ctx *specs.Context) {
				err := &Error{Code: "TEST_CODE", Message: "test message"}
				assert.Nil(ctx.T, err.Unwrap())
			})
		})

		s.When("NewError", func(s *specs.Spec) {
			s.It("sets code and message and no cause", func(ctx *specs.Context) {
				err := NewError("CODE", "message")
				assert.Equal(ctx.T, "CODE", err.Code)
				assert.Equal(ctx.T, "message", err.Message)
				assert.Nil(ctx.T, err.Cause)
			})
		})

		s.When("NewErrorWithCause", func(s *specs.Spec) {
			s.It("sets code message and cause", func(ctx *specs.Context) {
				cause := errors.New("original error")
				err := NewErrorWithCause("CODE", "message", cause)
				assert.Equal(ctx.T, "CODE", err.Code)
				assert.Equal(ctx.T, "message", err.Message)
				assert.Equal(ctx.T, cause, err.Cause)
			})
		})

		s.When("IsNotFound", func(s *specs.Spec) {
			s.It("covers all cases", func(ctx *specs.Context) {
				tests := []struct {
					name     string
					err      error
					expected bool
				}{
					{"nil error", nil, false},
					{"not found error", ErrNotFound, true},
					{"wrapped not found error", &Error{Code: ErrCodeNotFound, Message: "not found"}, true},
					{"different error code", ErrInternal, false},
					{"wrapped different error", &Error{Code: ErrCodeInternal, Message: "internal error"}, false},
					{"wrapped using errors.Is", errors.New("wrapped: not found"), false},
				}
				for _, tt := range tests {
					result := IsNotFound(tt.err)
					assert.Equal(ctx.T, tt.expected, result, "case: %s", tt.name)
				}
			})
		})

		s.When("IsInvalidState", func(s *specs.Spec) {
			s.It("covers all cases", func(ctx *specs.Context) {
				tests := []struct {
					name     string
					err      error
					expected bool
				}{
					{"nil error", nil, false},
					{"invalid state error", ErrInvalidState, true},
					{"wrapped invalid state error", &Error{Code: ErrCodeInvalidState, Message: "invalid state"}, true},
					{"different error code", ErrInternal, false},
					{"wrapped different error", &Error{Code: ErrCodeInternal, Message: "internal error"}, false},
				}
				for _, tt := range tests {
					result := IsInvalidState(tt.err)
					assert.Equal(ctx.T, tt.expected, result, "case: %s", tt.name)
				}
			})
		})

		s.When("IsConcurrency", func(s *specs.Spec) {
			s.It("covers all cases", func(ctx *specs.Context) {
				tests := []struct {
					name     string
					err      error
					expected bool
				}{
					{"nil error", nil, false},
					{"concurrency error", ErrConcurrency, true},
					{"wrapped concurrency error", &Error{Code: ErrCodeConcurrency, Message: "concurrency conflict"}, true},
					{"different error code", ErrInternal, false},
					{"wrapped different error", &Error{Code: ErrCodeInternal, Message: "internal error"}, false},
				}
				for _, tt := range tests {
					result := IsConcurrency(tt.err)
					assert.Equal(ctx.T, tt.expected, result, "case: %s", tt.name)
				}
			})
		})

		s.When("IsValidation", func(s *specs.Spec) {
			s.It("covers all cases", func(ctx *specs.Context) {
				tests := []struct {
					name     string
					err      error
					expected bool
				}{
					{"nil error", nil, false},
					{"validation error", ErrValidation, true},
					{"wrapped validation error", &Error{Code: ErrCodeValidation, Message: "validation failed"}, true},
					{"different error code", ErrInternal, false},
					{"wrapped different error", &Error{Code: ErrCodeInternal, Message: "internal error"}, false},
				}
				for _, tt := range tests {
					result := IsValidation(tt.err)
					assert.Equal(ctx.T, tt.expected, result, "case: %s", tt.name)
				}
			})
		})
	})
}
