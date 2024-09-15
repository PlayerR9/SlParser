package parser

import (
	dbg "github.com/PlayerR9/go-debug/assert"
	lxr "github.com/PlayerR9/grammar/lexer"
)

var (
	lexer *lxr.Lexer[TokenType]
)

func init() {
	var builder lxr.Builder[TokenType]

	builder.Register('*', TttMod, func(lexer *lxr.Lexer[TokenType]) (string, error) {
		_, err := lexer.NextRune()
		dbg.AssertErr(err, "lexer.NextRune()")

		return "*", nil
	})

	builder.Register('?', TttMod, func(lexer *lxr.Lexer[TokenType]) (string, error) {
		_, err := lexer.NextRune()
		dbg.AssertErr(err, "lexer.NextRune()")

		return "?", nil
	})

	builder.Register('+', TttMod, func(lexer *lxr.Lexer[TokenType]) (string, error) {
		_, err := lexer.NextRune()
		dbg.AssertErr(err, "lexer.NextRune()")

		return "+", nil
	})

	lexer = builder.Build()
}

// // Identifiers

// UPPERCASE_ID : [A-Z]+([_][A-Z]+)* ;
// LOWERCASE_ID : [a-z]+([A-Z][a-z]*)* ;

// // Operators

// MOD : [*?+];
// PIPE : '|';

// // SYMBOLS

// // Punctuation

// COLON : ':' ;
// SEMICOLON : ';' ;

// // Brackets

// OP_PAREN : '(' ;
// CL_PAREN : ')' ;

// // Whitespace

// WS : [ \t\r\n]+ -> skip ;
