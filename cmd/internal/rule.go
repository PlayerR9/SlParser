package internal

import (
	"fmt"
	"strings"
)

type Rule struct {
	Elems []string
}

func (r Rule) String() string {
	return "_ = is.AddRule(" + strings.Join(r.Elems, ", ") + ")"
}

func NewRule(lhs string, rhss []string) (*Rule, error) {
	if len(rhss) == 0 {
		return nil, fmt.Errorf("expected at least one rhss")
	}

	return &Rule{
		Elems: append([]string{lhs}, rhss...),
	}, nil
}
