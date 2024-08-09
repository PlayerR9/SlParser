package pkg

import (
	"fmt"

	ebnf "github.com/PlayerR9/EbnfParser/pkg"
)

func node_to_rule(root *ebnf.Node) *Rule {
	lhs := root.Data

	var rhss []string

	for c := root.FirstChild; c != nil; c = c.NextSibling {
		rhss = append(rhss, c.Data)
	}

	return NewRule(lhs, rhss, false)
}

func ExtractRules(root *ebnf.Node) ([]*Rule, error) {
	if root == nil {
		return nil, fmt.Errorf("expected %q, got nothing instead", ebnf.SourceNode.String())
	} else if root.Type != ebnf.SourceNode {
		return nil, fmt.Errorf("expected %q, got %q instead", ebnf.SourceNode.String(), root.Type.String())
	}

	var sub_roots []*ebnf.Node

	for c := root.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != ebnf.RuleNode {
			return nil, fmt.Errorf("expected %q, got %q instead", ebnf.RuleNode, c.Type)
		}

		sub_roots = append(sub_roots, c)
	}

	var rules []*Rule

	for _, sub_root := range sub_roots {
		rule := node_to_rule(sub_root)

		rules = append(rules, rule)
	}

	return rules, nil
}
