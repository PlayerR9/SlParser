package PageInterval

// ErrNoPagesInInterval represents an error where there are no pages in the interval.
type ErrNoPagesInInterval struct{}

// Error is a method of the error interface.
//
// Returns:
//
//   - string: The error message.
func (e *ErrNoPagesInInterval) Error() string {
	return "no pages in interval"
}

// NewErrNoPagesInInterval creates a new instance of ErrNoPagesInInterval.
//
// Returns:
//
//   - error: An error of type *ErrNoPagesInInterval.
func NewErrNoPagesInInterval() *ErrNoPagesInInterval {
	return &ErrNoPagesInInterval{}
}
