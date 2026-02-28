// Package observability defines minimal observability contracts used by the platform.
// Implementations (Prometheus, OpenTelemetry, structured loggers) live in kit-runtime.
// kit-core must not import any concrete metrics, tracing, or logging libraries.
package observability

import "context"

// Field represents a single key-value field for structured logging.
// Implementations in kit-runtime map this to their logger's field type.
type Field struct {
	Key   string
	Value interface{}
}

// Logger is the minimal logging contract for the platform.
// Concrete implementations (e.g. slog, zerolog, kit-logger) live in kit-runtime.
type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	WithContext(ctx context.Context) Logger
	With(fields ...Field) Logger
}

// MetricsRecorder is the minimal metrics recording contract.
// Concrete implementations (Prometheus, OTEL) live in kit-runtime.
type MetricsRecorder interface {
	Inc(name string, labels ...string)
	Observe(name string, value float64, labels ...string)
}
