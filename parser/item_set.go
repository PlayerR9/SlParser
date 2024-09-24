package parser

import (
	"slices"

	gr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/parser/internal"
	bck "github.com/PlayerR9/go-commons/backup"
	gcslc "github.com/PlayerR9/go-commons/slices"
	gcers "github.com/PlayerR9/go-errors"
	gcmap "github.com/PlayerR9/go-sets"
)

// ItemSet is an item set.
type ItemSet[T gr.TokenTyper] struct {
	// symbols is the list of symbols.
	symbols []T

	// rules is the list of rules.
	rules []*internal.Rule[T]

	// item_table is the item table.
	item_table map[T][]*internal.Item[T]
}

// NewItemSet creates a new item set.
//
// Returns:
//   - ItemSet: the new item set.
func NewItemSet[T gr.TokenTyper]() ItemSet[T] {
	return ItemSet[T]{
		item_table: make(map[T][]*internal.Item[T]),
	}
}

// PrintTable prints the item table of the item set.
//
// Each item is printed on a new line, with the format:
//
//	<item_string>
//
// The strings are separated by an empty line.
//
// Returns:
//   - []string: the lines of the item set.
func (is ItemSet[T]) PrintTable() []string {
	var lines []string

	for _, items := range is.item_table {
		for _, item := range items {
			lines = append(lines, item.String())
		}

		lines = append(lines, "")
	}

	return lines
}

// AddRule adds a new rule to the item set.
//
// Parameters:
//   - lhs: the left hand side of the rule.
//   - rhss: the right hand sides of the rule.
//
// Returns:
//   - error: if the rule could not be added.
//
// Errors:
//   - gcers.NilReceiver: if the receiver is nil.
//   - gcers.ErrInvalidParameter: if the rule does not have at least one right hand side.
func (is *ItemSet[T]) AddRule(lhs T, rhss ...T) error {
	if is == nil {
		return nil
	} else if len(rhss) == 0 {
		return gcers.NewErrInvalidParameter("rhss must have at least one element")
	}

	pos, ok := slices.BinarySearch(is.symbols, lhs)
	if !ok {
		is.symbols = slices.Insert(is.symbols, pos, lhs)
	}

	for _, rhs := range rhss {
		pos, ok := slices.BinarySearch(is.symbols, rhs)
		if !ok {
			is.symbols = slices.Insert(is.symbols, pos, rhs)
		}
	}

	rule, _ := internal.NewRule(lhs, rhss)
	is.rules = append(is.rules, rule)

	return nil
}

func (is *ItemSet[T]) make_items() {
	if is == nil {
		return
	}

	var item_list gcslc.Builder[*internal.Item[T]]

	for _, rhs := range is.symbols {
		for _, rule := range is.rules {
			indices := rule.IndicesOf(rhs)
			if len(indices) == 0 {
				continue
			}

			for _, idx := range indices {
				item, err := internal.NewItem(rule, idx)
				gcers.AssertErr(err, "internal.NewItem(rule, %d)", idx)

				item_list.Append(item)
			}
		}

		is.item_table[rhs] = item_list.Build()
		item_list.Reset()
	}
}

// ItemsWithLhsOf returns the items with the given left-hand side.
//
// Parameters:
//   - lhs: The left-hand side to search.
//
// Returns:
//   - []*Item[T]: The items with the given left-hand side. Nil if there are no items with the given left-hand side.
func (is ItemSet[T]) ItemsWithLhsOf(lhs T) []*internal.Item[T] {
	items, ok := is.item_table[lhs]
	if !ok || len(items) == 0 {
		return nil
	}

	return items
}

func get_lookahead_of[T gr.TokenTyper](is *ItemSet[T], item *internal.Item[T], seen *gcmap.SeenSet[*internal.Item[T]]) []T {
	seen.SetSeen(item)

	stack := []*internal.Item[T]{item}

	var lookaheads []T

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		rhs, ok := top.RhsAt(top.Pos + 1)
		if !ok {
			continue
		}

		if rhs.IsTerminal() {
			pos, ok := slices.BinarySearch(lookaheads, rhs)
			if ok {
				continue
			}

			lookaheads = slices.Insert(lookaheads, pos, rhs)

			continue
		}

		nexts := is.ItemsWithLhsOf(rhs)
		nexts = seen.FilterNotSeen(nexts)
		if len(nexts) == 0 {
			continue
		}
	}

	return lookaheads
}

func (is ItemSet[T]) make_lookahead() {
	seen := gcmap.NewSeenSet[*internal.Item[T]]()

	for _, items := range is.item_table {
		if len(items) == 0 {
			continue
		}

		for _, item := range items {
			lookaheads := get_lookahead_of(&is, item, seen)
			item.SetLookaheads(lookaheads)

			seen.Reset()
		}
	}
}

func (is *ItemSet[T]) init() {
	if is == nil {
		return
	}

	is.make_items()
	is.make_lookahead()
}

/*
func (is *ItemSet[T]) expected_of(symbol T, seen *SeenMap[*Item[T]]) (internal.Expecter, error) {
	if symbol.IsTerminal() {
		expected := internal.NewExpectTerminal(symbol)
		return expected, nil
	}

	nexts, ok := is.item_table[symbol]
	if !ok {
		return nil, errors.New("no nexts")
	}

	var sub_expecteds []internal.Expecter

	for _, next := range nexts {
		ok := seen.SetSeen(next)
		if !ok {
			continue
		}

		s, ok := next.RhsAt(0)
		dba.AssertOk(ok, "next.RhsAt(%d)", 0)

		expected, err := is.expected_of(s)
		if err != nil {
			return nil, fmt.Errorf("expected of %q: %w", s.String(), err)
		}

		sub_expecteds = append(sub_expecteds, expected)
	}

	expected := internal.NewExpectNonTerminal(symbol, sub_expecteds)

	return expected, nil
}

func (is *ItemSet[T]) determine_expecteds() {
	item := is.item_table[T(0)][0]

	var expected internal.Expecter

	next, ok := item.RhsByOffset(-1)
	if !ok {
		// TODO: Handle this case.
	} else {
		if next.IsTerminal() {
			expected = internal.NewExpectTerminal(next)
		} else {
			next_items, ok := is.item_table[next]
			if !ok {
				// TODO: Handle this case.
			} else {
				var sub_expecteds []internal.Expecter

				for _, ni := range next_items {
					s, ok := ni.RhsAt(0)
					dba.AssertOk(ok, "ni.RhsAt(%d)", 0)

				}

				// TODO: Handle this case.
			}
		}

	}

}
*/

// Build builds the parser.
//
// Returns:
//   - Parser: the parser. Never returns nil.
func (b ItemSet[T]) Build() *Parser[T] {
	p := &Parser[T]{}

	fn := func() *ActiveParser[T] {
		ap, err := NewActiveParser(p)
		gcers.AssertErr(err, "NewActiveParser(p)")

		return ap
	}

	p.seq = bck.Subject(fn)

	return p
}
