package slices

// Pointer is an interface that defines a method to check if a pointer is nil.
type Pointer interface {
	// IsNil checks if the pointer is nil.
	//
	// Returns:
	//   - bool: True if the pointer is nil, false otherwise.
	IsNil() bool
}
