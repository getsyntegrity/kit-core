/*
Package validation provides a robust data validation framework with custom rules, error handling, and comprehensive validation strategies.

This package offers production-ready validation capabilities including struct validation, field validation, custom validators,
conditional validation, and comprehensive error reporting for microservices applications.

## Features

### Struct Validation
Comprehensive struct validation with tag-based rules:

Example:

	type User struct {
		ID       string `validate:"required,ulid"`
		Name     string `validate:"required,min=2,max=50"`
		Email    string `validate:"required,email"`
		Age      int    `validate:"required,min=18,max=120"`
		Password string `validate:"required,min=8,contains=uppercase,contains=lowercase,contains=number"`
	}

	validator, err := validation.NewValidator()
	if err != nil {
		return err
	}

	user := User{
		ID:       "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Age:      30,
		Password: "SecurePass123",
	}

	if err := validator.Validate(user); err != nil {
		log.Printf("Validation failed: %v", err)
		return
	}

### Field Validation
Individual field validation with custom rules:

Example:

	validator, err := validation.NewValidator()
	if err != nil {
		return err
	}

	// Validate email
	email := "user@example.com"
	if err := validator.ValidateField(email, "email"); err != nil {
		log.Printf("Invalid email: %v", err)
		return
	}

	// Validate ULID
	id := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	if err := validator.ValidateField(id, "ulid"); err != nil {
		log.Printf("Invalid ULID: %v", err)
		return
	}

### Custom Validators
Extensible validation system with custom rules:

Example:

	// Register custom validator
	validator, err := validation.NewValidator()
	if err != nil {
		return err
	}
	validator.RegisterValidator("phone", func(value interface{}) error {
		phone, ok := value.(string)
		if !ok {
			return validation.NewError("phone", "invalid type")
		}

		// Phone validation logic
		if !regexp.MustCompile(`^\+?[1-9]\d{1,14}$`).MatchString(phone) {
			return validation.NewError("phone", "invalid phone format")
		}
		return nil
	})

	// Use custom validator
	type Contact struct {
		Phone string `validate:"required,phone"`
	}

### Conditional Validation
Conditional validation based on field values:

Example:

	type User struct {
		Type      string `validate:"required,oneof=individual business"`
		SSN       string `validate:"required_if=Type individual"`
		TaxID     string `validate:"required_if=Type business"`
		Company   string `validate:"required_if=Type business"`
	}

	validator, err := validation.NewValidator()
	if err != nil {
		return err
	}

	// Individual user
	individual := User{
		Type: "individual",
		SSN:  "123-45-6789",
	}
	if err := validator.Validate(individual); err != nil {
		log.Printf("Individual validation failed: %v", err)
		return
	}

	// Business user
	business := User{
		Type:    "business",
		TaxID:   "12-3456789",
		Company: "Acme Corp",
	}
	if err := validator.Validate(business); err != nil {
		log.Printf("Business validation failed: %v", err)
		return
	}

### Array Validation
Comprehensive array and slice validation:

Example:

	type Order struct {
		Items []OrderItem `validate:"required,min=1,dive"`
	}

	type OrderItem struct {
		ProductID string  `validate:"required,ulid"`
		Quantity  int     `validate:"required,min=1"`
		Price     float64 `validate:"required,min=0"`
	}

	order := Order{
		Items: []OrderItem{
			{ProductID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Quantity: 2, Price: 29.99},
			{ProductID: "01ARZ3NDEKTSV4RRFFQ69G5FBW", Quantity: 1, Price: 19.99},
		},
	}

	if err := validator.Validate(order); err != nil {
		log.Printf("Order validation failed: %v", err)
		return
	}

### Error Handling
Comprehensive error handling with detailed error information:

Example:

	if err := validator.Validate(user); err != nil {
		if validationErr, ok := err.(*validation.ValidationError); ok {
			for _, fieldErr := range validationErr.FieldErrors {
				log.Printf("Field: %s, Error: %s, Value: %v",
					fieldErr.Field, fieldErr.Message, fieldErr.Value)
			}
		}
		return
	}

### Common Validation Rules
Built-in validation rules for common use cases:

Example:

	rules := validation.NewCommonValidationRules()

	// Email validation
	isValidEmail := rules.IsValidEmail("user@example.com")

	// URL validation
	isValidURL := rules.IsValidURL("https://example.com")

	// ULID validation
	isValidULID := rules.IsValidULID("01ARZ3NDEKTSV4RRFFQ69G5FAV")

	// Phone validation
	isValidPhone := rules.IsValidPhone("+1234567890")

### Configuration
Flexible configuration via environment variables or programmatically:

Environment Variables:

	VALIDATION_STRICT_MODE=true
	VALIDATION_CUSTOM_RULES_ENABLED=true
	VALIDATION_ERROR_FORMAT=json
	VALIDATION_MAX_ERRORS=10

### Performance
Optimized for high-performance applications with efficient validation algorithms and caching.

### Thread Safety
All validation components are thread-safe and can be used concurrently.
*/
package validation
