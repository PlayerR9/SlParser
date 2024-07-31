package parser

import (
	"fmt"

	pkg "github.com/PlayerR9/SLParser/parser/pkg"
	ast "github.com/PlayerR9/grammar/ast"
	gr "github.com/PlayerR9/grammar/grammar"
	ulx "github.com/PlayerR9/grammar/lexer"
	uprx "github.com/PlayerR9/grammar/parser"
	luc "github.com/PlayerR9/lib_units/common"
)

var (
	// lexer is the lexer of the parser.
	lexer *pkg.Lexer

	// parser is the parser of the parser.
	parser *pkg.Parser
)

func init() {
	lexer = pkg.NewLexer()
	parser = pkg.NewParser()
}

// Parse parses the given data and returns the AST tree.
//
// Parameters:
//   - data: The data to parse.
//
// Returns:
//   - *ast.Node[NodeType]: The AST tree.
//   - error: An error if the parsing failed.
func Parse(data []byte) (*ast.Node[NodeType], error) {
	if len(data) == 0 {
		return nil, luc.NewErrInvalidParameter("data", luc.NewErrEmpty("slice of bytes"))
	}

	tokens, err := ulx.FullLex(lexer, data)
	if err != nil {
		// DEBUG: Print tokens:
		fmt.Println(string(ulx.PrintSyntaxError(data, tokens)))
		fmt.Println()

		return nil, fmt.Errorf("error while lexing: %w", err)
	}

	forest, err := uprx.FullParse(parser, tokens)
	if err != nil {
		for _, tree := range forest {
			fmt.Println(tree.String())
			fmt.Println()
		}

		fmt.Println()

		return nil, fmt.Errorf("error while parsing: %w", err)
	} else if len(forest) != 1 {
		for _, tree := range forest {
			fmt.Println(tree.String())
			fmt.Println()
		}

		fmt.Println()

		return nil, fmt.Errorf("expected 1 tree, got %d trees instead", len(forest))
	}

	nodes, err := AstBuilder.Apply(forest[0].Root())
	if err != nil {
		for _, node := range nodes {
			fmt.Println(ast.PrintAst(node))
			fmt.Println()
		}

		fmt.Println()

		return nil, fmt.Errorf("error while converting to AST: %w", err)
	} else if len(nodes) != 1 {
		return nil, fmt.Errorf("expected 1 node, got %d nodes instead", len(nodes))
	}

	return nodes[0], nil
}

var (
	// AstBuilder is the AST builder of the parser.
	AstBuilder *ast.Make[NodeType, pkg.TokenType]
)

