package grammar

import (
	"errors"
	"slices"

	faults "github.com/PlayerR9/go-fault"
	tr "github.com/PlayerR9/mygo-lib/CustomData/tree"
)

/////////////////////////////////////////////////////////

var (
	// EOFToken is the end of file token.
	EOFToken *tr.Node
)

func init() {
	EOFToken = NewToken(-1, EtEOF, "", nil)
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
func CheckNode(idx int, kind string, node *tr.Node, types ...string) error {
	if node == nil {
		fault := ErrNotAsExpected.New()

		faults.Faults.AddContext(fault, "Kind", kind)
		faults.Faults.AddContext(fault, "Index", idx)
		faults.Faults.AddContext(fault, "Types", types)
		faults.Faults.AddContext(fault, "Got", nil)

		return fault
	}

	var got string

	switch info := node.Info.(type) {
	case *NodeData:
		got = info.Type
	case *TokenData:
		got = info.Type
	default:
		return errors.New("node has the wrong data")
	}

	ok := slices.Contains(types, got)
	if ok {
		return nil
	}

	fault := ErrNotAsExpected.New()

	faults.Faults.AddContext(fault, "Kind", kind)
	faults.Faults.AddContext(fault, "Index", idx)
	faults.Faults.AddContext(fault, "Types", types)
	faults.Faults.AddContext(fault, "Got", &got)

	return fault
}
