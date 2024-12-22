package lexer

import (
	mtch "github.com/PlayerR9/SlParser/matcher"
	"github.com/PlayerR9/SlParser/mygo-lib/common"
	gch "github.com/PlayerR9/SlParser/mygo-lib/runes"
)

// MatchWord creates a matcher that matches a word.
//
// Parameters:
//   - str: The string to be matched.
//
// Returns:
//   - mtch.Matcher: A matcher that matches the given string.
//   - error: An error if the creation of the matcher fails.
//
// Errors:
//   - gch.ErrInvalidUtf8: If the string contains invalid utf-8 data.
func MatchWord(str string) (mtch.Matcher, error) {
	if str == "" {
		err := common.NewErrNilParam("str")
		return nil, err
	}

	chars, err := gch.StringToUtf8(str)
	if err != nil {
		return nil, err
	}

	m := mtch.Slice(chars)

	return m, nil
}