func init() {
	AstBuilder = ast.NewMake[NodeType, pkg.TokenType]()

	parts := ast.NewPartsBuilder[NodeType]()

	parts.Add(func(a *ast.Result[NodeType], prev any) (any, error) {
		luc.AssertNil(a, "a")

		root := luc.AssertConv[*gr.Token[pkg.TokenType]](prev, "prev")

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return nil, err
		}

		if len(children) != 1 {
			return nil, fmt.Errorf("expected 1 child, got %d instead", len(children))
		}

		return children[0], nil
	})

	parts.Add(func(a *ast.Result[NodeType], prev any) (any, error) {
		luc.AssertNil(a, "a")

		child := luc.AssertConv[*gr.Token[pkg.TokenType]](prev, "prev")

		// Identifier = uppercase_id .
		// Identifier = lowercase_id .

		data, err := ast.ExtractData(child)
		if err != nil {
			return nil, err
		}

		a.MakeNode(IdentifierNode, data)

		return nil, nil
	})

	AstBuilder.AddEntry(pkg.NtkIdentifier, parts.Build())
	parts.Reset()

	f1 := func(children []*gr.Token[pkg.TokenType]) ([]*ast.Node[NodeType], error) {
		switch len(children) {
		case 2:
			// OrExpr = Identifier pipe (OrExpr) .

			sub_nodes, err := AstBuilder.Apply(children[0])
			if err != nil {
				return sub_nodes, err
			}

			return sub_nodes, nil
		case 3:
			// OrExpr = Identifier pipe Identifier .

			var nodes []*ast.Node[NodeType]

			sub_nodes, err := AstBuilder.Apply(children[0])
			if len(sub_nodes) > 0 {
				nodes = append(nodes, sub_nodes...)
			}

			if err != nil {
				return nodes, err
			}

			sub_nodes, err = AstBuilder.Apply(children[2])
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

	parts.Add(func(a *ast.Result[NodeType], prev any) (any, error) {
		luc.AssertNil(a, "a")

		root := luc.AssertConv[*gr.Token[pkg.TokenType]](prev, "prev")

		a.MakeNode(OrNode, "")

		nodes, err := ast.LeftRecursive(root, pkg.NtkOrExpr, f1)
		a.AppendChildren(nodes)

		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	AstBuilder.AddEntry(pkg.NtkOrExpr, parts.Build())
	parts.Reset()

	parts.Add(func(a *ast.Result[NodeType], prev any) (any, error) {
		luc.AssertNil(a, "a")

		root := luc.AssertConv[*gr.Token[pkg.TokenType]](prev, "prev")

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return children, err
		}

		return children, nil
	})

	parts.Add(func(a *ast.Result[NodeType], prev any) (any, error) {
		luc.AssertNil(a, "a")

		children := luc.AssertConv[[]*gr.Token[pkg.TokenType]](prev, "prev")

		switch len(children) {
		case 1:
			// Rhs = Identifier .

			sub_nodes, err := AstBuilder.Apply(children[0])
			a.SetNodes(sub_nodes)

			if err != nil {
				return nil, err
			}

			if len(sub_nodes) != 1 {
				return nil, fmt.Errorf("expected 1 child, got %d instead", len(sub_nodes))
			}
		case 3:
			// Rhs = op_paren OrExpr cl_paren .

			a.MakeNode(OrNode, "")

			sub_nodes, err := AstBuilder.Apply(children[1])
			a.AppendChildren(sub_nodes)

			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("expected 1 or 3 children, got %d instead", len(children))
		}

		return nil, nil
	})

	AstBuilder.AddEntry(pkg.NtkRhs, parts.Build())
	parts.Reset()

	f2 := func(children []*gr.Token[pkg.TokenType]) ([]*ast.Node[NodeType], error) {
		nodes, err := AstBuilder.Apply(children[0])
		if err != nil {
			return nodes, err
		}

		return nodes, nil
	}

	parts.Add(func(a *ast.Result[NodeType], prev any) (any, error) {
		luc.AssertNil(a, "a")

		root := luc.AssertConv[*gr.Token[pkg.TokenType]](prev, "prev")

		// RhsCls = Rhs .
		// RhsCls = Rhs RhsCls .

		sub_nodes, err := ast.LeftRecursive(root, pkg.NtkRhsCls, f2)
		a.SetNodes(sub_nodes)

		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	AstBuilder.AddEntry(pkg.NtkRhsCls, parts.Build())
	parts.Reset()

	f3 := func(children []*gr.Token[pkg.TokenType]) ([]*ast.Node[NodeType], error) {
		switch len(children) {
		case 3:
			// RuleLine = newline tab dot  .

			// do nothing

			return nil, nil
		case 4:
			// RuleLine = newline tab pipe RhsCls (RuleLine) .

			n := ast.NewNode(OrNode, "")

			sub_nodes, err := AstBuilder.Apply(children[3])
			n.AddChildren(sub_nodes)

			if err != nil {
				return []*ast.Node[NodeType]{n}, err
			}

			return []*ast.Node[NodeType]{n}, nil
		default:
			return nil, fmt.Errorf("expected 3 or 4 children, got %d children instead", len(children))
		}
	}

	parts.Add(func(a *ast.Result[NodeType], prev any) (any, error) {
		luc.AssertNil(a, "a")

		root := luc.AssertConv[*gr.Token[pkg.TokenType]](prev, "prev")

		sub_nodes, err := ast.LeftRecursive(root, pkg.NtkRuleLine, f3)
		a.SetNodes(sub_nodes)

		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	AstBuilder.AddEntry(pkg.NtkRuleLine, parts.Build())
	parts.Reset()

	parts.Add(func(a *ast.Result[NodeType], prev any) (any, error) {
		luc.AssertNil(a, "a")

		root := luc.AssertConv[*gr.Token[pkg.TokenType]](prev, "prev")

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return children, err
		}

		return children, nil
	})

	parts.Add(func(a *ast.Result[NodeType], prev any) (any, error) {
		luc.AssertNil(a, "a")

		children := luc.AssertConv[[]*gr.Token[pkg.TokenType]](prev, "prev")

		lhs, err := ast.ExtractData(children[0])
		if err != nil {
			return nil, err
		}

		switch len(children) {
		case 4:
			// Rule = uppercase_id equal RhsCls dot .

			a.MakeNode(RuleNode, lhs)

			sub_nodes, err := AstBuilder.Apply(children[2])
			a.AppendChildren(sub_nodes)

			if err != nil {
				return nil, err
			}
		case 6:
			// Rule = uppercase_id newline tab equal RhsCls RuleLine .

			n := ast.NewNode(OrNode, "")

			sub_nodes, err := AstBuilder.Apply(children[4])
			n.AddChildren(sub_nodes)

			a.AppendNodes([]*ast.Node[NodeType]{n})

			if err != nil {
				a.TransformNodes(RuleNode, lhs)
				return nil, err
			}

			sub_nodes, err = AstBuilder.Apply(children[5])
			a.AppendNodes(sub_nodes)
			a.TransformNodes(RuleNode, lhs)

			if err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("expected 4 or 6 children, got %d children instead", len(children))
		}

		return nil, nil
	})

	AstBuilder.AddEntry(pkg.NtkRule, parts.Build())
	parts.Reset()

	parts.Add(func(a *ast.Result[NodeType], prev any) (any, error) {
		luc.AssertNil(a, "a")

		root := luc.AssertConv[*gr.Token[pkg.TokenType]](prev, "prev")

		// Source = Source1 EOF .

		a.MakeNode(SourceNode, "")

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return nil, err
		}

		if len(children) != 2 {
			return nil, fmt.Errorf("expected 2 children, got %d children instead", len(children))
		}

		return children[0], nil
	})

	f4 := func(children []*gr.Token[pkg.TokenType]) ([]*ast.Node[NodeType], error) {
		nodes, err := AstBuilder.Apply(children[0])
		if err != nil {
			return nodes, err
		}

		return nodes, nil
	}

	parts.Add(func(a *ast.Result[NodeType], prev any) (any, error) {
		luc.AssertNil(a, "a")

		child := luc.AssertConv[*gr.Token[pkg.TokenType]](prev, "prev")

		sub_nodes, err := ast.LeftRecursive(child, pkg.NtkSource1, f4)
		a.AppendChildren(sub_nodes)

		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	AstBuilder.AddEntry(pkg.NtkSource, parts.Build())
}
