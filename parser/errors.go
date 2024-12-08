package parser

import (
	"errors"
	"strings"
)

/////////////////////////////////////////////////////////

// NewUnsupportedValue creates a new errors.Err error with the code InvalidOperation and the
// message "value of <expected> is not a supported <kind> type".
//
// Parameters:
//   - kind: The kind of the value. Ignored if not provided.
//   - expected: The expected value. Ignored if not provided.
//
// Returns:
//   - *errors.Err[ErrorCode]: The new error. Never returns nil.
func NewUnsupportedValue(kind, expected string) error {
	var builder strings.Builder

	builder.WriteString("value ")

	if expected != "" {
		builder.WriteString("of ")
		builder.WriteString(expected)
	}

	builder.WriteString(" is not  ")

	if kind != "" {
		builder.WriteString("a supported ")
		builder.WriteString(kind)
		builder.WriteString(" type")
	} else {
		builder.WriteString("supported")
	}

	return errors.New(builder.String())
}
