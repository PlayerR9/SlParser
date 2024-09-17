package util

import (
	"strings"
)

// ErrorCode is an error code.
type ErrorCode interface {
	~int

	// String returns the string representation of the error code.
	//
	// Returns:
	//   - string: the string representation of the error code.
	String() string
}

// Err is a generic yet advanced error type.
type Err[T ErrorCode] struct {
	// Code is the error code.
	Code T

	// Reason is the error message.
	Reason error

	// Suggestions is a list of suggestions.
	Suggestions []string
}

// Error implements the error interface.
//
// Message:
//
//	"{code}: {reason}".
func (e Err[T]) Error() string {
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
func (e Err[T]) Unwrap() error {
	return e.Reason
}

// ChangeReason changes the reason for the error.
//
// Parameters:
//   - new_reason: the new reason for the error.
func (e *Err[T]) ChangeReason(new_reason error) {
	if e == nil {
		return
	}

	e.Reason = new_reason
}

// NewErr creates a new Err error.
//
// Parameters:
//   - code: the error code.
//   - reason: the reason for the error.
//
// Returns:
//   - *Err: the error. Never returns nil.
func NewErr[T ErrorCode](code T, reason error) *Err[T] {
	return &Err[T]{
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
func (e *Err[T]) AddSuggestion(sentences ...string) {
	if e == nil {
		return
	}

	paragraph := strings.Join(sentences, " ")

	e.Suggestions = append(e.Suggestions, paragraph)
}
