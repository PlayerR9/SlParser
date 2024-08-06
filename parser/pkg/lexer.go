package pkg

import (
	"fmt"
	"regexp"
	"strconv"
	"unicode"

	grlx "github.com/PlayerR9/SLParser/util/lexer"
	gr "github.com/PlayerR9/grammar/grammar"
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

var (
	// Lexer is the lexer of the grammar.
	Lexer *grlx.Lexer[TokenType]
)

func init() {
	f := func(lexer *grlx.Lexer[TokenType]) (*gr.Token[TokenType], error) {
		// luc.Assert(len(l.input_stream) > 0, "l.input_stream is empty")

		c, err := lexer.Peek()
		if grlx.IsExhausted(err) {
			return nil, fmt.Errorf("expected character, got nothing instead")
		} else if err != nil {
			return nil, err
		}

		token_type, ok := single_token_map[c]
		if ok {
			tk := gr.NewToken(token_type, string(c), lexer.At(), nil)

			_, _ = lexer.Next()

			return tk, nil
		}

		var tk *gr.Token[TokenType]

		switch c {
		case ' ', '\t':
			// Do nothing

			_, _ = lexer.Next()
		case '\r':
			_, err := lexer.MatchChars([]rune{'\r', '\n'})
			if err != nil {
				return nil, err
			}

			tk = gr.NewToken(TtkNewline, "\n", lexer.At(), nil)
		default:
			at := lexer.At()

			var data string
			var s TokenType

			if !unicode.IsLetter(c) {
				tmp, ok := lexer.MatchRegex(parse_newlines)
				if !ok {
					return nil, fmt.Errorf("invalid character: %s", strconv.QuoteRune(c))
				}

				data = tmp
				s = TtkNewline
			} else if unicode.IsUpper(c) {
				tmp, ok := lexer.MatchRegex(parse_uppercase_id)
				if !ok {
					return nil, fmt.Errorf("invalid uppercase identifier: %s", strconv.QuoteRune(c))
				}

				data = tmp
				s = TtkUppercaseID
			} else {
				tmp, ok := lexer.MatchRegex(parse_lowercase_id)
				if !ok {
					return nil, fmt.Errorf("invalid lowercase identifier: %s", strconv.QuoteRune(c))
				}

				data = tmp
				s = TtkLowercaseID
			}

			tk = gr.NewToken(s, data, at, nil)
		}

		return tk, nil
	}

	Lexer = grlx.NewLexer(f)
}
