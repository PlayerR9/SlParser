package common

// Resetter is an interface for objects that can be reset.
type Resetter interface {
	// Reset resets the object.
	//
	// Returns:
	//   - error: An error if the object could not be reset.
	//
	// Errors:
	//   - common.ErrNilReceiver: If the receiver is nil.
	//   - any other error: Implementation-specific.
	Reset() error
}
