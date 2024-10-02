package ast

import (
	"errors"
	"iter"
	"slices"

	gers "github.com/PlayerR9/go-errors"
)

// Traversor is a struct that holds the functions for traversing an ast.
type Traversor[N interface {
	Child() iter.Seq[N]

	Noder
}, I interface {
	NextInfos() ([]I, error)

	Infoer[N]
}] struct {
	// InitFn is the function that initializes the info.
	//
	// Parameters:
	//   - node: The node the info is about.
	//   - frames: The frames of the node. Used for stack traces.
	//
	// Returns:
	//   - I: The info.
	//   - error: The error that might occur during the initialization, such as
	//   the node being nil.
	InitFn func(node N, frames []string) (I, error)

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
func (t Traversor[N, I]) ReverseDFS(root N) error {
	if root.IsNil() {
		return gers.NewErrNilParameter("Traversor.ReverseDFS()", "root")
	}

	if t.DoFn == nil {
		t.DoFn = func(node N, info I) error {
			return nil
		}
	}

	root_inf, err := t.InitFn(root, nil)
	if err != nil {
		return err
	}

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

		children, err := top.NextInfos()
		if err != nil {
			inner_err = err
			continue
		}

		if len(children) > 0 {
			slices.Reverse(children)

			stack = append(stack, children...)
		}

		top.See()
	}

	if inner_err == nil {
		return nil
	}

	outer_err := gers.New(BadSyntaxTree, inner_err.Error())

	for frame := range last_top.Frame() {
		outer_err.AddFrame(frame)
	}

	return outer_err
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
func ReverseDFS[N interface {
	Child() iter.Seq[N]

	Noder
}, I interface {
	NextInfos() ([]I, error)

	Infoer[N]
}](do_fn func(node N, info I) error) func(info I) error {
	if do_fn == nil {
		do_fn = func(node N, info I) error {
			return nil
		}
	}

	return func(info I) error {
		if info.IsNil() {
			return gers.NewErrNilParameter("ast.ReverseDFS()", "info")
		}

		stack := []I{info}

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

				inner_err = do_fn(top.Node(), last_top)
				continue
			}

			// Add its children to the stack.

			children, err := top.NextInfos()
			if err != nil {
				inner_err = err
				continue
			}

			if len(children) > 0 {
				slices.Reverse(children)

				stack = append(stack, children...)
			}

			top.See()
		}

		if inner_err == nil {
			return nil
		}

		outer_err := gers.New(BadSyntaxTree, inner_err.Error())

		for frame := range last_top.Frame() {
			outer_err.AddFrame(frame)
		}

		return outer_err
	}
}
