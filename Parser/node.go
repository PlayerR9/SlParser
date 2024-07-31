package Parser

// NodeType represents the type of a node in the AST tree.
type NodeType int

const (
	SourceNode NodeType = iota
	RuleNode
	IdentifierNode
	OrNode
)

// String implements the NodeTyper interface.
func (t NodeType) String() string {
	return [...]string{
		"Source",
		"Rule",
		"Identifier",
		"OR",
	}[t]
}
