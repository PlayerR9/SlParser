package internal

import (
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
	builder.Register('\r', func(stream lexer.RuneStreamer, char rune) (TokenType, error) {
		gers.AssertNotNil(stream, "stream")

		char, err := stream.NextRune()
		if err == io.EOF {
			return EtInvalid, lexer.NewErrGotNothing('\r', '\n')
		} else if err != nil {
			return EtInvalid, err
		}

		if char != '\n' {
			return EtInvalid, lexer.NewErrGotUnexpected('\r', '\n', char)
		}

		err = lexer.ApplyMany(stream, lexer.FragNewline)
		if err != nil && err != lexer.NotFound {
			return EtInvalid, err
		}

		return TtNewline, nil
	})

	builder.Register('\n', func(stream lexer.RuneStreamer, char rune) (TokenType, error) {
		gers.AssertNotNil(stream, "stream")

		err := lexer.ApplyMany(stream, lexer.FragNewline)
		if err != nil && err != lexer.NotFound {
			return EtInvalid, err
		}

		return TtNewline, nil
	})

	// LIST_COMPREHENSION : 'sq = [x * x for x in range(10)]' ;
	// PRINT_STMT : 'sq' ;
	word_fn := lexer.FragWord(" = [x * x for x in range(10)]")

	builder.Register('s', func(stream lexer.RuneStreamer, char rune) (TokenType, error) {
		gers.AssertNotNil(stream, "stream")

		char, err := stream.NextRune()
		if err == io.EOF {
			return EtInvalid, lexer.NewErrGotNothing('s', 'q')
		} else if err != nil {
			return EtInvalid, err
		}

		if char != 'q' {
			return EtInvalid, lexer.NewErrGotUnexpected('s', 'q', char)
		}

		next, err := stream.NextRune()
		if err == io.EOF {
			return TtPrintStmt, nil
		} else if err != nil {
			return EtInvalid, err
		}

		if next != ' ' {
			err := stream.UnreadRune()
			gers.AssertErr(err, "lexer.UnreadRune()")

			return TtPrintStmt, nil
		}

		err = word_fn(stream)
		if err != nil {
			return EtInvalid, err
		}

		return TtListComprehension, nil
	})

	Lexer = builder.Build()
}
