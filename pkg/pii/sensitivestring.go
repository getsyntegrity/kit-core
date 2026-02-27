package pii

import (
	"log/slog"
)

// SensitiveString is a string that should be redacted when logged, but not when marshaled to JSON.
type SensitiveString string

// LogValue returns a safe representation of the sensitive string for logging.
func (p SensitiveString) LogValue() slog.Value {
	return slog.StringValue(Redact(string(p)))
}
