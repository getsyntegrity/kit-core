package repository

import "fmt"

// Error represents a repository-specific error.
type Error struct {
	Code    string
	Message string
	Err     error
}

// Error returns the error message.
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	return e.Err
}

// NewError creates a new repository error.
func NewError(code, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common repository error codes.
const (
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeAlreadyExists = "ALREADY_EXISTS"
	ErrCodeInvalidInput  = "INVALID_INPUT"
	ErrCodeDatabase      = "DATABASE_ERROR"
	ErrCodeConnection    = "CONNECTION_ERROR"
)

// NotFoundError creates a not found error.
func NotFoundError(entityType, id string) *Error {
	return NewError(
		ErrCodeNotFound,
		fmt.Sprintf("%s with id %s not found", entityType, id),
		nil,
	)
}

// AlreadyExistsError creates an already exists error.
func AlreadyExistsError(entityType, id string) *Error {
	return NewError(
		ErrCodeAlreadyExists,
		fmt.Sprintf("%s with id %s already exists", entityType, id),
		nil,
	)
}

// InvalidInputError creates an invalid input error.
func InvalidInputError(message string) *Error {
	return NewError(ErrCodeInvalidInput, message, nil)
}

// DatabaseError creates a database error.
func DatabaseError(message string, err error) *Error {
	return NewError(ErrCodeDatabase, message, err)
}

// ConnectionError creates a connection error.
func ConnectionError(message string, err error) *Error {
	return NewError(ErrCodeConnection, message, err)
}
