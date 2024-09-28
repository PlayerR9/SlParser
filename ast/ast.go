package ast

import (
	"fmt"

	gr "github.com/PlayerR9/SlParser/grammar"
	gers "github.com/PlayerR9/go-errors"
)

// ToNodeFn is a function that converts a token to an ast node.
//
// Parameters:
//   - tree: The parse tree to convert. Assumed to be non-nil.
//
// Returns:
//   - N: The ast node.
//   - error: if an error occurred.
type ToNodeFn[N interface {
	AddChildren(children []N)

	Noder
}, T gr.TokenTyper] func(tree *gr.ParseTree[T]) (N, error)

// TransitionFn is a function that converts a token into an intermediate
// stage of a set of nodes.
//
// Parameters:
//   - tree: The parse tree to convert. Assumed to be non-nil.
//
// Returns:
//   - []N: The intermediate nodes.
//   - error: if an error occurred.
type TransitionFn[N interface {
	AddChildren(children []N)

	Noder
}, T gr.TokenTyper] func(tree *gr.ParseTree[T]) ([]N, error)

// AstMaker is an ast maker.
type AstMaker[N interface {
	AddChildren(children []N)

	Noder
}, T gr.TokenTyper] struct {
	transformations map[T]ToNodeFn[N, T]
	transitions     map[T]TransitionFn[N, T]
}

// NewAstMaker is a function that creates a new ast maker.
//
// Returns:
//   - AstMaker: The created ast maker. Never returns nil.
func NewAstMaker[N interface {
	AddChildren(children []N)

	Noder
}, T gr.TokenTyper]() AstMaker[N, T] {
	return AstMaker[N, T]{
		transformations: make(map[T]ToNodeFn[N, T]),
		transitions:     make(map[T]TransitionFn[N, T]),
	}
}

// AddTransformation is a function that adds a transformation.
//
// Parameters:
//   - type_: The type of the token.
//   - fn: The transformation function.
//
// If 'fn' is nil, the transformation will be removed. Previously registered
// transformations with the same type will be overwritten.
func (am AstMaker[N, T]) AddTransformation(type_ T, fn ToNodeFn[N, T]) {
	if fn == nil {
		delete(am.transformations, type_)
	} else {
		am.transformations[type_] = fn
	}
}

// AddTransition is a function that adds a transition.
//
// Parameters:
//   - type_: The type of the token.
//   - fn: The transition function.
//
// If 'fn' is nil, the transition will be removed. Previously registered
// transitions with the same type will be overwritten.
func (am AstMaker[N, T]) AddTransition(type_ T, fn TransitionFn[N, T]) {
	if fn == nil {
		delete(am.transitions, type_)
	} else {
		am.transitions[type_] = fn
	}
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

	var fn ToNodeFn[N, T]

	if len(am.transformations) > 0 {
		val, ok := am.transformations[type_]
		if !ok {
			fn = nil
		} else {
			fn = val
		}
	}

	if fn == nil {
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

	node.SetPosition(root.Pos())

	return node, nil
}

// Apply is a function that applies a transition function.
//
// Parameters:
//   - root: The root token.
//
// Returns:
//   - []N: The intermediate nodes.
//   - error: if an error occurred.
func (am AstMaker[N, T]) Apply(root *gr.ParseTree[T]) ([]N, error) {
	if root == nil {
		return nil, gers.NewErrNilParameter("root")
	}

	type_ := root.Type()

	var fn TransitionFn[N, T]

	if len(am.transitions) > 0 {
		val, ok := am.transitions[type_]
		if !ok {
			fn = nil
		} else {
			fn = val
		}
	}

	if fn == nil {
		err := NewUnregisteredType(type_, type_.String())
		return nil, err
	}

	nodes, err := fn(root)
	return nodes, err
}
