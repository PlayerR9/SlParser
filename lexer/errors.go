package lexer

import (
	"strconv"
	"strings"

	gcstr "github.com/PlayerR9/go-commons/strings"
)

// ErrUnexpectedChar is an error that occurs when an unexpected character is encountered.
type ErrUnexpectedChar struct {
	// Expected is the expected character.
	Expecteds []rune

	// Previous is the previous character.
	Previous rune

	// Got is the current character.
	Got *rune
}

// Error implements the error interface.
//
// Message:
//
//	"expected {expected} after {previous}, got {got} instead".
func (e ErrUnexpectedChar) Error() string {
	var got string

	if e.Got == nil {
		got = "nothing"
	} else {
		got = strconv.QuoteRune(*e.Got)
	}

	var builder strings.Builder

	builder.WriteString("expected ")

	if len(e.Expecteds) == 0 {
		builder.WriteString("nothing")
	} else {
		elems := gcstr.SliceOfRunes(e.Expecteds)
		gcstr.QuoteStrings(elems)

		builder.WriteString(gcstr.EitherOrString(elems))
	}

	builder.WriteString(" after ")
	builder.WriteString(strconv.QuoteRune(e.Previous))
	builder.WriteString(", got ")
	builder.WriteString(got)
	builder.WriteString(" instead")

	return builder.String()
}

// NewErrUnexpectedChar creates a new ErrUnexpectedChar error.
//
// Parameters:
//   - previous: the previous character.
//   - expecteds: the expected characters.
//   - got: the current character.
//
// Returns:
//   - *ErrUnexpectedChar: the error. Never returns nil.
func NewErrUnexpectedChar(previous rune, expecteds []rune, got *rune) *ErrUnexpectedChar {
	return &ErrUnexpectedChar{
		Expecteds: expecteds,
		Previous:  previous,
		Got:       got,
	}
}
