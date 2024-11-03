package parser

import (
	slgr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/mygo-lib/common"
	tr "github.com/PlayerR9/mygo-lib/trees"
)

// ParseTree the parse tree output by the parser.
type ParseTree struct {
	// tree is the underlying parse tree.
	tree tr.Tree[*slgr.Token]
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
func NewParseTree(tk *slgr.Token) (*ParseTree, error) {
	if tk == nil {
		return nil, common.NewErrNilParam("tk")
	}

	tree := tr.New(tk)

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
//   - *slgr.Token: The root token of the parse tree.
func (pt ParseTree) Root() *slgr.Token {
	return pt.tree.Root()
}

// Slice returns a slice of all the tokens in the parse tree in a pre-order view.
//
// Returns:
//   - []*slgr.Token: A slice of all the tokens in the parse tree.
func (pt ParseTree) Slice() []*slgr.Token {
	slice := make([]*slgr.Token, 0, pt.tree.Size())

	visit_fn := func(t *slgr.Token) error {
		slice = append(slice, t)
		return nil
	}

	_ = tr.PreorderView(pt.tree, visit_fn)

	return slice
}
