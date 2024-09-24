package internal

import (
	"fmt"
	"slices"
	"strings"

	kdd "github.com/PlayerR9/SlParser/kdd"
	gers "github.com/PlayerR9/go-errors"
)

func extract_symbols_rec(root *kdd.Node, seen map[string]struct{}, symbols []*kdd.Node) []*kdd.Node {
	gers.AssertNotNil(root, "root")
	gers.AssertNotNil(seen, "seen")

	if root.Type == kdd.RhsNode {
		_, ok := seen[root.Data]
		if !ok {
			symbols = append(symbols, root)
			seen[root.Data] = struct{}{}
		}

		return symbols
	}

	for c := root.FirstChild; c != nil; c = c.NextSibling {
		symbols = extract_symbols_rec(c, seen, symbols)
	}

	return symbols
}

func sort(tokens []*kdd.Node) error {
	buckets := make(map[TokenType][]*kdd.Node, 3)
	buckets[ExtraTk] = make([]*kdd.Node, 0)
	buckets[TerminalTk] = make([]*kdd.Node, 0)
	buckets[NonterminalTk] = make([]*kdd.Node, 0)

	for i, tk := range tokens {
		type_, err := TypeOf(tk)
		if err != nil {
			return fmt.Errorf("at index %d: %w", i, err)
		}

		prev, ok := buckets[type_]
		if !ok {
			return fmt.Errorf("bucket %q not found", type_.String())
		}

		buckets[type_] = append(prev, tk)
	}

	for type_, bucket := range buckets {
		slices.SortStableFunc(bucket, func(a, b *kdd.Node) int {
			return strings.Compare(a.Data, b.Data)
		})

		buckets[type_] = bucket
	}

	i := 0

	tks := buckets[ExtraTk]
	for _, tk := range tks {
		tokens[i] = tk
		i++
	}

	tks = buckets[TerminalTk]
	for _, tk := range tks {
		tokens[i] = tk
		i++
	}

	tks = buckets[NonterminalTk]
	for _, tk := range tks {
		tokens[i] = tk
		i++
	}

	return nil
}

func ExtractSymbols(root *kdd.Node) []*kdd.Node {
	if root == nil {
		return nil
	}

	seen := make(map[string]struct{})

	symbols := extract_symbols_rec(root, seen, nil)

	sort(symbols)

	return symbols
}

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
