package parser

import (
	"errors"
	"fmt"
	"io"

	gr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/parser/internal"
	bck "github.com/PlayerR9/SlParser/util/go-commons/backup"
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

	// err is the error.
	err error
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
		err:    nil,
	}, nil
}

func (ap *ActiveParser[T]) Align(history *bck.History[*Item[T]]) bool {
	dba.AssertNotNil(ap, "ap")
	dba.AssertNotNil(history, "history")

	ap.shift()
	if ap.HasError() {
		return false
	}

	bck.Align(history, ap)

	return !ap.HasError()
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
func (ap *ActiveParser[T]) shift() {
	dba.AssertNotNil(ap, "ap")

	if len(ap.tokens) == 0 {
		ap.err = io.EOF

		return
	}

	tk := ap.tokens[0]
	ap.tokens = ap.tokens[1:]

	tree, err := gr.NewTree(tk)
	dba.AssertErr(err, "grammar.NewTree(tk)")

	ap.stack.Push(tree)
}

// reduce is a helper function that reduces an item.
//
// Parameters:
//   - it: the item to reduce. Assumed to be non-nil.
func (ap *ActiveParser[T]) reduce(it *Item[T]) {
	dba.AssertNotNil(ap, "ap")
	dba.AssertNotNil(it, "it")

	var prev *T

	for rhs := range it.BackwardRhs() {
		top, ok := ap.stack.Pop()
		if !ok {
			ap.err = NewErrUnexpectedToken([]T{rhs}, prev, nil)

			return
		}

		type_ := top.Type()
		if type_ != rhs {
			ap.err = NewErrUnexpectedToken([]T{rhs}, prev, &type_)

			return
		}

		prev = &type_
	}

	popped := ap.stack.Popped()
	ap.stack.Accept()

	lhs := it.Lhs()

	tree, err := gr.Combine(lhs, popped)
	dba.AssertErr(err, "grammar.Combine(%s, popped)", lhs.String())

	ap.stack.Push(tree)
}

// ApplyEvent applies an event to the active parser. Does nothing if item or
// the receiver are nil.
//
// Parameters:
//   - item: The item to apply the event to.
//
// Returns:
//   - bool: True if the active parser has accepted, false otherwise.
func (ap *ActiveParser[T]) ApplyEvent(item *Item[T]) bool {
	if ap == nil || item == nil {
		return false
	}

	switch item.act {
	case internal.ActShift:
		ap.shift()
	case internal.ActReduce:
		ap.reduce(item)
	case internal.ActAccept:
		ap.reduce(item)
	default:
		ap.err = fmt.Errorf("unexpected action: %v", item.act)
	}

	if ap.HasError() {
		ap.stack.Refuse()

		return false
	}

	return item.act == internal.ActAccept
}

func (ap *ActiveParser[T]) DetermineNextEvents() []*Item[T] {
	dba.AssertNotNil(ap, "ap")

	defer ap.stack.Refuse()

	top, ok := ap.stack.Pop()
	if !ok {
		ap.err = errors.New("End of Input was reached but no accepting state was found")

		return nil
	}

	lookahead := top.Lookahead()
	type_ := top.Type()

	fn, ok := ap.global.ParseFnOf(type_)
	if !ok || fn == nil {
		ap.err = fmt.Errorf("no rule for %q", type_.String())

		return nil
	}

	events, err := fn(ap, top, lookahead)
	if err != nil {
		ap.err = err

		return nil
	}

	if len(events) == 0 {
		ap.err = fmt.Errorf("no action for %q", type_.String())

		return nil
	}

	return events
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

func (ap ActiveParser[T]) HasError() bool {
	return ap.err != nil
}

func (ap ActiveParser[T]) Error() error {
	return ap.err
}
