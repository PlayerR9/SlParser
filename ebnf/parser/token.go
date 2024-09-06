package parser

//go:generate stringer -type=TokenType

type TokenType int

const (
	EttEOF TokenType = iota

	TttMod
)

func (t TokenType) String() string {
	return [...]string{
		"EOF",
	}[t]
}

func (t TokenType) IsTerminal() bool {
	return t <= EttEOF
}
