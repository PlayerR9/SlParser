package internal

import (
	"strconv"

	assert "github.com/PlayerR9/go-verify"
	"github.com/PlayerR9/mygo-data/sets"
)

// Rule is a grammar rule.
type Rule struct {
	// lhs is the left-hand side of the rule.
	lhs string

	// rhss is the right-hand side of the rule.
	rhss []string
}

// NewRule creates a new Rule.
//
// Parameters:
//   - lhs: The left-hand side of the rule.
//   - rhss: The right-hand side of the rule.
//
// Returns:
//   - Rule: The new Rule.
func NewRule(lhs string, rhss []string) Rule {
	return Rule{
		lhs:  lhs,
		rhss: rhss,
	}
}

// Size returns the number of right-hand sides in the rule.
//
// Returns:
//   - uint: The number of right-hand sides in the rule.
func (r Rule) Size() uint {
	return uint(len(r.rhss))
}

// Lhs returns the left-hand side of the rule.
//
// Returns:
//   - string: The left-hand side of the rule.
func (r Rule) Lhs() string {
	return r.lhs
}

// RhsAt returns the right-hand side at the given index.
//
// Parameters:
//   - idx: The index of the right-hand side to retrieve.
//
// Returns:
//   - string: The right-hand side at the given index.
//   - bool: True if the retrieval was successful, false otherwise.
func (r Rule) RhsAt(idx uint) (string, bool) {
	if idx >= uint(len(r.rhss)) {
		return "", false
	}

	rhs := r.rhss[idx]

	return rhs, true
}

// Symbols returns an ordered set of symbols which are used in the rule.
//
// The ordered set contains the left-hand side and all the right-hand sides.
//
// Returns:
//   - *sets.OrderedSet[string]: The ordered set of symbols. Never returns nil.
func (r Rule) Symbols() *sets.OrderedSet[string] {
	symbols := new(sets.OrderedSet[string])

	err := symbols.Insert(r.lhs)
	assert.Err(err, "symbols.Insert(%s)", strconv.Quote(r.lhs))

	for _, rhs := range r.rhss {
		err := symbols.Insert(rhs)
		assert.Err(err, "symbols.Insert(%s)", strconv.Quote(rhs))
	}

	return symbols
}

// IndicesOf returns all the indices of a target symbol in the right-hand side.
//
// Parameters:
//   - target: The target symbol to search for.
//
// Returns:
//   - []uint: The indices of the target symbol in the right-hand side or nil if target is not found.
func (r Rule) IndicesOf(target string) []uint {
	if len(r.rhss) == 0 {
		return nil
	}

	var count uint

	for _, rhs := range r.rhss {
		if rhs == target {
			count++
		}
	}

	if count == 0 {
		return nil
	}

	slice := make([]uint, 0, count)

	for i, rhs := range r.rhss {
		if rhs == target {
			slice = append(slice, uint(i))
		}
	}

	return slice
}

// Rhss returns a copy of the right-hand side of the rule.
//
// Returns:
//   - []string: A copy of the right-hand side of the rule.
func (r Rule) Rhss() []string {
	if len(r.rhss) == 0 {
		return nil
	}

	slice := make([]string, len(r.rhss))
	copy(slice, r.rhss)

	return slice
}
