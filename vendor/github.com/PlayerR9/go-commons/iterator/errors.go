package iterator

import (
	"strings"

	gcers "github.com/PlayerR9/go-commons/errors"
)

// ErrIteration is the error that is returned when an iteration error occurs.
type ErrIteration struct {
	// Idx is the index of the element where the error occurred.
	Idx int

	// Reason is the error that occurred.
	Reason error
}

// Error implements the error interface.
//
// Message: "iteration error"
func (e *ErrIteration) Error() string {
	var builder strings.Builder

	builder.WriteString("error while iterating over the ")
	builder.WriteString(gcers.GetOrdinalSuffix(e.Idx + 1))
	builder.WriteString(" element")

	if e.Reason != nil {
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap implements the error interface.
func (e *ErrIteration) Unwrap() error {
	return e.Reason
}

// NewErrIteration creates a new error that occurs when an iteration error occurs.
//
// Parameters:
//   - idx: The index of the element where the error occurred.
//   - reason: The error that occurred.
//
// Returns:
//   - error: An error if an iteration error occurs. Never returns nil.
func NewErrIteration(idx int, reason error) *ErrIteration {
	return &ErrIteration{
		Idx:    idx,
		Reason: reason,
	}
}

// ChangeReason changes the reason of the error.
//
// Parameters:
//   - reason: The reason for the error.
func (e *ErrIteration) ChangeReason(reason error) {
	e.Reason = reason
}
