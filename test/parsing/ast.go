package parsing

import (
	"fmt"

	ast "github.com/PlayerR9/SlParser/ast"
	gr "github.com/PlayerR9/SlParser/grammar"
)

//go:generate stringer -type=NodeType -linecomment

//go:generate go run github.com/PlayerR9/SlParser/cmd -o=node.go

type NodeType int

const (
	SouceNode NodeType = iota
)

var (
	AstMaker ast.AstMaker[*Node, TokenType]
)

func init() {
	builder := ast.NewBuilder[*Node, TokenType]()

	builder.Register(NttSource, func(tk *gr.Token[TokenType]) (*Node, error) {
		// Token[T][79:NttSource]
		//  ├── Token[T][79:TttNewline]
		//  └── Token[T][80:NttSource1]
		//  │   └── Token[T][80:NttStatement]
		//  │   │   └── Token[T][80:TttListComprehension ("sq = [x * x for x in range(10)]")]
		//  │   ├── Token[T][111:TttNewline]
		//  │   └── Token[T][112:NttSource1]
		//  │       └── Token[T][112:NttStatement]
		//  │           └── Token[T][112:TttPrintStmt ("sq")]
		//  └── Token[T][-1:EttEOF]

		children := tk.GetChildren()
		if len(children) != 3 {
			return nil, fmt.Errorf("expected 3 children, got %d instead", len(children))
		}

		if children[0].Type != TttNewline {
			return nil, fmt.Errorf("first child expected to be %q, got %q instead", children[0].Type.String())
		}

		ast.LhsToAst(children[1], NttSource1, func(children []*gr.Token[TokenType]) (*Node, error) {
			switch len(children) {
			case 2:
				// Token[T][80:NttStatement]
				// ├──Token[T][80:TttListComprehension ("sq = [x * x for x in range(10)]")]
				// └── Token[T][111:TttNewline]

				children[0]

				children[1]
			case 1:
				// Token[T][112:NttSource1]
				// └── Token[T][112:NttStatement]
				//     └── Token[T][112:TttPrintStmt ("sq")]
			default:
				return nil, fmt.Errorf("expected either 1 or 2 children, got %d instead", len(children))
			}
		})

		if children[2].Type != EttEOF {
			return nil, fmt.Errorf("third child expected to be %q, got %q instead", children[2].Type.String())
		}
	})

	// Token[T][79:NttSource]
	//  ├── Token[T][79:TttNewline]
	//  └── Token[T][80:NttSource1]
	//  │   └── Token[T][80:NttStatement]
	//  │   │   └── Token[T][80:TttListComprehension ("sq = [x * x for x in range(10)]")]
	//  │   ├── Token[T][111:TttNewline]
	//  │   └── Token[T][112:NttSource1]
	//  │       └── Token[T][112:NttStatement]
	//  │           └── Token[T][112:TttPrintStmt ("sq")]
	//  └── Token[T][-1:EttEOF]

	AstMaker = builder.Build()
}
