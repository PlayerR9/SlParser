package internal

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	slgr "github.com/PlayerR9/SlParser/grammar"
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

/////////////////////////////////////////////////////////

func (item Item) String() string {
	lhs := item.rule.Lhs()
	rhss := item.rule.Rhss()
	act := item.act.String()
	lookaheads := collections.String("*sets.OrderedSet[string]", item.lookaheads)

	elems := append([]string{lhs, "="}, rhss...)
	elems = append(elems, ":", "(", act, ")", "-->", lookaheads)

	elems = slices.Insert(elems, int(item.pos)+1, "#")

	str := strings.Join(elems, " ")
	return str
}

func NewItem(rule *Rule, pos uint) (*Item, error) {
	if rule == nil {
		return nil, common.NewErrNilParam("rule")
	}

	size := rule.Size()
	if pos < 0 || pos > size {
		return nil, common.NewErrBadParam("pos", fmt.Sprintf("be in range [0, %d]", size))
	}

	var act ActionType

	if pos == size {
		rhs, ok := rule.RhsAt(pos - 1)
		if !ok {
			return nil, fmt.Errorf("failed to get rhs at %d", pos-1)
		}

		if rhs == slgr.EtEOF {
			act = ActAccept
		} else {
			act = ActReduce
		}
	} else {
		act = ActShift
	}

	return &Item{
		rule:       rule,
		pos:        pos,
		act:        act,
		lookaheads: nil,
	}, nil
}

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

func (item Item) RhsAt(idx uint) (string, bool) {
	// assert.Cond(item.rule != nil, "item.rule must not be nil")

	return item.rule.RhsAt(idx)
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

func (item Item) Action() ActionType {
	return item.act
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

func (item Item) IndicesOf(symbol string) []uint {
	// assert.Cond(item.rule != nil, "item.rule must not be nil")

	return item.rule.IndicesOf(symbol)
}
