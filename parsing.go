package SlParser

import (
	gr "github.com/PlayerR9/SlParser/grammar"
	lxr "github.com/PlayerR9/SlParser/lexer"
	prx "github.com/PlayerR9/SlParser/parser"
	gcers "github.com/PlayerR9/go-commons/errors"
)

// Lex is a function that lexes the given data.
//
// The last token is always EOF, even if nothing was lexed; thus,
// length of the returned list is always >= 1.
//
// Parameters:
//   - lexer: The lexer.
//   - data: The data.
//
// Returns:
//   - []*gr.Token[T]: The list of tokens.
//   - error: if an error occurred.
func Lex[T gr.TokenTyper](lexer *lxr.Lexer[T], data []byte) ([]*gr.Token[T], error) {
	defer lexer.Reset()

	var err error

	if lexer != nil {
		input_stream := lxr.NewStream().FromBytes(data)
		lexer.SetInputStream(input_stream)
		err = lexer.Lex()
	} else {
		err = gcers.NewErrNilParameter("lexer")
	}

	tokens := lexer.Tokens()
	return tokens, err
}

// LexString is a function that lexes the given string.
//
// The last token is always EOF, even if nothing was lexed; thus,
// length of the returned list is always >= 1.
//
// Parameters:
//   - lexer: The lexer.
//   - str: The string.
//
// Returns:
//   - []*gr.Token[T]: The list of tokens.
//   - error: if an error occurred.
func LexString[T gr.TokenTyper](lexer *lxr.Lexer[T], str string) ([]*gr.Token[T], error) {
	defer lexer.Reset()

	var err error

	if lexer != nil {
		input_stream := lxr.NewStream().FromString(str)
		lexer.SetInputStream(input_stream)
		err = lexer.Lex()
	} else {
		err = gcers.NewErrNilParameter("lexer")
	}

	tokens := lexer.Tokens()
	return tokens, err
}

// Parse is a function that parses the given tokens.
//
// Parameters:
//   - parser: The parser.
//   - tokens: The tokens.
//
// Returns:
//   - []*gr.Token[T]: The list of tokens.
//   - error: if an error occurred.
func Parse[T gr.TokenTyper](parser *prx.Parser[T], tokens []*gr.Token[T]) ([]*gr.Token[T], error) {
	if parser == nil {
		return nil, gcers.NewErrNilParameter("parser")
	}

	defer parser.Reset()

	parser.SetTokens(tokens)
	err := parser.Parse()
	forest := parser.Forest()

	return forest, err
}
