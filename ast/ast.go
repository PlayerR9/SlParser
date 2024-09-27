package ast

import (
	"fmt"

	gr "github.com/PlayerR9/SlParser/grammar"
	gers "github.com/PlayerR9/go-errors"
)

// ToAstFunc is a function that converts a token to an ast node.
//
// Parameters:
//   - tree: The parse tree to convert. Assumed to be non-nil.
//
// Returns:
//   - N: The ast node.
//   - error: if an error occurred.
type ToAstFunc[N interface {
	AddChildren(children []N)

	Noder
}, T gr.TokenTyper] func(tree *gr.ParseTree[T]) (N, error)

// AstMaker is an ast maker.
type AstMaker[N interface {
	AddChildren(children []N)

	Noder
}, T gr.TokenTyper] map[T]ToAstFunc[N, T]

// FnOf is a function that gets the ast function of a type.
//
// Parameters:
//   - type_: The type of the token.
//
// Returns:
//   - ToAstFunc[N, T]: The ast function.
//   - bool: true if the ast function was found, false otherwise.
func (am AstMaker[N, T]) FnOf(type_ T) (ToAstFunc[N, T], bool) {
	if len(am) == 0 {
		return nil, false
	}

	fn, ok := am[type_]
	return fn, ok
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
	zero := *new(N)

	if root == nil {
		return zero, gers.NewErrNilParameter("root")
	}

	type_ := root.Type()

	fn, ok := am.FnOf(type_)
	if !ok || fn == nil {
		err := NewUnregisteredType(type_, type_.String())
		return zero, err
	}

	node, err := fn(root)
	if err != nil {
		return zero, err
	}

	if node.IsNil() {
		return zero, fmt.Errorf("the returned node must not be nil")
	}

	return node, nil
}
