package parser

import (
	"slices"

	gr "github.com/PlayerR9/SlParser/grammar"
	gcers "github.com/PlayerR9/go-commons/errors"
	gcslc "github.com/PlayerR9/go-commons/slices"
)

// CheckTop is a function that checks if the top of the stack is in the allowed list.
//
// Parameters:
//   - parser: The parser to check.
//   - allowed: The list of allowed tokens.
//
// Returns:
//   - *gr.Token[T]: The top of the stack.
//   - bool: True if the top of the stack is in the allowed list, false otherwise.
//
// If the receiver is nil, then it returns nil and false.
//
// If no allowed tokens are provided, then it returns the top of the stack and false.
func CheckTop[T gr.TokenTyper](parser *Parser[T], allowed ...T) (*gr.Token[T], bool) {
	if parser == nil {
		return nil, false
	}

	top, ok := parser.Pop()
	if !ok || len(allowed) == 0 {
		return top, false
	}

	for _, a := range allowed {
		if top.Type == a {
			return top, true
		}
	}

	return top, false
}

// CheckLookahead is a function that checks if the lookahead is in the allowed list.
//
// Parameters:
//   - lookahead: The lookahead to check.
//   - allowed: The list of allowed tokens.
//
// Returns:
//   - bool: True if the lookahead is in the allowed list, false otherwise.
//
// If the receiver is nil or no allowed tokens are provided, then it returns false.
func CheckLookahead[T gr.TokenTyper](lookahead *gr.Token[T], allowed ...T) bool {
	if lookahead == nil || len(allowed) == 0 {
		return false
	}

	for _, a := range allowed {
		if lookahead.Type == a {
			return true
		}
	}

	return false
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

// filter_items_without_rhs returns the items without the rhs with the given offset.
//
// Parameters:
//   - items: the items to filter.
//   - offset: the offset of the rhs to filter.
//   - rhs: the rhs to filter.
//
// Returns:
//   - []*Item[T]: the items without the rhs with the given offset.
func filter_items_without_rhs[T gr.TokenTyper](items []*Item[T], offset int, rhs T) []*Item[T] {
	fn := func(item *Item[T]) bool {
		tmp, ok := item.RhsAt(item.pos - offset)
		return !ok || tmp == rhs
	}

	items = gcslc.SliceFilter(items, fn)

	return items
}

// UnambiguousRule is a function that creates a function that parses an unambiguous rule.
//
// An unambiguous rule is one that has only one possible outcome or if it can be determined
// by only popping values from the stack.
//
// Parameters:
//   - items: The items to parse.
//
// Returns:
//   - ParseFn: The parse function.
//
// Will panic if items is empty or all items are nil.
func UnambiguousRule[T gr.TokenTyper](items ...*Item[T]) ParseFn[T] {
	items = gcslc.FilterNilValues(items)
	switch len(items) {
	case 0:
		return func(parser *Parser[T], top1, lookahead *gr.Token[T]) ([]*Item[T], error) {
			panic(gcers.NewErrInvalidParameter("items", gcers.NewErrEmpty(items)))
		}
	case 1:
		return func(parser *Parser[T], top1, lookahead *gr.Token[T]) ([]*Item[T], error) {
			return items, nil
		}
	default:
		max := len(items) // Ensure no infinite loop occurs

		return func(parser *Parser[T], top1, lookahead *gr.Token[T]) ([]*Item[T], error) {
			prev := top1

			for offset := 0; offset < max && len(items) > 1; offset++ {
				all_rhs := get_rhs_with_offset(items, offset)
				if len(all_rhs) == 0 {
					break
				}

				top, ok := parser.Pop()
				if !ok {
					return nil, NewErrUnexpectedToken(all_rhs, &prev.Type, nil)
				}

				items := filter_items_without_rhs(items, offset, top.Type)
				if len(items) == 0 {
					return nil, NewErrUnexpectedToken(all_rhs, &prev.Type, &top.Type)
				}

				prev = top
			}

			return items, nil
		}
	}
}
