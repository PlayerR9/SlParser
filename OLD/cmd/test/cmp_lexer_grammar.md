// Code generated by SlParser.
package test

import (
	"github.com/PlayerR9/grammar/grammar"
	"github.com/PlayerR9/grammar/lexer"
)

var (
	// internal_lexer is the lexer of the grammar.
	internal_lexer lexer.Lexer[token_type]
)

func init() {
	lex_one := func(l *lexer.Lexer[token_type]) (*grammar.Token[token_type], error) {
		// Lex here anything that matcher doesn't handle...

		panic("Implement me!")
	}

	internal_lexer.WithLexFunc(lex_one)

	// Add here your custom matcher rules.
}