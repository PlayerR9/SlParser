package parser

import (
	"fmt"
	"sync"

	"github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/parser/internal"
	tr "github.com/PlayerR9/mygo-lib/CustomData/tree"
)

/////////////////////////////////////////////////////////

// baseParser is the base implementation of the Parser interface.
type baseParser struct {
	// decision_table is the decision table of the parser. This must be
	// read-only.
	decision_table map[string][]*internal.Item

	// mu protects the fields above.
	mu sync.RWMutex
}

// ItemsOf implements the Parser interface.
func (p *baseParser) ItemsOf(type_ string) ([]*internal.Item, bool) {
	if p == nil {
		return nil, false
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	items, ok := p.decision_table[type_]
	return items, ok
}

// Parse implements the Parser interface.
func (p *baseParser) Parse(tokens []*tr.Node) *Iterator {
	if p == nil || len(tokens) == 0 {
		return nil
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	init_fn := func() (*Active, error) {
		slice := make([]*tr.Node, 0, len(tokens))

		for _, tk := range tokens {
			tkd := grammar.MustGet[*grammar.TokenData](tk)

			tk_copy := grammar.NewToken(tkd.Pos, tkd.Data, tkd.Type, nil)
			slice = append(slice, tk_copy)
		}

		if len(slice) >= 2 {
			for i, tk := range slice[:len(slice)-1] {
				tkd := grammar.MustGet[*grammar.TokenData](tk)
				tkd.Lookahead = slice[i+1]
			}
		}

		active, err := NewActive(p, slice)
		return active, err
	}

	itr, err := NewIterator(init_fn)
	if err != nil {
		panic(fmt.Sprintf("failed to create iterator: %v", err))
	}

	return itr
}
