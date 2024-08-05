// Code generated by EbnfParser.
package test

import (
	gr "github.com/PlayerR9/grammar/grammar"
	grlx "github.com/PlayerR9/grammar/lexer"
)

// Lexer is the lexer of the grammar.
type Lexer struct {
	// input_stream is the input stream of the lexer.
	input_stream []byte

	// tokens is the tokens of the lexer.
	tokens []*gr.Token[TokenType]

	// at is the position of the lexer in the input stream.
	at int
}

// NewLexer creates a new lexer.
//
// Returns:
//   - *Lexer: The new lexer. Never returns nil.
func NewLexer() *Lexer {
	return &Lexer{}
}

// SetInputStream sets the input stream of the lexer.
//
// Parameters:
//   - data: The input stream of the lexer.
func (l *Lexer) SetInputStream(data []byte) {
	l.input_stream = data
}

// Reset resets the lexer.
//
// This utility function allows to reset the information contained in the lexer
// so that it can be used multiple times.
func (l *Lexer) Reset() {
	l.tokens = l.tokens[:0]
	l.at = 0
}

// IsDone checks if the lexer is done.
//
// Returns:
//   - bool: True if the lexer is done, false otherwise.
func (l *Lexer) IsDone() bool {
	return len(l.input_stream) == 0
}

// LexOne lexes the next token of the lexer.
//
// Returns:
//   - *gr.Token[S]: The token of the lexer.
//   - error: An error if the lexer encounters an error while lexing the next token.
//
// If the returned token is nil, then it is marked as 'to skip' and, as a result,
// not added to the list of tokens.
func (l *Lexer) LexOne() (*gr.Token[TokenType], error) {
	// Lex here...
	
	panic("Implement me!")
}

// FullLex is just a wrapper around the Grammar.FullLex function.
//
// Parameters:
//   - data: The input stream of the lexer.
//
// Returns:
//   - []*Token[TokenType]: The tokens of the lexer.
//   - error: An error if the lexer encounters an error while lexing the input stream.
func FullLex(data []byte) ([]*gr.Token[TokenType], error) {
	lexer := NewLexer()

	tokens, err := grlx.FullLex(lexer, data)
	if err != nil {
		return tokens, err
	}

	return tokens, nil
}
