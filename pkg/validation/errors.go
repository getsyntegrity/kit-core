package validation

import (
	"fmt"
	"strings"
)

// ErrorType represents the type of validation error.
type ErrorType string

const (
	// ErrorTypeRequired indicates a required field is missing.
	ErrorTypeRequired ErrorType = "required"

	// ErrorTypeInvalid indicates an invalid value.
	ErrorTypeInvalid ErrorType = "invalid"

	// ErrorTypeTooShort indicates a value is too short.
	ErrorTypeTooShort ErrorType = "too_short"

	// ErrorTypeTooLong indicates a value is too long.
	ErrorTypeTooLong ErrorType = "too_long"

	// ErrorTypeInvalidFormat indicates an invalid format.
	ErrorTypeInvalidFormat ErrorType = "invalid_format"

	// ErrorTypeInvalidRange indicates a value is out of range.
	ErrorTypeInvalidRange ErrorType = "invalid_range"

	// ErrorTypeInvalidEnum indicates an invalid enum value.
	ErrorTypeInvalidEnum ErrorType = "invalid_enum"

	// ErrorTypeCustom indicates a custom validation error.
	ErrorTypeCustom ErrorType = "custom"
)

// ErrorBuilder provides utility functions for building validation errors.
type ErrorBuilder struct{}

// NewErrorBuilder creates a new validation error builder.
func NewErrorBuilder() *ErrorBuilder {
	return &ErrorBuilder{}
}

// BuildRequiredError creates a required field error.
func (b *ErrorBuilder) BuildRequiredError(field string) GenericValidationError {
	return GenericValidationError{
		Field:   field,
		Message: fmt.Sprintf("%s is required", field),
	}
}

// BuildInvalidEmailError creates an invalid email error.
func (b *ErrorBuilder) BuildInvalidEmailError(field string, value interface{}) GenericValidationError {
	return GenericValidationError{
		Field:   field,
		Value:   value,
		Message: fmt.Sprintf("%s must be a valid email address", field),
	}
}

// BuildInvalidURLError creates an invalid URL error.
func (b *ErrorBuilder) BuildInvalidURLError(field string, value interface{}) GenericValidationError {
	return GenericValidationError{
		Field:   field,
		Value:   value,
		Message: fmt.Sprintf("%s must be a valid URL", field),
	}
}

// BuildInvalidULIDError creates an invalid ULID error.
func (b *ErrorBuilder) BuildInvalidULIDError(field string, value interface{}) GenericValidationError {
	return GenericValidationError{
		Field:   field,
		Value:   value,
		Message: fmt.Sprintf("%s must be a valid ULID", field),
	}
}

// BuildInvalidDomainError creates an invalid domain error.
func (b *ErrorBuilder) BuildInvalidDomainError(field string, value interface{}) GenericValidationError {
	return GenericValidationError{
		Field:   field,
		Value:   value,
		Message: fmt.Sprintf("%s must be a valid domain name", field),
	}
}

// BuildTooShortError creates a too short error.
func (b *ErrorBuilder) BuildTooShortError(field string, value interface{}, minLength int) GenericValidationError {
	return GenericValidationError{
		Field:   field,
		Value:   value,
		Message: fmt.Sprintf("%s must be at least %d characters long", field, minLength),
	}
}

// BuildTooLongError creates a too long error.
func (b *ErrorBuilder) BuildTooLongError(field string, value interface{}, maxLength int) GenericValidationError {
	return GenericValidationError{
		Field:   field,
		Value:   value,
		Message: fmt.Sprintf("%s must be at most %d characters long", field, maxLength),
	}
}

// BuildInvalidEnumError creates an invalid enum error.
func (b *ErrorBuilder) BuildInvalidEnumError(field string, value interface{}, validValues []string) GenericValidationError {
	return GenericValidationError{
		Field:   field,
		Value:   value,
		Message: fmt.Sprintf("%s must be one of: %s", field, strings.Join(validValues, ", ")),
	}
}

