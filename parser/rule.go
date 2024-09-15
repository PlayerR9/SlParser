package parser

import (
	"errors"
	"iter"

	gr "github.com/PlayerR9/SlParser/grammar"
)

// A rule is a production of a grammar.
type Rule[T gr.TokenTyper] struct {
	// lhs is the left hand side of the rule.
	lhs T

	// rhss are the right hand sides of the rule.
	rhss []T
}

// NewRule creates a new rule.
//
// Parameters:
//   - lhs: the left hand side of the rule.
//   - rhss: the right hand sides of the rule.
//
// Returns:
//   - *Rule: the new rule.
//   - error: if the rule does not have at least one right hand side.
func NewRule[T gr.TokenTyper](lhs T, rhss ...T) (*Rule[T], error) {
	if len(rhss) == 0 {
		return nil, errors.New("at least one right hand side is required")
	}

	return &Rule[T]{
		lhs:  lhs,
		rhss: rhss,
	}, nil
}

func (r Rule[T]) size() int {
	return len(r.rhss)
}

func (r Rule[T]) rhs_at(pos int) (T, bool) {
	if pos < 0 || pos >= len(r.rhss) {
		return T(0), false
	}

	return r.rhss[pos], true
}

func (r Rule[T]) backward_rhs() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := len(r.rhss) - 1; i >= 0; i-- {
			if !yield(r.rhss[i]) {
				break
			}
		}
	}
}

func (r Rule[T]) get_lhs() T {
	return r.lhs
}
