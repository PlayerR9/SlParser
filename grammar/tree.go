package grammar

import (
	"fmt"
	"strings"
)

// TreeNode is an interface that represents a node in a tree.
type TreeNode interface {
	fmt.Stringer
}

// recTreeToString writes a string representation of a tree to the given
// strings.Builder. The string representation is a recursive stringification
// of the tree, with each node indented under its parent. The indent string
// is used to indent each level of the tree.
//
// Parameters:
//   - b: The strings.Builder to write the string representation to.
//   - indent: The string to use to indent each level of the tree.
//   - node: The root of the tree to stringify.
func recTreeToString[T interface {
	GetChildren() []T

	TreeNode
}](b *strings.Builder, indent string, node T) {
	_, _ = b.WriteString(indent)
	_, _ = b.WriteString(node.String())
	_, _ = b.WriteRune('\n')

	children := node.GetChildren()

	for _, child := range children {
		recTreeToString(b, indent+"   ", child)
	}
}

// TreeToString takes a tree and returns a string representation of it. The
// string representation is a recursive stringification of the tree, with each
// node indented under its parent. The stringification is done using the
// TreeNode interface, which is implemented by all nodes in the tree.
//
// Parameters:
//   - root: The root of the tree to stringify.
//
// Returns:
//   - string: A string representation of the tree.
func TreeToString[T interface {
	GetChildren() []T

	TreeNode
}](root T) string {
	var builder strings.Builder

	recTreeToString(&builder, "", root)

	return builder.String()
}
