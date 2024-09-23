package internal

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
