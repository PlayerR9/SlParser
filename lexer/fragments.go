package lexer

import (
	"errors"
	"fmt"
	"io"
	"strings"

	gcch "github.com/PlayerR9/go-commons/runes"
	dba "github.com/PlayerR9/go-debug/assert"
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
//   - string: the fragment. Empty string if nothing was lexed.
//   - error: if an error occurred.
//
// It can only either return NotFound or an error.
type LexFragment func(stream RuneStreamer) (string, error)

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
	FragNewline = func(lexer RuneStreamer) (string, error) {
		char, err := lexer.NextRune()
		if err == io.EOF {
			return "", NotFound
		} else if err != nil {
			return "", err
		}

		if char == '\n' {
			return "\n", nil
		}

		if char == '\r' {
			next, err := lexer.NextRune()
			if err == nil && next == '\n' {
				return "\r\n", nil
			}

			if err == io.EOF {
				return "", NewErrGotNothing('\r', '\n')
			} else if err != nil {
				return "", err
			}

			return "", NewErrGotUnexpected('\r', '\n', next)
		}

		err = lexer.UnreadRune()
		dba.AssertErr(err, "lexer.UnreadRune()")

		return "", NotFound
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
		return func(lexer RuneStreamer) (string, error) {
			return "", NewErrNoGroupSpecified()
		}
	}

	return func(lexer RuneStreamer) (string, error) {
		char, err := lexer.NextRune()
		if err == io.EOF {
			return "", NotFound
		} else if err != nil {
			return "", err
		}

		if is_fn(char) {
			return string(char), nil
		}

		err = lexer.UnreadRune()
		dba.AssertErr(err, "lexer.UnreadRune()")

		return "", NotFound
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
		return func(lexer RuneStreamer) (string, error) {
			return "", fmt.Errorf("invalid word: %w", err)
		}
	}

	return func(lexer RuneStreamer) (string, error) {
		prev := chars[0]

		for _, char := range chars[1:] {
			c, err := lexer.NextRune()
			if err == io.EOF {
				return "", NewErrGotNothing(prev, char)
			} else if err != nil {
				return "", err
			}

			if c != char {
				return "", NewErrGotUnexpected(prev, char, c)
			}

			prev = char
		}

		return word, nil
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
		fn = func(lexer RuneStreamer) (string, error) {
			var builder strings.Builder

			builder.WriteRune(prev)

			for {
				char, err := lexer.NextRune()
				if err == io.EOF {
					break
				} else if err != nil {
					return builder.String(), err
				}

				builder.WriteRune(char)

				if char == until {
					break
				}
			}

			return builder.String(), nil
		}
	} else {
		fn = func(lexer RuneStreamer) (string, error) {
			var builder strings.Builder

			builder.WriteRune(prev)

			for {
				char, err := lexer.NextRune()
				if err == io.EOF {
					return builder.String(), NewErrGotNothing(prev, until)
				} else if err != nil {
					return builder.String(), err
				}

				builder.WriteRune(char)

				if char == until {
					break
				}
			}

			return builder.String(), nil
		}
	}

	return fn
}
