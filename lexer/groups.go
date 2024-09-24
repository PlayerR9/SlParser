package lexer

import (
	"fmt"
	"unicode"
)

// GroupFn is a function that checks if a character belongs to a group.
//
// Parameters:
//   - char: the character to check.
//
// Returns:
//   - bool: true if the character belongs to the group, false otherwise.
type GroupFn func(char rune) bool

var (
	// GroupWs is the group of whitespace characters that does not include newlines.
	// (i.e. ' ', '\t')
	GroupWs GroupFn

	// GroupWsNl is the group of whitespace characters that includes newlines.
	// (i.e. ' ', '\t', '\n', '\r')
	GroupWsNl GroupFn

	// GroupUpper is the group of uppercase characters.
	// (i.e. 'A', 'B', 'C', ...)
	GroupUpper GroupFn

	// GroupLower is the group of lowercase characters.
	// (i.e. 'a', 'b', 'c', ...)
	GroupLower GroupFn
)

func init() {
	GroupWs = func(char rune) bool {
		return char == ' ' || char == '\t'
	}

	GroupWsNl = func(char rune) bool {
		return char == ' ' || char == '\t' || char == '\n' || char == '\r'
	}

	GroupUpper = unicode.IsUpper

	GroupLower = unicode.IsLower
}

func MakeGroup(from, to rune) (GroupFn, error) {
	if !unicode.IsLetter(from) {
		if !unicode.IsDigit(from) {
			return nil, fmt.Errorf("from must be a letter or a digit")
		}

		if !unicode.IsDigit(to) {
			return nil, fmt.Errorf("to must be a digit")
		}

		// NUMERIC

		if from > to {
			from, to = to, from
		}

		var fn GroupFn

		if from == to {
			fn = func(char rune) bool {
				return char == from
			}
		} else {
			fn = func(char rune) bool {
				return char >= from && char <= to
			}
		}

		return fn, nil
	}

	// LETTER

	if unicode.IsUpper(from) {
		if !unicode.IsUpper(to) {
			return nil, fmt.Errorf("to must be an uppercase letter")
		}
	} else {
		if !unicode.IsLower(to) {
			return nil, fmt.Errorf("to must be a lowercase letter")
		}
	}

	if from > to {
		from, to = to, from
	}

	var fn GroupFn

	if from == to {
		fn = func(char rune) bool {
			return char == from
		}
	} else {
		fn = func(char rune) bool {
			return char >= from && char <= to
		}
	}

	return fn, nil
}

/*
func GroupFromString(str string) (GroupFn, error) {
	chars, err := gcch.StringToUtf8(str)
	if err != nil {
		return nil, err
	}

	if len(chars) == 0 {
		return nil, fmt.Errorf("expected '[', got nothing instead")
	}

	if len(chars) == 1 {
		if chars[0] == '[' {
			return nil, fmt.Errorf("missing close group ']'")
		} else if chars[0] == ']' {
			return nil, fmt.Errorf("missing open group '['")
		}

		return nil, fmt.Errorf("expected '[', got %q instead", chars[0])
	}

	if chars[len(chars)-1] != ']' {
		return nil, fmt.Errorf("expected ']', got %q instead", chars[len(chars)-1])
	}

	chars = chars[1 : len(chars)-1]

	if len(chars) == 0 {
		return nil, nil
	}

	from := chars[0]
}
*/
