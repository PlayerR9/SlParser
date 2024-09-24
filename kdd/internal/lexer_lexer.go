package internal

import (
	"io"
	"strings"
	"unicode"

	"github.com/PlayerR9/SlParser/lexer"
)

var (
	Lexer *lexer.Lexer[TokenType]
)

func init() {
	builder := lexer.NewBuilder[TokenType]()

	// TODO: Add here your own custom rules...

	// COLON : ':' ;
	builder.Register(':', func(stream lexer.RuneStreamer, char rune) (TokenType, string, error) {
		return TtColon, ":", nil
	})

	// SEMICOLON : ';' ;
	builder.Register(';', func(stream lexer.RuneStreamer, char rune) (TokenType, string, error) {
		return TtSemicolon, ";", nil
	})

	// WS : [ \t]+ -> skip ;
	builder.RegisterSkip(' ', lexer.FragWs(false,
		lexer.WithAllowOptional(false),
		lexer.WithLexMany(true),
	))

	// NEWLINE : ('\r'? '\n')+ ;
	newline_frag := lexer.FragNewline(
		lexer.WithLexMany(true),
	)
	builder.Register('\n', func(stream lexer.RuneStreamer, char rune) (TokenType, string, error) {
		str, err := newline_frag(stream)
		if err == lexer.NotFound {
			return TtNewline, "\n", nil
		} else if err != nil {
			return EtInvalid, "", err
		}

		return TtNewline, "\n" + str, nil
	})

	builder.Register('\r', func(stream lexer.RuneStreamer, char rune) (TokenType, string, error) {
		char, err := stream.NextRune()
		if err == io.EOF {
			return EtInvalid, "", lexer.NewErrGotNothing('\r', '\n')
		} else if err != nil {
			return EtInvalid, "", lexer.NewErrInvalidInputStream(err)
		}

		if char != '\n' {
			return EtInvalid, "", lexer.NewErrGotUnexpected('\r', '\n', char)
		}

		str, err := newline_frag(stream)
		if err == lexer.NotFound {
			return TtNewline, "\r\n", nil
		} else if err != nil {
			return EtInvalid, "", err
		}

		return TtNewline, "\r\n" + str, nil
	})

	builder.Default(func(stream lexer.RuneStreamer, char rune) (TokenType, string, error) {
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
		prev := char

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

			if unicode.IsLower(char) {
				// lowercases (LOWERCASE) ;
			} else {
				// cont (cont1 (UPPERCASE)) ;
			}

			// LOWERCASE_ID : (LOWERCASE) CONT ;
			// LOWERCASE_ID : (LOWERCASE) LOWERCASES ;
			// LOWERCASE_ID : (LOWERCASE) LOWERCASES CONT ;

			// fragment CONT : CONT1 ;
			// fragment CONT : CONT1 CONT ;

			// fragment CONT1 : UPPERCASE ;
			// fragment CONT1 : UPPERCASE LOWERCASES ;
		} else {
			// UPPERCASE_ID : (UPPERCASE) ;
			// UPPERCASE_ID : (UPPERCASE) UPPERCASE_ID1 ;
			// UPPERCASE_ID : (UPPERCASE) UNDERSCORE UPPERCASE_ID1 ;

			// fragment UPPERCASE_ID1 : UPPERCASES ;
			// fragment UPPERCASE_ID1 : UPPERCASES UNDERSCORE UPPERCASE_ID1 ;
		}

		// fragment LOWERCASES : LOWERCASE ;
		// fragment LOWERCASES : LOWERCASE LOWERCASES ;

		// fragment LOWERCASE : [a-z];
		// fragment UPPERCASE : [A-Z];

		// fragment UNDERSCORE : '_' ;
	})

	Lexer = builder.Build()
}
