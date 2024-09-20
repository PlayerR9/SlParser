package parser

import (
	"errors"
	"fmt"
	"io"

	gr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/parser/internal"
	gcers "github.com/PlayerR9/go-commons/errors"
	dba "github.com/PlayerR9/go-debug/assert"
)

// ActiveParser is a parser.
type ActiveParser[T gr.TokenTyper] struct {
	// global is the global parser.
	global *Parser[T]

	// tokens is the list of tokens.
	tokens []*gr.Token[T]

	// stack is the parser stack.
	stack *internal.Stack[T]
}

func NewActiveParser[T gr.TokenTyper](global *Parser[T]) (*ActiveParser[T], error) {
	if global == nil {
		return nil, gcers.NewErrNilParameter("global")
	}

	tokens := make([]*gr.Token[T], len(global.tokens))
	copy(tokens, global.tokens)

	var stack internal.Stack[T]

	return &ActiveParser[T]{
		global: global,
		tokens: tokens,
		stack:  &stack,
	}, nil
}

func (ap *ActiveParser[T]) Align(history *History[*Item[T]]) error {
	if ap == nil {
		return gcers.NilReceiver
	} else if history == nil {
		return gcers.NewErrNilParameter("history")
	}

	err := ap.shift() // initial shift
	if err != nil {
		return err
	}

	err = Align(history, ap)
	if err != nil {
		return err
	}

	return nil
}

// SetTokens sets the list of tokens.
//
// Parameters:
//   - tokens: the list of tokens.
//
// Does nothing if the receiver is nil.
func (p *ActiveParser[T]) SetTokens(tokens []*gr.Token[T]) {
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
func (p ActiveParser[T]) Pop() (*gr.ParseTree[T], bool) {
	tk, ok := p.stack.Pop()
	return tk, ok
}

// shift is a helper function that shifts a token.
//
// Returns:
//   - error: if an error occurred.
func (p *ActiveParser[T]) shift() error {
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
func (p ActiveParser[T]) reduce(it *Item[T]) error {
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

// ApplyEvent applies an event to the active parser. Does nothing if item or
// the receiver are nil.
//
// Parameters:
//   - item: The item to apply the event to.
//
// Returns:
//   - bool: True if the active parser has accepted, false otherwise.
//   - error: if an error occurred.
func (ap *ActiveParser[T]) ApplyEvent(item *Item[T]) (bool, error) {
	if ap == nil || item == nil {
		return false, nil
	}

	switch item.act {
	case internal.ActShift:
		err := ap.shift()
		if err != nil {
			return false, err
		}
	case internal.ActReduce:
		err := ap.reduce(item)
		if err != nil {
			ap.stack.Refuse()

			return false, err
		}
	case internal.ActAccept:
		err := ap.reduce(item)
		if err != nil {
			ap.stack.Refuse()

			return false, err
		}

		return true, nil
	default:
		ap.stack.Refuse()

		return false, fmt.Errorf("unexpected action: %v", item.act)
	}

	return false, nil
}

func (ap *ActiveParser[T]) DetermineNextEvents() ([]*Item[T], error) {
	if ap == nil {
		return nil, gcers.NilReceiver
	}

	defer ap.stack.Refuse()

	top, ok := ap.stack.Pop()
	if !ok {
		return nil, errors.New("End of Input was reached but no accepting state was found")
	}

	lookahead := top.Lookahead()
	type_ := top.Type()

	fn, ok := ap.global.ParseFnOf(type_)
	if !ok || fn == nil {
		return nil, fmt.Errorf("no rule for %q", type_.String())
	}

	events, err := fn(ap, top, lookahead)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("no action for %q", type_.String())
	}

	return events, nil
}

// Forest returns the forest.
//
// Returns:
//   - []*gr.ParseTree[T]: the forest.
func (p ActiveParser[T]) Forest() []*gr.ParseTree[T] {
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
func (p *ActiveParser[T]) Reset() {
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
