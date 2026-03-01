package validation

import (
	"regexp"
	"strings"

	"github.com/getsyntegrity/kit-core/pkg/idgen"
)

// CommonValidationRules provides a collection of common validation rules.
type CommonValidationRules struct{}

// NewCommonValidationRules creates a new instance of common validation rules.
func NewCommonValidationRules() *CommonValidationRules {
	return &CommonValidationRules{}
}

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	ulidRegex     = regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{26}$`)
)

// IsValidEmail validates email format.
func (c *CommonValidationRules) IsValidEmail(email string) bool {
	if email == "" {
		return false
	}
	return emailRegex.MatchString(email)
}

// IsValidULID validates ULID format.
func (c *CommonValidationRules) IsValidULID(ulid string) bool {
	if ulid == "" {
		return false
	}
	return idgen.IsValidULID(ulid)
}

// IsValidDomain validates domain name format.
func (c *CommonValidationRules) IsValidDomain(domain string) bool {
	if domain == "" {
		return false
	}
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

// IsValidPostalCode validates postal code format (basic).
func (c *CommonValidationRules) IsValidPostalCode(postalCode string) bool {
	if postalCode == "" {
		return false
	}
	postalCodeRegex := regexp.MustCompile(`^[A-Z0-9\s-]{3,10}$`)
	return postalCodeRegex.MatchString(strings.ToUpper(postalCode))
}
