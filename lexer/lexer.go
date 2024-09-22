package lexer

import (
	"errors"
	"io"

	gr "github.com/PlayerR9/SlParser/grammar"
	gcers "github.com/PlayerR9/errors/error"
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

	// pos is the current position in the input stream.
	// This is in bytes.
	pos int

	// next_pos is the next position in the input stream.
	// This is in bytes.
	next_pos int

	// last_read_size is the size of the last read rune.
	last_read_size int

	// state is the lexer state.
	state LexerState
}

// NextRune implements RuneStreamer interface.
func (l *Lexer[T]) NextRune() (rune, error) {
	if l == nil || l.input_stream == nil {
		return 0, io.EOF
	}

	c, size, err := l.input_stream.ReadRune()
	if err != nil {
		l.state.UpdateLastErr(err)

		return 0, err
	}

	l.next_pos += size
	l.last_read_size = size

	l.state.UpdateLastErr(nil)

	return c, nil
}

// UnreadRune implements RuneStreamer interface.
func (l *Lexer[T]) UnreadRune() error {
	if l == nil || l.input_stream == nil {
		return errors.New("no rune to unread")
	}

	err := l.input_stream.UnreadRune()
	if err != nil {
		return err
	}

	l.next_pos -= l.last_read_size
	l.last_read_size = 0

	return nil
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

// add_token is a helper function that adds a token to the list of tokens
// iff the token is not nil; updating the lexer state.
//
// Parameters:
//   - tk: the token to add.
func (l *Lexer[T]) add_token(tk *gr.Token[T]) {
	if tk != nil {
		l.tokens = append(l.tokens, tk)
	}

	l.last_read_size = 0
	l.pos = l.next_pos
}

// Lex lexes the input stream.
//
// Returns:
//   - error: an error of type *Err if the input stream could not be lexed.
func (l *Lexer[T]) Lex() error {
	if l == nil {
		return nil
	}

	if len(l.table) == 0 {
		if l.def_fn == nil {
			char, err := l.NextRune()
			if err == io.EOF {
				return nil
			} else if err != nil {
				l.state.UpdateLastCharRead(nil)

				return l.make_error()
			}

			l.state.UpdateLastCharRead(&char)

			return l.make_error()
		}

		for {
			char, err := l.NextRune()
			if err == io.EOF {
				break
			} else if err != nil {
				l.state.UpdateLastCharRead(nil)

				return l.make_error()
			}

			l.state.UpdateLastCharRead(&char)

			type_, data, err := l.def_fn(l, char)
			if err == nil {
				tk := gr.NewToken(type_, data, l.pos)

				l.add_token(tk)
			} else if err == SkipToken {
				l.add_token(nil)
			} else {
				_ = l.input_stream.UnreadRune()

				l.state.UpdateLastErr(err)

				return l.make_error()
			}
		}

		return nil
	}

	if l.def_fn == nil {
		for {
			char, err := l.NextRune()
			if err == io.EOF {
				break
			} else if err != nil {
				l.state.UpdateLastCharRead(nil)

				return l.make_error()
			}

			l.state.UpdateLastCharRead(&char)

			var type_ T
			var data string

			var tk *gr.Token[T]

			fn, ok := l.table[char]
			if !ok {
				return l.make_error()
			}

			type_, data, err = fn(l, char)

			if err == nil {
				tk = gr.NewToken(type_, data, l.pos)

				l.add_token(tk)
			} else if err != SkipToken {
				_ = l.input_stream.UnreadRune()

				l.state.UpdateLastErr(err)

				return l.make_error()
			} else {
				l.add_token(nil)
			}
		}
	}

	for {
		char, err := l.NextRune()
		if err == io.EOF {
			break
		} else if err != nil {
			l.state.UpdateLastCharRead(nil)

			return l.make_error()
		}

		l.state.UpdateLastCharRead(&char)

		var type_ T
		var data string

		var tk *gr.Token[T]

		fn, ok := l.table[char]
		if ok {
			type_, data, err = fn(l, char)
		} else {
			type_, data, err = l.def_fn(l, char)
		}

		if err == nil {
			tk = gr.NewToken(type_, data, l.pos)

			l.add_token(tk)
		} else if err != SkipToken {
			_ = l.input_stream.UnreadRune()

			l.state.UpdateLastErr(err)

			return l.make_error()
		} else {
			l.add_token(nil)
		}
	}

	return nil
}

// Tokens returns the list of tokens. The last token is always EOF.
//
// Returns:
//   - []*gr.Token[T]: the list of tokens.
func (l *Lexer[T]) Tokens() []*gr.Token[T] {
	eof := gr.NewToken(T(0), "", -1)

	if l == nil {
		return []*gr.Token[T]{eof}
	}

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

	l.state.Reset()

	l.pos = 0
	l.next_pos = 0
	l.last_read_size = 0
}

// make_error creates a new error.
//
// Returns:
//   - *Err: the last error. Never returns nil.
func (l Lexer[T]) make_error() *gcers.Err[ErrorCode] {
	pos := l.next_pos

	if l.state.last_char_read == nil {
		err := gcers.NewErr(gcers.FATAL, InvalidInputStream, l.state.last_err.Error())
		err.AddSuggestion("Input is most likely not a valid input for the current lexer.")
		err.AddFrame("lexer", "Lexer[T]")

		AddContext(err, "pos", pos)

		return err
	}

	last_read := *l.state.last_char_read

	_, ok := l.table[last_read]
	if !ok && l.def_fn == nil {
		err := gcers.NewErr(gcers.FATAL, UnrecognizedChar, l.state.last_err.Error())

		err.AddSuggestion(
			"Input provided cannot be lexed by the current lexer. You may want to check for typos in the input.",
		)
		err.AddSuggestion(
			"(Less likely) The lexer table is not configured correctly. Contact the developer and provide this error code.",
		)

		AddContext(err, "pos", pos)

		return err
	}

	err := gcers.NewErr(gcers.FATAL, BadWord, l.state.last_err.Error())
	err.AddSuggestion("You may want to check for typos in the input.")

	AddContext(err, "pos", pos)

	return err
}
