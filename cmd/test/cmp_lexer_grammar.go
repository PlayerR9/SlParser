// Code generated by SlParser.
package test

import (
	"github.com/PlayerR9/grammar/grammar"
	"github.com/PlayerR9/grammar/lexing"
)

var (
	// matcher is the matcher of the grammar.
	matcher lexing.Matcher[token_type]
)

func init() {
	// Add here your custom matcher rules.
}

var (
	// internal_lexer is the lexer of the grammar.
	internal_lexer *lexing.Lexer[token_type]
)

func init() {
	lex_one := func(l *lexing.Lexer[token_type]) (*grammar.Token[token_type], error) {
		// Lex here anything that matcher doesn't handle...

		panic("Implement me!")
	}

	internal_lexer = lexing.NewLexer(lex_one, matcher)
}
