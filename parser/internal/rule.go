package internal

import (
	"iter"

	gr "github.com/PlayerR9/SlParser/grammar"
	gcers "github.com/PlayerR9/errors"
)

// A rule is a production of a grammar.
type Rule[T gr.TokenTyper] struct {
	// lhs is the left hand side of the rule.
	lhs T

	// rhss are the right hand sides of the rule.
	rhss []T
}

func NewRule[T gr.TokenTyper](lhs T, rhss []T) (*Rule[T], error) {
	if len(rhss) == 0 {
		return nil, gcers.NewErrInvalidParameter("rhss must have at least one element")
	}

	return &Rule[T]{
		lhs:  lhs,
		rhss: rhss,
	}, nil
}

// Size returns the number of right hand sides of the rule.
//
// Returns:
//   - int: The size of the rule.
func (r Rule[T]) Size() int {
	return len(r.rhss)
}

func (r Rule[T]) HasRhsAt(pos int) bool {
	return pos >= 0 && pos < len(r.rhss)
}

// RhsAt returns the right hand side at the given position.
//
// Parameters:
//   - pos: The position of the right hand side.
//
// Returns:
//   - T: The right hand side at the given position.
//   - bool: True if the position is valid, false otherwise.
func (r Rule[T]) RhsAt(pos int) (T, bool) {
	if pos < 0 || pos >= len(r.rhss) {
		return T(0), false
	}

	return r.rhss[pos], true
}

// BackwardRhs returns the backward rhs of the rule.
//
// Returns:
//   - iter.Seq[T]: the backward rhs of the rule.
func (r Rule[T]) BackwardRhs() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := len(r.rhss) - 1; i >= 0; i-- {
			if !yield(r.rhss[i]) {
				break
			}
		}
	}
}

// ForwardRhs returns the forward rhs of the rule.
//
// Returns:
//   - iter.Seq[T]: the forward rhs of the rule.
func (r Rule[T]) ForwardRhs() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, rhs := range r.rhss {
			if !yield(rhs) {
				break
			}
		}
	}
}

// Lhs returns the left hand side of the rule.
//
// Returns:
//   - T: The left hand side of the rule.
func (r Rule[T]) Lhs() T {
	return r.lhs
}

// IndicesOf returns the ocurrence indices of the rhs in the rule.
//
// Parameters:
//   - rhs: The right-hand side to search.
//
// Returns:
//   - []int: The indices of the rhs in the rule.
func (r Rule[T]) IndicesOf(rhs T) []int {
	var indices []int

	for i, r := range r.rhss {
		if r == rhs {
			indices = append(indices, i)
		}
	}

	return indices
}
