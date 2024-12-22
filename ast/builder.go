package ast

import (
	slgr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/mygo-lib/common"
)

// ASTFn is a function that converts a token tree into a node tree.
//
// Parameters:
//   - tk: The root token of the tree. (Assumed to not be nil)
//
// Returns:
//   - N: The node resulting from the conversion of the root token.
//   - error: An error if the conversion process fails.
//
// Errors:
//   - any error: Implementation-specific error.
type ASTFn[N slgr.TreeNode] func(tk *slgr.Token) (N, error)

// ASTMaker is a map of functions that convert a token tree into a node tree.
//
// Parameters:
//   - tk: The root token of the tree. (Assumed to not be nil)
//
// Returns:
//   - N: The node resulting from the conversion of the root token.
//   - error: An error if the conversion process fails.
//
// Errors:
//   - any error: Implementation-specific error.
type ASTMaker[N slgr.TreeNode] map[string]ASTFn[N]

// Builder is a struct that builds an ASTMaker.
type Builder[N slgr.TreeNode] struct {
	// table is the map of functions that convert a token tree into a node tree.
	table map[string]ASTFn[N]
}

// Reset implements common.Resetter.
func (b *Builder[N]) Reset() error {
	if b == nil {
		return common.ErrNilReceiver
	}

	if len(b.table) == 0 {
		return nil
	}

	clear(b.table)
	b.table = nil

	return nil
}

// Add adds a function to the ASTMaker.
//
// Parameters:
//   - type_: The type of the token that the function will be associated with.
//   - fn: The function that converts a token tree into a node tree.
//
// Returns:
//   - error: An error if the receiver is nil or if the function is nil.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
//   - common.ErrBadParam: If the function is nil.
func (b *Builder[N]) Add(type_ string, fn ASTFn[N]) error {
	if b == nil {
		return common.ErrNilReceiver
	}

	if fn == nil {
		err := common.NewErrNilParam("fn")
		return err
	}

	if b.table == nil {
		b.table = make(map[string]ASTFn[N])
	}

	b.table[type_] = fn

	return nil
}

// Build returns the table of functions that convert a token tree into a node tree.
//
// The returned table is a copy of the internal table of the builder.
//
// Returns:
//   - ASTMaker[N]: The table of functions that convert a token tree into a node tree. Never
//     returns nil.
func (b Builder[N]) Build() ASTMaker[N] {
	if len(b.table) == 0 {
		table := make(ASTMaker[N], 0)
		return table
	}

	table := make(ASTMaker[N], len(b.table))

	for k, v := range b.table {
		table[k] = v
	}

	return table
}
