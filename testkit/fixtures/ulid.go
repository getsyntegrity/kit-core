// Package fixtures provides deterministic test data for specs.
package fixtures

// ValidULID returns a fixed valid ULID string for use in tests.
// Deterministic; no I/O or randomness.
func ValidULID() string {
	return "01ARZ3NDEKTSV4RRFFQ69G5FAV"
}

// AnotherValidULID returns a second fixed valid ULID for tests requiring distinct IDs.
func AnotherValidULID() string {
	return "01ARZ3NDEKTSV4RRFFQ69G5FAW"
}
