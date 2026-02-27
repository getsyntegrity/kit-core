package domain

import (
	"errors"
	"fmt"
)

// Error represents a domain-specific error.
type Error struct {
	Code    string
	Message string
	Cause   error
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	return e.Cause
}

// NewError creates a new domain error.
func NewError(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewErrorWithCause creates a new domain error with a cause.
func NewErrorWithCause(code, message string, cause error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Common domain error codes.
const (
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeInvalidState    = "INVALID_STATE"
	ErrCodeConcurrency     = "CONCURRENCY_CONFLICT"
	ErrCodeValidation      = "VALIDATION_FAILED"
	ErrCodeUnauthorized    = "UNAUTHORIZED"
	ErrCodeForbidden       = "FORBIDDEN"
	ErrCodeInternal        = "INTERNAL_ERROR"
	ErrCodeExternalService = "EXTERNAL_SERVICE_ERROR"
)

// Common domain errors.
var (
	ErrNotFound        = NewError(ErrCodeNotFound, "entity not found")
	ErrInvalidState    = NewError(ErrCodeInvalidState, "invalid state")
	ErrConcurrency     = NewError(ErrCodeConcurrency, "concurrency conflict")
	ErrValidation      = NewError(ErrCodeValidation, "validation failed")
	ErrUnauthorized    = NewError(ErrCodeUnauthorized, "unauthorized")
	ErrForbidden       = NewError(ErrCodeForbidden, "forbidden")
	ErrInternal        = NewError(ErrCodeInternal, "internal error")
	ErrExternalService = NewError(ErrCodeExternalService, "external service error")
)

// IsNotFound checks if an error is a not found error.
func IsNotFound(err error) bool {
	var domainErr *Error
	if errors.As(err, &domainErr) {
		return domainErr.Code == ErrCodeNotFound
	}
	return errors.Is(err, ErrNotFound)
}

// IsInvalidState checks if an error is an invalid state error.
func IsInvalidState(err error) bool {
	var domainErr *Error
	if errors.As(err, &domainErr) {
		return domainErr.Code == ErrCodeInvalidState
	}
	return errors.Is(err, ErrInvalidState)
}

// IsConcurrency checks if an error is a concurrency error.
func IsConcurrency(err error) bool {
	var domainErr *Error
	if errors.As(err, &domainErr) {
		return domainErr.Code == ErrCodeConcurrency
	}
	return errors.Is(err, ErrConcurrency)
}

// IsValidation checks if an error is a validation error.
func IsValidation(err error) bool {
	var domainErr *Error
	if errors.As(err, &domainErr) {
		return domainErr.Code == ErrCodeValidation
	}
	return errors.Is(err, ErrValidation)
}
