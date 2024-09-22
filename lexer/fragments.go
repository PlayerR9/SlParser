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
//   - lexer: the lexer. Assumed to be non-nil.
//
// Returns:
//   - string: the fragment. Empty string if nothing was lexed.
//   - error: if an error occurred.
type LexFragment func(lexer RuneStreamer) (string, error)

// FragNewline lexes a newline.
//
// Parameters:
//   - opts: the lexer options.
//
// Returns:
//   - LexFragment: a function that lexes a newline.
func FragNewline(opts ...LexOption) LexFragment {
	fn := func(lexer RuneStreamer) (string, error) {
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
				return "", fmt.Errorf("after %q, %w",
					'\r',
					fmt.Errorf("expected %q, got nothing instead", '\n'),
				)
			} else if err != nil {
				return "", err
			}

			return "", fmt.Errorf("expected %q, got %q instead",
				'\n',
				next,
			)
		}

		err = lexer.UnreadRune()
		dba.AssertErr(err, "lexer.UnreadRune()")

		return "", NotFound
	}

	return FragWithOptions(fn, opts...)
}

// FragWs lexes a whitespace.
//
// Parameters:
//   - include_newline: if true, the lexer will include newlines in the whitespace.
//   - opts: the lexer options.
//
// Returns:
//   - LexFragment: a function that lexes a whitespace.
func FragWs(include_newline bool, opts ...LexOption) LexFragment {
	var is_fn GroupFn

	if include_newline {
		is_fn = GroupWsNl
	} else {
		is_fn = GroupWs
	}

	return FragGroup(is_fn, opts...)
}

// FragGroup lexes a group.
//
// Parameters:
//   - is_fn: the function to lex the group.
//
// Returns:
//   - LexFragment: a function that lexes the group.
//
// If 'is_fn' is nil, a function that returns an error is returned.
func FragGroup(is_fn GroupFn, opts ...LexOption) LexFragment {
	if is_fn == nil {
		return func(lexer RuneStreamer) (string, error) {
			return "", errors.New("no group function provided")
		}
	}

	fn := func(lexer RuneStreamer) (string, error) {
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

	return FragWithOptions(fn, opts...)
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
func FragWord(word string, opts ...LexOption) LexFragment {
	chars, err := gcch.StringToUtf8(word)
	if err != nil {
		return func(lexer RuneStreamer) (string, error) {
			return "", fmt.Errorf("invalid word: %w", err)
		}
	}

	fn := func(lexer RuneStreamer) (string, error) {
		prev := chars[0]

		for _, char := range chars[1:] {
			c, err := lexer.NextRune()
			if err == io.EOF {
				return "", fmt.Errorf("after %q, %w", prev,
					fmt.Errorf("expected %q, got nothing instead", char),
				)
			} else if err != nil {
				return "", err
			}

			if c != char {
				return "", fmt.Errorf("after %q, %w", prev,
					fmt.Errorf("expected %q, got %q instead", char, c),
				)
			}

			prev = char
		}

		return word, nil
	}

	return FragWithOptions(fn, opts...)
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
					return builder.String(), fmt.Errorf("after %q, %w", prev,
						fmt.Errorf("expected %q, got nothing instead", until),
					)
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
