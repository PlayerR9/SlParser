package parser

import (
	"strconv"
	"strings"

	gr "github.com/PlayerR9/SlParser/grammar"
	gcstr "github.com/PlayerR9/go-commons/strings"
)

// ErrUnexpectedToken is an error that is returned when an unexpected token is found.
type ErrUnexpectedToken[T gr.TokenTyper] struct {
	// Expecteds is the expected tokens.
	Expecteds []T

	// After is the token after the expected token.
	After *T

	// Got is the token that was found.
	Got *T
}

// Error implements the error interface.
//
// Message:
//
//	"expected {expected} before {after}, got {got} instead".
func (e ErrUnexpectedToken[T]) Error() string {
	var after string

	if e.After == nil {
		after = " at the end"
	} else {
		after = " before " + strconv.Quote((*e.After).String())
	}

	var got string

	if e.Got == nil {
		got = "nothing"
	} else {
		got = strconv.Quote((*e.Got).String())
	}

	var builder strings.Builder

	builder.WriteString("expected ")

	if len(e.Expecteds) == 0 {
		builder.WriteString("nothing")
	} else {
		elems := gcstr.SliceOfStringer(e.Expecteds)
		gcstr.QuoteStrings(elems)
		builder.WriteString(gcstr.EitherOr(elems))
	}

	builder.WriteString(after)
	builder.WriteString(", got ")
	builder.WriteString(got)
	builder.WriteString(" instead")

	return builder.String()
}

// NewErrUnexpectedToken creates a new ErrUnexpectedToken error.
//
// Parameters:
//   - expecteds: the expected tokens.
//   - after: the token after the expected token.
//   - got: the token that was found.
//
// Returns:
//   - *ErrUnexpectedToken[T]: the error. Never returns nil.
func NewErrUnexpectedToken[T gr.TokenTyper](expecteds []T, after, got *T) *ErrUnexpectedToken[T] {
	return &ErrUnexpectedToken[T]{
		Expecteds: expecteds,
		After:     after,
		Got:       got,
	}
}
