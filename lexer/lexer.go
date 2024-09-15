package lexer

import (
	"fmt"
	"io"

	gr "github.com/PlayerR9/SlParser/grammar"
	dba "github.com/PlayerR9/go-debug/assert"
)

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
}

// PeekRune returns the current character without consuming it.
//
// Returns:
//   - rune: the current character.
//   - error: if an error occurred.
func (l *Lexer[T]) PeekRune() (rune, error) {
	if l == nil || l.input_stream == nil {
		return 0, io.EOF
	}

	c, _, err := l.input_stream.ReadRune()
	if err != nil {
		return 0, err
	}

	err = l.input_stream.UnreadRune()
	dba.AssertErr(err, "l.input_stream.UnreadRune()")

	return c, nil
}

// NextRune consumes the current character.
//
// Returns:
//   - rune: the current character.
//   - error: if an error occurred.
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

/* // UnreadRune unreads the current character.
//
// Returns:
//   - error: if an error occurred.
func (l *Lexer) UnreadRune() error {
	if l == nil || l.input_stream == nil {
		return errors.New("no rune to unread")
	}

	err := l.input_stream.UnreadRune()
	return err
} */

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
//
// Returns:
//   - error: if an error occurred.
func (l *Lexer[T]) Lex() error {
	dba.AssertNotNil(l, "l")

	if len(l.table) == 0 {
		if l.def_fn == nil {
			c, _, err := l.input_stream.ReadRune()
			if err == io.EOF {
				return nil
			} else if err != nil {
				return err
			}

			err = l.input_stream.UnreadRune()
			dba.AssertErr(err, "l.input_stream.UnreadRune()")

			return fmt.Errorf("unexpected character %q", c)
		}

		for {
			c, _, err := l.input_stream.ReadRune()
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}

			tk, err := l.def_fn(l, c)
			if err != nil {
				_ = l.input_stream.UnreadRune()
				return err
			}

			if tk != nil {
				l.tokens = append(l.tokens, tk)
			}
		}

		return nil
	}

	if l.def_fn == nil {
		for {
			c, _, err := l.input_stream.ReadRune()
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}

			fn, ok := l.table[c]
			if !ok {
				err := l.input_stream.UnreadRune()
				dba.AssertErr(err, "l.input_stream.UnreadRune()")

				return fmt.Errorf("unexpected character %q", c)
			}

			tk, err := fn(l, c)
			if err != nil {
				_ = l.input_stream.UnreadRune()
				return err
			}

			if tk != nil {
				l.tokens = append(l.tokens, tk)
			}
		}

		return nil
	}

	for {
		c, _, err := l.input_stream.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		var tk *gr.Token[T]

		fn, ok := l.table[c]
		if ok {
			tk, err = fn(l, c)
		} else {
			tk, err = l.def_fn(l, c)
		}

		if err != nil {
			_ = l.input_stream.UnreadRune()
			return err
		}

		if tk != nil {
			l.tokens = append(l.tokens, tk)
		}
	}

	return nil
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
}
