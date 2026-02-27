package pii

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedact(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{src: "test123", expected: "t***3"},
		{src: "test123@email.com", expected: "t***3@email.com"},
		{src: "test123@", expected: "t***3@"},
		{src: "", expected: ""},
		{src: "  ", expected: " *** "},
		{src: "@", expected: "@***@"},
		{src: "a", expected: "a***a"},
		{src: "ab", expected: "a***b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Redact(tt.src)
			assert.Equal(t, tt.expected, res)
		})
	}
}

func TestRedactAll(t *testing.T) {
	res := RedactAll([]string{"test123", "test123@email.com"})
	assert.Equal(t, []string{"t***3", "t***3@email.com"}, res)
}

func TestSensitiveString_LogValue(t *testing.T) {
	tests := []struct {
		name     string
		input    SensitiveString
		expected string
	}{
		{
			name:     "normal string",
			input:    SensitiveString("test123"),
			expected: "t***3",
		},
		{
			name:     "email string",
			input:    SensitiveString("test123@email.com"),
			expected: "t***3@email.com",
		},
		{
			name:     "empty string",
			input:    SensitiveString(""),
			expected: "",
		},
		{
			name:     "single character",
			input:    SensitiveString("a"),
			expected: "a***a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logValue := tt.input.LogValue()
			assert.Equal(t, tt.expected, logValue.String())
		})
	}
}

func BenchmarkRedact(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		res := Redact("tes123t@test.com")
		if res != "t***t@test.com" {
			b.Fail()
		}
	}
}
