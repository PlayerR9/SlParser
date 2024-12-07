package parser

import (
	gr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/parser/internal"
	"github.com/PlayerR9/mygo-lib/common"
)

type ConflictTable map[string][]*internal.Item

func (c ConflictTable) Lines() []string {
	if len(c) == 0 {
		return nil
	}

	var lines []string

	for symbol, items := range c {
		lines = append(lines, "", symbol+":")

		for _, item := range items {
			lines = append(lines, item.String())
		}
	}

	return lines[1:]
}

func (c ConflictTable) ItemsByLhs(lhs string) []*internal.Item {
	if len(c) == 0 {
		return nil
	}

	var result []*internal.Item

	for _, items := range c {
		for _, item := range items {
			ok := item.Lhs() == lhs
			if ok {
				result = append(result, item)
			}
		}
	}

	return result
}

func (c ConflictTable) DetermineLookaheadsOf(item *internal.Item) (*sets.OrderedSet[string], error) {
	if item == nil {
		return nil, common.NewErrNilParam("item")
	}

	lookahead, ok := item.GetRhsByOffset(0)
	if !ok {
		return nil, nil
	}

	symbols := new(sets.OrderedSet[string])

	ok = gr.IsTerminal(lookahead)
	if ok {
		_ = symbols.Insert(lookahead)

		return symbols, nil
	}

	stack := new(lls.ArrayStack[string])
	_ = stack.Push(lookahead)

	seen := new(sets.SeenSet[*Item])
	defer seen.Reset()

	for {
		top, err := stack.Pop()
		if err != nil {
			break
		}

		nexts := c.ItemsByLhs(top)

		for _, next := range nexts {
			ok := seen.Has(next)
			if ok {
				continue
			}

			_ = seen.Insert(next)

			la, ok := next.GetRhsAt(0)
			assert.Cond(ok, "next.GetRhsAt(0) = false")

			ok = gr.IsTerminal(la)
			if ok {
				_ = symbols.Insert(la)
			} else {
				_ = stack.Push(la)
			}
		}
	}

	return symbols, nil
}

func (c ConflictTable) Solve() {
	if len(c) == 0 {
		return
	}

	for symbol, items := range c {
		for _, item := range items {
			lookaheads, _ := c.DetermineLookaheadsOf(item)
			_ = item.SetLookaheads(lookaheads)
		}

		c[symbol] = items
	}
}
