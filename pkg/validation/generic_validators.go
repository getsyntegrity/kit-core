package validation

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// GenericValidatorFunc represents a generic validation function.
type GenericValidatorFunc[T any] func(ctx context.Context, input T) error

// Rule represents a validation rule with error message.
type Rule[T any] struct {
	Rule    func(input T) bool
	Message string
}

// Result represents the result of validation.
type Result struct {
	IsValid bool
	Errors  []GenericValidationError
}

// GenericValidationError represents a validation error.
type GenericValidationError struct {
	Field   string
	Message string
	Value   interface{}
}

func (e GenericValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// GenericValidator provides generic validation functions.
type GenericValidator[T any] struct {
	rules []Rule[T]
}

// NewGenericValidator creates a new generic validator.
func NewGenericValidator[T any]() *GenericValidator[T] {
	return &GenericValidator[T]{
		rules: make([]Rule[T], 0),
	}
}

// AddRule adds a validation rule.
func (v *GenericValidator[T]) AddRule(rule Rule[T]) *GenericValidator[T] {
	v.rules = append(v.rules, rule)
	return v
}

// Validate validates the input against all rules.
func (v *GenericValidator[T]) Validate(_ context.Context, input T) error {
	for _, rule := range v.rules {
		if !rule.Rule(input) {
			return fmt.Errorf("%s", rule.Message)
		}
	}
	return nil
}

// Common validation functions using high-order functions

// Required validates that a field is not empty.
func Required[T any](fieldExtractor func(T) string, fieldName string) Rule[T] {
	return Rule[T]{
		Rule: func(input T) bool {
			value := fieldExtractor(input)
			return strings.TrimSpace(value) != ""
		},
		Message: fmt.Sprintf("%s is required", fieldName),
	}
}

// MinLength validates minimum length.
func MinLength[T any](fieldExtractor func(T) string, fieldName string, minLength int) Rule[T] {
	return Rule[T]{
		Rule: func(input T) bool {
			value := fieldExtractor(input)
			return len(strings.TrimSpace(value)) >= minLength
		},
		Message: fmt.Sprintf("%s must be at least %d characters long", fieldName, minLength),
	}
}

// MaxLength validates maximum length.
func MaxLength[T any](fieldExtractor func(T) string, fieldName string, maxLength int) Rule[T] {
	return Rule[T]{
		Rule: func(input T) bool {
			value := fieldExtractor(input)
			return len(strings.TrimSpace(value)) <= maxLength
		},
		Message: fmt.Sprintf("%s must be at most %d characters long", fieldName, maxLength),
	}
}

// Email validates email format.
func Email[T any](fieldExtractor func(T) string, fieldName string) Rule[T] {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return Rule[T]{
		Rule: func(input T) bool {
			value := fieldExtractor(input)
			return emailRegex.MatchString(value)
		},
		Message: fmt.Sprintf("%s must be a valid email address", fieldName),
	}
}

// Domain validates domain format.
func Domain[T any](fieldExtractor func(T) string, fieldName string) Rule[T] {
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	return Rule[T]{
		Rule: func(input T) bool {
			value := fieldExtractor(input)
			return domainRegex.MatchString(value)
		},
		Message: fmt.Sprintf("%s must be a valid domain name", fieldName),
	}
}

// Custom validates with a custom function.
func Custom[T any](validator func(T) bool, message string) Rule[T] {
	return Rule[T]{
		Rule:    validator,
		Message: message,
	}
}

// Compose combines multiple validators.
func Compose[T any](validators ...Rule[T]) Rule[T] {
	return Rule[T]{
		Rule: func(input T) bool {
			for _, validator := range validators {
				if !validator.Rule(input) {
					return false
				}
			}
			return true
		},
		Message: "validation failed",
	}
}

// FieldValidator creates a high-order function to create field-specific validators.
func FieldValidator[T any, F any](fieldExtractor func(T) F) *FieldValidatorBuilder[T, F] {
	return &FieldValidatorBuilder[T, F]{
		fieldExtractor: fieldExtractor,
	}
}

// FieldValidatorBuilder helps build field-specific validators.
type FieldValidatorBuilder[T any, F any] struct {
	fieldExtractor func(T) F
}

// Required creates a required validator for the field.
func (b *FieldValidatorBuilder[T, F]) Required(fieldName string) Rule[T] {
	return Rule[T]{
		Rule: func(input T) bool {
			value := b.fieldExtractor(input)
			return !reflect.ValueOf(value).IsZero()
		},
		Message: fmt.Sprintf("%s is required", fieldName),
	}
}

// StringFieldValidatorBuilder for string fields.
type StringFieldValidatorBuilder[T any] struct {
	fieldExtractor func(T) string
}

// StringField creates a string field validator.
func StringField[T any](fieldExtractor func(T) string) *StringFieldValidatorBuilder[T] {
	return &StringFieldValidatorBuilder[T]{
		fieldExtractor: fieldExtractor,
	}
}

// Required validates that the string field is not empty.
func (b *StringFieldValidatorBuilder[T]) Required(fieldName string) Rule[T] {
	return Rule[T]{
		Rule: func(input T) bool {
			value := b.fieldExtractor(input)
			return strings.TrimSpace(value) != ""
		},
		Message: fmt.Sprintf("%s is required", fieldName),
	}
}

// MinLength validates minimum length.
func (b *StringFieldValidatorBuilder[T]) MinLength(fieldName string, minLength int) Rule[T] {
	return Rule[T]{
		Rule: func(input T) bool {
			value := b.fieldExtractor(input)
			return len(strings.TrimSpace(value)) >= minLength
		},
		Message: fmt.Sprintf("%s must be at least %d characters long", fieldName, minLength),
	}
}

// MaxLength validates maximum length.
func (b *StringFieldValidatorBuilder[T]) MaxLength(fieldName string, maxLength int) Rule[T] {
	return Rule[T]{
		Rule: func(input T) bool {
			value := b.fieldExtractor(input)
			return len(strings.TrimSpace(value)) <= maxLength
		},
		Message: fmt.Sprintf("%s must be at most %d characters long", fieldName, maxLength),
	}
}

// Email validates email format.
func (b *StringFieldValidatorBuilder[T]) Email(fieldName string) Rule[T] {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return Rule[T]{
		Rule: func(input T) bool {
			value := b.fieldExtractor(input)
			return emailRegex.MatchString(value)
		},
		Message: fmt.Sprintf("%s must be a valid email address", fieldName),
	}
}

// Domain validates domain format.
func (b *StringFieldValidatorBuilder[T]) Domain(fieldName string) Rule[T] {
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	return Rule[T]{
		Rule: func(input T) bool {
			value := b.fieldExtractor(input)
			return domainRegex.MatchString(value)
		},
		Message: fmt.Sprintf("%s must be a valid domain name", fieldName),
	}
}
