package ast

import (
	gr "github.com/PlayerR9/SlParser/grammar"
	gcers "github.com/PlayerR9/go-errors"
)

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
		return gcers.NewErrInvalidParameter("at must be non-negative")
	}

	if at >= len(children) {
		return NewBadSyntaxTree(at, type_, "")
	}

	tk := children[at]
	if tk == nil {
		return NewBadSyntaxTree(at, type_, "")
	}

	tk_type := tk.Type()

	if tk_type != type_ {
		return NewBadSyntaxTree(at, type_, tk_type.String())
	}

	return nil
}

// LhsDoFunc is a function that does the conversion.
//
// Parameters:
//   - children: The list of children.
//
// Returns:
//   - []N: The converted nodes.
//   - error: if an error occurred.
type LhsDoFunc[N interface {
	AddChildren(children []N)

	Noder
}, T gr.TokenTyper] func(children []*gr.ParseTree[T]) ([]N, error)

// LhsToAst is a function that converts a token to an ast node.
//
// Parameters:
//   - at: The position of the token.
//   - root: The root token. Assumed to be non-nil.
//   - lhs: The lhs token.
//   - do: The function that does the conversion.
//
// Returns:
//   - []N: The converted nodes.
//   - error: if an error occurred.
//
// Errors:
//   - *errors.ErrNilParameter: If 'root' or 'do' is nil.
//   - any other error returned by 'do'.
func LhsToAst[N interface {
	AddChildren(children []N)

	Noder
}, T gr.TokenTyper](at int, children []*gr.ParseTree[T], lhs T, do LhsDoFunc[N, T]) ([]N, error) {
	if do == nil {
		return nil, gcers.NewErrNilParameter("do")
	}

	err := CheckType(children, at, lhs)
	if err != nil {
		return nil, err
	}

	root := children[at]
	var nodes []N

	for root != nil && err == nil {
		children := root.GetChildren()
		if len(children) == 0 {
			break
		}

		last_child := children[len(children)-1]

		var sub_nodes []N

		if last_child.Type() == lhs {
			sub_nodes, err = do(children[:len(children)-1])
			root = last_child
		} else {
			sub_nodes, err = do(children)
			root = nil
		}

		sub_nodes = FilterNonNilNodes(sub_nodes)
		nodes = append(nodes, sub_nodes...)
	}

	return nodes, err
}

func FilterNonNilNodes[N Noder](nodes []N) []N {
	var top int

	for i := 0; i < len(nodes); i++ {
		node := nodes[i]

		if !node.IsNil() {
			nodes[top] = node
			top++
		}
	}

	return nodes[:top:top]
}
