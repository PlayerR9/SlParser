package internal

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	slgr "github.com/PlayerR9/SlParser/grammar"
	assert "github.com/PlayerR9/go-verify"
	"github.com/PlayerR9/mygo-data/collections"
	"github.com/PlayerR9/mygo-data/sets"
	"github.com/PlayerR9/mygo-lib/common"
)

// Item is the item of a parser.
type Item struct {
	// rule is the rule of the item. (non-nil)
	rule *Rule

	// pos is the position of the item in the rule's right-hand side.
	pos uint

	// act is the action of the item.
	act ActionType

	// lookaheads is the lookaheads of the item.
	lookaheads *sets.OrderedSet[string]
}

// String implements fmt.Stringer.
func (item Item) String() string {
	lhs := item.rule.Lhs()
	rhss := item.rule.Rhss()
	act := item.act.String()

	elems := append([]string{lhs, "="}, rhss...)
	elems = append(elems, ":", "(", act, ")")

	ok := item.lookaheads.IsEmpty()

	if !ok {
		lookaheads := collections.String("*sets.OrderedSet[string]", item.lookaheads)
		elems = append(elems, "-->", lookaheads)
	}

	elems = slices.Insert(elems, int(item.pos)+1, "#")

	str := strings.Join(elems, " ")
	return str
}

// NewItem creates a new Item from the given rule and position in the rule's
// right-hand side.
//
// Parameters:
//   - rule: The rule of the item. Must not be nil.
//   - pos: The position of the item in the rule's right-hand side. Must be in
//     range (0, rule.Size()].
//
// Returns:
//   - *Item: The newly created item.
//   - error: An error if the parameters are invalid or if the item could not be
//     created.
//
// Errors:
//   - common.ErrBadParam: If the rule is nil or the position is out of range.
func NewItem(rule *Rule, pos uint) (*Item, error) {
	if rule == nil {
		return nil, common.NewErrNilParam("rule")
	}

	size := rule.Size()
	if pos == 0 || pos > size {
		return nil, common.NewErrBadParam("pos", fmt.Sprintf("must be in range (0, %d]", size))
	}

	var act ActionType

	if pos != size {
		act = ActShift
	} else {
		rhs, ok := rule.RhsAt(pos - 1)
		assert.True(ok, "rule.RhsAt(%d)", pos-1)

		if rhs == slgr.EtEOF {
			act = ActAccept
		} else {
			act = ActReduce
		}
	}

	item := &Item{
		rule:       rule,
		pos:        pos,
		act:        act,
		lookaheads: nil,
	}

	return item, nil
}

// Action returns the action type associated with the item.
//
// Returns:
//   - ActionType: The action type of the item.
func (item Item) Action() ActionType {
	return item.act
}

// Rule returns the rule associated with the item.
//
// Returns:
//   - *Rule: The rule of the item. Never returns nil.
func (item Item) Rule() *Rule {
	return item.rule
}

// RhsAt returns the right-hand side at the given index.
//
// Parameters:
//   - idx: The index of the right-hand side to retrieve.
//
// Returns:
//   - string: The right-hand side at the given index.
//   - bool: True if the retrieval was successful, false otherwise.
func (item Item) RhsAt(idx uint) (string, bool) {
	rule := item.rule

	rhs, ok := rule.RhsAt(idx)
	return rhs, ok
}

// IndicesOf returns all the indices of a target symbol in the right-hand side.
//
// Parameters:
//   - target: The target symbol to search for.
//
// Returns:
//   - []uint: The indices of the target symbol in the right-hand side or nil if target is not found.
func (item Item) IndicesOf(target string) []uint {
	rule := item.rule

	indices := rule.IndicesOf(target)
	return indices
}

/////////////////////////////////////////////////////////

func (item Item) BackwardRhs() iter.Seq[string] {
	rhss := item.rule.Rhss()

	slices.Reverse(rhss)

	fn := func(yield func(string) bool) {
		for _, rhs := range rhss {
			ok := yield(rhs)
			if !ok {
				break
			}
		}
	}

	return fn
}

func (item Item) Lhs() string {
	// assert.Cond(item.rule != nil, "item.rule must not be nil")

	return item.rule.Lhs()
}

func (item Item) Pos() uint {
	return item.pos
}

func (item *Item) SetLookaheads(lookaheads *sets.OrderedSet[string]) error {
	if item == nil {
		return common.ErrNilReceiver
	}

	item.lookaheads = lookaheads

	return nil
}

func (item Item) HasLookahead(la string) bool {
	ok := item.lookaheads.Has(la)
	return ok
}

func (item Item) ExpectLookahead() bool {
	ok := item.lookaheads.IsEmpty()
	return !ok
}

// NextRhs returns the symbol after the current position in the RHS.
//
// Returns:
//   - string: The symbol after the current position in the RHS.
//   - bool: True if the retrieval was successful, false otherwise.
func (item Item) NextRhs() (string, bool) {
	// assert.Cond(item.rule != nil, "item.rule must not be nil")

	return item.rule.RhsAt(item.pos)
}
