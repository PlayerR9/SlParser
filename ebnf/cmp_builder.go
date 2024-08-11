package ebnf

import (
	"fmt"

	ast "github.com/PlayerR9/grammar/ast"
	gr "github.com/PlayerR9/grammar/grammar"
)

var (
	// ast_builder is the AST builder of the parser.
	ast_builder *ast.Make[*Node, token_type]
)

func init() {
	ast_builder = ast.NewMake[*Node, token_type]()

	parts := ast.NewPartsBuilder[*Node]()

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		// luc.AssertNil(a, "a")

		root := prev.(*gr.Token[token_type])
		// root := luc.AssertConv[*gr.Token[TokenType]](prev, "prev")

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return nil, err
		}

		if len(children) != 1 {
			return nil, fmt.Errorf("expected 1 child, got %d instead", len(children))
		}

		return children[0], nil
	})

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		// luc.AssertNil(a, "a")

		child := prev.(*gr.Token[token_type])
		// child := luc.AssertConv[*gr.Token[TokenType]](prev, "prev")

		// Identifier = uppercase_id .
		// Identifier = lowercase_id .

		data, err := ast.ExtractData(child)
		if err != nil {
			return nil, err
		}

		a.SetNode(NewNode(IdentifierNode, data, child.At))

		return nil, nil
	})

	ast_builder.AddEntry(ntk_Identifier, parts.Build())
	parts.Reset()

	f1 := func(children []*gr.Token[token_type]) ([]*Node, error) {
		switch len(children) {
		case 2:
			// OrExpr = Identifier pipe (OrExpr) .

			sub_nodes, err := ast_builder.ApplyToken(children[0])
			if err != nil {
				return sub_nodes, err
			}

			return sub_nodes, nil
		case 3:
			// OrExpr = Identifier pipe Identifier .

			var nodes []*Node

			sub_nodes, err := ast_builder.ApplyToken(children[0])
			if len(sub_nodes) > 0 {
				nodes = append(nodes, sub_nodes...)
			}

			if err != nil {
				return nodes, err
			}

			sub_nodes, err = ast_builder.ApplyToken(children[2])
			if len(sub_nodes) > 0 {
				nodes = append(nodes, sub_nodes...)
			}

			if err != nil {
				return nodes, err
			}

			return nodes, nil
		default:
			return nil, fmt.Errorf("expected 2 or 3 children, got %d instead", len(children))
		}
	}

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		// luc.AssertNil(a, "a")

		root := prev.(*gr.Token[token_type])
		// root := luc.AssertConv[*gr.Token[TokenType]](prev, "prev")

		a.SetNode(NewNode(OrNode, "", root.At))

		nodes, err := ast.LeftRecursive(root, ntk_OrExpr, f1)

		a_nodes := make([]ast.Noder, 0, len(nodes))

		for _, node := range nodes {
			a_nodes = append(a_nodes, node)
		}

		a.AppendChildren(a_nodes)

		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	ast_builder.AddEntry(ntk_OrExpr, parts.Build())
	parts.Reset()

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		// luc.AssertNil(a, "a")

		root := prev.(*gr.Token[token_type])
		// root := luc.AssertConv[*gr.Token[TokenType]](prev, "prev")

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return children, err
		}

		return children, nil
	})

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		// luc.AssertNil(a, "a")

		children := prev.([]*gr.Token[token_type])
		// children := luc.AssertConv[[]*gr.Token[TokenType]](prev, "prev")

		switch len(children) {
		case 1:
			// Rhs = Identifier .

			sub_nodes, err := ast_builder.ApplyToken(children[0])
			a.SetNodes(sub_nodes)

			if err != nil {
				return nil, err
			}

			if len(sub_nodes) != 1 {
				return nil, fmt.Errorf("expected 1 child, got %d instead", len(sub_nodes))
			}
		case 3:
			// Rhs = op_paren OrExpr cl_paren .

			a.SetNode(NewNode(OrNode, "", children[1].At))

			sub_nodes, err := ast_builder.ApplyToken(children[1])
			a_nodes := make([]ast.Noder, 0, len(sub_nodes))

			for _, node := range sub_nodes {
				a_nodes = append(a_nodes, node)
			}

			a.AppendChildren(a_nodes)

			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("expected 1 or 3 children, got %d instead", len(children))
		}

		return nil, nil
	})

	ast_builder.AddEntry(ntk_Rhs, parts.Build())
	parts.Reset()

	f2 := func(children []*gr.Token[token_type]) ([]*Node, error) {
		nodes, err := ast_builder.ApplyToken(children[0])
		if err != nil {
			return nodes, err
		}

		return nodes, nil
	}

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		// luc.AssertNil(a, "a")

		root := prev.(*gr.Token[token_type])
		// root := luc.AssertConv[*gr.Token[TokenType]](prev, "prev")

		// RhsCls = Rhs .
		// RhsCls = Rhs RhsCls .

		sub_nodes, err := ast.LeftRecursive(root, ntk_RhsCls, f2)
		a.SetNodes(sub_nodes)

		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	ast_builder.AddEntry(ntk_RhsCls, parts.Build())
	parts.Reset()

	f3 := func(children []*gr.Token[token_type]) ([]*Node, error) {
		switch len(children) {
		case 3:
			// RuleLine = newline tab dot  .

			// do nothing

			return nil, nil
		case 4:
			// RuleLine = newline tab pipe RhsCls (RuleLine) .

			n := NewNode(OrNode, "", children[3].At)

			sub_nodes, err := ast_builder.ApplyToken(children[3])
			a_nodes := make([]ast.Noder, 0, len(sub_nodes))

			for _, node := range sub_nodes {
				a_nodes = append(a_nodes, node)
			}

			n.AddChildren(a_nodes)

			if err != nil {
				return []*Node{n}, err
			}

			return []*Node{n}, nil
		default:
			return nil, fmt.Errorf("expected 3 or 4 children, got %d children instead", len(children))
		}
	}

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		// luc.AssertNil(a, "a")

		root := prev.(*gr.Token[token_type])
		// root := luc.AssertConv[*gr.Token[TokenType]](prev, "prev")

		sub_nodes, err := ast.LeftRecursive(root, ntk_RuleLine, f3)
		a.SetNodes(sub_nodes)

		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	ast_builder.AddEntry(ntk_RuleLine, parts.Build())
	parts.Reset()

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		// luc.AssertNil(a, "a")

		root := prev.(*gr.Token[token_type])
		// root := luc.AssertConv[*gr.Token[TokenType]](prev, "prev")

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return children, err
		}

		return children, nil
	})

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		// luc.AssertNil(a, "a")

		children := prev.([]*gr.Token[token_type])
		// children := luc.AssertConv[[]*gr.Token[TokenType]](prev, "prev")

		lhs, err := ast.ExtractData(children[0])
		if err != nil {
			return nil, err
		}

		switch len(children) {
		case 4:
			// Rule = uppercase_id equal RhsCls dot .

			a.SetNode(NewNode(RuleNode, lhs, children[2].At))

			sub_nodes, err := ast_builder.ApplyToken(children[2])
			a_nodes := make([]ast.Noder, 0, len(sub_nodes))

			for _, node := range sub_nodes {
				a_nodes = append(a_nodes, node)
			}

			a.AppendChildren(a_nodes)

			if err != nil {
				return nil, err
			}
		case 6:
			// Rule = uppercase_id newline tab equal RhsCls RuleLine .

			n := NewNode(OrNode, "", children[4].At)

			sub_nodes, err := ast_builder.ApplyToken(children[4])
			a_nodes := make([]ast.Noder, 0, len(sub_nodes))

			for _, node := range sub_nodes {
				a_nodes = append(a_nodes, node)
			}

			n.AddChildren(a_nodes)

			a.AppendNodes([]*Node{n})

			if err != nil {
				_ = a.DoForEach(func(n *Node) error {
					n.Data = lhs
					n.Type = RuleNode

					return nil
				})

				return nil, err
			}

			sub_nodes, err = ast_builder.ApplyToken(children[5])
			a.AppendNodes(sub_nodes)

			_ = a.DoForEach(func(n *Node) error {
				n.Data = lhs
				n.Type = RuleNode

				return nil
			})

			if err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("expected 4 or 6 children, got %d children instead", len(children))
		}

		return nil, nil
	})

	ast_builder.AddEntry(ntk_Rule, parts.Build())
	parts.Reset()

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		// luc.AssertNil(a, "a")

		root := prev.(*gr.Token[token_type])
		// root := luc.AssertConv[*gr.Token[TokenType]](prev, "prev")

		// Source = Source1 EOF .

		a.SetNode(NewNode(SourceNode, "", root.At))

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return nil, err
		}

		if len(children) != 2 {
			return nil, fmt.Errorf("expected 2 children, got %d children instead", len(children))
		}

		return children[0], nil
	})

	f4 := func(children []*gr.Token[token_type]) ([]*Node, error) {
		nodes, err := ast_builder.ApplyToken(children[0])
		if err != nil {
			return nodes, err
		}

		return nodes, nil
	}

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		// luc.AssertNil(a, "a")

		child := prev.(*gr.Token[token_type])
		// child := luc.AssertConv[*gr.Token[TokenType]](prev, "prev")

		sub_nodes, err := ast.LeftRecursive(child, ntk_Source1, f4)
		a_nodes := make([]ast.Noder, 0, len(sub_nodes))

		for _, node := range sub_nodes {
			a_nodes = append(a_nodes, node)
		}

		a.AppendChildren(a_nodes)

		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	ast_builder.AddEntry(ntk_Source, parts.Build())
}
