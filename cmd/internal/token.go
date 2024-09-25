package internal

import (
	"unicode"
	"unicode/utf8"
)

//go:generate stringer -type=TokenType -linecomment

type TokenType int

const (
	InvalidTk     TokenType = iota - 1 // Invalid
	ExtraTk                            // Et
	TerminalTk                         // Tt
	NonterminalTk                      // Nt
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

func IsCandidateForAst(type_ TokenType, data string) bool {
	if type_ != NonterminalTk || data == "" {
		return false
	}

	r, _ := utf8.DecodeLastRuneInString(data)
	return !unicode.IsDigit(r) && unicode.IsLetter(r)
}
