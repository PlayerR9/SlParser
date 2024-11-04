package ast

import (
	"errors"
	"fmt"

	slgr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/mygo-lib/common"
)

type AST[N any] interface {
	Make(token *slgr.Token) ([]N, error)
}

type baseAST[N any] struct {
	table map[string]ToAstFunc[N]
}

func (b *baseAST[N]) Make(token *slgr.Token) ([]N, error) {
	if token == nil {
		return nil, nil
	} else if b == nil {
		return nil, common.NewErrNilParam("b")
	}

	type_ := token.Type

	if len(b.table) == 0 {
		return nil, fmt.Errorf("token type (%q) is not supported", type_)
	}

	fn, ok := b.table[type_]
	if !ok {
		return nil, fmt.Errorf("token type (%q) is not supported", type_)
	} else if fn == nil {
		return nil, errors.New("fn must not be nil")
	}

	nodes, err := fn(token)
	return nodes, err
}
