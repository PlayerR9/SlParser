package lexer

import (
	"errors"
	"io"

	gr "github.com/PlayerR9/SlParser/grammar"
	util "github.com/PlayerR9/SlParser/util"
	dba "github.com/PlayerR9/go-debug/assert"
)

// RuneStreamer is a rune streamer.
type RuneStreamer interface {
	// NextRune consumes the current character.
	//
	// Returns:
	//   - rune: the current character.
	//   - error: if an error occurred.
	//
	// Errors:
	//   - io.EOF: if the end of the stream is reached.
	//   - any other error if the stream could not be read.
	NextRune() (rune, error)

	// UnreadRune unreads the current character.
	//
	// Returns:
	//   - error: if an error occurred.
	UnreadRune() error
}

var (
	// SkipToken is an error that is returned when a token is skipped.
	//
	// Readers must return this error as is and not wrap it as callers check for this error using
	// ==.
	SkipToken error
)

func init() {
	SkipToken = errors.New("skip token")
}

// Lexer is a lexer.
type Lexer[T gr.TokenTyper] struct {
	// input_stream is the input stream.
	input_stream io.RuneScanner

	// tokens is the list of tokens.
	tokens []*gr.Token[T]

	// table is the lexer table.
	table map[rune]LexFunc[T]

	// def_fn is the default lexer function.
	def_fn LexFunc[T]

	// err is the last error.
	err *util.Err[ErrorCode]
}

// NextRune implements RuneStreamer interface.
func (l *Lexer[T]) NextRune() (rune, error) {
	if l == nil || l.input_stream == nil {
		return 0, io.EOF
	}

	c, _, err := l.input_stream.ReadRune()
	if err != nil {
		return 0, err
	}

	return c, nil
}

// UnreadRune implements RuneStreamer interface.
func (l *Lexer[T]) UnreadRune() error {
	if l == nil || l.input_stream == nil {
		return errors.New("no rune to unread")
	}

	err := l.input_stream.UnreadRune()
	return err
}

// SetInputStream sets the input stream.
//
// Parameters:
//   - input_stream: the input stream.
//
// Does nothing if the receiver is nil.
func (l *Lexer[T]) SetInputStream(input_stream io.RuneScanner) {
	if l == nil {
		return
	}

	l.input_stream = input_stream
}

// Lex lexes the input stream.
func (l *Lexer[T]) Lex() {
	dba.AssertNotNil(l, "l")

	if len(l.table) == 0 {
		if l.def_fn == nil {
			c, _, err := l.input_stream.ReadRune()
			if err == nil {
				err = l.input_stream.UnreadRune()
				dba.AssertErr(err, "l.input_stream.UnreadRune()")

				l.err = NewErrUnrecognizedChar(c)
			} else if err != io.EOF {
				l.err = NewErrInvalidInputStream(err)
			}

			return
		}

		for l.err == nil {
			c, _, err := l.input_stream.ReadRune()
			if err == io.EOF {
				break
			} else if err != nil {
				l.err = NewErrInvalidInputStream(err)
				break
			}

			type_, data, err := l.def_fn(l, c)
			if err == nil {
				tk := gr.NewTerminalToken(type_, data)
				l.tokens = append(l.tokens, tk)
			} else if err != SkipToken {
				_ = l.input_stream.UnreadRune()

				l.err = NewErrInvalidInputStream(err)
			}
		}

		return
	}

	if l.def_fn == nil {
		for l.err == nil {
			c, _, err := l.input_stream.ReadRune()
			if err == io.EOF {
				break
			} else if err != nil {
				l.err = NewErrInvalidInputStream(err)
				break
			}

			fn, ok := l.table[c]
			if !ok {
				err := l.input_stream.UnreadRune()
				dba.AssertErr(err, "l.input_stream.UnreadRune()")

				l.err = NewErrUnrecognizedChar(c)
				break
			}

			type_, data, err := fn(l, c)
			if err == nil {
				tk := gr.NewTerminalToken(type_, data)
				l.tokens = append(l.tokens, tk)
			} else if err != SkipToken {
				_ = l.input_stream.UnreadRune()

				l.err = NewErrInvalidInputStream(err)
			}
		}

		return
	}

	for l.err == nil {
		c, _, err := l.input_stream.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			l.err = NewErrInvalidInputStream(err)

			break
		}

		var type_ T
		var data string

		var tk *gr.Token[T]

		fn, ok := l.table[c]
		if ok {
			type_, data, err = fn(l, c)
		} else {
			type_, data, err = l.def_fn(l, c)
		}

		if err == nil {
			tk = gr.NewTerminalToken(type_, data)
			l.tokens = append(l.tokens, tk)
		} else if err != SkipToken {
			_ = l.input_stream.UnreadRune()

			l.err = NewErrInvalidInputStream(err)
		}
	}
}

// Tokens returns the list of tokens.
//
// Returns:
//   - []*gr.Token[T]: the list of tokens.
func (l Lexer[T]) Tokens() []*gr.Token[T] {
	eof := gr.NewTerminalToken(T(0), "")

	tokens := make([]*gr.Token[T], len(l.tokens), len(l.tokens)+1)
	copy(tokens, l.tokens)

	tokens = append(tokens, eof)

	for i := 0; i < len(tokens)-1; i++ {
		tokens[i].Lookahead = tokens[i+1]
	}

	return tokens
}

// Reset resets the lexer and makes it reusable.
func (l *Lexer[T]) Reset() {
	if l == nil {
		return
	}

	l.input_stream = nil

	if len(l.tokens) > 0 {
		for i := 0; i < len(l.tokens); i++ {
			l.tokens[i] = nil
		}

		l.tokens = l.tokens[:0]
	}

	l.err = nil
}

// Error returns the last error.
//
// Returns:
//   - *util.Err[ErrorCode]: the last error. Nil if no error occurred.
func (l Lexer[T]) Error() *util.Err[ErrorCode] {
	if l.err == nil {
		return nil
	}

	return l.err
}
