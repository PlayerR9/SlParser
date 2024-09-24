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

	// fragment LOWERCASES : LOWERCASE ;
	// fragment LOWERCASES : LOWERCASE LOWERCASES ;
	frag_must_lowercases := lexer.FragGroup(lexer.GroupLower)

	frag_may_lowercases := lexer.FragGroup(lexer.GroupLower)

	// fragment CONT1 : UPPERCASE ;
	// fragment CONT1 : UPPERCASE LOWERCASES ;
	frag_cont1 := func(stream lexer.RuneStreamer) error {
		gers.AssertNotNil(stream, "stream")

		char, err := stream.NextRune()
		if err == io.EOF {
			return lexer.NewErrBadGroup("uppercase", nil)
		} else if err != nil {
			return err
		}

		if !unicode.IsUpper(char) {
			err := stream.UnreadRune()
			gers.AssertErr(err, "stream.UnreadRune()")

			return lexer.NotFound
		}

		err = lexer.ApplyMany(stream, frag_may_lowercases)
		if err == lexer.NotFound {
			return nil
		} else if err != nil {
			return err
		}

		return nil
	}

	// fragment CONT : CONT1 ;
	// fragment CONT : CONT1 CONT ;
	frag_conts := frag_cont1
	frag_may_conts := frag_cont1

	frag_uppercases := lexer.FragGroup(lexer.GroupUpper)

	// fragment UPPERCASE_ID1 : UPPERCASES ;
	// fragment UPPERCASE_ID1 : UPPERCASES UNDERSCORE UPPERCASE_ID1 ;

	frag_uppercase_id1 := func(stream lexer.RuneStreamer) error {
		gers.AssertNotNil(stream, "stream")

		for {
			err := lexer.ApplyMany(stream, frag_uppercases)
			if err == lexer.NotFound {
				return lexer.NewErrBadGroup("uppercase", nil)
			} else if err != nil {
				return err
			}

			char, err := stream.NextRune()
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}

			if char != '_' {
				_ = stream.UnreadRune()

				break
			}
		}

		return nil
	}

	builder.Default(func(stream lexer.RuneStreamer, char rune) (TokenType, error) {
		gers.AssertNotNil(stream, "stream")

		if !unicode.IsLetter(char) {
			return EtInvalid, lexer.NewErrBadGroup("letter", &char)
		}

		if unicode.IsLower(char) {
			char, err := stream.NextRune()
			if err == io.EOF {
				// nothing
				// LOWERCASE_ID : (LOWERCASE) ;

				return TtLowercaseId, nil
			} else if err != nil {
				return EtInvalid, err
			}

			if !unicode.IsLetter(char) {
				return EtInvalid, lexer.NewErrBadGroup("letter", &char)
			}

			if unicode.IsLower(char) {
				// lowercases (LOWERCASE) ;

				err := lexer.ApplyMany(stream, frag_must_lowercases)
				if err == lexer.NotFound {
					return EtInvalid, lexer.NewErrBadGroup("lowercase letter", &char)
				} else if err != nil {
					return EtInvalid, err
				}

				// LOWERCASE_ID : (LOWERCASE) LOWERCASES ;
				// LOWERCASE_ID : (LOWERCASE) LOWERCASES CONT ;

				err = lexer.ApplyMany(stream, frag_may_conts)
				if err != nil && err != lexer.NotFound {
					return EtInvalid, err
				}
			} else {
				// LOWERCASE_ID : (LOWERCASE) CONT ;

				err := lexer.ApplyMany(stream, frag_conts)
				if err == lexer.NotFound {
					return EtInvalid, lexer.NewErrBadGroup("uppercase letter", &char)
				} else if err != nil {
					return EtInvalid, err
				}
			}

			return TtLowercaseId, nil
		} else {
			// UPPERCASE_ID : (UPPERCASE) ;
			// UPPERCASE_ID : (UPPERCASE) UPPERCASE_ID1 ;
			// UPPERCASE_ID : (UPPERCASE) UNDERSCORE UPPERCASE_ID1 ;

			var must bool

			char, err := stream.NextRune()
			if err == io.EOF {
				must = false
			} else if err != nil {
				return EtInvalid, err
			}

			if char == '_' {
				must = true
			} else {
				_ = stream.UnreadRune()

				must = false
			}

			err = frag_uppercase_id1(stream)
			if err == nil {
				return TtUppercaseId, nil
			} else if err != lexer.NotFound {
				return EtInvalid, err
			} else if must {
				return EtInvalid, lexer.NewErrBadGroup("uppercase letter", &char)
			} else {
				return TtUppercaseId, nil
			}
		}

		// fragment LOWERCASE : [a-z];
		// fragment UPPERCASE : [A-Z];

		// fragment UNDERSCORE : '_' ;
	})

	Lexer = builder.Build()
}
