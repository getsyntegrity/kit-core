package fixtures

// ValidTenantID returns a fixed tenant ID for use in tests.
// Deterministic; no I/O.
func ValidTenantID() string {
	return "01HXYZ00000000000000000001"
}

// ValidUserID returns a fixed user ID for use in tests.
func ValidUserID() string {
	return "01HXYZ00000000000000000002"
}

// ValidCorrelationID returns a fixed correlation ID for tests.
func ValidCorrelationID() string {
	return "corr-01"
}

// ValidCausationID returns a fixed causation ID for tests.
func ValidCausationID() string {
	return "caus-01"
}
