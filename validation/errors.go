package validation

// ErrorType represents the type of validation error.
type ErrorType string

const (
	ErrorTypeRequired      ErrorType = "required"
	ErrorTypeInvalid       ErrorType = "invalid"
	ErrorTypeTooShort      ErrorType = "too_short"
	ErrorTypeTooLong       ErrorType = "too_long"
	ErrorTypeInvalidFormat ErrorType = "invalid_format"
	ErrorTypeInvalidRange  ErrorType = "invalid_range"
	ErrorTypeInvalidEnum   ErrorType = "invalid_enum"
	ErrorTypeCustom        ErrorType = "custom"
)

// GenericValidationError represents a single validation error.
type GenericValidationError struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

// Error returns the error message.
func (e GenericValidationError) Error() string {
	return e.Message
}

// GenericValidationErrors represents multiple generic validation errors.
type GenericValidationErrors struct {
	Errors []GenericValidationError `json:"errors"`
}
