package parser

import (
	tr "github.com/PlayerR9/mygo-lib/CustomData/tree"
	"github.com/PlayerR9/mygo-lib/common"
)

/////////////////////////////////////////////////////////

// ParseTree the parse tree output by the parser.
type ParseTree struct {
	// tree is the underlying parse tree.
	tree *tr.Tree
}

// String implements the fmt.Stringer interface.
func (pt ParseTree) String() string {
	return pt.tree.String()
}

// NewParseTree creates a new parse tree from a token.
//
// Parameters:
//   - tk: The root token of the parse tree. Must not be nil.
//
// Returns:
//   - *ParseTree: The new parse tree.
//   - error: An error if the token is nil.
func NewParseTree(tk *tr.Node) (*ParseTree, error) {
	if tk == nil {
		return nil, common.NewErrNilParam("tk")
	}

	tree := tr.NewTree(tk)

	return &ParseTree{
		tree: tree,
	}, nil
}

// Equals checks whether the given parse tree is equal to the current parse tree.
//
// Parameters:
//   - other: The parse tree to compare with. May be nil.
//
// Returns:
//   - bool: True if the parse trees are equal, false otherwise.
func (pt ParseTree) Equals(other *ParseTree) bool {
	return other != nil && tr.Equals(pt.tree, other.tree)
}

// Root returns the root token of the parse tree.
//
// Returns:
//   - *tr.Node: The root token of the parse tree.
func (pt ParseTree) Root() *tr.Node {
	return pt.tree.Root()
}

// Slice returns a slice of all the tokens in the parse tree in a pre-order view.
//
// Returns:
//   - []*tr.Node: A slice of all the tokens in the parse tree.
func (pt ParseTree) Slice() []*tr.Node {
	slice := make([]*tr.Node, 0, pt.tree.Size())

	visit_fn := func(t *tr.Node) error {
		slice = append(slice, t)
		return nil
	}

	_ = tr.View.Preorder(pt.tree, visit_fn)

	return slice
}
