package lexer

//go:generate stringer -type=ErrorCode

type ErrorCode int

const (
	// UnrecognizedChar occurs when an unrecognized character is encountered.
	//
	// Example:
	// 	let L = lexer with integers as its lexing table.
	// 	Lex(L, "a")
	UnrecognizedChar ErrorCode = iota

	// InvalidInputStream occurs when the input stream is invalid.
	//
	// Example:
	// 	let is = input stream of non-utf8 characters.
	InvalidInputStream

	// BadWord occurs when a word is invalid.
	//
	// Example:
	// 	let L = lexer with integers as its lexing table.
	// 	Lex(L, "01")
	BadWord
)

func (e ErrorCode) Int() int {
	return int(e)
}
