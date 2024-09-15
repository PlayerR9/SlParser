package parser

import (
	gr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/parser/internal"
)

type ParseFn[T gr.TokenTyper] func(parser *Parser[T], top1 *gr.Token[T], lookahead *gr.Token[T]) ([]*Item[T], error)

type Builder[T gr.TokenTyper] struct {
	table map[T]ParseFn[T]
}

func NewBuilder[T gr.TokenTyper]() Builder[T] {
	return Builder[T]{
		table: make(map[T]ParseFn[T]),
	}
}

func (b *Builder[T]) Register(rhs T, fn ParseFn[T]) {
	if b == nil || fn == nil {
		return
	}

	b.table[rhs] = fn
}

func (b Builder[T]) Build() Parser[T] {
	var table map[T]ParseFn[T]

	if len(b.table) > 0 {
		table = make(map[T]ParseFn[T], len(b.table))
		for k, v := range b.table {
			table[k] = v
		}
	}

	var stack internal.Stack[T]

	return Parser[T]{
		table: table,
		stack: &stack,
	}
}

func (b *Builder[T]) Reset() {
	if b == nil {
		return
	}

	if len(b.table) > 0 {
		for k := range b.table {
			b.table[k] = nil
			delete(b.table, k)
		}

		b.table = make(map[T]ParseFn[T])
	}
}
