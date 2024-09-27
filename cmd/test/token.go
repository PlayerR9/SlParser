// Code generated by SlParser. Do not edit.
package test

import (
	"github.com/PlayerR9/SlParser/parser"
)

type TokenType int

const (
	EtInvalid TokenType = iota -1
	EtEof
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
	
	_ = is.AddRule(NtSource, NtSource1, EtEof)
	_ = is.AddRule(NtSource1, NtRule)
	_ = is.AddRule(NtSource1, NtRule, TtNewline, NtSource1)
	_ = is.AddRule(NtRule, TtLowercaseId, TtColon, NtRule1, TtSemicolon)
	_ = is.AddRule(NtRule1, NtRhs)
	_ = is.AddRule(NtRule1, NtRhs, NtRule1)
	_ = is.AddRule(NtRhs, TtUppercaseId)
	_ = is.AddRule(NtRhs, TtLowercaseId)

	Parser = parser.Build(&is)
}