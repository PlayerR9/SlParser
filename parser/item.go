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

// String implements the fmt.Stringer interface.
func (item Item[T]) String() string {
	if item.rule == nil {
		return ""
	}

	var elems []string

	lhs := item.rule.Lhs()

	elems = append(elems, lhs.String(), "->")

	var i int
	for rhs := range item.rule.ForwardRhs() {
		elems = append(elems, rhs.String())

		if i == item.pos {
			elems = append(elems, "#")
		}

		i++
	}

	return strings.Join(elems, " ")
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

// ForwardRhs returns the forward rhs of the item.
//
// Returns:
//   - iter.Seq[T]: The forward rhs of the item.
func (item Item[T]) ForwardRhs() iter.Seq[T] {
	dba.AssertNotNil(item.rule, "item.rule")

	return item.rule.ForwardRhs()
}

// BackwardRhs returns the backward rhs of the item.
//
// Returns:
//   - iter.Seq[T]: the backward rhs of the item.
func (item Item[T]) BackwardRhs() iter.Seq[T] {
	dba.AssertNotNil(item.rule, "item.rule")

	return item.rule.BackwardRhs()
}

// Lhs returns the left hand side of the item.
//
// Returns:
//   - T: the left hand side of the item.
func (item Item[T]) Lhs() T {
	dba.AssertNotNil(item.rule, "item.rule")

	return item.rule.Lhs()
}

// RhsAt returns the rhs at the given index.
//
// Returns:
//   - T: the rhs at the given index.
//   - bool: true if the index is valid, false otherwise.
func (item Item[T]) RhsAt(idx int) (T, bool) {
	if item.rule == nil {
		return T(-1), false
	}

	return item.rule.RhsAt(idx)
}

// Pos returns the position of the item in the rule.
//
// Returns:
//   - int: the position of the item in the rule.
func (item Item[T]) Pos() int {
	return item.pos
}

// set_lookaheads is a helper function for the lookaheads of the item.
//
// Parameters:
//   - lookaheads: the list of lookaheads.
func (item *Item[T]) set_lookaheads(lookaheads []T) {
	if item == nil {
		return
	}

	item.lookaheads = lookaheads
}

func (item Item[T]) HasRhsByOffset(offset int) bool {
	idx := item.pos - offset

	return item.rule != nil && item.rule.HasRhsAt(idx)
}

func (item Item[T]) RhsByOffset(offset int) (T, bool) {
	if item.rule == nil {
		return T(-1), false
	}

	idx := item.pos - offset

	rhs, ok := item.rule.RhsAt(idx)
	return rhs, ok
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
