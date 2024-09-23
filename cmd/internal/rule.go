package internal

import "strings"

type Rule struct {
	Elems []string
}

func (r Rule) String() string {
	return "_ = is.AddRule(" + strings.Join(r.Elems, ", ") + ")"
}

func NewRule(lhs *Token, rhss []*Token) *Rule {
	elems := make([]string, 1, len(rhss)+1)
	elems[0] = lhs.String()

	for _, rhs := range rhss {
		elems = append(elems, rhs.String())
	}

	return &Rule{
		Elems: elems,
	}
}
