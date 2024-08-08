package pkg

import (
	"fmt"
	"regexp"
	"strconv"
	"unicode"

	gr "github.com/PlayerR9/grammar/grammar"
	grlx "github.com/PlayerR9/grammar/lexer"
)

var (
	parse_lowercase_id *regexp.Regexp
	parse_uppercase_id *regexp.Regexp

	matcher *grlx.Matcher[TokenType]

	single_token_map map[rune]TokenType
)

func init() {
	parse_lowercase_id = regexp.MustCompile(`^[a-z]+([_][a-z]+)*[0-9]*`)
	parse_uppercase_id = regexp.MustCompile(`^([A-Z][a-z]*)+[0-9]*`)

	matcher = grlx.NewMatcher[TokenType]()

	_ = matcher.AddToMatch(TtkDot, ".")
	_ = matcher.AddToMatch(TtkOpParen, "(")
	_ = matcher.AddToMatch(TtkClParen, ")")
	_ = matcher.AddToMatch(TtkPipe, "|")
	_ = matcher.AddToMatch(TtkEqualSign, "=")
}

var (
	// Lexer is the lexer of the grammar.
	Lexer *grlx.Lexer[TokenType]
)

func init() {
	f := func(lexer *grlx.Lexer[TokenType]) (*gr.Token[TokenType], error) {
		// luc.Assert(len(l.input_stream) > 0, "l.input_stream is empty")

		match, err := matcher.Match(lexer)
		if err != nil {
			return nil, err
		}

		if match.IsValidMatch() {
			symbol, data := match.GetMatch()

			return gr.NewToken(symbol, data, lexer.Pos(), nil), nil
		}

		c, _, err := lexer.ReadRune()
		if err != nil {
			return nil, err
		}

		// lowercase ([_] lowercase)* -> lowercase id
		// lowercase ([_] lowercase)* digit -> lowercase id
		// lowercase :
		// 	[a-z]
		// 	| [a-z] lowercase
		// 	;
		// digit :
		// 	[0-9]
		// 	| [0-9] digit
		// 	;
		// [A-Z] lowercase? ([A-Z] lowercase?)* -> uppercase id
		// [A-Z] lowercase? ([A-Z] lowercase?)* digit -> uppercase id
		// digit :
		// 	[0-9]
		// 	| [0-9] digit
		// 	;
		// newline :
		// 	[\r]?[\n]
		// 	| [\r]?[\n] newline
		// 	;
		// whitespace :
		// 	[ \t]
		// 	| [ \t] whitespace
		// 	;

		c, err := lexer.Peek()
		if err != nil {
			return nil, err
		}

		at := lexer.At()

		token_type, ok := single_token_map[c]
		if ok {
			tk := gr.NewToken(token_type, string(c), at, nil)

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

			tk = gr.NewToken(TtkNewline, "\n", at, nil)
		default:
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
