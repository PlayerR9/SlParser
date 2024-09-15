package ast

import (
	"strconv"
	"strings"

	gr "github.com/PlayerR9/SlParser/grammar"
)

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
