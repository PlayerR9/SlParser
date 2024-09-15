package ast

import (
	"errors"
	"fmt"

	gr "github.com/PlayerR9/SlParser/grammar"
	gcers "github.com/PlayerR9/go-commons/errors"
)

type AstMaker[N interface {
	SetChildren(children []N)
}, T gr.TokenTyper] struct {
	table          map[T]ToAstFunc[N, T]
	make_fake_node func(root *gr.Token[T]) N
}

func (am AstMaker[N, T]) Convert(root *gr.Token[T]) (N, error) {
	if root == nil {
		return *new(N), gcers.NewErrNilParameter("root")
	}

	type_ := root.Type

	var node N
	var err error

	fn, ok := am.table[type_]
	if !ok {
		err = fmt.Errorf("type is not registered")
	} else {
		node, err = fn(root)
	}

	if err != nil {
		if am.make_fake_node != nil {
			node = TransformFakeNode[N](root, am.make_fake_node)
		}

		err = NewErrIn(type_, err)
	}

	return node, nil
}

func LhsToAst[N interface {
	SetChildren(children []N)
}, T gr.TokenTyper](root *gr.Token[T], lhs T, do func(children []*gr.Token[T]) (N, error)) ([]N, error) {
	if do == nil {
		return nil, gcers.NewErrNilParameter("do")
	} else if root == nil {
		return nil, gcers.NewErrNilParameter("root")
	} else if root.Type != lhs {
		return nil, fmt.Errorf("expected %q, got %s instead", lhs.String(), root.Type.String())
	}

	var nodes []N

	for root != nil {
		children := root.Children

		if len(children) == 0 {
			return nil, errors.New("expected at least one child")
		}

		last_child := children[len(children)-1]

		var node N
		var err error

		if last_child.Type == lhs {
			node, err = do(children[:len(children)-1])
			root = last_child
		} else {
			node, err = do(children)
			root = nil
		}

		if err != nil {
			return nodes, err
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

/* func MakeAst(am AstMaker, tk *gr.Token[T]) (*Node, error) {

} */

// TransformFakeNode transforms a node into a fake AST node.
//
// Parameters:
//   - tk: the token to transform.
//
// Returns:
//   - *Node: the transformed node.
//
// This function transforms a node into a fake AST node. It does this by creating a new node with the correct type and data,
// and then setting the children of the new node to be the transformed children of the fake node.
func TransformFakeNode[N interface {
	SetChildren(children []N)
}, T gr.TokenTyper](tk *gr.Token[T], fn func(tk *gr.Token[T]) N) N {
	if tk == nil {
		return *new(N)
	}

	node := fn(tk)

	// node := NewNode(FakeNode, tk.Type.String()+" : "+tk.Data)

	var subnodes []N

	for child := range tk.Child() {
		n := TransformFakeNode[N](child, fn)
		subnodes = append(subnodes, n)
	}

	node.SetChildren(subnodes)

	return node
}
