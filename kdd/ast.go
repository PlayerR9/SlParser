package kdd

import (
	"fmt"

	"github.com/PlayerR9/SlParser/ast"
	"github.com/PlayerR9/SlParser/grammar"
	gers "github.com/PlayerR9/go-errors"
	"github.com/PlayerR9/go-errors/assert"
)

// NodeType is the type of a node.
type NodeType int

const (
	/*InvalidNode represents an invalid node.
	Node[InvalidNode]
	*/
	InvalidNode NodeType = iota - 1 // Invalid

	/*RhsNode represents the terminal symbol.
	Node[RhsNode (<id>)]
	*/
	RhsNode // Rhs

	/*RuleNode represents a single rule.
	Node[RuleNode]
	 ├── RhsNode (<id>) // This is the LHS of the rule.
	 ├── RhsNode (<id>) // This is the RHS of the rule.
	 └── ...
	*/
	RuleNode // Rule

	/*SourceNode is the collection of all rules in the grammar.
	Node[SourceNode]
	 ├── Node[RuleNode]
	 └── ...
	*/
	SourceNode // Source
)

// rule : LOWERCASE_ID COLON rhs+ SEMICOLON ;

var (
	ast_maker ast.AstMaker[*Node, TokenType]
)

func init() {
	ast_maker = ast.NewAstMaker[*Node, TokenType]()

	// TODO: Add here your own custom rules...

	// rhs : UPPERCASE_ID ;
	// rhs : LOWERCASE_ID ;
	ast_maker.AddTransformation(NtRhs, func(tk *grammar.ParseTree[TokenType]) (*Node, error) {
		assert.NotNil(tk, "tk")

		field := assert.New(
			NewField(TtUppercaseId, TtLowercaseId),
		)

		rule := assert.New(
			NewRule(NtRhs, false, field),
		)

		sub_nodes, err := rule.ApplyField(tk.GetChildren())
		if err != nil {
			return nil, err
		} else if len(sub_nodes) != 1 {
			return nil, fmt.Errorf("expected one child, got %d instead", len(sub_nodes))
		}

		node := NewNode(RhsNode, sub_nodes[0].Data)
		return node, nil
	})

	ast_maker.AddTransition(NtRule1, func(tree *grammar.ParseTree[TokenType]) ([]*Node, error) {
		assert.NotNil(tree, "tree")

		children := tree.GetChildren()

		var sub_nodes []*Node
		var err error

		switch len(children) {
		case 1:
			// rule1 : rhs ;

			field := assert.New(
				NewField(NtRhs),
			)

			rule := assert.New(
				NewRule(NtRule1, false, field),
			)
			rule.AddExpected(0, RhsNode)

			sub_nodes, err = rule.ApplyField(children)
		case 2:
			// rule1 : rhs rule1 ;

			field1 := assert.New(
				NewField(NtRhs),
			)

			field2 := assert.New(
				NewField(NtRule1),
			)

			rule := assert.New(
				NewRule(NtRule1, true, field1, field2),
			)
			rule.AddExpected(0, RhsNode)

			sub_nodes, err = rule.ApplyField(children)
		default:
			return nil, fmt.Errorf("expected 1 or 2 children, got %d instead", len(children))
		}

		return sub_nodes, err
	})

	ast_maker.AddTransformation(NtRule, func(tk *grammar.ParseTree[TokenType]) (*Node, error) {
		assert.NotNil(tk, "tk")

		children := tk.GetChildren()

		// rule : LOWERCASE_ID COLON rule1 SEMICOLON ;
		err := ast.CheckType(children, 0, TtLowercaseId)
		if err != nil {
			return nil, err
		}

		lhs := NewNode(RhsNode, children[0].Data())
		lhs.SetPosition(children[0].Pos())

		node := NewNode(RuleNode, "")
		node.AddChild(lhs)

		err = ast.CheckType(children, 1, TtColon)
		if err != nil {
			return nil, err
		}

		err = ast.CheckType(children, 3, TtSemicolon)
		if err != nil {
			return nil, err
		}

		sub_children, err := ast.LhsToAst(2, children, NtRule1, rule1)
		if err != nil {
			return nil, err
		}

		node.AddChildren(sub_children)

		return node, nil
	})

	ast_maker.AddTransition(NtSource1, func(tree *grammar.ParseTree[TokenType]) ([]*Node, error) {
		assert.NotNil(tree, "tree")

		children := tree.GetChildren()

		var node *Node

		switch len(children) {
		case 1:
			// source1 : rule ;

			rule := assert.New(NewRule(NtSource1, true, NtRule))
			rule.AddExpected(0, RuleNode)

			sub_rules, err := rule.ApplyField(children)
			if err != nil {
				return nil, err
			}

			node = sub_rules[0]
		case 3:
			// source1 : rule NEWLINE source1 ;

			rule := assert.New(NewRule(NtSource1, true, NtRule, TtNewline, NtSource1))
			rule.AddExpected(0, RuleNode)

			sub_rules, err := rule.ApplyField(children)
			if err != nil {
				return nil, err
			}

			node = sub_rules[0]
		default:
			return nil, fmt.Errorf("expected one or three children, got %d instead", len(children))
		}

		return []*Node{node}, nil
	})

	ast_maker.AddTransformation(NtSource, func(tk *grammar.ParseTree[TokenType]) (*Node, error) {
		if tk == nil {
			return nil, gers.NewErrNilParameter("ast_maker.AddTransformation()", "tk")
		}

		// source : source1 EOF ;
		rule := assert.New(
			NewRule(NtSource, false, NtSource1, EtEOF),
		)

		sub_nodes, err := rule.ApplyField(tk.GetChildren())
		if err != nil {
			return nil, err
		}

		tmp, err := ast.LhsToAst(0, children, NtSource1, source1)
		if err != nil {
			return nil, err
		}

		node := NewNode(SourceNode, "")
		node.AddChildren(sub_nodes)

		return node, nil
	})
}
