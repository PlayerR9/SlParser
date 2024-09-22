package ast

import (
	"strconv"
	"strings"

	gr "github.com/PlayerR9/SlParser/grammar"
	gcers "github.com/PlayerR9/errors/error"
	gcstr "github.com/PlayerR9/go-commons/strings"
	"github.com/dustin/go-humanize"
)

//go:generate stringer -type=ErrorCode

type ErrorCode int

const (
	UnregisteredType ErrorCode = iota
	BadSyntaxTree
)

func NewUnregisteredType[T gr.TokenTyper](type_ T, in string) *gcers.Err[ErrorCode] {
	err := gcers.NewErr(gcers.FATAL, UnregisteredType, "type "+type_.String()+"is not registered")
	err.AddFrame("", in)

	return err
}

func NewBadSyntaxTree[T gr.TokenTyper](at int, type_ T, got string) *gcers.Err[ErrorCode] {
	if got != "" {
		got = strconv.Quote(got)
	}

	msg := gcstr.ExpectedValue("type", Quote(type_), got)

	err := gcers.NewErr(gcers.FATAL, BadSyntaxTree, msg)
	err.AddFrame("", humanize.Ordinal(at+1)+" child")

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
