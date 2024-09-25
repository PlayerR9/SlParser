package kdd

import (
	gerr "github.com/PlayerR9/go-errors/error"
)

//go:generate stringer -type=ErrorCode

type ErrorCode int

const (
	// InvalidSyntax occurs when the AST is invalid or not
	// as it should be.
	InvalidSyntax ErrorCode = iota
)

// Int implements the error.ErrorCoder interface.
func (e ErrorCode) Int() int {
	return int(e)
}

// NewErrSyntax returns a new error.Err error representing a
// syntax error.
//
// Parameters:
//   - msg: The reason why the syntax is wrong.
//
// Returns:
//   - *error.Err: A pointer to the newly created error. Never returns nil.
//
// If the message is not specified, the string "the AST is invalid is used instead".
func NewErrSyntax(msg string) *gerr.Err {
	if msg == "" {
		msg = "the AST is invalid"
	}

	err := gerr.New(InvalidSyntax, msg)

	return err
}

func CheckNode(node *Node) error {
	if node == nil {
		return nil // TODO: Decide this.
	}

	switch node.Type {
	case SourceNode:
	default:
	}

	return nil
}
