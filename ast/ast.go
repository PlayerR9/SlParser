package ast

import (
	slgr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/mygo-lib/common"
)

// ASTMaker is an AST maker.
type ASTMaker[T any] struct {
	// table is the AST maker table.
	table map[string]ToASTFn[T]
}

// Make makes an AST node from a token.
//
// Parameters:
//   - token: The token to make an AST node from.
//
// Returns:
//   - []T: The AST nodes.
//   - error: An error if the evaluation failed.
func (b *ASTMaker[T]) Make(token *slgr.Token) ([]T, error) {
	if b == nil {
		return nil, common.ErrNilReceiver
	} else if token == nil {
		return nil, common.NewErrNilParam("token")
	}

	type_ := token.Type

	fn, ok := b.table[type_]
	if !ok || fn == nil {
		return nil, NewErrUnsupportedType(true, type_)
	}

	nodes, err := fn(token)
	return nodes, err
}
