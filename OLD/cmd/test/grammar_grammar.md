// Code generated by SlParser. Do not edit.
package test

import (
	"github.com/PlayerR9/grammar"
)

var (
	// Parser is the complete parser of the grammar.
	Parser grammar.ParsingFunc[token_type]
)

func init() {
	Parser.Init(
		internal_lexer,
		internal_parser,
		ast_builder,
	)
}