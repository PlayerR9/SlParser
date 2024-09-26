package internal

import (
	"fmt"
	"strings"
)

// Rule is a production of a grammar.
type Rule struct {
	// Elems is the elements of the rule.
	Elems []string
}

// String implements the fmt.Stringer interface.
func (r Rule) String() string {
	return "_ = is.AddRule(" + strings.Join(r.Elems, ", ") + ")"
}

// NewRule creates a new rule.
//
// Parameters:
//   - lhs: the left hand side of the rule.
//   - rhss: the right hand sides of the rule.
//
// Returns:
//   - *Rule: the new rule.
//   - error: the error if any.
func NewRule(lhs string, rhss []string) (*Rule, error) {
	if len(rhss) == 0 {
		return nil, fmt.Errorf("expected at least one rhss")
	}

	return &Rule{
		Elems: append([]string{lhs}, rhss...),
	}, nil
}
