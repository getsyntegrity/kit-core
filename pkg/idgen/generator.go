// Package idgen provides ID generation utilities and implementations
// for creating unique identifiers across different systems.
// All IDs are generated as ULIDs (Universally Unique Lexicographically Sortable Identifiers).
package idgen

import (
	"crypto/rand"
	"io"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	mtx    sync.Mutex
	reader io.Reader
)

func init() {
	reader = ulid.Monotonic(rand.Reader, 0)
}

func newULID() ulid.ULID {
	mtx.Lock()
	defer mtx.Unlock()
	return ulid.MustNew(ulid.Timestamp(time.Now()), reader)
}

// NewULID returns new ULID string.
func NewULID() string {
	return newULID().String()
}

// MustNewULID is an alias for NewULID for consistency with other Must* functions.
func MustNewULID() string {
	return NewULID()
}

// ParseULID parses a ULID string and returns the ULID value.
func ParseULID(s string) (ulid.ULID, error) {
	return ulid.Parse(s)
}

// MustParseULID parses a ULID string and returns the ULID value. Panics if invalid.
func MustParseULID(s string) ulid.ULID {
	return ulid.MustParse(s)
}

// IsValidULID checks if a string is a valid ULID format.
func IsValidULID(s string) bool {
	_, err := ulid.Parse(s)
	return err == nil
}
