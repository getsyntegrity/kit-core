// Package validation provides generic validation utilities for struct validation.
//
// OWNERSHIP RULES (CRITICAL):
//   - This package MUST remain domain-agnostic and service-agnostic.
//   - Do NOT add domain-specific validators (e.g., tenant types, approval statuses).
//   - Do NOT encode business rules or domain invariants.
//   - Domain-specific validation belongs in services/*, not in libs/*.
//
// Generic validators (e.g., email, domain, ULID) are acceptable.
// Service-specific validators (e.g., tenant_type, tenant_status) are FORBIDDEN.
package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator provides robust validation with custom rules.
type Validator struct {
	validate *validator.Validate
}

// Error represents a validation error with field details.
type Error struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// Errors represents multiple validation errors.
type Errors struct {
	Errors []Error `json:"errors"`
}

// Error returns a formatted error message.
func (v Errors) Error() string {
	if len(v.Errors) == 0 {
		return "no validation errors"
	}

	var messages []string
	for _, err := range v.Errors {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, "; ")
}

// NewValidator creates a new validator with custom rules.
// Returns an error if any custom validator registration fails.
func NewValidator() (*Validator, error) {
	validatorInstance := validator.New()

	// Register custom validators
	if err := registerCustomValidators(validatorInstance); err != nil {
		return nil, err
	}

	// Register custom tag name function
	validatorInstance.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{
		validate: validatorInstance,
	}, nil
}

// Validate validates a struct and returns detailed errors.
func (v *Validator) Validate(structToValidate interface{}) error {
	err := v.validate.Struct(structToValidate)
	if err == nil {
		return nil
	}

	var validationErrors Errors

	for _, err := range err.(validator.ValidationErrors) {
		validationError := Error{
			Field: err.Field(),
			Tag:   err.Tag(),
			Value: fmt.Sprintf("%v", err.Value()),
		}

		// Generate user-friendly error messages
		validationError.Message = v.generateErrorMessage(err)
		validationErrors.Errors = append(validationErrors.Errors, validationError)
	}

	return validationErrors
}

// ValidateVar validates a single field.
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	err := v.validate.Var(field, tag)
	if err == nil {
		return nil
	}

	var validationErrors Errors

	for _, err := range err.(validator.ValidationErrors) {
		validationError := Error{
			Field:   err.Field(),
			Tag:     err.Tag(),
			Value:   fmt.Sprintf("%v", err.Value()),
			Message: v.generateErrorMessage(err),
		}
		validationErrors.Errors = append(validationErrors.Errors, validationError)
	}

	return validationErrors
}

// generateErrorMessage generates user-friendly error messages.
func (v *Validator) generateErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", err.Field(), err.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", err.Field(), err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param())
	case "domain":
		return fmt.Sprintf("%s must be a valid domain name", err.Field())
	case "tenant_type":
		return fmt.Sprintf("%s must be one of: CUSTOMER, SYSTEM, PARTNER", err.Field())
	case "tenant_status":
		return fmt.Sprintf("%s must be one of: ACTIVE, INACTIVE, SUSPENDED, DELETED", err.Field())
	case "positive":
		return fmt.Sprintf("%s must be a positive number", err.Field())
	case "ulid":
		return fmt.Sprintf("%s must be a valid ULID", err.Field())
	default:
		return fmt.Sprintf("%s failed on the %s tag", err.Field(), err.Tag())
	}
}

// registerCustomValidators registers custom validation rules.
// Returns the first registration error if any.
//
// NOTE: Domain-specific validators (tenant_type, tenant_status) are DEPRECATED
// and will be removed. Services should implement their own domain-specific validation.
func registerCustomValidators(validatorInstance *validator.Validate) error {
	// Generic validators (domain-agnostic)
	if err := validatorInstance.RegisterValidation("domain", validateDomain); err != nil {
		return fmt.Errorf("register domain validation: %w", err)
	}
	if err := validatorInstance.RegisterValidation("ulid", validateULID); err != nil {
		return fmt.Errorf("register ULID validation: %w", err)
	}

	// DEPRECATED: Domain-specific validators - DO NOT ADD MORE
	if err := validatorInstance.RegisterValidation("tenant_type", validateTenantType); err != nil {
		return fmt.Errorf("register tenant type validation: %w", err)
	}
	if err := validatorInstance.RegisterValidation("tenant_status", validateTenantStatus); err != nil {
		return fmt.Errorf("register tenant status validation: %w", err)
	}

	if err := validatorInstance.RegisterValidation("positive", validatePositive); err != nil {
		return fmt.Errorf("register positive number validation: %w", err)
	}
	return nil
}

// validateDomain validates domain names.
func validateDomain(fl validator.FieldLevel) bool {
	domain := fl.Field().String()
	if domain == "" {
		return true // Let required handle empty values
	}

	// Basic domain validation regex
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	return domainRegex.MatchString(domain)
}

// validateTenantType validates tenant types.
//
// DEPRECATED: This function encodes domain-specific business rules (tenant types)
// and violates libs/* ownership rules. This will be removed in Phase 2 extraction.
// Services should implement their own tenant type validation.
func validateTenantType(fl validator.FieldLevel) bool {
	tenantType := fl.Field().String()
	validTypes := []string{"CUSTOMER", "SYSTEM", "PARTNER"}

	for _, validType := range validTypes {
		if tenantType == validType {
			return true
		}
	}
	return false
}

// validateTenantStatus validates tenant statuses.
//
// DEPRECATED: This function encodes domain-specific business rules (tenant statuses)
// and violates libs/* ownership rules. This will be removed in Phase 2 extraction.
// Services should implement their own tenant status validation.
func validateTenantStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	validStatuses := []string{"ACTIVE", "INACTIVE", "SUSPENDED", "DELETED"}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

// validateULID validates ULID format.
func validateULID(fl validator.FieldLevel) bool {
	ulidStr := fl.Field().String()
	if ulidStr == "" {
		return true // Let required handle empty values
	}
	return ulidRegex.MatchString(ulidStr)
}

// validatePositive validates positive numbers.
func validatePositive(fieldLevel validator.FieldLevel) bool {
	switch fieldLevel.Field().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fieldLevel.Field().Int() > 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fieldLevel.Field().Uint() > 0
	case reflect.Float32, reflect.Float64:
		return fieldLevel.Field().Float() > 0
	default:
		return false
	}
}
