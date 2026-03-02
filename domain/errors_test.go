package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name: "error without cause",
			err: &Error{
				Code:    "TEST_CODE",
				Message: "test message",
			},
			expected: "TEST_CODE: test message",
		},
		{
			name: "error with cause",
			err: &Error{
				Code:    "TEST_CODE",
				Message: "test message",
				Cause:   errors.New("original error"),
			},
			expected: "TEST_CODE: test message (caused by: original error)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	err := &Error{
		Code:    "TEST_CODE",
		Message: "test message",
		Cause:   originalErr,
	}

	assert.Equal(t, originalErr, err.Unwrap())
}

func TestError_Unwrap_NilCause(t *testing.T) {
	err := &Error{
		Code:    "TEST_CODE",
		Message: "test message",
	}

	assert.Nil(t, err.Unwrap())
}

func TestNewError(t *testing.T) {
	err := NewError("CODE", "message")

	assert.Equal(t, "CODE", err.Code)
	assert.Equal(t, "message", err.Message)
	assert.Nil(t, err.Cause)
}

func TestNewErrorWithCause(t *testing.T) {
	cause := errors.New("original error")
	err := NewErrorWithCause("CODE", "message", cause)

	assert.Equal(t, "CODE", err.Code)
	assert.Equal(t, "message", err.Message)
	assert.Equal(t, cause, err.Cause)
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "not found error",
			err:      ErrNotFound,
			expected: true,
		},
		{
			name: "wrapped not found error",
			err: &Error{
				Code:    ErrCodeNotFound,
				Message: "not found",
			},
			expected: true,
		},
		{
			name:     "different error code",
			err:      ErrInternal,
			expected: false,
		},
		{
			name: "wrapped different error",
			err: &Error{
				Code:    ErrCodeInternal,
				Message: "internal error",
			},
			expected: false,
		},
		{
			name:     "wrapped using errors.Is",
			err:      errors.New("wrapped: not found"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotFound(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsInvalidState(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "invalid state error",
			err:      ErrInvalidState,
			expected: true,
		},
		{
			name: "wrapped invalid state error",
			err: &Error{
				Code:    ErrCodeInvalidState,
				Message: "invalid state",
			},
			expected: true,
		},
		{
			name:     "different error code",
			err:      ErrInternal,
			expected: false,
		},
		{
			name: "wrapped different error",
			err: &Error{
				Code:    ErrCodeInternal,
				Message: "internal error",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsInvalidState(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsConcurrency(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "concurrency error",
			err:      ErrConcurrency,
			expected: true,
		},
		{
			name: "wrapped concurrency error",
			err: &Error{
				Code:    ErrCodeConcurrency,
				Message: "concurrency conflict",
			},
			expected: true,
		},
		{
			name:     "different error code",
			err:      ErrInternal,
			expected: false,
		},
		{
			name: "wrapped different error",
			err: &Error{
				Code:    ErrCodeInternal,
				Message: "internal error",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsConcurrency(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidation(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "validation error",
			err:      ErrValidation,
			expected: true,
		},
		{
			name: "wrapped validation error",
			err: &Error{
				Code:    ErrCodeValidation,
				Message: "validation failed",
			},
			expected: true,
		},
		{
			name:     "different error code",
			err:      ErrInternal,
			expected: false,
		},
		{
			name: "wrapped different error",
			err: &Error{
				Code:    ErrCodeInternal,
				Message: "internal error",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidation(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
