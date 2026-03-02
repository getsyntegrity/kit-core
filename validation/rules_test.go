package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommonValidationRules(t *testing.T) {
	r := NewCommonValidationRules()
	assert.NotNil(t, r)
}

func TestCommonValidationRules_IsValidEmail(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, r.IsValidEmail(tt.email))
		})
	}
}

func TestCommonValidationRules_IsValidULID(t *testing.T) {
	r := NewCommonValidationRules()
	// Valid ULID: 26 chars Crockford base32
	validULID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	tests := []struct {
		name     string
		ulid     string
		expected bool
	}{
		{"empty", "", false},
		{"valid", validULID, true},
		{"too short", "01ARZ3NDEKTSV4RRFFQ69G5FA", false},
		{"invalid char", "01ARZ3NDEKTSV4RRFFQ69G5FAX ", false}, // space invalid
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, r.IsValidULID(tt.ulid))
		})
	}
}

func TestCommonValidationRules_IsValidDomain(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, r.IsValidDomain(tt.domain))
		})
	}
}

func TestCommonValidationRules_IsValidName(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, r.IsValidName(tt.input))
		})
	}
}

func TestCommonValidationRules_IsValidUsername(t *testing.T) {
	r := NewCommonValidationRules()
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty", "", false},
		{"too short", "ab", false},
		{"valid min", "abc", true},
		{"valid max", "a1234567890123456789", true}, // 20 chars
		{"too long", "a123456789012345678901", false},
		{"with underscore", "user_name", true},
		{"with hyphen", "user-name", true},
		{"invalid char", "user@name", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, r.IsValidUsername(tt.input))
		})
	}
}

func TestCommonValidationRules_IsValidPostalCode(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, r.IsValidPostalCode(tt.input))
		})
	}
}
