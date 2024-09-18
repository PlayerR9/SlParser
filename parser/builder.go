package parser

import (
	"fmt"
	"slices"

	gr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/parser/internal"
	dba "github.com/PlayerR9/go-debug/assert"
)

// ParseFn is a function that parses a production.
//
// Parameters:
//   - parser: The parser. Assumed to be non-nil.
//   - top1: The top token of the production. Assumed to be non-nil.
//   - lookahead: The lookahead token.
//
// Returns:
//   - []*Item: The list of items.
//   - error: if an error occurred.
type ParseFn[T gr.TokenTyper] func(parser *Parser[T], top1 *gr.Token[T], lookahead *gr.Token[T]) ([]*Item[T], error)

// Builder is a parser builder.
type Builder[T gr.TokenTyper] struct {
	// table is the parser table.
	table map[T]ParseFn[T]
}

// get_rhs_with_offset is a helper function that returns the rhs with the given offset.
//
// Returns:
//   - []T: the rhs with the given offset.
//
// The returned list is sorted and unique.
func get_rhs_with_offset[T gr.TokenTyper](items []*Item[T], offset int) []T {
	var all_rhs []T

	for _, item := range items {
		if item == nil {
			continue
		}

		rhs, ok := item.RhsAt(item.pos - offset)
		if !ok {
			continue
		}

		pos, ok := slices.BinarySearch(all_rhs, rhs)
		if !ok {
			all_rhs = slices.Insert(all_rhs, pos, rhs)
		}
	}

	return all_rhs
}

func apply_items_filter[T gr.TokenTyper](sols *internal.SolWithLevel[*Item[T]], type_ T, offset int, items []*Item[T]) []*Item[T] {
	dba.AssertNotNil(sols, "sols")

	fn := func(item *Item[T]) bool {
		rhs, ok := item.RhsByOffset(offset)
		if !ok {
			sols.AddSolution(offset-1, item)
		}

		return ok && type_ == rhs
	}

	items = SliceFilter(items, fn)
	return items
}

// register_unambiguous registers a new parser function.
//
// An unambiguous rule is one that has only one possible outcome or if it can be determined
// by only popping values from the stack.
//
// Parameters:
//   - rhs: the right hand side of the production.
//   - items: the list of items. (Tthis is assumed to be longer than 1.)
//
// If the receiver is nil or 'items' is empty or all items are nil, then nothing is registered.
func register_unambiguous[T gr.TokenTyper](items []*Item[T]) ParseFn[T] {
	dba.Assert(len(items) > 1, "len(items) > 1")

	fn := func(parser *Parser[T], top1, lookahead *gr.Token[T]) ([]*Item[T], error) {
		items_left := make([]*Item[T], len(items))
		copy(items_left, items)

		var sols internal.SolWithLevel[*Item[T]]
		offset := 1

		prev := top1.Type
		var last_got *T

		for {
			top, ok := parser.Pop()
			if !ok {
				break
			}

			last_got = &top.Type

			items_left = apply_items_filter(&sols, top.Type, offset, items_left)
			if len(items_left) == 0 {
				break
			}

			offset++
			prev = top.Type
		}

		if len(items_left) > 0 {
			_ = apply_items_filter(&sols, T(-1), offset, items_left)
		}

		solutions := sols.Solutions()
		if len(solutions) > 0 {
			return solutions, nil
		}

		var expecteds []T

		for _, item := range items {
			rhs, ok := item.RhsAt(item.pos - (offset - 1))
			if !ok {
				continue
			}

			pos, ok := slices.BinarySearch(expecteds, rhs)
			if !ok {
				expecteds = slices.Insert(expecteds, pos, rhs)
			}
		}

		return nil, NewErrUnexpectedToken(expecteds, &prev, last_got)
	}

	return fn
}

// NewBuilder creates a new parser builder.
//
// Returns:
//   - Builder: the builder.
func NewBuilder[T gr.TokenTyper](is *ItemSet[T]) Builder[T] {
	table := make(map[T]ParseFn[T])

	if is == nil {
		return Builder[T]{
			table: table,
		}
	}

	is.init()

	fmt.Println(is.PrintTable())

	for lhs, items := range is.item_table {
		var fn ParseFn[T]

		switch len(items) {
		case 0:
			fn = func(_ *Parser[T], top1, _ *gr.Token[T]) ([]*Item[T], error) {
				return nil, fmt.Errorf("no rule for %q", top1.Type.String())
			}
		case 1:
			fn = func(parser *Parser[T], top1, lookahead *gr.Token[T]) ([]*Item[T], error) {
				return items, nil
			}
		default:
			fn = register_unambiguous(items)
		}

		table[lhs] = fn
	}

	return Builder[T]{
		table: table,
	}
}

// Register registers a new parser function.
//
// Parameters:
//   - rhs: the right hand side of the production.
//   - fn: the parser function.
//
// If the receiver or 'fn' are nil, then nothing is registered.
//
// Previous functions are overwritten.
func (b *Builder[T]) Register(rhs T, fn ParseFn[T]) {
	if b == nil || fn == nil {
		return
	}

	b.table[rhs] = fn
}

// Build builds the parser.
//
// Returns:
//   - Parser: the parser.
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

// Reset resets the builder.
//
// Does nothing if the receiver is nil.
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
