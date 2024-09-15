package parser

import (
	"strconv"
	"strings"

	gr "github.com/PlayerR9/SlParser/grammar"
)

type ErrUnexpectedToken[T gr.TokenTyper] struct {
	Expected T
	After    *T
	Got      *T
}

func (e ErrUnexpectedToken[T]) Error() string {
	var after string

	if e.After == nil {
		after = " at the end"
	} else {
		after = " before " + strconv.Quote((*e.After).String())
	}

	var got string

	if e.Got == nil {
		got = "nothing"
	} else {
		got = strconv.Quote((*e.Got).String())
	}

	var builder strings.Builder

	builder.WriteString("expected ")
	builder.WriteString(strconv.Quote(e.Expected.String()))
	builder.WriteString(after)
	builder.WriteString(", got ")
	builder.WriteString(got)
	builder.WriteString(" instead")

	return builder.String()
}

func NewErrUnexpectedToken[T gr.TokenTyper](expected T, after, got *T) *ErrUnexpectedToken[T] {
	return &ErrUnexpectedToken[T]{
		Expected: expected,
		After:    after,
		Got:      got,
	}
}
