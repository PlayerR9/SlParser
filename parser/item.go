package parser

import (
	"errors"
	"iter"

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

	size := rule.size()
	if pos < 0 || pos > size {
		return nil, gcers.NewErrInvalidParameter("pos", errors.New("value is out of range"))
	}

	var act actioner

	if pos < size {
		act = &shift_action{}
	} else {
		rhs, ok := rule.rhs_at(pos - 1)
		dba.AssertOk(ok, "rhs_at(%d)", pos-1)

		if rhs == T(0) {
			act = &accept_action{}
		} else {
			act = &reduce_action{}
		}
	}

	return &Item[T]{
		rule: rule,
		act:  act,
	}, nil
}

// backward_rhs returns the backward rhs of the item.
//
// Returns:
//   - iter.Seq[T]: the backward rhs of the item.
func (i Item[T]) backward_rhs() iter.Seq[T] {
	return i.rule.backward_rhs()
}

// lhs returns the left hand side of the item.
//
// Returns:
//   - T: the left hand side of the item.
func (i Item[T]) lhs() T {
	return i.rule.get_lhs()
}
