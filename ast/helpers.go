package ast

import (
	"slices"

	slgr "github.com/PlayerR9/SlParser/grammar"
	gslc "github.com/PlayerR9/mygo-lib/slices"
)

// CheckNode checks that the given node's type is one of the given types. If
// the node is nil, an error is returned. If the node is not nil and the type
// of the node is not one of the given types, an error is also returned.
//
// Parameters:
//   - kind: The kind of AST node to check.
//   - node: The node to check.
//   - types: The types to check against.
//
// Returns:
//   - error: Returns an error if the type of the node is not as expected.
func CheckNode(kind string, node Noder, types ...string) error {
	if node == nil {
		return gslc.NewErrNotAsExpected(true, kind, nil, types...)
	}

	got := node.GetType()

	ok := slices.Contains(types, got)
	if ok {
		return nil
	}

	return gslc.NewErrNotAsExpected(true, kind, &got, types...)
}

// CheckToken checks that the given token's type is one of the given types. If
// the token is nil, an error is returned. If the token is not nil and the type
// of the token is not one of the given types, an error is also returned.
//
// Parameters:
//   - kind: The kind of token to check.
//   - token: The token to check.
//   - types: The types to check against.
//
// Returns:
//   - error: Returns an error if the type of the token is not as expected.
func CheckToken(kind string, token *slgr.Token, types ...string) error {
	if token == nil {
		return gslc.NewErrNotAsExpected(true, kind, nil, types...)
	}

	got := token.Type

	ok := slices.Contains(types, got)
	if ok {
		return nil
	}

	return gslc.NewErrNotAsExpected(true, kind, &got, types...)
}
