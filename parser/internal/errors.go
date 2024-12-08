package internal

import (
	"errors"
	"slices"
	"strconv"
	"strings"
)

// NewUnsupportedValue creates an error indicating an unsupported value type.
//
// Parameters:
//   - kind: The kind of the value. If provided, it specifies the unsupported type.
//   - expected: The expected value. If provided, it specifies the value that is not supported.
//
// Returns:
//   - error: The error with a message indicating the unsupported value type.
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

// ErrNotAsExpected occurs when a string is not as expected.
type ErrNotAsExpected struct {
	// Quote if true, the strings will be quoted before being printed.
	Quote bool

	// Kind is the kind of the string that is not as expected.
	Kind string

	// Expecteds are the strings that were expecteds.
	Expecteds []string

	// Got is the actual string.
	Got string
}

// Error implements the error interface.
func (e ErrNotAsExpected) Error() string {
	var kind string

	if e.Kind != "" {
		kind = e.Kind + " to be "
	}

	var got string

	if e.Got == "" {
		got = "nothing"
	} else if e.Quote {
		got = strconv.Quote(e.Got)
	} else {
		got = e.Got
	}

	var builder strings.Builder

	builder.WriteString("expected ")
	builder.WriteString(kind)

	if len(e.Expecteds) > 0 {
		var elems []string

		if !e.Quote {
			elems = e.Expecteds
		} else {
			elems = make([]string, 0, len(e.Expecteds))

			for _, elem := range e.Expecteds {
				str := strconv.Quote(elem)
				elems = append(elems, str)
			}
		}

		builder.WriteString(EitherOrString(elems))
	} else {
		builder.WriteString("something")
	}

	builder.WriteString(", got ")
	builder.WriteString(got)

	return builder.String()
}

// NewErrNotAsExpected creates a new ErrNotAsExpected error.
//
// Parameters:
//   - quote: Whether or not to quote the strings in the error message.
//   - kind: The kind of thing that was not as expected. This is used in the error message.
//   - got: The actual value. If empty, "nothing" is used in the error message.
//   - expecteds: The expected values. If empty, "something" is used in the error message.
//
// Returns:
//   - error: The new error. Never returns nil.
//
// Format:
//
//	"expected <kind> to be <expected>, got <got>"
//
// Where:
//   - <kind>: The kind of thing that was not as expected. This is used in the error message.
//   - <expected>: The expected values. This is used in the error message.
//   - <got>: The actual value. This is used in the error message. If nil, "nothing" is used instead.
//
// Duplicate values are automatically removed and the list of expected values is sorted in ascending order.
func NewErrNotAsExpected(quote bool, kind string, got string, expecteds ...string) error {
	unique := make([]string, 0, len(expecteds))

	for _, expected := range expecteds {
		pos, ok := slices.BinarySearch(unique, expected)
		if ok {
			continue
		}

		unique = slices.Insert(unique, pos, expected)
	}

	unique = unique[:len(unique):len(unique)]

	return &ErrNotAsExpected{
		Quote:     quote,
		Kind:      kind,
		Expecteds: unique,
		Got:       got,
	}
}
