package ast

import (
	"fmt"
	"iter"

	gers "github.com/PlayerR9/go-errors"
	"github.com/PlayerR9/go-errors/assert"
)

// NodeTyper is an interface representing a node type.
type NodeTyper interface {
	~int

	// String returns the string representation of the node type.
	//
	// Returns:
	//   - string: the string representation of the node type.
	String() string
}

// CheckASTWithLimit checks if the given ast node is valid, up to a given
// limit depth. If limit is negative, it will check all the way down to the
// leaves. On the other hand, if limit is 0, it will only check if the node
// is nil or not.
//
// Parameters:
//   - node: The node to check.
//   - limit: The maximum depth to check. If negative, it will check all the
//     way down to the leaves.
//
// Returns:
//   - error: If the node is not valid, an error describing the problem will
//     be returned. Otherwise, nil is returned.
//
// See documentation to learn what counts as a valid node.
type CheckASTWithLimit[N interface {
	Child() iter.Seq[N]

	Noder
}] func(node N, limit int) error

// CheckNodeFn is a function that checks if the given node is valid.
//
// Parameters:
//   - node: The node to check. Assumed to be non-nil.
//
// Returns:
//   - error: An error if the node is invalid. Otherwise, nil.
type CheckNodeFn[N interface {
	Child() iter.Seq[N]

	Noder
}] func(node N) error

// MakeCheckFn is a function that creates a CheckASTWithLimit function
// with the given check function.
//
// Parameters:
//   - check_fn: The check function.
//
// Returns:
//   - CheckASTWithLimit: The created function. Never returns nil.
//
// If the check_fn is nil, a function that returns an error will be returned.
func MakeCheckFn[N interface {
	Child() iter.Seq[N]
	GetType() T

	Noder
}, T NodeTyper](table map[T]CheckNodeFn[N]) CheckASTWithLimit[N] {
	var do_fn func(node N, _ *CheckerInfo[N]) error

	if len(table) == 0 {
		do_fn = func(node N, _ *CheckerInfo[N]) error {
			if node.IsNil() {
				return gers.NewErrNilParameter("ast.CheckASTWithLimit()", "node")
			}

			type_ := node.GetType()

			return fmt.Errorf("unknown node type: %s", type_.String())
		}
	} else {
		do_fn = func(node N, _ *CheckerInfo[N]) error {
			if node.IsNil() {
				return gers.NewErrNilParameter("ast.CheckASTWithLimit()", "node")
			}

			type_ := node.GetType()

			check_fn, ok := table[type_]
			if !ok {
				return fmt.Errorf("unknown node type: %s", type_.String())
			}

			return check_fn(node)
		}
	}

	fn := ReverseDFS(do_fn)

	return func(node N, limit int) error {
		if node.IsNil() {
			return gers.NewErrNilParameter("ast.CheckASTWithLimit()", "node")
		}

		info := assert.New(
			NewCheckerInfo(node, nil, limit),
		)

		return fn(info)
	}
}
