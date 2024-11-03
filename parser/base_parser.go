package parser

import (
	"fmt"
	"sync"

	slgr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/parser/internal"
)

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
func (p *baseParser) Parse(tokens []*slgr.Token) *Iterator {
	if p == nil || len(tokens) == 0 {
		return nil
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	init_fn := func() (*Active, error) {
		slice := make([]*slgr.Token, 0, len(tokens))

		for _, tk := range tokens {
			tk_copy := &slgr.Token{
				Type: tk.Type,
				Data: tk.Data,
				Pos:  tk.Pos,
			}

			slice = append(slice, tk_copy)
		}

		if len(slice) >= 2 {
			for i, tk := range slice[:len(slice)-1] {
				tk.Lookahead = slice[i+1]
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
