package lexing

import (
	"io"
	"unicode"

	"github.com/PlayerR9/grammar/lexing"
)

var (
	CatDecimal LexFunc
)

func init() {

	cat_decimal = func(scanner io.RuneScanner) ([]rune, error) {
		// [0-9]

		c, _, err := scanner.ReadRune()
		if err != nil {
			return nil, err
		}

		if !unicode.IsDigit(c) {
			_ = scanner.UnreadRune()

			return nil, lexing.Done
		}

		return []rune{c}, lexing.Done
	}

	cat_uppercase = func(scanner io.RuneScanner) ([]rune, error) {
		// [A-Z]

		c, _, err := scanner.ReadRune()
		if err != nil {
			return nil, err
		}

		if !unicode.IsUpper(c) {
			_ = scanner.UnreadRune()

			return nil, lexing.Done
		}

		return []rune{c}, lexing.Done
	}

	cat_lowercase = func(scanner io.RuneScanner) ([]rune, error) {
		// [a-z]

		c, _, err := scanner.ReadRune()
		if err != nil {
			return nil, err
		}

		if !unicode.IsLower(c) {
			_ = scanner.UnreadRune()

			return nil, lexing.Done
		}

		return []rune{c}, lexing.Done
	}
}
