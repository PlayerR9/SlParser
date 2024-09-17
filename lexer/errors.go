package lexer

import (
	"fmt"

	util "github.com/PlayerR9/SlParser/util"
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
)

// NewErrUnrecognizedChar creates a new error for an unrecognized character.
//
// Parameters:
//   - char: the unrecognized character.
//
// Returns:
//   - *util.Err[ErrorCode]: the error. Never returns nil.
func NewErrUnrecognizedChar(char rune) *util.Err[ErrorCode] {
	err := util.NewErr(UnrecognizedChar, fmt.Errorf("character (%q) is not a recognized character", char))
	err.AddSuggestion(
		"1. Input provided cannot be lexed by the current lexer.",
		"You may want to check for typos in the input.",
	)
	err.AddSuggestion(
		"2. (Less likely) The lexer table is not configured correctly.",
		"Contact the developer and provide this error code.",
	)

	return err
}

// NewErrInvalidInputStream creates a new error for an invalid input stream.
//
// Parameters:
//   - reason: the underlying error.
//
// Returns:
//   - *util.Err[ErrorCode]: the error. Never returns nil.
func NewErrInvalidInputStream(reason error) *util.Err[ErrorCode] {
	err := util.NewErr(InvalidInputStream, reason)
	err.AddSuggestion("Input is most likely not a valid input for the current lexer.")

	return err
}
