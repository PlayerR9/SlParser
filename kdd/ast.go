package kdd

import (
	"fmt"

	"github.com/PlayerR9/SlParser/ast"
	"github.com/PlayerR9/SlParser/grammar"
	gers "github.com/PlayerR9/go-errors"
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
	ast_maker = make(ast.AstMaker[*Node, TokenType])

	// TODO: Add here your own custom rules...

	// rhs : UPPERCASE_ID ;
	// rhs : LOWERCASE_ID ;
	ast_maker[NtRhs] = func(tk *grammar.ParseTree[TokenType]) (*Node, error) {
		children := tk.GetChildren()

		if len(children) != 1 {
			return nil, fmt.Errorf("expected one child, got %d instead", len(children))
		}

		type_ := children[0].Type()
		gers.AssertNotNil(type_, "type_")

		node := NewNode(tk.Pos(), RhsNode, children[0].Data())
		return node, nil
	}

	// rule1 : rhs ;
	// rule1 : rhs rule1 ;
	rule1 := func(children []*grammar.ParseTree[TokenType]) (*Node, error) {
		if len(children) != 1 {
			return nil, fmt.Errorf("expected one child, got %d instead", len(children))
		}

		node, err := ast_maker.Convert(children[0])
		if err != nil {
			return nil, err
		} else if node.Type != RhsNode {
			return nil, fmt.Errorf("expected RhsNode, got %s instead", node.Type.String())
		}

		return node, nil
	}

	ast_maker[NtRule] = func(tk *grammar.ParseTree[TokenType]) (*Node, error) {
		children := tk.GetChildren()

		// rule : LOWERCASE_ID COLON rule1 SEMICOLON ;
		err := ast.CheckType(children, 0, TtLowercaseId)
		if err != nil {
			return nil, err
		}

		lhs := NewNode(children[0].Pos(), RhsNode, children[0].Data())

		node := NewNode(tk.Pos(), RuleNode, "")
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
	}

	source1 := func(children []*grammar.ParseTree[TokenType]) (*Node, error) {
		var node *Node

		switch len(children) {
		case 1:
			// source1 : rule ;

			var err error

			node, err = ast_maker.Convert(children[0])
			if err != nil {
				return nil, err
			} else if node.Type != RuleNode {
				return nil, fmt.Errorf("expected RuleNode, got %s instead", node.Type.String())
			}
		case 2:
			// source1 : rule NEWLINE source1 ;

			err := ast.CheckType(children, 1, TtNewline)
			if err != nil {
				return nil, err
			}

			node, err = ast_maker.Convert(children[0])
			if err != nil {
				return nil, err
			} else if node.Type != RuleNode {
				return nil, fmt.Errorf("expected RuleNode, got %s instead", node.Type.String())
			}
		default:
			return nil, fmt.Errorf("expected one or two children, got %d instead", len(children))
		}

		return node, nil
	}

	ast_maker[NtSource] = func(tk *grammar.ParseTree[TokenType]) (*Node, error) {
		children := tk.GetChildren()
		if len(children) != 2 {
			return nil, fmt.Errorf("expected two children, got %d instead", len(children))
		}

		err := ast.CheckType(children, 1, EtEOF)
		if err != nil {
			return nil, err
		}

		// source : source1 EOF ;

		tmp, err := ast.LhsToAst(0, children, NtSource1, source1)
		if err != nil {
			return nil, err
		}

		node := NewNode(tk.Pos(), SourceNode, "")
		node.AddChildren(tmp)

		return node, nil
	}
}
