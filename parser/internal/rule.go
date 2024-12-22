package internal

// Rule is a rule in the grammar.
type Rule struct {
	// lhs is the left-hand side of the rule.
	lhs string

	// rhss is the right-hand side of the rule.
	rhss []string
}

// NewRule creates a new rule with the given left-hand side and right-hand side symbols.
//
// Parameters:
//   - lhs: The left-hand side of the rule.
//   - rhss: The right-hand side of the rule.
//
// Returns:
//   - *Rule: The newly created rule. Never returns nil.
func NewRule(lhs string, rhss []string) *Rule {
	rule := &Rule{
		lhs:  lhs,
		rhss: rhss,
	}

	return rule
}

// Rhss returns a copy of the right-hand side symbols of the rule.
//
// Returns:
//   - []string: A copy of the right-hand side symbols of the rule. Returns nil
//     if the rule has no right-hand side symbols.
func (r Rule) Rhss() []string {
	if len(r.rhss) == 0 {
		return nil
	}

	rhss := make([]string, len(r.rhss))
	copy(rhss, r.rhss)

	return rhss
}

// Lhs returns the left-hand side of the rule.
//
// Returns:
//   - string: The left-hand side of the rule. Never returns nil.
func (r Rule) Lhs() string {
	return r.lhs
}
