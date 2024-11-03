package internal

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	slgr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/mygo-lib/common"
	gslc "github.com/PlayerR9/mygo-lib/slices"
)

type ActionType int

const (
	ActAccept ActionType = iota // (ACCEPT)
	ActReduce                   // (REDUCE)
	ActShift                    // (SHIFT)
)

type Item struct {
	rule       *Rule
	pos        int
	act        ActionType
	lookaheads []string
}

func (item Item) String() string {
	// assert.Cond(item.rule != nil, "item.rule must not be nil")

	var builder strings.Builder

	builder.WriteString(item.rule.Lhs())
	builder.WriteString(" -> ")

	var i int

	for rhs := range item.rule.Rhs() {
		if i == item.pos {
			builder.WriteString("# ")
		}

		builder.WriteString(rhs + " ")

		i++
	}

	if i == item.pos {
		builder.WriteString("# ")
	}

	builder.WriteString(": ")
	builder.WriteString(item.act.String())
	builder.WriteString(" -- ")
	builder.WriteString(strings.Join(item.lookaheads, ", "))

	return builder.String()
}

func NewItem(rule *Rule, pos int) (*Item, error) {
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
	// assert.Cond(item.rule != nil, "item.rule must not be nil")

	return item.rule.BackwardRhs()
}

func (item Item) Lhs() string {
	// assert.Cond(item.rule != nil, "item.rule must not be nil")

	return item.rule.Lhs()
}

func (item Item) RhsAt(idx int) (string, bool) {
	// assert.Cond(item.rule != nil, "item.rule must not be nil")

	return item.rule.RhsAt(idx)
}

func (item Item) Pos() int {
	return item.pos
}

func (item *Item) SetLookaheads(lookaheads []string) error {
	if lookaheads == nil {
		return nil
	} else if item == nil {
		return common.ErrNilReceiver
	}

	_, _ = gslc.Merge(&item.lookaheads, lookaheads)

	return nil
}

func (item Item) HasLookahead(la string) bool {
	_, ok := slices.BinarySearch(item.lookaheads, la)
	return ok
}

func (item Item) ExpectLookahead() bool {
	return len(item.lookaheads) != 0
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

func (item Item) IndicesOf(symbol string) []int {
	// assert.Cond(item.rule != nil, "item.rule must not be nil")

	return item.rule.IndicesOf(symbol)
}
