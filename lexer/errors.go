package lexer

import (
	"fmt"

	gcerr "github.com/PlayerR9/go-errors/error"
)

//go:generate stringer -type=ErrorCode

type ErrorCode int

const (
	// UnrecognizedChar occurs when an unrecognized character is encountered.
	//
	// Example:
	// 	let L = lexer with integers as its lexing table.
	// 	Lex(L, "a")
	UnrecognizedChar ErrorCode = iota

	// InvalidInputStream occurs when the input stream is invalid.
	//
	// Example:
	// 	let is = input stream of non-utf8 characters.
	InvalidInputStream

	// BadWord occurs when a word is invalid.
	//
	// Example:
	// 	let L = lexer with integers as its lexing table.
	// 	Lex(L, "01")
	BadWord

	// BadGroup occurs when a group that was expected was not found or it was
	// found but was invalid.
	//
	// Example:
	// 	let L = lexer with integers as its lexing table.
	// 	Lex(L, "a")
	BadGroup
)

func (e ErrorCode) Int() int {
	return int(e)
}

// NewErrInvalidInputStream returns a new error when the lexer encounters an invalid input stream.
//
// This error has code InvalidInputStream and the reason as its message. Yet, if it is nil,
// "something went wrong" is used.
//
// Parameters:
//   - reason: The reason for the error.
//
// Returns:
//   - *gcerr.Err: The error. Never returns nil.
func NewErrInvalidInputStream(reason error) *gcerr.Err {
	var msg string

	if reason == nil {
		msg = "something went wrong"
	} else {
		msg = reason.Error()
	}

	err := gcerr.New(InvalidInputStream, msg)
	return err
}

// NewErrGotNothing returns a new error when the lexer encounters nothing after a character
// when, in reality, it should have lexed something.
//
// This error has code BadWord and message "expected <expected> after <prev>, got nothing instead".
//
// Parameters:
//   - prev: The previous character.
//   - expected: The expected character.
//
// Returns:
//   - *gcerr.Err: The error. Never returns nil.
func NewErrGotNothing(prev, expected rune) *gcerr.Err {
	msg := fmt.Sprintf("expected %q after %q, got nothing instead", expected, prev)

	err := gcerr.New(BadWord, msg)
	return err
}

// NewErrGotUnexpected returns a new error when the lexer encounters an unexpected character.
//
// This error has code BadWord and message "expected <expected> before <after>, got <got> instead".
//
// Parameters:
//   - after: The character after the unexpected character.
//   - expected: The expected character.
//   - got: The unexpected character.
//
// Returns:
//   - *gcerr.Err: The error. Never returns nil.
func NewErrGotUnexpected(after, expected, got rune) *gcerr.Err {
	msg := fmt.Sprintf("expected %q before %q, got %q instead", expected, after, got)

	err := gcerr.New(BadWord, msg)
	return err
}

// NewErrBadGroup returns a new error when the lexer encounters a group that was expected
// but was not found or it was found but was invalid.
//
// This error has code BadGroup and message "expected group <expected>, got <got> instead".
//
// Parameters:
//   - expected: The expected group.
//   - got: The unexpected group.
//
// Returns:
//   - *gcerr.Err: The error. Never returns nil.
func NewErrBadGroup(expected string, got *rune) *gcerr.Err {
	var msg string

	if got == nil {
		msg = fmt.Sprintf("expected group %q, got nothing instead", expected)
	} else {
		msg = fmt.Sprintf("expected group %q, got %q instead", expected, *got)
	}

	err := gcerr.New(BadGroup, msg)

	return err
}

// NewErrNoGroupSpecified returns a new error when the lexer encounters a group that was expected
// but was not found or it was found but was invalid.
//
// This error has code BadGroup and message "no group was specified".
//
// Returns:
//   - *gcerr.Err: The error. Never returns nil.
func NewErrNoGroupSpecified() *gcerr.Err {
	err := gcerr.New(BadGroup, "no group was specified")
	return err
}
