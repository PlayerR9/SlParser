package internal

import (
	"fmt"
	"io"

	"github.com/PlayerR9/SlParser/lexer"
	"github.com/PlayerR9/SlParser/parser"
	gers "github.com/PlayerR9/go-errors"
)

//go:generate stringer -type=TokenType

type TokenType int

const (
	EtInvalid TokenType = iota - 1
	EtEOF
	TtListComprehension
	TtPrintStmt
	TtNewline
	NtSource
	NtSource1
	NtStatement
)

func (t TokenType) IsTerminal() bool {
	return t <= TtNewline
}

var (
	Lexer  *lexer.Lexer[TokenType]
	Parser *parser.Parser[TokenType]
)

func init() {
	is := parser.NewItemSet[TokenType]()

	_ = is.AddRule(NtSource, TtNewline, NtSource1, EtEOF)
	_ = is.AddRule(NtSource1, NtStatement)
	_ = is.AddRule(NtSource1, NtStatement, TtNewline, NtSource1)
	_ = is.AddRule(NtStatement, TtListComprehension)
	_ = is.AddRule(NtStatement, TtPrintStmt)

	Parser = parser.Build(&is)

	builder := lexer.NewBuilder[TokenType]()

	// COMMENT : '#' .*? '\n' -> skip ;
	comment_fn := lexer.FragUntil('#', '\n', true)
	builder.RegisterSkip('#', comment_fn)

	// NEWLINE : ('\r'? '\n')+ ;
	builder.Register('\r', func(stream lexer.RuneStreamer, char rune) (TokenType, string, error) {
		char, err := stream.NextRune()
		if err == io.EOF {
			return EtInvalid, "", fmt.Errorf("after %q, %w", '\r',
				fmt.Errorf("expected %q, got nothing instead", '\n'),
			)
		} else if err != nil {
			return EtInvalid, "", err
		}

		if char != '\n' {
			return EtInvalid, "", fmt.Errorf("after %q, %w", '\r',
				fmt.Errorf("expected %q, got %q instead", '\n', char),
			)
		}

		str, err := lexer.ApplyMany(stream, lexer.FragNewline)
		if err != nil && err != lexer.NotFound {
			return EtInvalid, "", err
		}

		return TtNewline, str, nil
	})

	builder.Register('\n', func(stream lexer.RuneStreamer, char rune) (TokenType, string, error) {
		str, err := lexer.ApplyMany(stream, lexer.FragNewline)
		if err != nil && err != lexer.NotFound {
			return EtInvalid, "", err
		}

		return TtNewline, str, nil
	})

	// LIST_COMPREHENSION : 'sq = [x * x for x in range(10)]' ;
	// PRINT_STMT : 'sq' ;
	word_fn := lexer.FragWord(" = [x * x for x in range(10)]")

	builder.Register('s', func(stream lexer.RuneStreamer, char rune) (TokenType, string, error) {
		char, err := stream.NextRune()
		if err == io.EOF {
			return EtInvalid, "", fmt.Errorf("after %q, %w", 's',
				fmt.Errorf("expected %q, got nothing instead", 'q'),
			)
		} else if err != nil {
			return EtInvalid, "", err
		}

		if char != 'q' {
			return EtInvalid, "", fmt.Errorf("after %q, %w", 's',
				fmt.Errorf("expected %q, got %q instead", 'q', char),
			)
		}

		next, err := stream.NextRune()
		if err == io.EOF {
			return TtPrintStmt, "sq", nil
		} else if err != nil {
			return EtInvalid, "", err
		}

		if next != ' ' {
			err := stream.UnreadRune()
			gers.AssertErr(err, "lexer.UnreadRune()")

			return TtPrintStmt, "sq", nil
		}

		word, err := word_fn(stream)
		if err != nil {
			return EtInvalid, "", err
		}

		return TtListComprehension, "sq" + word, nil
	})

	Lexer = builder.Build()
}
