package internal

import (
	"io"
	"unicode"

	"github.com/PlayerR9/SlParser/lexer"
	gers "github.com/PlayerR9/go-errors"
)

var (
	Lexer *lexer.Lexer[TokenType]
)

func init() {
	builder := lexer.NewBuilder[TokenType]()

	// TODO: Add here your own custom rules...

	// COLON : ':' ;
	builder.RegisterChar(':', TtColon)

	// SEMICOLON : ';' ;
	builder.RegisterChar(';', TtSemicolon)

	// WS : [ \t]+ -> skip ;
	builder.RegisterSkip(' ', lexer.FragWs(false))

	// NEWLINE : ('\r'? '\n')+ ;
	builder.Register('\n', func(stream lexer.RuneStreamer, char rune) (TokenType, error) {
		gers.AssertNotNil(stream, "stream")

		err := lexer.ApplyMany(stream, lexer.FragNewline)

		if err == nil || err == lexer.NotFound {
			return TtNewline, nil
		} else {
			return EtInvalid, err
		}
	})

	builder.Register('\r', func(stream lexer.RuneStreamer, char rune) (TokenType, error) {
		gers.AssertNotNil(stream, "stream")

		char, err := stream.NextRune()
		if err == io.EOF {
			return EtInvalid, lexer.NewErrGotNothing('\n', '\r')
		} else if err != nil {
			return EtInvalid, err
		}

		if char != '\n' {
			return EtInvalid, lexer.NewErrGotUnexpected('\r', '\n', char)
		}

		err = lexer.ApplyMany(stream, lexer.FragNewline)

		if err == nil || err == lexer.NotFound {
			return TtNewline, nil
		} else {
			return EtInvalid, err
		}
	})

	// fragment NONLOWER : [A-Z0-9] ;
	frag_nonlower := func(stream lexer.RuneStreamer) error {
		gers.AssertNotNil(stream, "stream")

		char, err := stream.NextRune()
		if err == io.EOF {
			return lexer.NotFound
		} else if err != nil {
			return err
		}

		if unicode.IsUpper(char) || unicode.IsDigit(char) {
			return nil
		}

		err = stream.UnreadRune()
		gers.AssertErr(err, "stream.UnreadRune()")

		return lexer.NotFound
	}

	// UPPERCASE_ID1 : UNDERSCORE NONLOWER+ ;
	frag_uppercase_id1 := func(stream lexer.RuneStreamer) error {
		gers.AssertNotNil(stream, "stream")

		char, err := stream.NextRune()
		if err == io.EOF {
			return lexer.NotFound
		} else if err != nil {
			return err
		}

		if char != '_' {
			err := stream.UnreadRune()
			gers.AssertErr(err, "stream.UnreadRune()")

			return lexer.NotFound
		}

		err = lexer.ApplyMany(stream, frag_nonlower)
		if err == lexer.NotFound {
			return lexer.NewErrBadGroup("uppercase letter or digit", &char)
		} else if err != nil {
			return err
		}

		return nil
	}

	// fragment ANY : [A-Za-z0-9] ;
	frag_any := func(stream lexer.RuneStreamer) error {
		gers.AssertNotNil(stream, "stream")

		char, err := stream.NextRune()
		if err == io.EOF {
			return lexer.NotFound
		} else if err != nil {
			return err
		}

		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			return nil
		}

		err = stream.UnreadRune()
		gers.AssertErr(err, "stream.UnreadRune()")

		return lexer.NotFound
	}

	builder.Default(func(stream lexer.RuneStreamer, char rune) (TokenType, error) {
		gers.AssertNotNil(stream, "stream")

		if !unicode.IsLetter(char) {
			return EtInvalid, lexer.NewErrBadGroup("letter", &char)
		}

		if unicode.IsLower(char) {
			// LOWERCASE_ID : LOWERCASE ANY* ;

			err := lexer.ApplyMany(stream, frag_any)
			if err == nil || err == lexer.NotFound {
				return TtLowercaseId, nil
			} else {
				return EtInvalid, lexer.NewErrBadGroup("lowercase and uppercase letter or digit", &char)
			}
		} else {
			err := stream.UnreadRune() // Push back the 'char' passed as argument.
			gers.AssertErr(err, "stream.UnreadRune()")

			// UPPERCASE_ID : UPPERCASE+ ;
			// UPPERCASE_ID : UPPERCASE+ UPPERCASE_ID1+ ;

			err = lexer.ApplyMany(stream, lexer.FragUppercase)
			if err == lexer.NotFound {
				return EtInvalid, lexer.NewErrBadGroup("uppercase letter", &char)
			} else if err != nil {
				return EtInvalid, err
			}

			err = lexer.ApplyMany(stream, frag_uppercase_id1)
			if err == nil || err == lexer.NotFound {
				return TtUppercaseId, nil
			} else {
				return EtInvalid, err
			}
		}
	})

	Lexer = builder.Build()
}
