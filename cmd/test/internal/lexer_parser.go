
package internal

import (
	"github.com/PlayerR9/SlParser/lexer"
)

var (
	Lexer *lexer.Lexer[TokenType]
)

func init() {
	builder := lexer.NewBuilder[TokenType]()

	// TODO: Add here your own custom rules...
	
	Lexer = builder.Build()
}