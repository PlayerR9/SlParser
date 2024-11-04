package ast

import (
	slgr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/mygo-lib/common"
)

// ToASTFn is a function that converts a token to a node.
//
// Parameters:
//   - token: the token to convert. Assumed to not be nil.
//
// Returns:
//   - []N: the nodes created by the token.
//   - error: an error if the token cannot be converted.
type ToASTFn[N any] func(token *slgr.Token) ([]N, error)

// Builder is a builder for ASTs.
type Builder[N any] struct {
	// table is the underlying table of the builder.
	table map[string]ToASTFn[N]
}

// Register registers a token type with the builder. Does nothing if
// the function is nil.
//
// Parameters:
//   - type_: the token type to register.
//   - fn: the function to convert tokens to nodes.
//
// Returns:
//   - error: an error if the receiver is nil.
func (b *Builder[N]) Register(type_ string, fn ToASTFn[N]) error {
	if fn == nil {
		return nil
	} else if b == nil {
		return common.NewErrNilParam("b")
	}

	if b.table == nil {
		b.table = make(map[string]ToASTFn[N])
	}

	b.table[type_] = fn

	return nil
}

// Build builds the AST.
//
// Returns:
//   - *ASTMaker[N]: the AST maker. Never returns nil.
func (b Builder[N]) Build() *ASTMaker[N] {
	if len(b.table) == 0 {
		return &ASTMaker[N]{
			table: make(map[string]ToASTFn[N]),
		}
	}

	table := make(map[string]ToASTFn[N], len(b.table))

	for k, v := range b.table {
		table[k] = v
	}

	return &ASTMaker[N]{
		table: table,
	}
}

// Reset resets the builder for reuse.
func (b *Builder[N]) Reset() {
	if b == nil {
		return
	}

	if len(b.table) > 0 {
		clear(b.table)
		b.table = nil
	}
}
