package internal

import (
	"fmt"

	kdd "github.com/PlayerR9/SlParser/kdd"
	gers "github.com/PlayerR9/go-errors"
)

func ExtractRules(root *kdd.Node) ([]*Rule, error) {
	if root == nil {
		return nil, nil
	}

	if root.Type != kdd.SourceNode {
		return nil, fmt.Errorf("expected SourceNode, got %s instead", root.Type.String())
	}

	var rules []*Rule

	for c := root.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != kdd.RuleNode {
			return nil, fmt.Errorf("expected RuleNode, got %s instead", c.Type.String())
		}

		var rhss []string

		children := c.GetChildren()
		if len(children) < 2 {
			return nil, fmt.Errorf("expected at least two children, got %d instead", len(children))
		}

		lhs := children[0].Data
		type_ := children[0].Type

		if type_ != kdd.RhsNode {
			return nil, fmt.Errorf("lhs expected to be RHS, got %s instead", type_.String())
		}

		for i := 1; i < len(children); i++ {
			data := children[i].Data
			type_ := children[i].Type

			if type_ != kdd.RhsNode {
				return nil, fmt.Errorf("lhs expected to be RHS, got %s instead", type_.String())
			}

			rhss = append(rhss, data)
		}

		rule, err := NewRule(lhs, rhss)
		gers.AssertErr(err, "NewRule(lhs, rhss)")
		gers.AssertNotNil(rule, "rule")

		rules = append(rules, rule)
	}

	return rules, nil
}
