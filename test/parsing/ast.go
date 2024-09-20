package parsing

import (
	ast "github.com/PlayerR9/SlParser/ast"
	gr "github.com/PlayerR9/SlParser/grammar"
	util "github.com/PlayerR9/SlParser/util/go-commons/errors"
)

//go:generate stringer -type=NodeType -linecomment

// go:generate go run github.com/PlayerR9/SlParser/cmd -o=node.go

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
	AstMaker *ast.AstMaker[*Node, TokenType]
)

func init() {
	builder := ast.NewBuilder[*Node, TokenType]()

	builder.Register(NttSource, func(tk *gr.ParseTree[TokenType]) (*Node, error) {
		// Token[T][79:NttSource]
		//  ├── Token[T][79:TttNewline]
		//  ├── Token[T][80:NttSource1]
		//  └── Token[T][-1:EttEOF]

		children := tk.GetChildren()
		if len(children) != 3 {
			return nil, util.NewErrValue("children", 3, len(children), true)
		}

		err := ast.CheckType(children, 0, TttNewline)
		if err != nil {
			return nil, err
		}

		err = ast.CheckType(children, 2, EttEOF)
		if err != nil {
			return nil, err
		}

		fn := func(children []*gr.ParseTree[TokenType]) (*Node, error) {
			var node *Node

			switch len(children) {
			case 1:
				// Token[T][112:NttStatement]
				// └── Token[T][112:TttPrintStmt ("sq")]

				tmp, err := AstMaker.Convert(children[0])
				if err != nil {
					return nil, err
				}

				node = tmp
			case 2:
				// Token[T][80:NttStatement]
				// ├── Token[T][80:TttListComprehension ("sq = [x * x for x in range(10)]")]
				// └── Token[T][111:TttNewline]

				err := ast.CheckType(children, 1, TttNewline)
				if err != nil {
					return nil, err
				}

				tmp, err := AstMaker.Convert(children[0])
				if err != nil {
					return nil, err
				}

				node = tmp
			default:
				return nil, util.NewErrValues("children", []int{1, 2}, len(children), false)
			}

			return node, nil
		}

		subnodes, err := ast.LhsToAst(children[1], NttSource1, fn)
		if err != nil {
			return nil, err
		}

		node := NewNode(tk.Pos(), SourceNode, "")
		node.AddChildren(subnodes)

		return node, nil
	})

	builder.Register(NttStatement, func(tk *gr.ParseTree[TokenType]) (*Node, error) {
		// Token[T][80:NttStatement]
		//  └── Token[T][80:TttListComprehension ("sq = [x * x for x in range(10)]")]

		// Token[T][112:NttStatement]
		//  └── Token[T][112:TttPrintStmt ("sq")]

		children := tk.GetChildren()
		if len(children) != 1 {
			return nil, util.NewErrValue("children", 1, len(children), true)
		}

		node, err := AstMaker.Convert(children[0])
		if err != nil {
			return nil, err
		}

		return node, nil
	})

	builder.Register(TttListComprehension, func(tk *gr.ParseTree[TokenType]) (*Node, error) {
		// Token[T][80:TttListComprehension ("sq = [x * x for x in range(10)]")]

		children := tk.GetChildren()
		if len(children) != 0 {
			return nil, util.NewErrValue("children", 0, len(children), true)
		}

		node := NewNode(tk.Pos(), ListComprehensionNode, tk.Data())

		return node, nil
	})

	builder.Register(TttPrintStmt, func(tk *gr.ParseTree[TokenType]) (*Node, error) {
		// Token[T][112:TttPrintStmt ("sq")]

		children := tk.GetChildren()
		if len(children) != 0 {
			return nil, util.NewErrValue("children", 0, len(children), true)
		}

		node := NewNode(tk.Pos(), PrintStmtNode, tk.Data())

		return node, nil
	})

	AstMaker = builder.Build()
}
