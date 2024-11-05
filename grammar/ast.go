package grammar

import (
	"github.com/PlayerR9/mygo-lib/common"
)

// astT is for internal use.
type astT struct{}

var (
	// AST is the namespace that allows for the creation of AST nodes.
	AST astT
)

func init() {
	AST = astT{}
}

// ToASTFn is a function that converts a token to a node.
//
// Parameters:
//   - token: the token to convert. Assumed to not be nil.
//
// Returns:
//   - []*Node: the nodes created by the token.
//   - error: an error if the token cannot be converted.
type ToASTFn func(token *Token) ([]*Node, error)

// Make makes an AST node from a token.
//
// Parameters:
//   - ast: The AST maker.
//   - token: The token to make an AST node from.
//
// Returns:
//   - []*Node: The AST nodes.
//   - error: An error if the evaluation failed.
func (astT) Make(ast map[string]ToASTFn, token *Token) ([]*Node, error) {
	if ast == nil {
		return nil, common.NewErrNilParam("ast")
	} else if token == nil {
		return nil, common.NewErrNilParam("token")
	}

	type_ := token.Type

	fn, ok := ast[type_]
	if !ok || fn == nil {
		return nil, NewErrUnsupportedType(true, type_)
	}

	nodes, err := fn(token)
	return nodes, err
}
