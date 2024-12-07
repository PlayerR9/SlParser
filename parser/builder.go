package parser

import (
	slgr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/parser/internal"
	sets "github.com/PlayerR9/mygo-data/sets"
	"github.com/PlayerR9/mygo-lib/common"
	gslc "github.com/PlayerR9/mygo-lib/slices"
)

type Builder struct {
	rules []*internal.Rule
}

func (b *Builder) AddRule(lhs string, rhss ...string) error {
	if b == nil {
		return common.ErrNilReceiver
	}

	rule := internal.NewRule(lhs, rhss)
	b.rules = append(b.rules, rule)

	return nil
}

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

		for _, symbol := range tmp {
			_ = symbols.Insert(symbol)
		}
	}

	return symbols
}

func (b Builder) Build() map[string][]*internal.Item {
	symbols := b.determineSymbols()

	table := make(map[string][]*internal.Item)

	table := make(map[string][]*internal.Item, len(all_symbols))

	var builder gslc.Builder[*internal.Item]

	for _, symbol := range all_symbols {
		for _, rule := range b.rules {
			indices := rule.IndicesOf(symbol)

			for _, idx := range indices {
				item, _ := internal.NewItem(rule, idx+1)
				// assert.Err(err, "internal.NewItem(rule, %d)", idx+1)
				// assert.Condf(item != nil, "item must not be nil")

				_ = builder.Append(item)
				// assert.Err(err, "builder.Append(item)")
			}
		}

		// assert.Cond(len(items) > 0, "len(items) must be greater than 0")

		table[symbol] = builder.Build()
		builder.Reset()
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

func (b *Builder) Reset() {
	if b == nil {
		return
	}

	if len(b.rules) > 0 {
		for k := range b.rules {
			b.rules[k] = nil
		}

		b.rules = b.rules[:0]
	}
}
