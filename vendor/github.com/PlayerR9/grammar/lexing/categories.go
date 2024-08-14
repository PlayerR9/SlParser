package lexing

import (
	"fmt"
	"io"
	"unicode"

	dbg "github.com/PlayerR9/go-debug/assert"
)

var (
	// CatDecimal is the category of decimal digits (i.e., from 0 to 9).
	CatDecimal LexFunc

	// CatUppercase is the category of uppercase letters (i.e., from A to Z).
	CatUppercase LexFunc

	// CatLowercase is the category of lowercase letters (i.e., from a to z).
	CatLowercase LexFunc
)

func init() {
	CatDecimal = func(scanner io.RuneScanner) ([]rune, error) {
		// [0-9]

		c, _, err := scanner.ReadRune()
		if err != nil {
			return nil, err
		}

		if !unicode.IsDigit(c) {
			_ = scanner.UnreadRune()

			return nil, nil
		}

		return []rune{c}, nil
	}

	CatUppercase = func(scanner io.RuneScanner) ([]rune, error) {
		// [A-Z]

		c, _, err := scanner.ReadRune()
		if err != nil {
			return nil, err
		}

		if !unicode.IsUpper(c) {
			_ = scanner.UnreadRune()

			return nil, nil
		}

		return []rune{c}, nil
	}

	CatLowercase = func(scanner io.RuneScanner) ([]rune, error) {
		// [a-z]

		c, _, err := scanner.ReadRune()
		if err != nil {
			return nil, err
		}

		if !unicode.IsLower(c) {
			_ = scanner.UnreadRune()

			return nil, nil
		}

		return []rune{c}, nil
	}
}

// IsFunc is a function that checks whether a character is in a category.
//
// Parameters:
//   - c: The character to check.
//
// Returns:
//   - bool: True if the character is in the category, false otherwise.
type IsFunc func(c rune) bool

func LexCategory(scanner io.RuneScanner, allow_optional, allow_many bool, is_f IsFunc) ([]rune, error) {
	var chars []rune

	if allow_many {
		for {
			c, _, err := scanner.ReadRune()
			if err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			if !is_f(c) {
				err = scanner.UnreadRune()
				dbg.AssertErr(err, "scanner.UnreadRune()")

				break
			}

			chars = append(chars, c)
		}
	} else {
		c, _, err := scanner.ReadRune()
		if err == io.EOF {
			if !allow_optional {
				return nil, fmt.Errorf("expected a character")
			} else {

			}
		}

		if err != nil {
			return nil, err
		}

		if !is_f(c) {
			_ = scanner.UnreadRune()

			return nil, nil
		}

		chars = append(chars, c)
	}

	if !allow_optional {
		if len(chars) == 0 {
			return nil, nil
		}
	}

	return chars, nil
}
