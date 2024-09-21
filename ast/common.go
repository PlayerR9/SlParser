package ast

import (
	"errors"
	"fmt"

	gr "github.com/PlayerR9/SlParser/grammar"
	gcers "github.com/PlayerR9/go-commons/errors"
)

/* // TransformFakeNode transforms a node into a fake AST node.
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
	AddChildren(children []N)
}, T gr.TokenTyper](tk *gr.ParseTree[T], fn func(tk *gr.ParseTree[T]) N) N {
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

	node.AddChildren(subnodes)

	return node
} */

// CheckType is a helper function that checks the type of the token at the given
// position.
//
// Parameters:
//   - children: The list of children.
//   - at: The position of the token.
//   - type_: The type of the token.
//
// Returns:
//   - error: if an error occurred.
//
// Errors:
//   - *errors.ErrInvalidParameter: If 'at' is less than 0.
//   - *errors.ErrValue: If the token at the given position is nil or
//     'at' is out of range.
func CheckType[T gr.TokenTyper](children []*gr.ParseTree[T], at int, type_ T) error {
	if at < 0 {
		return gcers.NewErrInvalidParameter("at", gcers.NewErrGTE(0))
	}

	pos_str := gcers.GetOrdinalSuffix(at+1) + " child"

	if at >= len(children) {
		return gcers.NewErrValue(pos_str, type_, nil, true)
	}

	tk := children[at]
	if tk == nil {
		return gcers.NewErrValue(pos_str, type_, nil, true)
	} else if tk.Type() != type_ {
		return gcers.NewErrValue(pos_str, type_, tk.Type(), true)
	}

	return nil
}

// LhsToAst is a function that converts a token to an ast node.
//
// Parameters:
//   - root: The root token. Assumed to be non-nil.
//   - lhs: The lhs token.
//   - do: The function that does the conversion.
func LhsToAst[N interface {
	AddChildren(children []N)
}, T gr.TokenTyper](root *gr.ParseTree[T], lhs T, do func(children []*gr.ParseTree[T]) (N, error)) ([]N, error) {
	if do == nil {
		return nil, gcers.NewErrNilParameter("do")
	} else if root == nil {
		return nil, gcers.NewErrNilParameter("root")
	} else if root.Type() != lhs {
		return nil, fmt.Errorf("expected %q, got %s instead", lhs.String(), root.Type().String())
	}

	var nodes []N

	for root != nil {
		children := root.GetChildren()

		if len(children) == 0 {
			return nil, errors.New("expected at least one child")
		}

		last_child := children[len(children)-1]

		var node N
		var err error

		if last_child.Type() == lhs {
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
