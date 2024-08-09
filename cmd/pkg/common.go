package pkg

import (
	"fmt"

	prx "github.com/PlayerR9/SLParser/parser"
	ast "github.com/PlayerR9/grammar/ast"
)

func node_to_rule(root *ast.Node[prx.NodeType]) *Rule {
	lhs := root.Data

	var rhss []string

	for c := root.FirstChild; c != nil; c = c.NextSibling {
		rhss = append(rhss, c.Data)
	}

	return NewRule(lhs, rhss, false)
}

func ExtractRules(root *ast.Node[prx.NodeType]) ([]*Rule, error) {
	if root == nil {
		return nil, fmt.Errorf("expected %q, got nothing instead", prx.SourceNode.String())
	} else if root.Type != prx.SourceNode {
		return nil, fmt.Errorf("expected %q, got %q instead", prx.SourceNode.String(), root.Type.String())
	}

	var sub_roots []*ast.Node[prx.NodeType]

	for c := root.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != prx.RuleNode {
			return nil, fmt.Errorf("expected %q, got %q instead", prx.RuleNode, c.Type)
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
