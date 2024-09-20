package internal

import (
	"strings"

	gr "github.com/PlayerR9/SlParser/grammar"
)

// Expecter is an expecter.
type Expecter interface {
	// String is the string representation of the expecter.
	//
	// Returns:
	//   - string: the string representation of the expecter.
	String() string
}

// ExpectTerminal is an expecter that expects a terminal symbol.
type ExpectTerminal[T gr.TokenTyper] struct {
	// symbol is the terminal symbol.
	symbol T
}

// String implements the Expecter interface.
func (e ExpectTerminal[T]) String() string {
	return e.symbol.String()
}

// NewExpectTerminal creates a new expecter that expects a terminal symbol.
//
// Parameters:
//   - symbol: the terminal symbol.
//
// Returns:
//   - *ExpectTerminal[T]: the new expecter. Never returns nil.
func NewExpectTerminal[T gr.TokenTyper](symbol T) *ExpectTerminal[T] {
	return &ExpectTerminal[T]{
		symbol: symbol,
	}
}

// ExpectNonTerminal is an expecter that expects a non-terminal symbol.
type ExpectNonTerminal[T gr.TokenTyper] struct {
	// reduce is the non-terminal symbol that is expected.
	reduce T

	// expecteds is a list of lookaheads.
	expecteds []Expecter
}

// String implements the Expecter interface.
func (e ExpectNonTerminal[T]) String() string {
	var builder strings.Builder

	builder.WriteString("expected(")
	builder.WriteString(e.reduce.String())
	builder.WriteString(") = {")

	if len(e.expecteds) > 0 {
		elems := make([]string, 0, len(e.expecteds))

		for _, e := range e.expecteds {
			elems = append(elems, e.String())
		}

		builder.WriteString(strings.Join(elems, ", "))
	}

	builder.WriteRune('}')

	return builder.String()
}

// NewExpectNonTerminal creates a new expecter that expects a non-terminal symbol.
//
// Parameters:
//   - reduce: the non-terminal symbol that is expected.
//   - expecteds: a list of lookaheads.
//
// Returns:
//   - *ExpectNonTerminal[T]: the new expecter. Never returns nil.
func NewExpectNonTerminal[T gr.TokenTyper](reduce T, expecteds []Expecter) *ExpectNonTerminal[T] {
	return &ExpectNonTerminal[T]{
		reduce:    reduce,
		expecteds: expecteds,
	}
}
