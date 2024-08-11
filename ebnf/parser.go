package ebnf

import (
	grammar "github.com/PlayerR9/grammar"
)

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

var (
	Parser *grammar.Parser[*Node, token_type]
)

func init() {
	Parser = grammar.NewParser(
		internal_lexer,
		internal_parser,
		ast_builder,
	)
}
