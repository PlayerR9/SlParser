package grammar

import (
	gers "github.com/PlayerR9/go-errors"
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

func NewBadParseTree(msg string) *gers.Err {
	return gers.New(BadParseTree, msg)
}
