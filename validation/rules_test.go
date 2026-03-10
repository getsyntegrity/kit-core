package validation

import (
	"testing"

	"github.com/pablogore/go-specs/specs"
	"github.com/stretchr/testify/assert"
)

func TestValidationRules(t *testing.T) {
	specs.Describe(t, "validation rules", func(s *specs.Spec) {
		s.When("NewCommonValidationRules", func(s *specs.Spec) {
			s.It("returns non-nil rules", func(ctx *specs.Context) {
				r := NewCommonValidationRules()
				assert.NotNil(ctx.T, r)
			})
		})

		s.When("CommonValidationRules_IsValidEmail", func(s *specs.Spec) {
			s.It("covers all email cases", func(ctx *specs.Context) {
				r := NewCommonValidationRules()
				tests := []struct {
					name     string
					email    string
					expected bool
				}{
					{"empty", "", false},
					{"valid simple", "a@b.co", true},
					{"valid with plus", "user+tag@example.com", true},
					{"valid with dot", "user.name@example.com", true},
					{"no at", "userexample.com", false},
					{"no domain", "user@", false},
					{"invalid tld", "user@example", false},
				}
				for _, tt := range tests {
					assert.Equal(ctx.T, tt.expected, r.IsValidEmail(tt.email), "case: %s", tt.name)
				}
			})
		})

		s.When("CommonValidationRules_IsValidULID", func(s *specs.Spec) {
			s.It("covers all ULID cases", func(ctx *specs.Context) {
				r := NewCommonValidationRules()
				validULID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
				tests := []struct {
					name     string
					ulid     string
					expected bool
				}{
					{"empty", "", false},
					{"valid", validULID, true},
					{"too short", "01ARZ3NDEKTSV4RRFFQ69G5FA", false},
					{"invalid char", "01ARZ3NDEKTSV4RRFFQ69G5FAX ", false},
				}
				for _, tt := range tests {
					assert.Equal(ctx.T, tt.expected, r.IsValidULID(tt.ulid), "case: %s", tt.name)
				}
			})
		})

		s.When("CommonValidationRules_IsValidDomain", func(s *specs.Spec) {
			s.It("covers all domain cases", func(ctx *specs.Context) {
				r := NewCommonValidationRules()
				tests := []struct {
					name     string
					domain   string
					expected bool
				}{
					{"empty", "", false},
					{"valid single", "example.com", true},
					{"valid sub", "sub.example.com", true},
					{"valid with hyphen", "my-domain.com", true},
					{"starts with hyphen", "-example.com", false},
					{"invalid char", "example!.com", false},
				}
				for _, tt := range tests {
					assert.Equal(ctx.T, tt.expected, r.IsValidDomain(tt.domain), "case: %s", tt.name)
				}
			})
		})

		s.When("CommonValidationRules_IsValidName", func(s *specs.Spec) {
			s.It("covers all name cases", func(ctx *specs.Context) {
				r := NewCommonValidationRules()
				tests := []struct {
					name     string
					input    string
					expected bool
				}{
					{"empty", "", false},
					{"letters and spaces", "John Doe", true},
					{"with hyphen", "Mary-Jane", true},
					{"with apostrophe", "O'Brien", true},
					{"numbers", "John2", false},
					{"special", "John@Doe", false},
				}
				for _, tt := range tests {
					assert.Equal(ctx.T, tt.expected, r.IsValidName(tt.input), "case: %s", tt.name)
				}
			})
		})

		s.When("CommonValidationRules_IsValidUsername", func(s *specs.Spec) {
			s.It("covers all username cases", func(ctx *specs.Context) {
				r := NewCommonValidationRules()
				tests := []struct {
					name     string
					input    string
					expected bool
				}{
					{"empty", "", false},
					{"too short", "ab", false},
					{"valid min", "abc", true},
					{"valid max", "a1234567890123456789", true},
					{"too long", "a123456789012345678901", false},
					{"with underscore", "user_name", true},
					{"with hyphen", "user-name", true},
					{"invalid char", "user@name", false},
				}
				for _, tt := range tests {
					assert.Equal(ctx.T, tt.expected, r.IsValidUsername(tt.input), "case: %s", tt.name)
				}
			})
		})

		s.When("CommonValidationRules_IsValidPostalCode", func(s *specs.Spec) {
			s.It("covers all postal code cases", func(ctx *specs.Context) {
				r := NewCommonValidationRules()
				tests := []struct {
					name     string
					input    string
					expected bool
				}{
					{"empty", "", false},
					{"US zip", "12345", true},
					{"UK style", "SW1A 1AA", true},
					{"with hyphen", "12345-6789", true},
					{"lowercase", "sw1a1aa", true},
					{"too short", "12", false},
					{"too long", "12345678901", false},
				}
				for _, tt := range tests {
					assert.Equal(ctx.T, tt.expected, r.IsValidPostalCode(tt.input), "case: %s", tt.name)
				}
			})
		})
	})
}
