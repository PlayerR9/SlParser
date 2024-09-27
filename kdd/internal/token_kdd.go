package internal

import (
	"github.com/PlayerR9/SlParser/parser"
)

type TokenType int

const (
	EtInvalid TokenType = iota - 1
	EtEOF
	TtColon
	TtLowercaseId
	TtNewline
	TtSemicolon
	TtUppercaseId
	NtRhs
	NtRule
	NtRule1
	NtSource
	NtSource1
)

func (t TokenType) IsTerminal() bool {
	return t <= TtUppercaseId
}

var (
	Parser *parser.Parser[TokenType]
)

func init() {
	is := parser.NewItemSet[TokenType]()

	// source : source1 EOF ;
	// source1 : rule ;
	// source1 : rule NEWLINE source1 ;
	// rule : LOWERCASE_ID COLON rule1 SEMICOLON ;
	// rule1 : rhs ;
	// rule1 : rhs rule1 ;
	// rhs : UPPERCASE_ID ;
	// rhs : LOWERCASE_ID ;

	_ = is.AddRule(NtSource, NtSource1, EtEOF)
	_ = is.AddRule(NtSource1, NtRule)
	_ = is.AddRule(NtSource1, NtRule, TtNewline, NtSource1)
	_ = is.AddRule(NtRule, TtLowercaseId, TtColon, NtRule1, TtSemicolon)
	_ = is.AddRule(NtRule1, NtRhs)
	_ = is.AddRule(NtRule1, NtRhs, NtRule1)
	_ = is.AddRule(NtRhs, TtUppercaseId)
	_ = is.AddRule(NtRhs, TtLowercaseId)

	Parser = parser.Build(&is)
}
