package grammar

import "unicode/utf8"

/////////////////////////////////////////////////////////

const (
	// EtEOF is a special token type used to indicate the end of the input (or file).
	EtEOF string = "EtEOF"

	// EtToSkip is a special token type used to indicate the token must be skipped.
	EtToSkip string = "EtToSkip"
)

// IsTerminal checks if the token type is a terminal. Terminal token types
// start with "T" or "E".
//
// Parameters:
//   - rhs: The token type.
//
// Returns:
//   - bool: True if the token type is a terminal, false otherwise.
func IsTerminal(rhs string) bool {
	if rhs == "" {
		return false
	}

	c, _ := utf8.DecodeRuneInString(rhs)
	if c == 'T' || c == 'E' {
		return true
	} else {
		return false
	}
}
