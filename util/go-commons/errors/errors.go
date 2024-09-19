package util

// TODO: Remove this once go-commons is updated

import (
	"fmt"
	"strconv"
	"strings"
)

// ErrValue represents an error when a value is not expected.
type ErrValue struct {
	// Kind is the name of the thing that was expected.
	Kind string

	// Expected is the value that was expected.
	Expected any

	// Got is the value that was received.
	Got any

	// ShouldQuote is true if the expected and got values should be quoted,
	// false otherwise.
	ShouldQuote bool
}

// Error implements the error interface.
//
// Message:
//
//	"expected <kind> to be <expected>, got <got> instead"
func (e ErrValue) Error() string {
	var builder strings.Builder

	builder.WriteString("expected ")

	if e.Kind != "" {
		builder.WriteString(e.Kind)
		builder.WriteString(" to be ")
	}

	if e.Expected == nil {
		builder.WriteString("nothing")
	} else if e.ShouldQuote {
		fmt.Fprintf(&builder, "%q", e.Expected)
	} else {
		fmt.Fprintf(&builder, "%s", e.Expected)
	}

	builder.WriteString(", got ")

	if e.Got == nil {
		builder.WriteString("nothing")
	} else if e.ShouldQuote {
		fmt.Fprintf(&builder, "%q", e.Got)
	} else {
		fmt.Fprintf(&builder, "%s", e.Got)
	}

	builder.WriteString(" instead")

	return builder.String()
}

// NewErrValue creates a new ErrValue error.
//
// Parameters:
//   - kind: The name of the thing that was expected.
//   - expected: The value that was expected.
//   - got: The value that was received.
//   - should_quote: True if the expected and got values should be quoted,
//     false otherwise.
//
// Returns:
//   - *ErrValue: A pointer to the newly created ErrValue. Never returns nil.
func NewErrValue(kind string, expected, got any, should_quote bool) *ErrValue {
	return &ErrValue{
		Kind:        kind,
		Expected:    expected,
		Got:         got,
		ShouldQuote: should_quote,
	}
}

// ErrValues represents an error when multiple value are not expected.
type ErrValues[T any] struct {
	// Kind is the name of the thing that was expected.
	Kind string

	// Expecteds is the values that were expected.
	Expecteds []T

	// Got is the value that was received.
	Got any

	// ShouldQuote is true if the expected and got values should be quoted,
	// false otherwise.
	ShouldQuote bool
}

// Error implements the error interface.
//
// Message:
//
//	"expected <kind> to be <expected>, got <got> instead"
func (e ErrValues[T]) Error() string {
	var builder strings.Builder

	builder.WriteString("expected ")

	if e.Kind != "" {
		builder.WriteString(e.Kind)
		builder.WriteString(" to be ")
	}

	switch len(e.Expecteds) {
	case 0:
		builder.WriteString("nothing")
	case 1:
		if e.ShouldQuote {
			builder.WriteString(strconv.Quote(fmt.Sprintf("%v", e.Expecteds[0])))
		} else {
			fmt.Fprintf(&builder, "%v", e.Expecteds[0])
		}
	default:
		elems := make([]string, 0, len(e.Expecteds))

		if e.ShouldQuote {
			for i := 0; i < len(e.Expecteds); i++ {
				elems = append(elems, strconv.Quote(fmt.Sprintf("%v", e.Expecteds[i])))
			}
		} else {
			for i := 0; i < len(e.Expecteds); i++ {
				elems = append(elems, fmt.Sprintf("%v", e.Expecteds[i]))
			}
		}

		builder.WriteString("either ")
		builder.WriteString(strings.Join(elems[:len(elems)-1], ", "))
		builder.WriteString(" or ")
		builder.WriteString(elems[len(elems)-1])
	}

	builder.WriteString(", got ")

	if e.Got == nil {
		builder.WriteString("nothing")
	} else if e.ShouldQuote {
		fmt.Fprintf(&builder, "%q", e.Got)
	} else {
		fmt.Fprintf(&builder, "%s", e.Got)
	}

	builder.WriteString(" instead")

	return builder.String()
}

// NewErrValues creates a new ErrValues error.
//
// Parameters:
//   - kind: The name of the thing that was expected.
//   - expected: The values that were expected.
//   - got: The value that was received.
//   - should_quote: True if the expected and got values should be quoted,
//     false otherwise.
//
// Returns:
//   - *ErrValue: A pointer to the newly created ErrValue. Never returns nil.
func NewErrValues[T any](kind string, expected []T, got any, should_quote bool) *ErrValue {
	return &ErrValue{
		Kind:        kind,
		Expected:    expected,
		Got:         got,
		ShouldQuote: should_quote,
	}
}
