package grammar

import (
	"slices"

	faults "github.com/PlayerR9/go-fault"
)

var (
	// EOFToken is the end of file token.
	EOFToken *Token
)

func init() {
	EOFToken = NewToken(-1, EtEOF, "")
}

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
func CheckNode(idx int, kind string, node *Node, types ...string) faults.Fault {
	if node == nil {
		fault := ErrNotAsExpected.New()

		faults.Faults.AddContext(fault, "Kind", kind)
		faults.Faults.AddContext(fault, "Index", idx)
		faults.Faults.AddContext(fault, "Types", types)
		faults.Faults.AddContext(fault, "Got", nil)

		return fault
	}

	got := node.Type

	ok := slices.Contains(types, got)
	if ok {
		return nil
	}

	err := ErrNotAsExpected.New()

	faults.Faults.AddContext(err, "Kind", kind)
	faults.Faults.AddContext(err, "Index", idx)
	faults.Faults.AddContext(err, "Types", types)
	faults.Faults.AddContext(err, "Got", &got)

	return err
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
func CheckToken(idx int, kind string, token *Token, types ...string) error {
	if token == nil {
		err := ErrNotAsExpected.New()

		faults.Faults.AddContext(err, "Kind", kind)
		faults.Faults.AddContext(err, "Index", idx)
		faults.Faults.AddContext(err, "Types", types)
		faults.Faults.AddContext(err, "Got", nil)

		return err
	}

	got := token.Type

	ok := slices.Contains(types, got)
	if ok {
		return nil
	}

	err := ErrNotAsExpected.New()

	faults.Faults.AddContext(err, "Kind", kind)
	faults.Faults.AddContext(err, "Index", idx)
	faults.Faults.AddContext(err, "Types", types)
	faults.Faults.AddContext(err, "Got", &got)

	return err
}
