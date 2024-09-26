package ast

import (
	"errors"
	"iter"
	"slices"

	gers "github.com/PlayerR9/go-errors"
	gerr "github.com/PlayerR9/go-errors/error"
)

var (
	// IgnoreInfo is an error that is used to ignore info. Readers must return
	// this error as is and not wrap it as callers check for this error using
	// ==.
	IgnoreInfo error
)

func init() {
	IgnoreInfo = errors.New("ignore info")
}

// NewInfoFn is a function that creates an info for a node.
//
// Parameters:
//   - node: The node to create the info for.
//
// Returns:
//   - I: The info.
//   - error: An error if the info could not be created.
//
// Return the error IgnoreInfo to not add the info to the table.
type NewInfoFn[N interface {
	Child() iter.Seq[N]

	Noder
}, I interface {
	NextInfos() []I

	Infoer[N]
}] func(node N) (I, error)

// InfoTableOfFn is a function that creates an info table for a node.
//
// Parameters:
//   - root: The root node.
//
// Returns:
//   - map[N]I: The info table.
//   - error: An error if the info table could not be created.
type InfoTableOfFn[N interface {
	Child() iter.Seq[N]

	Noder
}, I interface {
	NextInfos() []I

	Infoer[N]
}] func(root N) (map[N]I, error)

// MakeInfoTable is a function that creates an info table for a node.
//
// Parameters:
//   - fn: The new info function.
//
// Returns:
//   - InfoTableOfFn: The info table function.
//
// Whenever the info function returns nil, the corresponding node will be removed from the table.
//
// If 'fn' is nil, then the function returns a function that returns an error.
func MakeInfoTable[N interface {
	Child() iter.Seq[N]

	Noder
}, I interface {
	NextInfos() []I

	Infoer[N]
}](fn NewInfoFn[N, I]) InfoTableOfFn[N, I] {
	if fn == nil {
		return func(root N) (map[N]I, error) {
			return nil, gers.NewErrNilParameter("fn")
		}
	}

	return func(root N) (map[N]I, error) {
		table := make(map[N]I)

		root_info := *new(I)
		root_info.Init(root, nil)

		stack := []I{root_info}

		var inner_err error
		var last_top I

		for len(stack) > 0 && inner_err == nil {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			last_top = top

			if top.IsNil() {
				inner_err = errors.New("found node to be nil")
				continue
			}

			node := top.Node()

			info, err := fn(node)
			if err == nil {
				table[node] = info
			} else if err != IgnoreInfo {
				inner_err = err
				continue
			}

			nexts := top.NextInfos()

			if len(nexts) > 0 {
				slices.Reverse(nexts)
				stack = append(stack, nexts...)
			}
		}

		if inner_err == nil {
			return table, nil
		}

		gers.AssertNotNil(last_top, "last_top")

		outer_err := gerr.New(BadSyntaxTree, inner_err.Error())

		for frame := range last_top.Frame() {
			outer_err.AddFrame("", frame)
		}

		return nil, outer_err
	}
}
