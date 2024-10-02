package internal

import (
	"fmt"
	"iter"
	"strings"

	gr "github.com/PlayerR9/SlParser/grammar"
	gcers "github.com/PlayerR9/go-errors"
	"github.com/PlayerR9/go-errors/assert"
)

// Item is an item in the parsing table.
type Item[T gr.TokenTyper] struct {
	// rule is the rule of the item.
	rule *Rule[T]

	// Act is the action of the item.
	Act ActionType

	// Pos is the position of the rhs in the rule.
	Pos int

	// lookaheads is the list of lookaheads.
	lookaheads []T

	// expected is the next element that is expected.
	expected Expecter
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

		if i == item.Pos {
			elems = append(elems, "#")
		}

		i++
	}

	elems = append(elems, ".", "("+item.Act.String()+")")

	if item.expected != nil {
		elems = append(elems, "--", item.expected.String())
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
		return nil, gcers.NewErrNilParameter("internal.NewItem()", "rule")
	}

	size := rule.Size()
	if pos < 0 || pos >= size {
		return nil, gcers.NewErrInvalidParameter("internal.NewItem()", fmt.Sprintf("pos of (%d) must be in range [0, %d)", pos, size))
	}

	var act ActionType

	if pos < size-1 {
		act = ActShift
	} else {
		rhs, ok := rule.RhsAt(pos)
		assert.Ok(ok, "rhs_at(%d)", pos)

		if rhs == T(0) {
			act = ActAccept
		} else {
			act = ActReduce
		}
	}

	return &Item[T]{
		rule: rule,
		Act:  act,
		Pos:  pos,
	}, nil
}

// Same as NewItem, but panics instead of returning an error.
//
// Never returns nil.
func MustNewItem[T gr.TokenTyper](rule *Rule[T], pos int) *Item[T] {
	assert.NotNil(rule, "rule")

	size := rule.Size()
	assert.CondF(pos >= 0 && pos < size, "pos must be in range [0, %d)", size)

	var act ActionType

	if pos < size-1 {
		act = ActShift
	} else {
		rhs, ok := rule.RhsAt(pos)
		assert.Ok(ok, "rhs_at(%d)", pos)

		if rhs == T(0) {
			act = ActAccept
		} else {
			act = ActReduce
		}
	}

	return &Item[T]{
		rule: rule,
		Act:  act,
	}
}

// ForwardRhs returns the forward rhs of the item.
//
// Returns:
//   - iter.Seq[T]: The forward rhs of the item.
func (item Item[T]) ForwardRhs() iter.Seq[T] {
	assert.NotNil(item.rule, "item.rule")

	return item.rule.ForwardRhs()
}

// BackwardRhs returns the backward rhs of the item.
//
// Returns:
//   - iter.Seq[T]: the backward rhs of the item.
func (item Item[T]) BackwardRhs() iter.Seq[T] {
	assert.NotNil(item.rule, "item.rule")

	return item.rule.BackwardRhs()
}

// Lhs returns the left hand side of the item.
//
// Returns:
//   - T: the left hand side of the item.
func (item Item[T]) Lhs() T {
	assert.NotNil(item.rule, "item.rule")

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

// SetLookaheads is a helper function for the lookaheads of the item.
//
// Parameters:
//   - lookaheads: the list of lookaheads.
func (item *Item[T]) SetLookaheads(lookaheads []T) {
	if item == nil {
		return
	}

	item.lookaheads = lookaheads
}

func (item Item[T]) HasRhsByOffset(offset int) bool {
	idx := item.Pos - offset

	return item.rule != nil && item.rule.HasRhsAt(idx)
}

func (item Item[T]) RhsByOffset(offset int) (T, bool) {
	if item.rule == nil {
		return T(-1), false
	}

	idx := item.Pos - offset

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
