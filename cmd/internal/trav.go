package internal

import (
	"fmt"

	kdd "github.com/PlayerR9/SlParser/kdd"
	gers "github.com/PlayerR9/go-errors"
)

func ExtractRules(infos map[*kdd.Node]*Info, root *kdd.Node) ([]*Rule, error) {
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

		children := c.GetChildren()
		if len(children) < 2 {
			return nil, fmt.Errorf("expected at least two children, got %d instead", len(children))
		}

		var ids []string

		for _, child := range children {
			info, ok := infos[child]
			if !ok || info == nil {
				return nil, fmt.Errorf("expected Info, got %s instead", child.Type.String())
			}

			type_ := child.Type

			if type_ != kdd.RhsNode {
				return nil, fmt.Errorf("id expected to be RHS, got %s instead", type_.String())
			}

			ids = append(ids, info.Literal)
		}

		rule, err := NewRule(ids[0], ids[1:])
		gers.AssertErr(err, "NewRule(lhs, rhss)")
		gers.AssertNotNil(rule, "rule")

		rules = append(rules, rule)
	}

	return rules, nil
}
