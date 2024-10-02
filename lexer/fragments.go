package lexer

import (
	"errors"
	"fmt"
	"io"

	gcch "github.com/PlayerR9/go-commons/runes"
	"github.com/PlayerR9/go-errors/assert"
)

var (
	// NotFound is an error that is returned when a fragment is not found. Unlike with other errors,
	// this one does not indicate an invalid lexing; it simply means that the fragment was not found.
	//
	// Readers must return this error as is and not wrap it as callers check for this error using
	// ==.
	NotFound error
)

func init() {
	NotFound = errors.New("not found")
}

// LexFragment is a function that lexes a fragment.
//
// Parameters:
//   - stream: The rune stream. Assumed to be non-nil.
//
// Returns:
//   - error: An error if the fragment was not found.
//
// It can only either return NotFound or an error.
type LexFragment func(stream RuneStreamer) error

var (
	// FragNewline lexes a newline.
	//
	// Parameters:
	//   - opts: the lexer options.
	//
	// Returns:
	//   - LexFragment: a function that lexes a newline.
	//
	// By default, the lexer does allow optional fragments and only lexes once.
	//   - Use WithAllowOptional(false) to disable optional fragments.
	//   - Use WithLexMany(true) to enable one or more fragments.
	FragNewline LexFragment
)

func init() {
	FragNewline = func(stream RuneStreamer) error {
		assert.NotNil(stream, "stream")

		char, err := stream.NextRune()
		if err == io.EOF {
			return NotFound
		} else if err != nil {
			return err
		}

		if char == '\n' {
			return nil
		} else if char != '\r' {
			err = stream.UnreadRune()
			assert.Err(err, "lexer.UnreadRune()")

			return NotFound
		}

		next, err := stream.NextRune()
		if err == nil && next == '\n' {
			return nil
		}

		if err == io.EOF {
			return NewErrGotNothing('\r', '\n')
		} else if err != nil {
			return err
		}

		return NewErrGotUnexpected('\r', '\n', next)
	}
}

// FragWs lexes a whitespace.
//
// Parameters:
//   - include_newline: if true, the lexer will include newlines in the whitespace.
//   - opts: the lexer options.
//
// Returns:
//   - LexFragment: a function that lexes a whitespace.
//
// By default, the lexer does allow optional fragments and only lexes once.
//   - Use WithAllowOptional(false) to disable optional fragments.
//   - Use WithLexMany(true) to enable one or more fragments.
func FragWs(include_newline bool) LexFragment {
	var is_fn GroupFn

	if include_newline {
		is_fn = GroupWsNl
	} else {
		is_fn = GroupWs
	}

	return FragGroup(is_fn)
}

// FragGroup lexes a group.
//
// Parameters:
//   - is_fn: the function to lex the group.
//   - opts: the lexer options.
//
// Returns:
//   - LexFragment: a function that lexes the group.
//
// If 'is_fn' is nil, a function that returns an error is returned.
//
// By default, the lexer does allow optional fragments and only lexes once.
//   - Use WithAllowOptional(false) to disable optional fragments.
//   - Use WithLexMany(true) to enable one or more fragments.
func FragGroup(is_fn GroupFn) LexFragment {
	if is_fn == nil {
		return func(lexer RuneStreamer) error {
			return NewErrNoGroupSpecified()
		}
	}

	return func(lexer RuneStreamer) error {
		assert.NotNil(lexer, "lexer")

		char, err := lexer.NextRune()
		if err == io.EOF {
			return NotFound
		} else if err != nil {
			return err
		}

		if is_fn(char) {
			return nil
		}

		err = lexer.UnreadRune()
		assert.Err(err, "lexer.UnreadRune()")

		return NotFound
	}
}

// FragWord lexes a word. However, the first character of the word is ignored
// as it assumes that said character is already lexed.
//
// Parameters:
//   - word: the word to lex.
//   - opts: the lexer options.
//
// Returns:
//   - LexFragment: a function that lexes the word.
//
// If 'word' is an invalid UTF-8 string, a function that returns an error is returned.
//
// If the word is not found in the lexer's input, a ErrUnexpectedChar error is returned.
//
// By default, the lexer does allow optional fragments and only lexes once.
//   - Use WithAllowOptional(false) to disable optional fragments.
//   - Use WithLexMany(true) to enable one or more fragments.
func FragWord(word string) LexFragment {
	chars, err := gcch.StringToUtf8(word)
	if err != nil {
		return func(lexer RuneStreamer) error {
			return fmt.Errorf("invalid word: %w", err)
		}
	}

	return func(lexer RuneStreamer) error {
		assert.NotNil(lexer, "lexer")

		prev := chars[0]

		for _, char := range chars[1:] {
			c, err := lexer.NextRune()
			if err == io.EOF {
				return NewErrGotNothing(prev, char)
			} else if err != nil {
				return err
			}

			if c != char {
				return NewErrGotUnexpected(prev, char, c)
			}

			prev = char
		}

		return nil
	}
}

// FragUntil lexes until a character is found.
//
// Parameters:
//   - prev: the previous character.
//   - until: the character to find.
//   - allow_eof: if true, the lexer will allow the end of the input.
//
// Returns:
//   - LexFragment: a function that lexes until the character is found.
func FragUntil(prev, until rune, allow_eof bool) LexFragment {
	var fn LexFragment

	if allow_eof {
		fn = func(lexer RuneStreamer) error {
			for {
				char, err := lexer.NextRune()
				if err == io.EOF {
					break
				} else if err != nil {
					return err
				}

				if char == until {
					break
				}
			}

			return nil
		}
	} else {
		fn = func(lexer RuneStreamer) error {
			for {
				char, err := lexer.NextRune()
				if err == io.EOF {
					return NewErrGotNothing(prev, until)
				} else if err != nil {
					return err
				}

				if char == until {
					break
				}
			}

			return nil
		}
	}

	return fn
}
