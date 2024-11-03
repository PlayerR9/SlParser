package ast

import (
	"iter"

	slgr "github.com/PlayerR9/SlParser/grammar"
)

type ToAstFunc[N interface {
	Child() iter.Seq[N]

	Noder
}] func(token *slgr.Token) ([]N, error)

type Builder[N interface {
	Child() iter.Seq[N]

	Noder
}] struct {
	table map[string]ToAstFunc[N]
}

func (b *Builder[N]) Register(tt string, fn ToAstFunc[N]) bool {
	if fn == nil {
		delete(b.table, tt)

		return true
	}

	if b == nil {
		return false
	}

	if b.table == nil {
		b.table = make(map[string]ToAstFunc[N])
	}

	b.table[tt] = fn

	return true
}

func (b Builder[N]) Build() AST[N] {
	if len(b.table) == 0 {
		return &baseAST[N]{
			table: make(map[string]ToAstFunc[N]),
		}
	}

	table := make(map[string]ToAstFunc[N], len(b.table))

	for k, v := range b.table {
		table[k] = v
	}

	return &baseAST[N]{
		table: table,
	}
}

func (b *Builder[N]) Reset() {
	if b == nil {
		return
	}

	if len(b.table) > 0 {
		clear(b.table)
		b.table = nil
	}
}
