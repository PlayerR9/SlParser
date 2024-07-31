package SLParser

import (
	"os"

	prx "github.com/PlayerR9/SLParser/parser"
	ast "github.com/PlayerR9/grammar/ast"
)

// ParseEbnf parses an EBNF file.
//
// Parameters:
//   - loc: The location of the EBNF file.
//
// Returns:
//   - *ast.Node[prx.NodeType]: The root node of the AST tree.
func ParseEbnf(loc string) (*ast.Node[prx.NodeType], error) {
	data, err := os.ReadFile(loc)
	if err != nil {
		return nil, err
	}

	root, err := prx.Parse(data)
	if err != nil {
		return root, err
	}

	return root, nil
}
