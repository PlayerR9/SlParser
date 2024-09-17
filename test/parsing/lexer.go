package parsing

import (
	"io"

	gr "github.com/PlayerR9/SlParser/grammar"
	lxr "github.com/PlayerR9/SlParser/lexer"
)

//go:generate stringer -type=TokenType

type TokenType int

const (
	EttEOF TokenType = iota

	TttListComprehension
	TttPrintStmt
	TttNewline
)

func (t TokenType) IsTerminal() bool {
	return t <= TttNewline
}

var (
	Lexer lxr.Lexer[TokenType]
)

func init() {
	builder := lxr.NewBuilder[TokenType]()

	// COMMENT : '#' .*? '\n' -> skip ;
	comment_fn := lxr.FragUntil[TokenType]('#', '\n', true)
	builder.RegisterSkip('#', comment_fn)

	builder.RegisterSkip('#', func(lexer *lxr.Lexer[TokenType]) (string, error) {
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
	newline_fn := lxr.FragNewline[TokenType](
		lxr.WithLexMany(true),
	)

	builder.Register('\r', func(lexer *lxr.Lexer[TokenType], char rune) (*gr.Token[TokenType], error) {
		_, err := newline_fn(lexer)
		if err != nil && err != lxr.NotFound {
			return nil, err
		}

		tk := gr.NewTerminalToken(TttNewline, "\n")
		return tk, nil
	})

	builder.Register('\n', func(lexer *lxr.Lexer[TokenType], char rune) (*gr.Token[TokenType], error) {
		_, err := newline_fn(lexer)
		if err != nil && err != lxr.NotFound {
			return nil, err
		}

		tk := gr.NewTerminalToken(TttNewline, "\n")
		return tk, nil
	})

	// LIST_COMPREHENSION : 'sq = [x * x for x in range(10)]' ;
	// PRINT_STMT : 'sq' ;
	builder.Register('s', func(lexer *lxr.Lexer[TokenType], char rune) (*gr.Token[TokenType], error) {
		char, err := lexer.NextRune()
		if err == io.EOF {
			return nil, lxr.NewErrUnexpectedChar('s', []rune{'q'}, nil)
		} else if err != nil {
			return nil, err
		}

		if char != 'q' {
			return nil, lxr.NewErrUnexpectedChar('s', []rune{'q'}, &char)
		}

		next, err := lexer.NextRune()
		if err == io.EOF {
			tk := gr.NewTerminalToken(TttPrintStmt, "sq")
			return tk, nil
		} else if err != nil {
			return nil, err
		}

		if next != ' ' {
			tk := gr.NewTerminalToken(TttPrintStmt, "sq")
			return tk, nil
		}

		word_fn := lxr.FragWord[TokenType](" = [x * x for x in range(10)]")

		word, err := word_fn(lexer)
		if err != nil {
			return nil, err
		}

		tk := gr.NewTerminalToken(TttListComprehension, "sq"+word)

		return tk, nil
	})

	Lexer = builder.Build()
}
