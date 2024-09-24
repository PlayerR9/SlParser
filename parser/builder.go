package parser

import (
	"fmt"
	"slices"

	gr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/parser/internal"
	bck "github.com/PlayerR9/go-commons/backup"
	gcslc "github.com/PlayerR9/go-commons/slices"
	gers "github.com/PlayerR9/go-errors"
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
type ParseFn[T gr.TokenTyper] func(parser *ActiveParser[T], top1 *gr.ParseTree[T], lookahead *gr.Token[T]) ([]*internal.Item[T], error)

// get_rhs_with_offset is a helper function that returns the rhs with the given offset.
//
// Returns:
//   - []T: the rhs with the given offset.
//
// The returned list is sorted and unique.
func get_rhs_with_offset[T gr.TokenTyper](items []*internal.Item[T], offset int) []T {
	var all_rhs []T

	for _, item := range items {
		if item == nil {
			continue
		}

		rhs, ok := item.RhsAt(item.Pos - offset)
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

func apply_items_filter[T gr.TokenTyper](sols *gcslc.SolWithLevel[*internal.Item[T]], type_ T, offset int, items []*internal.Item[T]) []*internal.Item[T] {
	gers.AssertNotNil(sols, "sols")

	fn := func(item *internal.Item[T]) bool {
		rhs, ok := item.RhsByOffset(offset)
		if !ok {
			sols.AddSolution(offset-1, item)
		}

		return ok && type_ == rhs
	}

	items = gcslc.FilterSlice(items, fn)
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
func register_unambiguous[T gr.TokenTyper](items []*internal.Item[T]) ParseFn[T] {
	gers.Assert(len(items) > 1, "len(items) > 1")

	fn := func(parser *ActiveParser[T], top1 *gr.ParseTree[T], lookahead *gr.Token[T]) ([]*internal.Item[T], error) {
		items_left := make([]*internal.Item[T], len(items))
		copy(items_left, items)

		var sols gcslc.SolWithLevel[*internal.Item[T]]
		offset := 1

		prev := top1.Type()
		var last_got *T

		for {
			top, ok := parser.Pop()
			if !ok {
				break
			}

			last_type := top.Type()

			last_got = &last_type

			items_left = apply_items_filter(&sols, top.Type(), offset, items_left)
			if len(items_left) == 0 {
				break
			}

			offset++
			prev = top.Type()
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
			rhs, ok := item.RhsAt(item.Pos - (offset - 1))
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
func Build[T gr.TokenTyper](is *ItemSet[T]) *Parser[T] {
	p := &Parser[T]{
		table: make(map[T]ParseFn[T]),
	}

	if is != nil {
		is.init()

		for lhs, items := range is.item_table {
			var fn ParseFn[T]

			switch len(items) {
			case 0:
				fn = func(_ *ActiveParser[T], top1 *gr.ParseTree[T], _ *gr.Token[T]) ([]*internal.Item[T], error) {
					return nil, fmt.Errorf("no rule for %q", top1.Type().String())
				}
			case 1:
				fn = func(parser *ActiveParser[T], top1 *gr.ParseTree[T], lookahead *gr.Token[T]) ([]*internal.Item[T], error) {
					return items, nil
				}
			default:
				fn = register_unambiguous(items)
			}

			p.table[lhs] = fn
		}
	}

	fn := func() *ActiveParser[T] {
		ap, err := NewActiveParser(p)
		gers.AssertErr(err, "NewActiveParser(p)")

		return ap
	}

	p.seq = bck.Subject(fn)

	return p
}
