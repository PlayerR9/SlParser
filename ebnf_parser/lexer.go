package ebnf_parser

import (
	"errors"
	"fmt"
	"unicode"

	utch "github.com/PlayerR9/MyGoLib/Utility/runes"
	gr "github.com/PlayerR9/SLParser/grammar"
)

type Lexer struct {
	chars []rune
}

const grammar string = `
equal = "=" .
dot = "." .
pipe = "|" .
newline = [ "\r" ] "\n" .
tab = "\t" .
ws = " " . -> skip
op_paren = "(" .
cl_paren = ")" .

uppercase_id = "A".."Z" { "a".."z" } .
lowercase_id = "a".."z" { "a".."z" } .

Source = Rule { newline Rule } EOF .

Rule = uppercase_id equal RhsCls dot .

Rule
	= uppercase_id newline equal tab RhsCls { newline tab pipe RhsCls } newline dot
	.

RhsCls
	= Rhs { Rhs }
	.

Rhs
	= uppercase_id
	| lowercase_id
	| op_paren Rhs pipe Rhs { pipe Rhs } cl_paren
	.
`

func (l *Lexer) lex_one() error {
	first := l.chars[0]

	switch first {
	case '=':
		// equal = "=" .
	case '.':
		// dot = "." .
	case '|':
		// pipe = "|" .
	case '\t':
		// tab = "\t" .
	case '\n':
		// newline = [ "\r" ] "\n" .
	case ' ':
		// ws = " " . -> skip
	case '(':
		// op_paren = "(" .
	case ')':
		// cl_paren = ")" .
	default:
		if !unicode.IsLetter(first) {
			return fmt.Errorf("invalid character: %q", first)
		}

		if unicode.IsLower(first) {

		} else {

		}
	}
}

func Lex(data []byte) ([]*gr.Token, error) {
	if len(data) == 0 {
		return nil, errors.New("no data provided")
	}

	chars, idx := utch.BytesToUtf8(data)
	if idx != -1 {
		return nil, utch.NewErrInvalidUTF8Encoding(idx)
	}

	l := &Lexer{
		chars: chars,
	}

}
