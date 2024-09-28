package grammar

import (
	gerr "github.com/PlayerR9/go-errors/error"
)

type ErrorCode int

const (
	// BadParseTree occurs when a parse tree is invalid.
	BadParseTree ErrorCode = iota
)

// Int implements the error.ErrorCoder interface.
func (e ErrorCode) Int() int {
	return int(e)
}

func NewBadParseTree(msg string) *gerr.Err {
	return gerr.New(BadParseTree, msg)
}
