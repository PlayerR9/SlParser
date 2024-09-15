package lexer

import (
	"io"
	"strings"

	gr "github.com/PlayerR9/SlParser/grammar"
	dba "github.com/PlayerR9/go-debug/assert"
)

// LexFragment is a function that lexes a fragment.
//
// Parameters:
//   - lexer: the lexer. Assumed to be non-nil.
//
// Returns:
//   - string: the fragment. Empty string if nothing was lexed.
//   - error: if an error occurred.
type LexFragment[T gr.TokenTyper] func(lexer *Lexer[T]) (string, error)

// FragNewline lexes a newline.
//
// Parameters:
//   - lexer: the lexer. Assumed to be non-nil.
//
// Returns:
//   - string: the newline. Empty string if no newline is found.
//   - error: if an error occurred.
func FragNewline[T gr.TokenTyper](lexer *Lexer[T]) (string, error) {
	if lexer == nil {
		return "", nil
	}

	char, err := lexer.PeekRune()
	if err == io.EOF {
		return "", nil
	} else if err != nil {
		return "", err
	}

	switch char {
	case '\n':
		_, err = lexer.NextRune()
		dba.AssertErr(err, "lexer.NextRune()")

		fallthrough
	case '\r':
		next, err := lexer.NextRune()
		if err == io.EOF {
			return "", NewErrUnexpectedChar('\r', []rune{'\n'}, nil)
		} else if err != nil {
			return "", err
		}

		if next != '\n' {
			return "", NewErrUnexpectedChar('\r', []rune{'\n'}, &next)
		}
	default:
		return "", nil
	}

	return "\n", nil
}

// FragWs lexes a whitespace.
//
// Parameters:
//   - lexer: the lexer. Assumed to be non-nil.
//
// Returns:
//   - string: the whitespace. Empty string if nothing was lexed.
//   - error: if an error occurred.
func FragWs[T gr.TokenTyper](lexer *Lexer[T]) (string, error) {
	if lexer == nil {
		return "", nil
	}

	char, err := lexer.PeekRune()
	if err == io.EOF {
		return "", nil
	} else if err != nil {
		return "", err
	}

	if char != ' ' && char != '\t' {
		return "", nil
	}

	_, err = lexer.NextRune()
	dba.AssertErr(err, "lexer.NextRune()")

	return " ", nil
}

// FragGroup lexes a group.
//
// Parameters:
//   - lexer: the lexer. Assumed to be non-nil.
//   - fn: the function to lex the group. Assumed to be non-nil.
//
// Returns:
//   - string: the group. Empty string if nothing was lexed.
//   - error: if an error occurred.
func FragGroup[T gr.TokenTyper](lexer *Lexer[T], fn func(char rune) bool) (string, error) {
	if lexer == nil || fn == nil {
		return "", nil
	}

	char, err := lexer.PeekRune()
	if err == io.EOF {
		return "", nil
	} else if err != nil {
		return "", err
	}

	if !fn(char) {
		return "", nil
	}

	_, err = lexer.NextRune()
	dba.AssertErr(err, "lexer.NextRune()")

	return string(char), nil
}

// LexMany lexes the fragment one or more times.
//
// Parameters:
//   - lexer: the lexer. Assumed to be non-nil.
//   - fn: the function to lex the fragment. Assumed to be non-nil.
//
// Returns:
//   - string: the fragment. Empty string if nothing was lexed.
//   - error: if an error occurred.
func LexMany[T gr.TokenTyper](lexer *Lexer[T], fn LexFragment[T]) (string, error) {
	if lexer == nil || fn == nil {
		return "", nil
	}

	var builder strings.Builder

	for {
		res, err := fn(lexer)
		if err != nil {
			return builder.String(), err
		} else if res == "" {
			break
		}

		builder.WriteString(res)
	}

	return builder.String(), nil
}
