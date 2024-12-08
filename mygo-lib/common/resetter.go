package common

// Resetter is an interface for types that can reset their internal state.
type Resetter interface {
	// Reset 'resets' the type's internal state in order to make it ready
	// for reuse.
	//
	// Returns:
	//   - error: An error if the type could not be reset.
	//
	// Errors:
	//   - ErrNilReceiver: If the receiver is nil.
	//   - any other error: Implementation-dependent error.
	Reset() error
}
