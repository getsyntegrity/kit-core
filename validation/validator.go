package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/getsyntegrity/kit-core/idgen"
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
func NewValidator() *Validator {
	validatorInstance := validator.New()
	registerCustomValidators(validatorInstance)
	validatorInstance.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &Validator{validate: validatorInstance}
}

// Validate validates a struct and returns detailed errors.
func (v *Validator) Validate(structToValidate interface{}) error {
	err := v.validate.Struct(structToValidate)
	if err == nil {
		return nil
	}
	var validationErrors Errors
	for _, e := range err.(validator.ValidationErrors) {
		validationErrors.Errors = append(validationErrors.Errors, Error{
			Field:   e.Field(),
			Tag:     e.Tag(),
			Value:   fmt.Sprintf("%v", e.Value()),
			Message: v.generateErrorMessage(e),
		})
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
	for _, e := range err.(validator.ValidationErrors) {
		validationErrors.Errors = append(validationErrors.Errors, Error{
			Field:   e.Field(),
			Tag:     e.Tag(),
			Value:   fmt.Sprintf("%v", e.Value()),
			Message: v.generateErrorMessage(e),
		})
	}
	return validationErrors
}

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
	case "ulid":
		return fmt.Sprintf("%s must be a valid ULID", err.Field())
	case "positive":
		return fmt.Sprintf("%s must be a positive number", err.Field())
	default:
		return fmt.Sprintf("%s failed on the %s tag", err.Field(), err.Tag())
	}
}

func registerCustomValidators(validatorInstance *validator.Validate) {
	_ = validatorInstance.RegisterValidation("domain", validateDomain)
	_ = validatorInstance.RegisterValidation("ulid", validateULID)
	_ = validatorInstance.RegisterValidation("positive", validatePositive)
}

func validateDomain(fl validator.FieldLevel) bool {
	domain := fl.Field().String()
	if domain == "" {
		return true
	}
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	return domainRegex.MatchString(domain)
}

func validateULID(fl validator.FieldLevel) bool {
	ulidStr := fl.Field().String()
	if ulidStr == "" {
		return true
	}
	return idgen.IsValidULID(ulidStr)
}

func validatePositive(fl validator.FieldLevel) bool {
	switch fl.Field().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fl.Field().Int() > 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fl.Field().Uint() > 0
	case reflect.Float32, reflect.Float64:
		return fl.Field().Float() > 0
	default:
		return false
	}
}
