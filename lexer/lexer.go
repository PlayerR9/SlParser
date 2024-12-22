package lexer

import (
	"io"
	"unicode/utf8"

	slgr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/mygo-lib/common"
	gch "github.com/PlayerR9/SlParser/mygo-lib/runes"
)

// LexOneFn is the function used to lex one token from the input data.
//
// Parameters:
//   - scanner: The input data to be lexed.
//
// Returns:
//   - *slgr.Token: The lexed token, or nil if the lexing process fails.
//   - error: An error if the lexing process fails.
type LexOneFn func(scanner io.RuneScanner) (*slgr.Token, error)

// Lexer is a lexer that can be used to lex input data into a list of tokens.
type Lexer struct {
	// chars is a list of characters that have not been lexed yet.
	chars []rune

	// tokens is a list of lexed tokens.
	tokens []*slgr.Token

	// last_read is the last rune that was read from the input data.
	last_read *rune

	// lex_one_fn is the function used to lex one token from the input data.
	lex_one_fn LexOneFn
}

// Write implements io.Writer.
func (l *Lexer) Write(data []byte) (int, error) {
	if l == nil {
		return 0, common.ErrNilReceiver
	}

	if len(data) == 0 {
		return 0, nil
	}

	chars, err := gch.BytesToUtf8(data)
	if err != nil {
		return 0, err
	}

	l.chars = append(l.chars, chars...)

	return len(data), nil
}

// ReadRune implements io.RuneScanner.
func (l *Lexer) ReadRune() (rune, int, error) {
	if l == nil {
		return 0, 0, common.ErrNilReceiver
	}

	if len(l.chars) == 0 {
		return 0, 0, io.EOF
	}

	c := l.chars[0]
	l.chars = l.chars[1:]

	l.last_read = &c

	size := utf8.RuneLen(c)

	return c, size, nil
}

// UnreadRune implements io.RuneScanner.
func (l *Lexer) UnreadRune() error {
	if l == nil {
		return common.ErrNilReceiver
	}

	if l.last_read == nil {
		return ErrCannotUnread
	}

	l.chars = append([]rune{*l.last_read}, l.chars...)
	l.last_read = nil

	return nil
}

// Reset implements common.Resetter.
func (l *Lexer) Reset() error {
	if l == nil {
		return common.ErrNilReceiver
	}

	if len(l.chars) > 0 {
		clear(l.chars)
		l.chars = nil
	}

	if len(l.tokens) > 0 {
		clear(l.tokens)
		l.tokens = nil
	}

	l.last_read = nil

	return nil
}

// GetTokens returns a copy of the tokens that have been lexed.
//
// The function returns a copy of the tokens that have been lexed so far. If
// no tokens have been lexed, the function returns nil.
//
// Returns:
//   - []*slgr.Token: A copy of the tokens that have been lexed, or nil if no
//     tokens have been lexed.
func (l Lexer) GetTokens() []*slgr.Token {
	if len(l.tokens) == 0 {
		return nil
	}

	tokens := make([]*slgr.Token, len(l.tokens))
	copy(tokens, l.tokens)

	return tokens
}

// Lex lexes input data into a list of tokens using the lexing function.
//
// The function reads runes from the input data and applies the lexing
// function to convert them into tokens. The process continues until
// the end of the input is reached or an error occurs.
//
// Returns:
//   - error: An error if the lexing process fails or if the receiver
//     is nil.
func (l *Lexer) Lex() error {
	if l == nil {
		return common.ErrNilReceiver
	}

	for {
		tk, err := l.lex_one_fn(l)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if tk != nil {
			l.tokens = append(l.tokens, tk)
		}
	}

	return nil
}

// Lex takes input data, lexes it, and returns the list of tokens.
//
// Parameters:
//   - lexer: The lexer to be used to lex the input data.
//   - data: The input data to be lexed.
//
// Returns:
//   - []*slgr.Token: The list of lexed tokens. If an error occurs while lexing, the
//     returned list may be empty.
//   - error: An error if the lexing process fails.
func Lex(lexer *Lexer, data []byte) ([]*slgr.Token, error) {
	if lexer == nil {
		return nil, common.ErrNilReceiver
	}

	defer lexer.Reset()

	if len(data) > 0 {
		n, err := lexer.Write(data)
		if err == nil && n != len(data) {
			err = io.ErrShortWrite
		}

		if err != nil {
			tokens := lexer.GetTokens()
			return tokens, err
		}
	}

	err := lexer.Lex()

	tokens := lexer.GetTokens()
	return tokens, err
}
