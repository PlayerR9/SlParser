package ast

import (
	"fmt"
	"strconv"
	"strings"

	gr "github.com/PlayerR9/SlParser/grammar"
	gcstr "github.com/PlayerR9/go-commons/strings"
	gcers "github.com/PlayerR9/go-errors/error"
	"github.com/dustin/go-humanize"
)

//go:generate stringer -type=ErrorCode

// ErrorCode is the error code of an error.
type ErrorCode int

const (
	// UnregisteredType occurs when a type is not registered.
	UnregisteredType ErrorCode = iota

	// BadSyntaxTree occurs when a syntax tree is invalid.
	BadSyntaxTree
)

// Int implements the error.ErrorCoder interface.
func (e ErrorCode) Int() int {
	return int(e)
}

// NewUnregisteredType creates a new UnregisteredType error.
//
// Parameters:
//   - type_: The type that is not registered.
//   - in: The input that caused the error.
//
// Returns:
//   - *gcers.Err: The error. Never returns nil.
func NewUnregisteredType[T gr.TokenTyper](type_ T, in string) *gcers.Err {
	msg := fmt.Sprintf("type %q is not registered", type_.String())

	err := gcers.New(UnregisteredType, msg)
	err.AddFrame(in)

	return err
}

// NewBadSyntaxTree creates a new BadSyntaxTree error.
//
// Parameters:
//   - at: The position of the token.
//   - type_: The type of the token.
//   - got: The unexpected value.
//
// Returns:
//   - *gcers.Err: The error. Never returns nil.
func NewBadSyntaxTree[T gr.TokenTyper](at int, type_ T, got string) *gcers.Err {
	if got != "" {
		got = strconv.Quote(got)
	}

	msg := gcstr.ExpectedValue("type", gcstr.Quote(type_), got)

	err := gcers.New(BadSyntaxTree, msg)
	err.AddFrame(humanize.Ordinal(at+1) + " child")

	return err
}

//////////////////////////////////////////////

type ErrIn[T gr.TokenTyper] struct {
	Type   T
	Reason error
}

func (e ErrIn[T]) Error() string {
	var reason string

	if e.Reason == nil {
		reason = "something went wrong"
	} else {
		reason = e.Reason.Error()
	}

	var builder strings.Builder

	builder.WriteString("in rule ")
	builder.WriteString(strconv.Quote(e.Type.String()))
	builder.WriteString(": ")
	builder.WriteString(reason)

	return builder.String()
}

func (e ErrIn[T]) Unwrap() error {
	return e.Reason
}

func NewErrIn[T gr.TokenTyper](type_ T, reason error) *ErrIn[T] {
	return &ErrIn[T]{
		Type:   type_,
		Reason: reason,
	}
}

func (e *ErrIn[T]) ChangeReason(new_reason error) {
	if e == nil {
		return
	}

	e.Reason = new_reason
}
