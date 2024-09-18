package util

// TODO: Remove this once go-commons is updated

import (
	"fmt"
	"strconv"
	"strings"
)

type ErrValue struct {
	Kind        string
	Expected    any
	Got         any
	ShouldQuote bool
}

func (e ErrValue) Error() string {
	var builder strings.Builder

	builder.WriteString("expected ")

	if e.Kind != "" {
		builder.WriteString(e.Kind)
		builder.WriteString(" to be ")
	}

	if e.Expected == nil {
		builder.WriteString("nothing")
	} else if e.ShouldQuote {
		fmt.Fprintf(&builder, "%q", e.Expected)
	} else {
		fmt.Fprintf(&builder, "%s", e.Expected)
	}

	builder.WriteString(", got ")

	if e.Got == nil {
		builder.WriteString("nothing")
	} else if e.ShouldQuote {
		fmt.Fprintf(&builder, "%q", e.Got)
	} else {
		fmt.Fprintf(&builder, "%s", e.Got)
	}

	builder.WriteString(" instead")

	return builder.String()
}

func NewErrValue(kind string, expected, got any, should_quote bool) *ErrValue {
	return &ErrValue{
		Kind:        kind,
		Expected:    expected,
		Got:         got,
		ShouldQuote: should_quote,
	}
}

type ErrValues[T any] struct {
	Kind        string
	Expecteds   []T
	Got         any
	ShouldQuote bool
}

func (e ErrValues[T]) Error() string {
	var builder strings.Builder

	builder.WriteString("expected ")

	if e.Kind != "" {
		builder.WriteString(e.Kind)
		builder.WriteString(" to be ")
	}

	switch len(e.Expecteds) {
	case 0:
		builder.WriteString("nothing")
	case 1:
		if e.ShouldQuote {
			builder.WriteString(strconv.Quote(fmt.Sprintf("%v", e.Expecteds[0])))
		} else {
			fmt.Fprintf(&builder, "%v", e.Expecteds[0])
		}
	default:
		elems := make([]string, 0, len(e.Expecteds))

		if e.ShouldQuote {
			for i := 0; i < len(e.Expecteds); i++ {
				elems = append(elems, strconv.Quote(fmt.Sprintf("%v", e.Expecteds[i])))
			}
		} else {
			for i := 0; i < len(e.Expecteds); i++ {
				elems = append(elems, fmt.Sprintf("%v", e.Expecteds[i]))
			}
		}

		builder.WriteString("either ")
		builder.WriteString(strings.Join(elems[:len(elems)-1], ", "))
		builder.WriteString(" or ")
		builder.WriteString(elems[len(elems)-1])
	}

	builder.WriteString(", got ")

	if e.Got == nil {
		builder.WriteString("nothing")
	} else if e.ShouldQuote {
		fmt.Fprintf(&builder, "%q", e.Got)
	} else {
		fmt.Fprintf(&builder, "%s", e.Got)
	}

	builder.WriteString(" instead")

	return builder.String()
}

func NewErrValues[T any](kind string, expected []T, got any, should_quote bool) *ErrValue {
	return &ErrValue{
		Kind:        kind,
		Expected:    expected,
		Got:         got,
		ShouldQuote: should_quote,
	}
}
