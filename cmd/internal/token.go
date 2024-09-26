package internal

import (
	"unicode"
	"unicode/utf8"
)

// Token is a token.
type Token struct {
	// Type is the type of the token.
	Type TokenType

	// Data is the data of the token.
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
