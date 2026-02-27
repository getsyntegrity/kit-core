package validation

import (
	"regexp"
	"strings"
)

// Common validation rules that can be reused across services

// CommonValidationRules provides a collection of common validation rules.
type CommonValidationRules struct{}

// NewCommonValidationRules creates a new instance of common validation rules.
func NewCommonValidationRules() *CommonValidationRules {
	return &CommonValidationRules{}
}

// Email validation regex.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// URL validation regex.
var urlRegex = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)

// Phone validation regex (basic international format).
var phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

// ULID validation regex.
// ULID format: 26 characters, base32 encoded (0-9, A-Z excluding I, L, O, U)
var ulidRegex = regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{26}$`)

// Common validation functions

// IsValidEmail validates email format.
func (c *CommonValidationRules) IsValidEmail(email string) bool {
	if email == "" {
		return false
	}
	return emailRegex.MatchString(email)
}

// IsValidURL validates URL format.
func (c *CommonValidationRules) IsValidURL(url string) bool {
	if url == "" {
		return false
	}
	return urlRegex.MatchString(url)
}

// IsValidPhone validates phone number format.
func (c *CommonValidationRules) IsValidPhone(phone string) bool {
	if phone == "" {
		return false
	}
	return phoneRegex.MatchString(phone)
}

// IsValidULID validates ULID format.
// ULID (Universally Unique Lexicographically Sortable Identifier) is a 26-character
// base32-encoded identifier that is sortable and contains timestamp information.
func (c *CommonValidationRules) IsValidULID(ulid string) bool {
	if ulid == "" {
		return false
	}
	return ulidRegex.MatchString(ulid)
}

// IsValidDomain validates domain name format.
func (c *CommonValidationRules) IsValidDomain(domain string) bool {
	if domain == "" {
		return false
	}

	// Basic domain validation regex
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	return domainRegex.MatchString(domain)
}

// IsValidName validates name format (letters, spaces, hyphens, apostrophes).
func (c *CommonValidationRules) IsValidName(name string) bool {
	if name == "" {
		return false
	}

	nameRegex := regexp.MustCompile(`^[a-zA-ZÀ-ÿ\s\-']+$`)
	return nameRegex.MatchString(name)
}

// IsValidUsername validates username format (alphanumeric, underscores, hyphens).
func (c *CommonValidationRules) IsValidUsername(username string) bool {
	if username == "" {
		return false
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`)
	return usernameRegex.MatchString(username)
}

// IsValidPassword validates password strength.
func (c *CommonValidationRules) IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	// At least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// At least one lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// At least one digit
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	// At least one special character
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	return hasUpper && hasLower && hasDigit && hasSpecial
}

// IsValidIPAddress validates IP address format.
func (c *CommonValidationRules) IsValidIPAddress(ipAddress string) bool {
	if ipAddress == "" {
		return false
	}

	// IPv4 regex
	ipv4Regex := regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)
	// IPv6 regex (simplified)
	ipv6Regex := regexp.MustCompile(`^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`)

	return ipv4Regex.MatchString(ipAddress) || ipv6Regex.MatchString(ipAddress)
}

// IsValidDate validates date format (YYYY-MM-DD).
func (c *CommonValidationRules) IsValidDate(date string) bool {
	if date == "" {
		return false
	}

	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return dateRegex.MatchString(date)
}

// IsValidTime validates time format (HH:MM:SS).
func (c *CommonValidationRules) IsValidTime(time string) bool {
	if time == "" {
		return false
	}

	timeRegex := regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`)
	return timeRegex.MatchString(time)
}

// IsValidDateTime validates datetime format (YYYY-MM-DD HH:MM:SS).
func (c *CommonValidationRules) IsValidDateTime(datetime string) bool {
	if datetime == "" {
		return false
	}

	datetimeRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2} ([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`)
	return datetimeRegex.MatchString(datetime)
}

// IsValidCurrency validates currency code format (ISO 4217).
func (c *CommonValidationRules) IsValidCurrency(currency string) bool {
	if currency == "" {
		return false
	}

	currencyRegex := regexp.MustCompile(`^[A-Z]{3}$`)
	return currencyRegex.MatchString(currency)
}

// IsValidLanguage validates language code format (ISO 639-1).
func (c *CommonValidationRules) IsValidLanguage(language string) bool {
	if language == "" {
		return false
	}

	languageRegex := regexp.MustCompile(`^[a-z]{2}$`)
	return languageRegex.MatchString(language)
}

// IsValidCountry validates country code format (ISO 3166-1 alpha-2).
func (c *CommonValidationRules) IsValidCountry(country string) bool {
	if country == "" {
		return false
	}

	countryRegex := regexp.MustCompile(`^[A-Z]{2}$`)
	return countryRegex.MatchString(country)
}

// IsValidPostalCode validates postal code format (basic).
func (c *CommonValidationRules) IsValidPostalCode(postalCode string) bool {
	if postalCode == "" {
		return false
	}

	// Basic postal code validation (alphanumeric, 3-10 characters)
	postalCodeRegex := regexp.MustCompile(`^[A-Z0-9\s-]{3,10}$`)
	return postalCodeRegex.MatchString(strings.ToUpper(postalCode))
}

// IsValidCreditCard validates credit card number format (Luhn algorithm).
func (c *CommonValidationRules) IsValidCreditCard(cardNumber string) bool {
	if cardNumber == "" {
		return false
	}

	// Remove spaces and dashes
	cardNumber = regexp.MustCompile(`[\s-]`).ReplaceAllString(cardNumber, "")

	// Check if it's numeric and has reasonable length
	if !regexp.MustCompile(`^\d{13,19}$`).MatchString(cardNumber) {
		return false
	}

	// Luhn algorithm
	sum := 0
	alternate := false

	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit := int(cardNumber[i] - '0')

		if alternate {
			digit *= 2
			if digit > 9 {
				digit = (digit % 10) + 1
			}
		}

		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}
