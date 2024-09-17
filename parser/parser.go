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
//   - grammar.Token[T]: the popped token.
//   - bool: true if the token was found, false otherwise.
func (p Parser[T]) Pop() (*gr.Token[T], bool) {
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

	p.stack.Push(tk)

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

	for rhs := range it.backward_rhs() {
		top, ok := p.stack.Pop()
		if !ok {
			return NewErrUnexpectedToken([]T{rhs}, prev, nil)
		} else if top.Type != rhs {
			return NewErrUnexpectedToken([]T{rhs}, prev, &top.Type)
		}

		prev = &top.Type
	}

	popped := p.stack.Popped()
	p.stack.Accept()

	tk, err := gr.NewNonTerminalToken(it.lhs(), popped)
	dba.AssertErr(err, "grammar.NewNonTerminalToken(%s, popped)", it.lhs().String())

	p.stack.Push(tk)

	return nil
}

// Parse parses the list of tokens.
//
// Parameters:
//   - tokens: the list of tokens.
//
// Returns:
//   - error: if an error occurred.
func (p *Parser[T]) Parse() error {
	if p == nil {
		return nil
	}

	err := p.shift() // initial shift
	if err != nil {
		return err
	}

	if len(p.table) == 0 {
		top1, ok := p.stack.Pop()
		dba.AssertOk(ok, "p.pop()")

		p.stack.Refuse()

		return fmt.Errorf("no rule for %q", top1.Type.String())
	}

	for {
		top1, ok := p.stack.Pop()
		if !ok {
			break
		}

		la := top1.Lookahead

		fn, ok := p.table[top1.Type]
		if !ok {
			p.stack.Refuse()

			return fmt.Errorf("no rule for %q", top1.Type.String())
		}

		items, err := fn(p, top1, la)
		p.stack.Refuse()

		if err != nil {
			return err
		}

		if len(items) > 1 {
			fmt.Printf("WARNING: ambiguous grammar, found %d actions\n", len(items))
		} else if len(items) == 0 {
			return fmt.Errorf("no action for %q", top1.Type.String())
		}

		it := items[0]

		switch act := it.act.(type) {
		case *shift_action:
			if len(p.tokens) == 0 {
				return fmt.Errorf("unexpected end of input")
			}

			tk := p.tokens[0]
			p.tokens = p.tokens[1:]

			p.stack.Push(tk)
		case *reduce_action:
			err := p.reduce(it)
			if err != nil {
				p.stack.Refuse()

				return err
			}
		case *accept_action:
			err := p.reduce(it)
			if err != nil {
				p.stack.Refuse()

				return err
			}

			return nil
		default:
			p.stack.Refuse()

			return fmt.Errorf("unexpected action: %v", act)
		}
	}

	p.stack.Refuse()

	return fmt.Errorf("end of input but no accept action found")
}

// Forest returns the forest.
//
// Returns:
//   - []*gr.Token[T]: the forest.
func (p Parser[T]) Forest() []*gr.Token[T] {
	if p.stack.IsEmpty() {
		return nil
	}

	var forest []*gr.Token[T]

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
