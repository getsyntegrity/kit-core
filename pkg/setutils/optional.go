// Package setutils provides utility functions, constants, and optional value handling.
// This package contains common utilities used across the go-kit framework.
package setutils

import "errors"

// OptionalVal is a generic wrapper that represents an optional value.
type OptionalVal[T any] struct {
	Value T
	Valid bool
}

// NewOptionalVal creates a new OptionalVal with a valid value.
//
// val: The value to be wrapped in the OptionalVal.
//
// Returns: An OptionalVal with the provided value and Valid set to true.
func NewOptionalVal[T any](val T) OptionalVal[T] {
	return OptionalVal[T]{
		Valid: true,
		Value: val,
	}
}

// EmptyOptionalVal creates an OptionalVal with no valid value (empty).
//
// Returns: An OptionalVal with Valid set to false and Value set to its zero value.
func EmptyOptionalVal[T any]() OptionalVal[T] {
	var zero T
	return OptionalVal[T]{
		Valid: false,
		Value: zero,
	}
}

// IsValid returns true if the OptionalVal contains a valid value.
//
// Returns: True if the OptionalVal's Valid field is true, false otherwise.
func (o OptionalVal[T]) IsValid() bool {
	return o.Valid
}

// Get returns the value if it's valid, otherwise it returns the zero value for T and an error.
//
// Returns: The wrapped value if Valid is true, otherwise the zero value for T and an error.
func (o OptionalVal[T]) Get() (T, error) {
	if !o.Valid {
		var zero T
		return zero, errors.New("no valid value")
	}
	return o.Value, nil
}

// OrElse returns the value if it's valid, otherwise it returns the provided default value.
//
// defaultVal: The value to be returned if the OptionalVal's Valid field is false.
//
// Returns: The wrapped value if Valid is true, otherwise the provided default value.
func (o OptionalVal[T]) OrElse(defaultVal T) T {
	if o.Valid {
		return o.Value
	}
	return defaultVal
}
