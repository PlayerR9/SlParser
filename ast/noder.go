package ast

// Noder is an interface for nodes in the AST
type Noder interface {
	comparable

	// IsNil checks whether the node is nil.
	//
	// Returns:
	//   - bool: true if the node is nil, false otherwise.
	IsNil() bool

	// String returns the string representation of the node.
	//
	// Returns:
	//   - string: the string representation of the node.
	String() string

	// SetPosition sets the position of the node. Does nothing if the node is nil.
	//
	// Parameters:
	//   - pos: The position of the node.
	SetPosition(pos int)
}
