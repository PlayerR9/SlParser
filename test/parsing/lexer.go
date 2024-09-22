package parsing

import (
	"fmt"
	"io"

	lxr "github.com/PlayerR9/SlParser/lexer"
	gcers "github.com/PlayerR9/go-commons/errors"
	dba "github.com/PlayerR9/go-debug/assert"
)

//go:generate stringer -type=TokenType

type TokenType int

const (
	EttInvalid TokenType = iota - 1
	EttEOF

	TttListComprehension
	TttPrintStmt
	TttNewline

	NttSource
	NttSource1
	NttStatement
)

func (t TokenType) IsTerminal() bool {
	return t <= TttNewline
}

var (
	Lexer *lxr.Lexer[TokenType]
)

func init() {
	builder := lxr.NewBuilder[TokenType]()

	// COMMENT : '#' .*? '\n' -> skip ;
	comment_fn := lxr.FragUntil('#', '\n', true)
	builder.RegisterSkip('#', comment_fn)

	builder.RegisterSkip('#', func(lexer lxr.RuneStreamer) (string, error) {
		for {
			char, err := lexer.NextRune()
			if err == io.EOF {
				break
			} else if err != nil {
				return "", err
			}

			if char == '\n' {
				break
			}
		}

		return "", nil
	})

	// NEWLINE : ('\r'? '\n')+ ;
	newline_fn := lxr.FragNewline(
		lxr.WithLexMany(true),
	)

	builder.Register('\r', func(lexer lxr.RuneStreamer, char rune) (TokenType, string, error) {
		char, err := lexer.NextRune()
		if err == io.EOF {
			return EttInvalid, "", fmt.Errorf("after %q, %w", '\r', gcers.NewErrValue("character", '\n', nil, true))
		} else if err != nil {
			return EttInvalid, "", err
		}

		if char != '\n' {
			return EttInvalid, "", fmt.Errorf("after %q, %w", '\r', gcers.NewErrValue("character", '\n', char, true))
		}

		str, err := newline_fn(lexer)
		if err != nil && err != lxr.NotFound {
			return EttInvalid, "", err
		}

		return TttNewline, str, nil
	})

	builder.Register('\n', func(lexer lxr.RuneStreamer, char rune) (TokenType, string, error) {
		str, err := newline_fn(lexer)
		if err != nil && err != lxr.NotFound {
			return EttInvalid, "", err
		}

		return TttNewline, str, nil
	})

	// LIST_COMPREHENSION : 'sq = [x * x for x in range(10)]' ;
	// PRINT_STMT : 'sq' ;
	builder.Register('s', func(lexer lxr.RuneStreamer, char rune) (TokenType, string, error) {
		char, err := lexer.NextRune()
		if err == io.EOF {
			return EttInvalid, "", fmt.Errorf("after %q, %w", 's', gcers.NewErrValue("character", 'q', nil, true))
		} else if err != nil {
			return EttInvalid, "", err
		}

		if char != 'q' {
			return EttInvalid, "", fmt.Errorf("after %q, %w", 's', gcers.NewErrValue("character", 'q', char, true))
		}

		next, err := lexer.NextRune()
		if err == io.EOF {
			return TttPrintStmt, "sq", nil
		} else if err != nil {
			return EttInvalid, "", err
		}

		if next != ' ' {
			err := lexer.UnreadRune()
			dba.AssertErr(err, "lexer.UnreadRune()")

			return TttPrintStmt, "sq", nil
		}

		word_fn := lxr.FragWord(" = [x * x for x in range(10)]")

		word, err := word_fn(lexer)
		if err != nil {
			return EttInvalid, "", err
		}

		return TttListComprehension, "sq" + word, nil
	})

	Lexer = builder.Build()
}
