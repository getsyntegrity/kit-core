// Package setutils provides utility data structures and functions for the go-kit framework.
// This package contains common utilities like sets, optional values, and constants.
package setutils

// Set is a utility data structure that holds a unique collection of strings.
type Set struct {
	items map[string]struct{}
}

// NewSet creates a new Set instance.
func NewSet() *Set {
	return &Set{
		items: make(map[string]struct{}),
	}
}

// Add inserts an item into the set.
func (s *Set) Add(item string) {
	s.items[item] = struct{}{}
}

// Contains checks if an item is present in the set.
func (s *Set) Contains(item string) bool {
	_, exists := s.items[item]
	return exists
}

// Items returns all the items in the set.
func (s *Set) Items() []string {
	keys := make([]string, 0, len(s.items))
	for key := range s.items {
		keys = append(keys, key)
	}
	return keys
}
