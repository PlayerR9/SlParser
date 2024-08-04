package SLParser

import (
	prx "github.com/PlayerR9/SLParser/parser"
	ast "github.com/PlayerR9/grammar/ast"
)

// ParseEbnf parses an EBNF file.
//
// Parameters:
//   - data: The data to parse.
//
// Returns:
//   - *ast.Node[prx.NodeType]: The root node of the AST tree.
func ParseEbnf(data []byte) (*ast.Node[prx.NodeType], error) {
	root, err := prx.Parse(data)
	if err != nil {
		return root, err
	}

	return root, nil
}
