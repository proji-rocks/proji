// Package pointer provides functions to create pointers to basic types. This is mostly for convenience.
package pointer

// To returns a pointer to the given value. Purely for convenience.
func To[T any](v T) *T { return &v }
