package lexer

import (
	"errors"
	"io"

	slgr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/mygo-lib/common"
)

// Builder is a builder for lexers.
type Builder struct {
	// lex_one_fn is the function used to lex one token from the input data.
	lex_one_fn LexOneFn
}

// Reset implements common.Resetter.
func (b *Builder) Reset() error {
	if b == nil {
		return common.ErrNilReceiver
	}

	b.lex_one_fn = nil

	return nil
}

// SetLexOneFn sets the lexing function used by the lexer.
//
// Parameters:
//   - fn: The new lexing function. Must not be nil.
//
// Returns:
//   - error: An error if the receiver is nil or if the parameter is nil.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
//   - common.ErrBadParam: If the parameter is nil.
func (b *Builder) SetLexOneFn(fn LexOneFn) error {
	if b == nil {
		return common.ErrNilReceiver
	}

	if fn == nil {
		err := common.NewErrNilParam("fn")
		return err
	}

	b.lex_one_fn = fn

	return nil
}

// Build creates a new lexer using the values set on the builder.
//
// Returns:
//   - *Lexer: The newly created lexer. Never returns nil.
func (b Builder) Build() *Lexer {
	var fn LexOneFn

	if b.lex_one_fn == nil {
		fn = func(_ io.RuneScanner) (*slgr.Token, error) {
			err := errors.New("no lexing function provided")
			return nil, err
		}
	} else {
		fn = b.lex_one_fn
	}

	lexer := &Lexer{
		lex_one_fn: fn,
	}

	return lexer
}
