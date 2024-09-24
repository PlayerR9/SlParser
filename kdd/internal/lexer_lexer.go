package internal

import (
	"io"
	"strings"
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
	builder.Register(':', func(stream lexer.RuneStreamer, char rune) (TokenType, string, error) {
		gers.AssertNotNil(stream, "stream")

		return TtColon, ":", nil
	})

	// SEMICOLON : ';' ;
	builder.Register(';', func(stream lexer.RuneStreamer, char rune) (TokenType, string, error) {
		gers.AssertNotNil(stream, "stream")

		return TtSemicolon, ";", nil
	})

	// WS : [ \t]+ -> skip ;
	builder.RegisterSkip(' ', lexer.FragWs(false))

	// NEWLINE : ('\r'? '\n')+ ;
	builder.Register('\n', func(stream lexer.RuneStreamer, char rune) (TokenType, string, error) {
		gers.AssertNotNil(stream, "stream")

		str, err := lexer.ApplyMany(stream, lexer.FragNewline)
		str = "\n" + str

		if err == nil || err == lexer.NotFound {
			return TtNewline, str, nil
		} else {
			return EtInvalid, str, err
		}
	})

	builder.Register('\r', func(stream lexer.RuneStreamer, char rune) (TokenType, string, error) {
		gers.AssertNotNil(stream, "stream")

		var builder strings.Builder
		builder.WriteRune('\r')

		char, err := stream.NextRune()
		if err == io.EOF {
			return EtInvalid, builder.String(), lexer.NewErrGotNothing('\n', '\r')
		} else if err != nil {
			return EtInvalid, builder.String(), lexer.NewErrInvalidInputStream(err)
		}

		builder.WriteRune(char)

		if char != '\n' {
			return EtInvalid, builder.String(), lexer.NewErrGotUnexpected('\r', '\n', char)
		}

		str, err := lexer.ApplyMany(stream, lexer.FragNewline)
		builder.WriteString(str)

		if err == nil || err == lexer.NotFound {
			return TtNewline, builder.String(), nil
		} else {
			return EtInvalid, builder.String(), err
		}
	})

	// fragment LOWERCASES : LOWERCASE ;
	// fragment LOWERCASES : LOWERCASE LOWERCASES ;
	frag_must_lowercases := lexer.FragGroup(lexer.GroupLower)

	frag_may_lowercases := lexer.FragGroup(lexer.GroupLower)

	// fragment CONT1 : UPPERCASE ;
	// fragment CONT1 : UPPERCASE LOWERCASES ;
	frag_cont1 := func(stream lexer.RuneStreamer) (string, error) {
		gers.AssertNotNil(stream, "stream")

		char, err := stream.NextRune()
		if err == io.EOF {
			return "", lexer.NewErrBadGroup("uppercase", nil)
		} else if err != nil {
			return "", lexer.NewErrInvalidInputStream(err)
		}

		if !unicode.IsUpper(char) {
			return "", lexer.NewErrBadGroup("uppercase", &char)
		}

		str, err := lexer.ApplyMany(stream, frag_may_lowercases)
		if err == lexer.NotFound {
			return string(char), nil
		} else if err != nil {
			return "", err
		}

		return string(char) + str, nil
	}

	// fragment CONT : CONT1 ;
	// fragment CONT : CONT1 CONT ;
	frag_conts := frag_cont1
	frag_may_conts := frag_cont1

	frag_uppercases := lexer.FragGroup(lexer.GroupUpper)

	// fragment UPPERCASE_ID1 : UPPERCASES ;
	// fragment UPPERCASE_ID1 : UPPERCASES UNDERSCORE UPPERCASE_ID1 ;

	frag_uppercase_id1 := func(stream lexer.RuneStreamer) (string, error) {
		gers.AssertNotNil(stream, "stream")

		var builder strings.Builder

		for {
			str, err := lexer.ApplyMany(stream, frag_uppercases)
			if err == lexer.NotFound {
				return "", lexer.NewErrBadGroup("uppercase", nil)
			} else if err != nil {
				return "", err
			}

			builder.WriteString(str)

			char, err := stream.NextRune()
			if err == io.EOF {
				break
			} else if err != nil {
				return "", lexer.NewErrInvalidInputStream(err)
			}

			if char != '_' {
				_ = stream.UnreadRune()

				break
			}
		}

		if builder.Len() == 0 {
			return "", lexer.NotFound
		}

		return builder.String(), nil
	}

	builder.Default(func(stream lexer.RuneStreamer, char rune) (TokenType, string, error) {
		gers.AssertNotNil(stream, "stream")

		char, err := stream.NextRune()
		if err == io.EOF {
			return EtInvalid, "", lexer.NewErrBadGroup("letter", nil)
		} else if err != nil {
			return EtInvalid, "", lexer.NewErrInvalidInputStream(err)
		}

		if !unicode.IsLetter(char) {
			return EtInvalid, "", lexer.NewErrBadGroup("letter", &char)
		}

		var builder strings.Builder

		builder.WriteRune(char)

		if unicode.IsLower(char) {
			char, err := stream.NextRune()
			if err == io.EOF {
				// nothing
				// LOWERCASE_ID : (LOWERCASE) ;

				return TtLowercaseId, builder.String(), nil
			} else if err != nil {
				return EtInvalid, "", lexer.NewErrInvalidInputStream(err)
			}

			if !unicode.IsLetter(char) {
				return EtInvalid, "", lexer.NewErrBadGroup("letter", &char)
			}

			builder.WriteRune(char)

			if unicode.IsLower(char) {
				// lowercases (LOWERCASE) ;

				lower, err := lexer.ApplyMany(stream, frag_must_lowercases)
				if err == lexer.NotFound {
					return EtInvalid, "", lexer.NewErrBadGroup("lowercase letter", &char)
				} else if err != nil {
					return EtInvalid, "", err
				}

				builder.WriteString(lower)

				// LOWERCASE_ID : (LOWERCASE) LOWERCASES ;
				// LOWERCASE_ID : (LOWERCASE) LOWERCASES CONT ;

				str, err := lexer.ApplyMany(stream, frag_may_conts)
				if err == nil {
					builder.WriteString(str)
				} else if err != lexer.NotFound {
					return EtInvalid, "", err
				}
			} else {
				// LOWERCASE_ID : (LOWERCASE) CONT ;

				str, err := lexer.ApplyMany(stream, frag_conts)
				if err == lexer.NotFound {
					return EtInvalid, "", lexer.NewErrBadGroup("uppercase letter", &char)
				} else if err != nil {
					return EtInvalid, "", err
				}

				builder.WriteString(str)
			}

			return TtLowercaseId, builder.String(), nil
		} else {
			// UPPERCASE_ID : (UPPERCASE) ;
			// UPPERCASE_ID : (UPPERCASE) UPPERCASE_ID1 ;
			// UPPERCASE_ID : (UPPERCASE) UNDERSCORE UPPERCASE_ID1 ;

			var must bool

			char, err := stream.NextRune()
			if err == io.EOF {
				must = false
			} else if err != nil {
				return EtInvalid, "", lexer.NewErrInvalidInputStream(err)
			}

			if char == '_' {
				must = true

				builder.WriteRune(char)
			} else {
				_ = stream.UnreadRune()

				must = false
			}

			str, err := frag_uppercase_id1(stream)
			if err == nil {
				builder.WriteString(str)
			} else if err != lexer.NotFound {
				return EtInvalid, "", err
			} else if must {
				return EtInvalid, "", lexer.NewErrBadGroup("uppercase letter", &char)
			}

			return TtUppercaseId, builder.String(), nil
		}

		// fragment LOWERCASE : [a-z];
		// fragment UPPERCASE : [A-Z];

		// fragment UNDERSCORE : '_' ;
	})

	Lexer = builder.Build()
}
