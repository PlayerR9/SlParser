package lexer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrInputStreamExhausted is an error that occurs when the input stream
// is exhausted.
type ErrInputStreamExhausted struct{}

// Error implements the error interface.
//
// Format:
//
//	"input stream is exhausted"
func (l *ErrInputStreamExhausted) Error() string {
	return "input stream is exhausted"
}

// NewErrInputStreamExhausted creates a new error.
//
// Returns:
//   - *ErrInputStreamExhausted: The new error. Never returns nil.
func NewErrInputStreamExhausted() *ErrInputStreamExhausted {
	return &ErrInputStreamExhausted{}
}

// IsExhausted checks if the error is an ErrInputStreamExhausted.
//
// Returns:
//   - bool: True if the error is an ErrInputStreamExhausted, false otherwise.
func IsExhausted(err error) bool {
	if err == nil {
		return false
	}

	var exhausted_err *ErrInputStreamExhausted

	return errors.As(err, &exhausted_err)
}

// ErrInvalidUTF8Encoding is an error that occurs when an invalid UTF-8
// encoding is encountered.
type ErrInvalidUTF8Encoding struct {
	// At is the position of the invalid UTF-8 encoding.
	At int
}

// Error implements the error interface.
//
// Format:
//
//	"invalid UTF-8 encoding at <position>"
func (l *ErrInvalidUTF8Encoding) Error() string {
	return fmt.Sprintf("invalid UTF-8 encoding at %d", l.At)
}

// NewErrInvalidUTF8Encoding creates a new error.
//
// Parameters:
//   - at: The position of the invalid UTF-8 encoding.
//
// Returns:
//   - *ErrInvalidUTF8Encoding: The new error. Never returns nil.
func NewErrInvalidUTF8Encoding(at int) *ErrInvalidUTF8Encoding {
	return &ErrInvalidUTF8Encoding{
		At: at,
	}
}

// ErrUnexpectedRune is an error that occurs when an unexpected rune is
// encountered.
type ErrUnexpectedRune struct {
	// Expected is the expected rune.
	Expected rune

	// Prev is the rune that was encountered before the expected rune.
	Prev *rune

	// Got is the rune that was encountered.
	Got *rune
}

// Error implements the error interface.
//
// Format:
//
//	"expected <expected> <prev>, got <actual> instead"
func (e *ErrUnexpectedRune) Error() string {
	var prev string

	if e.Prev == nil {
		prev = "at the beginning"
	} else {
		prev = "after " + strconv.QuoteRune(*e.Prev)
	}

	var got string

	if e.Got == nil {
		got = "nothing"
	} else {
		got = strconv.QuoteRune(*e.Got)
	}

	var builder strings.Builder

	builder.WriteString("expected ")
	builder.WriteString(strconv.QuoteRune(e.Expected))
	builder.WriteString(" ")
	builder.WriteString(prev)
	builder.WriteString(", got ")
	builder.WriteString(got)
	builder.WriteString(" instead")

	return builder.String()
}

// NewErrUnexpectedRune creates a new error.
//
// Parameters:
//   - expected: The expected rune.
//   - prev: The rune that was encountered before the expected rune.
//   - got: The rune that was encountered.
//
// Returns:
//   - *ErrUnexpectedRune: The new error. Never returns nil.
func NewErrUnexpectedRune(expected rune, prev *rune, got *rune) *ErrUnexpectedRune {
	return &ErrUnexpectedRune{
		Expected: expected,
		Prev:     prev,
		Got:      got,
	}
}
