package ast

import (
	"fmt"

	gr "github.com/PlayerR9/SlParser/grammar"
	gcers "github.com/PlayerR9/errors"
)

type ToAstFunc[N interface {
	AddChildren(children []N)
}, T gr.TokenTyper] func(tk *gr.ParseTree[T]) (N, error)

// AstMaker is an ast maker.
type AstMaker[N interface {
	AddChildren(children []N)
}, T gr.TokenTyper] struct {
	// table is the ast table.
	table map[T]ToAstFunc[N, T]

	// make_fake_node is the function that makes the fake node.
	// make_fake_node func(root *gr.ParseTree[T]) N
}

// Convert is a function that converts a token to an ast node.
//
// Parameters:
//   - root: The root token. Assumed to be non-nil.
//
// Returns:
//   - N: The ast node.
//   - error: if an error occurred.
func (am AstMaker[N, T]) Convert(root *gr.ParseTree[T]) (N, error) {
	if root == nil {
		return *new(N), gcers.NewErrNilParameter("root")
	}

	type_ := root.Type()

	var node N
	var err error

	fn, ok := am.table[type_]
	if !ok {
		err = fmt.Errorf("type is not registered")
	} else {
		node, err = fn(root)
	}

	if err != nil {
		// if am.make_fake_node != nil {
		// 	node = TransformFakeNode[N](root, am.make_fake_node)
		// }

		err = NewErrIn(type_, err)
	}

	return node, err
}

/* func MakeAst(am AstMaker, tk *gr.Token[T]) (*Node, error) {

} */
