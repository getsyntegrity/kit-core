// Package ref provides reference utilities for Go values.
package ref

// Of get reference of any value, useful in cases when const is used: ref.Of("mystring"), ref.Of(int32(123)).
func Of[T any](val T) *T {
	return &val
}

// SafeDeref safely dereference value for a pointer.
func SafeDeref[T any](val *T) T {
	if val != nil {
		return *val
	}

	var def T
	return def
}

// NilIfNil returns nil if the given pointer is nil. Otherwise, it returns the value pointed to by the pointer.
func NilIfNil[T any](ptr *T) any {
	if ptr == nil {
		return nil
	}
	return *ptr
}

// DerefSlice dereferences a slice of pointers and returns a slice of values.
func DerefSlice[T any](slice []*T) []T {
	var result []T
	for _, item := range slice {
		if item != nil {
			result = append(result, *item)
		}
	}
	return result
}
