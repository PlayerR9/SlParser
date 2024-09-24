package grammar

import (
	gcslc "github.com/PlayerR9/go-commons/slices"
	dba "github.com/PlayerR9/go-debug/assert"
	gcers "github.com/PlayerR9/go-errors"
)

// Combine creates a new ParseTree by combining the subtrees.
//
// The root of the new tree is a Token with the given type and the position of the first subtree,
// and the lookahead of the last subtree. The children of the new tree are the
// subtrees.
//
// Parameters:
//   - type_: The type of the new tree.
//   - subtrees: The subtrees.
//
// Returns:
//   - *ParseTree[T]: The new tree.
//   - error: an error if there are no subtrees. Nil parse trees are ignored.
func Combine[T TokenTyper](type_ T, subtrees []*ParseTree[T]) (*ParseTree[T], error) {
	subtrees = gcslc.FilterZeroValues(subtrees)
	if len(subtrees) == 0 {
		return nil, gcers.NewErrInvalidParameter("cannot combine an empty list of subtrees")
	}

	last_tk := subtrees[len(subtrees)-1].Root()
	first_tk := subtrees[0].Root()

	root_tk := &Token[T]{
		Type:      type_,
		Lookahead: last_tk.Lookahead,
		Pos:       first_tk.Pos,
	}

	tree, err := NewTree(root_tk)
	dba.AssertErr(err, "NewTree(root_tk)")
	tree.SetChildren(subtrees)

	return tree, nil
}
