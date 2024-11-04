package ast

// Noder is the interface that all nodes must implement.
type Noder interface {
	// GetType returns the type of the node.
	//
	// Returns:
	//   - string: The type of the node.
	GetType() string
}
