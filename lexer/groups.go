package lexer

import (
	"fmt"
	"io"
	"unicode"

	gers "github.com/PlayerR9/go-errors"
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
)

func init() {
	GroupWs = func(char rune) bool {
		return char == ' ' || char == '\t'
	}

	GroupWsNl = func(char rune) bool {
		return char == ' ' || char == '\t' || char == '\n' || char == '\r'
	}
}

// FragUppercase is a fragment that checks if the current character is uppercase.
//
// Parameters:
//   - stream: the rune streamer.
//
// Returns:
//   - error: if an error occurred.
func FragUppercase(stream RuneStreamer) error {
	gers.AssertNotNil(stream, "stream")

	char, err := stream.NextRune()
	if err == io.EOF {
		return NotFound
	} else if err != nil {
		return err
	}

	if unicode.IsUpper(char) {
		return nil
	}

	err = stream.UnreadRune()
	gers.AssertErr(err, "stream.UnreadRune()")

	return NotFound
}

// FragLowercase is a fragment that checks if the current character is lowercase.
//
// Parameters:
//   - stream: the rune streamer.
//
// Returns:
//   - error: if an error occurred.
func FragLowercase(stream RuneStreamer) error {
	gers.AssertNotNil(stream, "stream")

	char, err := stream.NextRune()
	if err == io.EOF {
		return NotFound
	} else if err != nil {
		return err
	}

	if unicode.IsLower(char) {
		return nil
	}

	err = stream.UnreadRune()
	gers.AssertErr(err, "stream.UnreadRune()")

	return NotFound
}

// FragLetter is a fragment that checks if the current character is a letter.
//
// Parameters:
//   - stream: the rune streamer.
//
// Returns:
//   - error: if an error occurred.
func FragLetter(stream RuneStreamer) error {
	gers.AssertNotNil(stream, "stream")

	char, err := stream.NextRune()
	if err == io.EOF {
		return NotFound
	} else if err != nil {
		return err
	}

	if unicode.IsLetter(char) {
		return nil
	}

	err = stream.UnreadRune()
	gers.AssertErr(err, "stream.UnreadRune()")

	return NotFound
}

// FragDigit is a fragment that checks if the current character is a digit.
//
// Parameters:
//   - stream: the rune streamer.
//
// Returns:
//   - error: if an error occurred.
func FragDigit(stream RuneStreamer) error {
	gers.AssertNotNil(stream, "stream")

	char, err := stream.NextRune()
	if err == io.EOF {
		return NotFound
	} else if err != nil {
		return err
	}

	if unicode.IsDigit(char) {
		return nil
	}

	err = stream.UnreadRune()
	gers.AssertErr(err, "stream.UnreadRune()")

	return NotFound
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
