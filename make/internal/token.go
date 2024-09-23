package internal

import (
	"unicode"
	"unicode/utf8"
)

//go:generate stringer -type=TokenType -linecomment

type TokenType int

const (
	ExtraTk       TokenType = iota // Et
	TerminalTk                     // Tt
	NonterminalTk                  // Nt
)

type Token struct {
	Type TokenType
	Data string
}

func (t Token) String() string {
	return t.Type.String() + t.Data
}

func NewToken(type_ TokenType, data string) *Token {
	return &Token{
		Type: type_,
		Data: data,
	}
}

func (t Token) IsCandidateForAst() bool {
	if t.Type != NonterminalTk || t.Data == "" {
		return false
	}

	r, _ := utf8.DecodeLastRuneInString(t.Data)
	return !unicode.IsDigit(r) && unicode.IsLetter(r)
}
