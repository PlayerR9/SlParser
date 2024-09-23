package parsing

import (
	"fmt"

	ast "github.com/PlayerR9/SlParser/ast"
	gr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/test/parsing/internal"
)

//go:generate stringer -type=NodeType -linecomment

type NodeType int

const (
	/*SourceNode represents the root node.
	Node[SourceNode]
	 ├── // ...
	 └── any statement
	*/
	SourceNode NodeType = iota

	/*ListComprehensionNode represents a list comprehension.
	Node[ListComprehensionNode ("sq = [x * x for x in range(10)]")]
	*/
	ListComprehensionNode

	/*PrintStmtNode represents a print statement.
	Node[PrintStmtNode ("sq")]
	*/
	PrintStmtNode
)

var (
	ast_maker *ast.AstMaker[*Node, internal.TokenType]
)

func init() {
	builder := ast.NewBuilder[*Node, internal.TokenType]()

	builder.Register(internal.NtSource, func(tk *gr.ParseTree[internal.TokenType]) (*Node, error) {
		// Token[T][79:NttSource]
		//  ├── Token[T][79:TttNewline]
		//  ├── Token[T][80:NttSource1]
		//  └── Token[T][-1:EttEOF]

		children := tk.GetChildren()
		if len(children) != 3 {
			return nil, fmt.Errorf("expected 3 children, got %d instead", len(children))
		}

		err := ast.CheckType(children, 0, internal.TtNewline)
		if err != nil {
			return nil, err
		}

		err = ast.CheckType(children, 2, internal.EtEOF)
		if err != nil {
			return nil, err
		}

		fn := func(children []*gr.ParseTree[internal.TokenType]) (*Node, error) {
			var node *Node

			switch len(children) {
			case 1:
				// Token[T][112:NttStatement]
				// └── Token[T][112:TttPrintStmt ("sq")]

				tmp, err := ast_maker.Convert(children[0])
				if err != nil {
					return nil, err
				}

				node = tmp
			case 2:
				// Token[T][80:NttStatement]
				// ├── Token[T][80:TttListComprehension ("sq = [x * x for x in range(10)]")]
				// └── Token[T][111:TttNewline]

				err := ast.CheckType(children, 1, internal.TtNewline)
				if err != nil {
					return nil, err
				}

				tmp, err := ast_maker.Convert(children[0])
				if err != nil {
					return nil, err
				}

				node = tmp
			default:
				return nil, fmt.Errorf("expected 1 or 2 children, got %d instead", len(children))
			}

			return node, nil
		}

		subnodes, err := ast.LhsToAst(1, children, internal.NtSource1, fn)
		if err != nil {
			return nil, err
		}

		node := NewNode(tk.Pos(), SourceNode, "")
		node.AddChildren(subnodes)

		return node, nil
	})

	builder.Register(internal.NtStatement, func(tk *gr.ParseTree[internal.TokenType]) (*Node, error) {
		// Token[T][80:NttStatement]
		//  └── Token[T][80:TttListComprehension ("sq = [x * x for x in range(10)]")]

		// Token[T][112:NttStatement]
		//  └── Token[T][112:TttPrintStmt ("sq")]

		children := tk.GetChildren()
		if len(children) != 1 {
			return nil, fmt.Errorf("expected 1 child, got %d instead", len(children))
		}

		node, err := ast_maker.Convert(children[0])
		if err != nil {
			return nil, err
		}

		return node, nil
	})

	builder.Register(internal.TtListComprehension, func(tk *gr.ParseTree[internal.TokenType]) (*Node, error) {
		// Token[T][80:TttListComprehension ("sq = [x * x for x in range(10)]")]

		children := tk.GetChildren()
		if len(children) != 0 {
			return nil, fmt.Errorf("expected no children, got %d instead", len(children))
		}

		node := NewNode(tk.Pos(), ListComprehensionNode, tk.Data())

		return node, nil
	})

	builder.Register(internal.TtPrintStmt, func(tk *gr.ParseTree[internal.TokenType]) (*Node, error) {
		// Token[T][112:TttPrintStmt ("sq")]

		children := tk.GetChildren()
		if len(children) != 0 {
			return nil, fmt.Errorf("expected no children, got %d instead", len(children))
		}

		node := NewNode(tk.Pos(), PrintStmtNode, tk.Data())

		return node, nil
	})

	ast_maker = builder.Build()
}
