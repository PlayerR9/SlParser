package parser

import (
	"iter"
	"strings"

	gr "github.com/PlayerR9/SlParser/grammar"
	gcers "github.com/PlayerR9/go-commons/errors"
	dba "github.com/PlayerR9/go-debug/assert"
)

// Item is an item in the parsing table.
type Item[T gr.TokenTyper] struct {
	// rule is the rule of the item.
	rule *Rule[T]

	// act is the action of the item.
	act actioner

	// pos is the position of the rhs in the rule.
	pos int

	// lookaheads is the list of lookaheads.
	lookaheads []T
}

func (i Item[T]) String() string {
	var builder strings.Builder

	lhs := i.rule.Lhs()

	builder.WriteString(lhs.String())
	builder.WriteString(" ->")

	var j int
	for rhs := range i.rule.ForwardRhs() {
		builder.WriteRune(' ')
		builder.WriteString(rhs.String())

		if j == i.pos {
			builder.WriteString(" #")
		}

		j++
	}

	return builder.String()
}

// NewItem creates a new item.
//
// Parameters:
//   - rule: the rule of the item.
//   - pos: the position of the item in the rule.
//
// Returns:
//   - *Item: the new item.
//   - error: if the position is out of range.
func NewItem[T gr.TokenTyper](rule *Rule[T], pos int) (*Item[T], error) {
	if rule == nil {
		return nil, gcers.NewErrNilParameter("rule")
	}

	size := rule.Size()
	if pos < 0 || pos >= size {
		return nil, gcers.NewErrInvalidParameter("pos", gcers.NewErrOutOfBounds(pos, 0, size))
	}

	var act actioner

	if pos < size-1 {
		act = &shift_action{}
	} else {
		rhs, ok := rule.RhsAt(pos)
		dba.AssertOk(ok, "rhs_at(%d)", pos)

		if rhs == T(0) {
			act = &accept_action{}
		} else {
			act = &reduce_action{}
		}
	}

	return &Item[T]{
		rule: rule,
		act:  act,
		pos:  pos,
	}, nil
}

// Same as NewItem, but panics instead of returning an error.
//
// Never returns nil.
func MustNewItem[T gr.TokenTyper](rule *Rule[T], pos int) *Item[T] {
	if rule == nil {
		panic(gcers.NewErrNilParameter("rule"))
	}

	size := rule.Size()
	if pos < 0 || pos >= size {
		panic(gcers.NewErrInvalidParameter("pos", gcers.NewErrOutOfBounds(pos, 0, size)))
	}

	var act actioner

	if pos < size-1 {
		act = &shift_action{}
	} else {
		rhs, ok := rule.RhsAt(pos)
		dba.AssertOk(ok, "rhs_at(%d)", pos)

		if rhs == T(0) {
			act = &accept_action{}
		} else {
			act = &reduce_action{}
		}
	}

	return &Item[T]{
		rule: rule,
		act:  act,
	}
}

// backward_rhs returns the backward rhs of the item.
//
// Returns:
//   - iter.Seq[T]: the backward rhs of the item.
func (i Item[T]) backward_rhs() iter.Seq[T] {
	return i.rule.BackwardRhs()
}

// lhs returns the left hand side of the item.
//
// Returns:
//   - T: the left hand side of the item.
func (i Item[T]) lhs() T {
	return i.rule.Lhs()
}

// RhsAt returns the rhs at the given index.
//
// Returns:
//   - T: the rhs at the given index.
//   - bool: true if the index is valid, false otherwise.
func (i Item[T]) RhsAt(idx int) (T, bool) {
	return i.rule.RhsAt(idx)
}

// Pos returns the position of the item in the rule.
//
// Returns:
//   - int: the position of the item in the rule.
func (i Item[T]) Pos() int {
	return i.pos
}

// set_lookaheads is a helper function for the lookaheads of the item.
//
// Parameters:
//   - lookaheads: the list of lookaheads.
func (item *Item[T]) set_lookaheads(lookaheads []T) {
	item.lookaheads = lookaheads
}

/*
func LookaheadsOf[T gr.TokenTyper](items ...*Item[T]) {
	item := items[0]

	rhs, ok := item.RhsAt(item.pos + 1)
	if !ok {

	} else {

	}
}
*/
