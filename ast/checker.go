package ast

import (
	"iter"

	"github.com/PlayerR9/SlParser/ast/internal"
	gers "github.com/PlayerR9/go-errors"
)

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

	internal.Noder
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

	internal.Noder
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

	internal.Noder
}](check_fn CheckNodeFn[N]) CheckASTWithLimit[N] {
	if check_fn == nil {
		return func(node N, limit int) error {
			return gers.NewErrInvalidParameter("check_fn")
		}
	}

	trav := Traversor[N, *internal.CheckerInfo[N]]{
		InitFn: nil,
		DoFn: func(node N, info *internal.CheckerInfo[N]) error {
			gers.AssertNotNil(info, "info")

			return check_fn(node)
		},
	}

	return func(node N, limit int) error {
		trav.InitFn = func() *internal.CheckerInfo[N] {
			info := internal.NewCheckerInfo[N](limit)
			return &info
		}
		return trav.ReverseDFS(node)
	}
}
