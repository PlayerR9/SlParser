package parser

import (
	"errors"
	"strconv"

	slgr "github.com/PlayerR9/SlParser/grammar"
	gslc "github.com/PlayerR9/SlParser/mygo-lib/slices"
	"github.com/PlayerR9/SlParser/parser/internal"
	assert "github.com/PlayerR9/go-verify"
	sets "github.com/PlayerR9/mygo-data/sets"
	"github.com/PlayerR9/mygo-lib/common"
)

// Builder is the builder for a Parser.
//
// An empty builder can be created with the `var b Builder` syntax or with the
// `b := new(Builder)` constructor.
type Builder struct {
	// rules is a list of rules.
	rules []internal.Rule
}

// Reset implements common.Resetter.
func (b *Builder) Reset() error {
	if b == nil {
		return common.ErrNilReceiver
	}

	if len(b.rules) == 0 {
		return nil
	}

	clear(b.rules)
	b.rules = nil

	return nil
}

// AddRule appends a new grammar rule to the Builder's list of rules.
//
// Parameters:
//   - lhs: The left-hand side of the rule.
//   - rhss: The right-hand side(s) of the rule. Must have at least one element.
//
// Returns:
//   - error: An error if the rule could not be added.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
//   - an error if the rule has no rhs.
func (b *Builder) AddRule(lhs string, rhss ...string) error {
	if b == nil {
		return common.ErrNilReceiver
	}

	if len(rhss) == 0 {
		return errors.New("rule must have at least one rhs")
	}

	rule := internal.NewRule(lhs, rhss)
	b.rules = append(b.rules, rule)

	return nil
}

/////////////////////////////////////////////////////////

func lookaheadsOf(table map[string][]*internal.Item, item *internal.Item) *sets.OrderedSet[string] {
	lookaheads := new(sets.OrderedSet[string])

	next, ok := item.NextRhs()
	if !ok {
		return nil
	}

	ok = slgr.IsTerminal(next)

	if ok {
		_ = lookaheads.Insert(next)

		return lookaheads
	}

	seen := make(map[string]interface{})

	stack := []string{next}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		ct := (ConflictTable)(table)

		others := ct.ItemsByLhs(top)

		seen[top] = struct{}{}

		for _, other := range others {
			rhs, _ := other.RhsAt(0)
			// assert.True(ok, "other.RhsAt(0)")

			if slgr.IsTerminal(rhs) {
				_ = lookaheads.Insert(rhs)
			} else {
				_, ok := seen[rhs]
				if !ok {
					stack = append(stack, rhs)

					seen[rhs] = struct{}{}
				}
			}
		}
	}

	return lookaheads
}

func (b Builder) determineSymbols() *sets.OrderedSet[string] {
	if len(b.rules) == 0 {
		return nil
	}

	symbols := new(sets.OrderedSet[string])

	for _, rule := range b.rules {
		tmp := rule.Symbols()

		for _, symbol := range tmp.Slice() {
			err := symbols.Insert(symbol)
			assert.Err(err, "symbols.Insert(%s)", strconv.Quote(symbol))
		}
	}

	return symbols
}

func (b Builder) Build() map[string][]*internal.Item {
	symbols := b.determineSymbols()

	table := make(map[string][]*internal.Item, symbols.Size())

	var builder gslc.Builder[*internal.Item]

	for _, symbol := range symbols.Slice() {
		for _, rule := range b.rules {
			indices := rule.IndicesOf(symbol)

			for _, idx := range indices {
				item, err := internal.NewItem(&rule, idx+1)
				assert.Err(err, "internal.NewItem(rule, %d)", idx+1)
				assert.Cond(item != nil, "item != nil")

				err = builder.Append(item)
				assert.Err(err, "builder.Append(item)")
			}
		}

		items := builder.Build()
		table[symbol] = items

		err := builder.Reset()
		assert.Err(err, "builder.Reset()")
	}

	for _, items := range table {
		for _, item := range items {
			las := lookaheadsOf(table, item)

			_ = item.SetLookaheads(las)
			// assert.Err(err, "item.SetLookaheads(las)")
		}
	}

	// for _, items := range table {
	// 	for _, item := range items {
	// 		fmt.Println(item.String())
	// 	}

	// 	fmt.Println()
	// }

	// _, err := is.WriteItems(os.Stdout)
	// assert.Err(err, "is.WriteItems(os.Stdout)")

	return table
}
