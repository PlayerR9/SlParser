package pkg

import (
	"fmt"
	"regexp"
	"strconv"
	"unicode"
	"unicode/utf8"

	gr "github.com/PlayerR9/grammar/grammar"
	ulpx "github.com/PlayerR9/grammar/lexer"
	utch "github.com/PlayerR9/lib_units/runes"
)

var (
	parse_lowercase_id *regexp.Regexp
	parse_uppercase_id *regexp.Regexp
	parse_newlines     *regexp.Regexp

	single_token_map map[rune]TokenType
)

func init() {
	parse_lowercase_id = regexp.MustCompile(`^[a-z]+([_][a-z]+)*[0-9]*`)
	parse_uppercase_id = regexp.MustCompile(`^([A-Z][a-z]*)+[0-9]*`)
	parse_newlines = regexp.MustCompile(`((\r)?\n)+`)

	single_token_map = map[rune]TokenType{
		'.': TtkDot,
		'(': TtkOpParen,
		')': TtkClParen,
		'|': TtkPipe,
		'=': TtkEqualSign,
	}
}

// Lexer is a struct that represents a lexer.
type Lexer struct {
	// input_stream is the input stream of the lexer.
	input_stream []byte

	// tokens is the tokens of the lexer.
	tokens []*gr.Token[TokenType]

	// at is the position of the lexer in the input stream.
	at int
}

// SetInputStream implements the Grammar.Lexer interface.
func (l *Lexer) SetInputStream(data []byte) {
	l.input_stream = data
}

// Reset implements the Grammar.Lexer interface.
func (l *Lexer) Reset() {
	l.tokens = l.tokens[:0]
	l.at = 0
}

// IsDone implements the Grammar.Lexer interface.
func (l *Lexer) IsDone() bool {
	return len(l.input_stream) == 0
}

// LexOne implements the Grammar.Lexer interface.
func (l *Lexer) LexOne() (*gr.Token[TokenType], error) {
	// luc.Assert(len(l.input_stream) > 0, "l.input_stream is empty")

	c, size := utf8.DecodeRune(l.input_stream)
	if c == utf8.RuneError {
		return nil, utch.NewErrInvalidUTF8Encoding(l.at)
	}

	token_type, ok := single_token_map[c]
	if ok {
		tk := gr.NewToken(token_type, string(c), l.at, nil)

		l.input_stream = l.input_stream[size:]
		l.at += size

		return tk, nil
	}

	var tk *gr.Token[TokenType]

	switch c {
	case ' ', '\t':
		// Do nothing
	case '\r':
		l.input_stream = l.input_stream[size:]
		l.at += size

		if len(l.input_stream) == 0 {
			return nil, fmt.Errorf("expected '\\n' after '\\r', got nothing instead")
		}

		c, size = utf8.DecodeRune(l.input_stream)
		if c == utf8.RuneError {
			return nil, utch.NewErrInvalidUTF8Encoding(l.at)
		}

		if c != '\n' {
			return nil, fmt.Errorf("expected '\\n' after '\\r', got %s instead", strconv.QuoteRune(c))
		}

		tmp := gr.NewToken(TtkNewline, "\n", l.at, nil)

		tk = tmp
	default:
		var match []byte

		if !unicode.IsLetter(c) {
			match = parse_newlines.Find(l.input_stream)
			if len(match) == 0 {
				return nil, fmt.Errorf("invalid character: %s", strconv.QuoteRune(c))
			}

			tmp := gr.NewToken(TtkNewline, string(match), l.at, nil)

			tk = tmp
		} else {
			if unicode.IsUpper(c) {
				match = parse_uppercase_id.Find(l.input_stream)
				if len(match) == 0 {
					return nil, fmt.Errorf("invalid uppercase identifier: %s", strconv.QuoteRune(c))
				}

				tmp := gr.NewToken(TtkUppercaseID, string(match), l.at, nil)

				tk = tmp
			} else {
				match = parse_lowercase_id.Find(l.input_stream)
				if len(match) == 0 {
					return nil, fmt.Errorf("invalid lowercase identifier: %s", strconv.QuoteRune(c))
				}

				tmp := gr.NewToken(TtkLowercaseID, string(match), l.at, nil)

				tk = tmp
			}
		}

		size = len(match)
	}

	l.input_stream = l.input_stream[size:]
	l.at += size

	return tk, nil
}

// NewLexer creates a new lexer.
//
// Returns:
//   - *Lexer: The new lexer. Never returns nil.
func NewLexer() *Lexer {
	return &Lexer{}
}

// FullLex is just a wrapper around the Grammar.FullLex function.
//
// Parameters:
//   - data: The input stream of the lexer.
//
// Returns:
//   - []*Token[T]: The tokens of the lexer.
//   - error: An error if the lexer encounters an error while lexing the input stream.
func FullLex(data []byte) ([]*gr.Token[TokenType], error) {
	lexer := NewLexer()

	tokens, err := ulpx.FullLex(lexer, data)
	if err != nil {
		return tokens, err
	}

	return tokens, nil
}
