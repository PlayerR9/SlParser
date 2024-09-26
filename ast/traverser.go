package ast

import (
	"errors"
	"iter"
	"slices"

	gerr "github.com/PlayerR9/go-errors/error"
)

// Traversor is a struct that holds the functions for traversing an ast.
type Traversor[N interface {
	Child() iter.Seq[N]

	Noder
}, I interface {
	NextInfos() []I

	Infoer[N]
}] struct {
	// InitFn is the function that initializes the info.
	//
	// Returns:
	//   - I: The info.
	InitFn func() I

	// DoFn is the function that does the traversal.
	//
	// Parameters:
	//   - node: The node to traverse. Assume non-nil.
	//   - info: The info of the node. Assume non-nil.
	//
	// Returns:
	//   - error: The error that might occur during the traversal.
	DoFn func(node N, info I) error
}

// ReverseDFS is a function that traverses the ast in reverse depth first order.
//
// This means that the children are traversed before the parent; effectively making
// it a traversal that starts from the left-most leaf and goes to the right-most;
// and then from the right-most parent to the left-most parent.
//
// Parameters:
//   - root: The root node of the tree.
//
// Returns:
//   - error: The error that might occur during the traversal.
func (t Traversor[N, I]) MakeReverseDFS() func(root N) error {
	var root_inf I

	if t.InitFn == nil {
		root_inf = *new(I)
	} else {
		root_inf = t.InitFn()
	}

	var fn func(root N) error

	if t.DoFn == nil {
		fn = func(root N) error {
			root_inf.Init(root, nil)

			stack := []I{root_inf}

			var inner_err error
			var last_top I

			for len(stack) > 0 && inner_err == nil {
				top := stack[len(stack)-1]
				last_top = top

				if top.IsNil() {
					inner_err = errors.New("node found to be nil")
					continue
				}

				if top.IsSeen() {
					stack = stack[:len(stack)-1]
					continue
				}

				// Add its children to the stack.

				children := top.NextInfos()
				if len(children) > 0 {
					slices.Reverse(children)

					stack = append(stack, children...)
				}

				top.See()
			}

			if inner_err == nil {
				return nil
			}

			outer_err := gerr.New(BadSyntaxTree, inner_err.Error())

			for frame := range last_top.Frame() {
				outer_err.AddFrame(frame)
			}

			return outer_err
		}
	} else {
		fn = func(root N) error {
			root_inf.Init(root, nil)

			stack := []I{root_inf}

			var inner_err error
			var last_top I

			for len(stack) > 0 && inner_err == nil {
				top := stack[len(stack)-1]
				last_top = top

				if top.IsNil() {
					inner_err = errors.New("node found to be nil")
					continue
				}

				if top.IsSeen() {
					stack = stack[:len(stack)-1]

					inner_err = t.DoFn(top.Node(), last_top)
					continue
				}

				// Add its children to the stack.

				children := top.NextInfos()
				if len(children) > 0 {
					slices.Reverse(children)

					stack = append(stack, children...)
				}

				top.See()
			}

			if inner_err == nil {
				return nil
			}

			outer_err := gerr.New(BadSyntaxTree, inner_err.Error())

			for frame := range last_top.Frame() {
				outer_err.AddFrame(frame)
			}

			return outer_err
		}
	}

	return fn
}

// ReverseDFS is a function that traverses the ast in reverse depth first order.
//
// This means that the children are traversed before the parent; effectively making
// it a traversal that starts from the left-most leaf and goes to the right-most;
// and then from the right-most parent to the left-most parent.
//
// Parameters:
//   - root: The root node of the tree.
//
// Returns:
//   - error: The error that might occur during the traversal.
func (t Traversor[N, I]) ReverseDFS(root N) error {
	var root_inf I

	if t.InitFn == nil {
		root_inf = *new(I)
	} else {
		root_inf = t.InitFn()
	}

	if t.DoFn == nil {
		t.DoFn = func(node N, info I) error {
			return nil
		}
	}

	root_inf.Init(root, nil)

	stack := []I{root_inf}

	var inner_err error
	var last_top I

	for len(stack) > 0 && inner_err == nil {
		top := stack[len(stack)-1]
		last_top = top

		if top.IsNil() {
			inner_err = errors.New("node found to be nil")
			continue
		}

		if top.IsSeen() {
			stack = stack[:len(stack)-1]

			inner_err = t.DoFn(top.Node(), last_top)
			continue
		}

		// Add its children to the stack.

		children := top.NextInfos()
		if len(children) > 0 {
			slices.Reverse(children)

			stack = append(stack, children...)
		}

		top.See()
	}

	if inner_err == nil {
		return nil
	}

	outer_err := gerr.New(BadSyntaxTree, inner_err.Error())

	for frame := range last_top.Frame() {
		outer_err.AddFrame(frame)
	}

	return outer_err
}