// BuildInvalidRangeError creates an invalid range error.
func (b *ErrorBuilder) BuildInvalidRangeError(field string, value interface{}, minVal interface{}, maxVal interface{}) GenericValidationError {
	var message string
	switch {
	case minVal != nil && maxVal != nil:
		message = fmt.Sprintf("%s must be between %v and %v", field, minVal, maxVal)
	case minVal != nil:
		message = fmt.Sprintf("%s must be at least %v", field, minVal)
	case maxVal != nil:
		message = fmt.Sprintf("%s must be at most %v", field, maxVal)
	}

	return GenericValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

// BuildCustomError creates a custom validation error.
func (b *ErrorBuilder) BuildCustomError(field string, value interface{}, message string) GenericValidationError {
	return GenericValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

// GenericValidationErrors represents multiple generic validation errors.
type GenericValidationErrors struct {
	Errors []GenericValidationError `json:"errors"`
}

// ErrorsBuilder provides utility functions for building validation error collections.
type ErrorsBuilder struct{}

// NewErrorsBuilder creates a new validation errors builder.
func NewErrorsBuilder() *ErrorsBuilder {
	return &ErrorsBuilder{}
}

// NewValidationErrors creates a new GenericValidationErrors instance.
func (b *ErrorsBuilder) NewValidationErrors() *GenericValidationErrors {
	return &GenericValidationErrors{
		Errors: []GenericValidationError{},
	}
}

// AddError adds a new validation error to the collection.
func (b *ErrorsBuilder) AddError(errors *GenericValidationErrors, err GenericValidationError) {
	errors.Errors = append(errors.Errors, err)
}

// AddRequiredError adds a required field error.
func (b *ErrorsBuilder) AddRequiredError(errors *GenericValidationErrors, field string) {
	builder := NewErrorBuilder()
	b.AddError(errors, builder.BuildRequiredError(field))
}

// AddInvalidError adds an invalid value error.
func (b *ErrorsBuilder) AddInvalidError(errors *GenericValidationErrors, field string, value interface{}, message string) {
	builder := NewErrorBuilder()
	b.AddError(errors, builder.BuildCustomError(field, value, message))
}

// AddFormatError adds an invalid format error.
func (b *ErrorsBuilder) AddFormatError(errors *GenericValidationErrors, field string, value interface{}, expectedFormat string) {
	builder := NewErrorBuilder()
	b.AddError(errors, builder.BuildCustomError(field, value, fmt.Sprintf("%s must be in format: %s", field, expectedFormat)))
}

// AddRangeError adds a range error.
func (b *ErrorsBuilder) AddRangeError(errors *GenericValidationErrors, field string, value interface{}, minVal interface{}, maxVal interface{}) {
	builder := NewErrorBuilder()
	b.AddError(errors, builder.BuildInvalidRangeError(field, value, minVal, maxVal))
}

// AddEnumError adds an enum error.
func (b *ErrorsBuilder) AddEnumError(errors *GenericValidationErrors, field string, value interface{}, validValues []string) {
	builder := NewErrorBuilder()
	b.AddError(errors, builder.BuildInvalidEnumError(field, value, validValues))
}

// Utility functions for working with ValidationErrors

// HasErrors returns true if there are validation errors.
func HasErrors(errors *GenericValidationErrors) bool {
	return errors != nil && len(errors.Errors) > 0
}

// GetErrorsByField returns all errors for a specific field.
func GetErrorsByField(errors *GenericValidationErrors, field string) []GenericValidationError {
	var fieldErrors []GenericValidationError
	if errors == nil {
		return fieldErrors
	}

	for _, err := range errors.Errors {
		if err.Field == field {
			fieldErrors = append(fieldErrors, err)
		}
	}
	return fieldErrors
}

// GetFirstError returns the first validation error.
func GetFirstError(errors *GenericValidationErrors) *GenericValidationError {
	if errors == nil || len(errors.Errors) == 0 {
		return nil
	}
	return &errors.Errors[0]
}

// GetErrorMessages returns all error messages as a slice.
func GetErrorMessages(errors *GenericValidationErrors) []string {
	var messages []string
	if errors == nil {
		return messages
	}

	for _, err := range errors.Errors {
		messages = append(messages, err.Message)
	}
	return messages
}

// GetErrorMessagesByField returns error messages for a specific field.
func GetErrorMessagesByField(errors *GenericValidationErrors, field string) []string {
	var messages []string
	fieldErrors := GetErrorsByField(errors, field)

	for _, err := range fieldErrors {
		messages = append(messages, err.Message)
	}
	return messages
}
