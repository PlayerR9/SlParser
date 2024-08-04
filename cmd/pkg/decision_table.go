package pkg

import (
	"strings"

	lus "github.com/PlayerR9/lib_units/slices"
)

type DecisionTable struct {
	symbols []string
	rules   []*Rule

	table map[string][]*Item
}

func (dt *DecisionTable) String() string {
	var values []string

	for _, symbol := range dt.symbols {
		items, _ := dt.table[symbol]
		// luc.AssertF(ok, "symbol %q not found in table", symbol)

		elems := make([]string, 0, len(items))

		for _, item := range items {
			elems = append(elems, item.String())
		}

		values = append(values, strings.Join(elems, "\n"))
	}

	return strings.Join(values, "\n\n")
}

// make_symbols makes the symbols of the grammar.
//
// This function returns only unique symbols and sorts them.
func (dt *DecisionTable) make_symbols() {
	var all_symbols []string

	for _, rule := range dt.rules {
		symbols := rule.GetSymbols()
		all_symbols = append(all_symbols, symbols...)
	}

	dt.symbols = lus.OrderedUniquefy(all_symbols)
}

// make_items_per_symbol makes the items of the grammar per symbol.
//
// Parameters:
//   - symbol: The symbol to make the items per.
func (dt *DecisionTable) make_items_per_symbol(symbol string) {
	index_map := make(map[int][]int)

	for i := 0; i < len(dt.rules); i++ {
		rule := dt.rules[i]

		indices := rule.GetIndicesOfRhs(symbol)

		if len(indices) > 0 {
			index_map[i] = indices
		}
	}

	if len(index_map) == 0 {
		dt.table[symbol] = []*Item{}

		return
	}

	var items []*Item

	for key, indices := range index_map {
		for _, idx := range indices {
			act := DetermineAction(idx, symbol)

			item, _ := NewItem(dt.rules[key], idx, act)
			// luc.AssertErr(err, "NewItem(rule, %d, %q)", idx, act)

			items = append(items, item)
		}
	}

	dt.table[symbol] = items
}

// make_items makes the items of the grammar.
func (dt *DecisionTable) make_items() {
	dt.table = make(map[string][]*Item)

	for _, symbol := range dt.symbols {
		dt.make_items_per_symbol(symbol)
	}
}

// NewDecisionTable is a constructor for a DecisionTable.
//
// Parameters:
//   - rules: The rules of the grammar.
//
// Returns:
//   - *DecisionTable: The new DecisionTable. Nil iff the rules are empty.
func NewDecisionTable(rules []*Rule) *DecisionTable {
	if len(rules) == 0 {
		return nil
	}

	dt := &DecisionTable{
		rules: rules,
	}

	dt.make_symbols()
	dt.make_items()

	return dt
}
