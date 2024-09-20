package parser

import (
	"fmt"
	"io"

	gr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/parser/internal"
	gcers "github.com/PlayerR9/go-commons/errors"
	dba "github.com/PlayerR9/go-debug/assert"
)

// Parser is a parser.
type Parser[T gr.TokenTyper] struct {
	// tokens is the list of tokens.
	tokens []*gr.Token[T]

	// stack is the parser stack.
	stack *internal.Stack[T]

	// table is the parser table.
	table map[T]ParseFn[T]
}

// SetTokens sets the list of tokens.
//
// Parameters:
//   - tokens: the list of tokens.
//
// Does nothing if the receiver is nil.
func (p *Parser[T]) SetTokens(tokens []*gr.Token[T]) {
	if p == nil {
		return
	}

	p.tokens = tokens
}

// Pop pops a token from the stack.
//
// Returns:
//   - grammar.ParseTree[T]: the popped token.
//   - bool: true if the token was found, false otherwise.
func (p Parser[T]) Pop() (*gr.ParseTree[T], bool) {
	tk, ok := p.stack.Pop()
	return tk, ok
}

// shift is a helper function that shifts a token.
//
// Returns:
//   - error: if an error occurred.
func (p *Parser[T]) shift() error {
	if p == nil {
		return gcers.NilReceiver
	}

	if len(p.tokens) == 0 {
		return io.EOF
	}

	tk := p.tokens[0]
	p.tokens = p.tokens[1:]

	tree, err := gr.NewTree(tk)
	dba.AssertErr(err, "grammar.NewTree(tk)")

	p.stack.Push(tree)

	return nil
}

// reduce is a helper function that reduces an item.
//
// Parameters:
//   - it: the item to reduce.
//
// Returns:
//   - error: if an error occurred.
func (p Parser[T]) reduce(it *Item[T]) error {
	var prev *T

	for rhs := range it.BackwardRhs() {
		top, ok := p.stack.Pop()
		if !ok {
			return NewErrUnexpectedToken([]T{rhs}, prev, nil)
		}

		type_ := top.Type()
		if type_ != rhs {
			return NewErrUnexpectedToken([]T{rhs}, prev, &type_)
		}

		prev = &type_
	}

	popped := p.stack.Popped()
	p.stack.Accept()

	tree, err := gr.Combine(it.Lhs(), popped)
	dba.AssertErr(err, "grammar.Combine(%s, popped)", it.Lhs().String())

	p.stack.Push(tree)

	return nil
}

// Parse parses the list of tokens.
//
// Parameters:
//   - tokens: the list of tokens.
//
// Returns:
//   - error: if an error occurred.
func (p *Parser[T]) Parse() (*ActiveParser[T], error) {
	if p == nil {
		return nil, nil
	}

	if p.Size() == 0 {
		top1, ok := p.stack.Pop()
		dba.AssertOk(ok, "p.pop()")

		p.stack.Refuse()

		return nil, fmt.Errorf("no rule for %q", top1.Type().String())
	}

	ap, err := NewActiveParser(p)
	dba.AssertErr(err, "NewActiveParser(p)")

	var aps []*ActiveParser[T]

	var history History[*Item[T]]

	possible, err := Execute(&history, ap)
	if len(possible) > 0 {
		for _, path := range possible {
			ap, err := NewActiveParser(p)
			dba.AssertErr(err, "NewActiveParser(p)")

			err = ap.Align(path)
			if err != nil {
				return nil, err
			}

			aps = append(aps, ap)
		}
	}

	if err == nil {
		return ap, nil
	}

	return ap, nil
}

// Forest returns the forest.
//
// Returns:
//   - []*gr.ParseTree[T]: the forest.
func (p Parser[T]) Forest() []*gr.ParseTree[T] {
	if p.stack.IsEmpty() {
		return nil
	}

	var forest []*gr.ParseTree[T]

	for {
		top, ok := p.stack.Pop()
		if !ok {
			break
		}

		forest = append(forest, top)
	}

	// slices.Reverse(forest)

	return forest
}

// Reset resets the parser; allowing it to be reused.
func (p *Parser[T]) Reset() {
	if p == nil {
		return
	}

	if len(p.tokens) > 0 {
		for i := 0; i < len(p.tokens); i++ {
			p.tokens[i] = nil
		}

		p.tokens = p.tokens[:0]
	}

	p.stack.Reset()
}

// Size returns the size of the table.
//
// Returns:
//   - int: the size of the table.
func (p Parser[T]) Size() int {
	return len(p.table)
}

// ParseFnOf returns the parse function for the given symbol.
//
// Parameters:
//   - symbol: the symbol.
//
// Returns:
//   - ParseFn[T]: the parse function.
//   - bool: true if the symbol was found, false otherwise.
func (p Parser[T]) ParseFnOf(symbol T) (ParseFn[T], bool) {
	if p.table == nil {
		return nil, false
	}

	fn, ok := p.table[symbol]
	if !ok {
		return nil, false
	}

	return fn, true
}
