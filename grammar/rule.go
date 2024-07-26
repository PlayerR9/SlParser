package grammar

import (
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// Rule is a rule of the grammar.
type Rule[T TokenTyper] struct {
	// lhs is the left-hand side of the rule.
	lhs T

	// rhss are the right-hand sides of the rule.
	rhss []T
}

// String implements the fmt.Stringer interface.
//
// Format:
//
//	RHS(n) RHS(n-1) ... RHS(1) -> LHS .
func (r *Rule[T]) String() string {
	var values []string

	for _, rhs := range r.rhss {
		values = append(values, rhs.GoString())
	}

	values = append(values, "->")
	values = append(values, r.lhs.GoString())
	values = append(values, ".")

	return strings.Join(values, " ")
}

// Iterator implements the common.Iterater interface.
func (r *Rule[T]) Iterator() uc.Iterater[T] {
	return uc.NewSimpleIterator(r.rhss)
}

// NewRule is a constructor for a Rule.
//
// Parameters:
//   - lhs: The left-hand side of the rule.
//   - rhss: The right-hand sides of the rule.
//
// Returns:
//   - *Rule: The created rule.
//
// Returns nil iff the rhss is empty.
func NewRule[T TokenTyper](lhs T, rhss []T) *Rule[T] {
	if len(rhss) == 0 {
		return nil
	}

	return &Rule[T]{
		lhs:  lhs,
		rhss: rhss,
	}
}

// GetLhs returns the left-hand side of the rule.
//
// Returns:
//   - T: The left-hand side of the rule.
func (r *Rule[T]) GetLhs() T {
	return r.lhs
}

// GetIndicesOfRhs returns the ocurrence indices of the rhs in the rule.
//
// Parameters:
//   - rhs: The right-hand side to search.
//
// Returns:
//   - []int: The indices of the rhs in the rule.
func (r *Rule[T]) GetIndicesOfRhs(rhs T) []int {
	var indices []int

	for i := 0; i < len(r.rhss); i++ {
		if r.rhss[i] == rhs {
			indices = append(indices, i)
		}
	}

	return indices
}

// GetRhss returns the right-hand sides of the rule.
//
// Returns:
//   - []T: The right-hand sides of the rule.
func (r *Rule[T]) GetRhss() []T {
	return r.rhss
}

// Size returns the number of right-hand sides of the rule.
//
// Returns:
//   - int: The "size" of the rule.
func (r *Rule[T]) Size() int {
	return len(r.rhss)
}
