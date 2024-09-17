package lexer

import (
	"strings"
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
)

// Err is a generic yet advanced error type.
type Err struct {
	// Code is the error code.
	Code ErrorCode

	// Reason is the error message.
	Reason error

	// Suggestions is a list of suggestions.
	Suggestions []string

	// Pos is the position of the error.
	Pos int
}

// Error implements the error interface.
//
// Message:
//
//	"{code}: {reason}".
func (e Err) Error() string {
	var builder strings.Builder

	builder.WriteString(e.Code.String())
	builder.WriteString(": ")

	if e.Reason == nil {
		builder.WriteString("something went wrong")
	} else {
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap unwraps the error.
//
// Returns:
//   - error: the reason for the error.
func (e Err) Unwrap() error {
	return e.Reason
}

// ChangeReason changes the reason for the error.
//
// Parameters:
//   - new_reason: the new reason for the error.
func (e *Err) ChangeReason(new_reason error) {
	if e == nil {
		return
	}

	e.Reason = new_reason
}

// NewErr creates a new Err error.
//
// Parameters:
//   - code: the error code.
//   - pos: the position of the error.
//   - reason: the reason for the error.
//
// Returns:
//   - *Err: the error. Never returns nil.
func NewErr(code ErrorCode, pos int, reason error) *Err {
	return &Err{
		Pos:    pos,
		Code:   code,
		Reason: reason,
	}
}

// AddSuggestion adds a suggestion to the error.
//
// Parameters:
//   - sentences: the sentences of the suggestion.
//
// Sentences are joined with spaces.
func (e *Err) AddSuggestion(sentences ...string) {
	if e == nil {
		return
	}

	paragraph := strings.Join(sentences, " ")

	e.Suggestions = append(e.Suggestions, paragraph)
}
