package pkg

import (
	"fmt"
	"io"
	"unicode"

	gr "github.com/PlayerR9/grammar/grammar"
	grlx "github.com/PlayerR9/grammar/lexer"
)

var (
	matcher *grlx.Matcher[TokenType]

	lex_whitespace grlx.LexFunc
	lex_digit      grlx.LexFunc
	lex_lowercase  grlx.LexFunc
	lex_newlines   grlx.LexFunc

	frag_uppercases grlx.LexFunc
	frag_lowercases grlx.LexFunc
)

func init() {
	matcher = grlx.NewMatcher[TokenType]()

	_ = matcher.AddToMatch(TtkDot, ".")
	_ = matcher.AddToMatch(TtkOpParen, "(")
	_ = matcher.AddToMatch(TtkClParen, ")")
	_ = matcher.AddToMatch(TtkPipe, "|")
	_ = matcher.AddToMatch(TtkEqualSign, "=")

	lex_whitespace = func(scanner io.RuneScanner) ([]rune, error) {
		// [ \t]+

		c, _, err := scanner.ReadRune()
		if err != nil {
			return nil, err
		}

		if c != ' ' && c != '\t' {
			_ = scanner.UnreadRune()

			return nil, grlx.Done
		}

		return []rune{c}, nil
	}

	lex_digit = func(scanner io.RuneScanner) ([]rune, error) {
		// [0-9]+

		c, _, err := scanner.ReadRune()
		if err != nil {
			return nil, err
		}

		if !unicode.IsDigit(c) {
			_ = scanner.UnreadRune()

			return nil, grlx.Done
		}

		return []rune{c}, nil
	}

	lex_newlines = func(scanner io.RuneScanner) ([]rune, error) {
		// ([\r]?[\n])+

		c1, _, err := scanner.ReadRune()
		if err != nil {
			return nil, err
		}

		if c1 == '\n' {
			return []rune{c1}, nil
		}

		if c1 != '\r' {
			_ = scanner.UnreadRune()

			return nil, grlx.Done
		}

		c2, _, err := scanner.ReadRune()
		if err == io.EOF {
			return nil, grlx.NewErrUnexpectedRune(&c1, nil, '\n')
		} else if err != nil {
			return nil, err
		}

		if c2 != '\n' {
			_ = scanner.UnreadRune()

			return nil, grlx.NewErrUnexpectedRune(&c1, &c2, '\n')
		}

		return []rune{c2}, nil
	}

	lex_lowercase = func(scanner io.RuneScanner) ([]rune, error) {
		// [a-z]+

		c, _, err := scanner.ReadRune()
		if err != nil {
			return nil, err
		}

		if !unicode.IsLower(c) {
			_ = scanner.UnreadRune()

			return nil, grlx.Done
		}

		return []rune{c}, nil
	}

	frag_uppercases = func(scanner io.RuneScanner) ([]rune, error) {
		// ([A-Z] | [A-Z] lowercase)+

		c, _, err := scanner.ReadRune()
		if err != nil {
			return nil, err
		}

		if !unicode.IsUpper(c) {
			_ = scanner.UnreadRune()

			return nil, grlx.Done
		}

		chars, err := grlx.RightLex(scanner, lex_lowercase)
		if err != nil {
			return nil, err
		}

		return append([]rune{c}, chars...), nil
	}

	frag_lowercases = func(scanner io.RuneScanner) ([]rune, error) {
		// (lowercase | lowercase [_])+

		chars, err := grlx.RightLex(scanner, lex_lowercase)
		if err != nil {
			return nil, err
		}

		c, _, err := scanner.ReadRune()
		if err == io.EOF {
			return chars, nil
		} else if err != nil {
			return nil, err
		}

		if c != '_' {
			_ = scanner.UnreadRune()

			return chars, nil
		}

		chars = append(chars, c)

		return chars, nil
	}
}

var (
	// Lexer is the lexer of the grammar.
	Lexer *grlx.Lexer[TokenType]
)

func match_rules(lexer *grlx.Lexer[TokenType]) (*gr.Token[TokenType], error) {
	chars, err := grlx.RightLex(lexer, lex_whitespace)
	if err != nil {
		return nil, err
	}

	if len(chars) != 0 {
		return nil, nil
	}

	chars, err = grlx.RightLex(lexer, lex_newlines)
	if err != nil {
		return nil, err
	}

	if len(chars) != 0 {
		return gr.NewToken(TtkNewline, "\n", lexer.Pos(), nil), nil
	}

	chars, err = grlx.RightLex(lexer, frag_uppercases)
	if err != nil {
		return nil, err
	}

	if len(chars) != 0 {
		// do digits

		digit, err := grlx.RightLex(lexer, lex_digit)
		if err != nil {
			return nil, err
		}

		chars = append(chars, digit...)

		return gr.NewToken(TtkUppercaseID, string(chars), lexer.Pos(), nil), nil
	}

	chars, err = grlx.RightLex(lexer, frag_lowercases)
	if err != nil {
		return nil, err
	}

	if len(chars) != 0 {
		// do digits

		digit, err := grlx.RightLex(lexer, lex_digit)
		if err != nil {
			return nil, err
		}

		chars = append(chars, digit...)

		return gr.NewToken(TtkLowercaseID, string(chars), lexer.Pos(), nil), nil
	}

	return nil, fmt.Errorf("no match found at %d", lexer.Pos())
}

func init() {
	f := func(lexer *grlx.Lexer[TokenType]) (*gr.Token[TokenType], error) {
		// luc.Assert(len(l.input_stream) > 0, "l.input_stream is empty")

		match, _ := matcher.Match(lexer)

		if match.IsValidMatch() {
			symbol, data := match.GetMatch()

			return gr.NewToken(symbol, data, lexer.Pos(), nil), nil
		}

		tk, err := match_rules(lexer)
		if err != nil {
			return nil, err
		}

		return tk, nil
	}

	Lexer = grlx.NewLexer(f)
}
