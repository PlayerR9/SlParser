
package internal

import (
	"github.com/PlayerR9/SlParser/lexer"
	"github.com/PlayerR9/SlParser/parser"
)

//go:generate stringer -type=TokenType

type TokenType int

const (
	EtInvalid TokenType = iota -1
	EtEOF
	TtListComprehension
	TtNewline
	TtPrintStmt
	NtSource
	NtSource1
	NtStatement
)

func (t TokenType) IsTerminal() bool {
	return t <= TtPrintStmt
}

var (
	Lexer *lexer.Lexer[TokenType]
	Parser *parser.Parser[TokenType]
)

func init() {
	is := parser.NewItemSet[TokenType]()
	
	_ = is.AddRule(NtSource, TtNewline, NtSource1, EtEOF)
	_ = is.AddRule(NtSource1, NtStatement)
	_ = is.AddRule(NtSource1, NtStatement, TtNewline, NtSource1)
	_ = is.AddRule(NtStatement, TtListComprehension)
	_ = is.AddRule(NtStatement, TtPrintStmt)

	Parser = parser.Build(&is)

	builder := lexer.NewBuilder[TokenType]()

	// TODO: Add here your own custom rules...
	
	Lexer = builder.Build()
}