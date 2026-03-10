package validation

// ErrorType represents the type of validation error.
type ErrorType string

const (
	//ErrorTypeRequired required field missing or empty
	ErrorTypeRequired ErrorType = "required"
	//ErrorTypeInvalid invalid value
	ErrorTypeInvalid ErrorType = "invalid"
	//ErrorTypeTooShort value below minimum length
	ErrorTypeTooShort ErrorType = "too_short"
	//ErrorTypeTooLong value above maximum length
	ErrorTypeTooLong ErrorType = "too_long"
	//ErrorTypeInvalidFormat value does not match format
	ErrorTypeInvalidFormat ErrorType = "invalid_format"
	//ErrorTypeInvalidRange value out of allowed range
	ErrorTypeInvalidRange ErrorType = "invalid_range"
	//ErrorTypeInvalidEnum value not in allowed set
	ErrorTypeInvalidEnum ErrorType = "invalid_enum"
	//ErrorTypeCustom custom validation error
	ErrorTypeCustom ErrorType = "custom"
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
