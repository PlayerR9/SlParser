package grammar

import "strings"

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
	return strings.HasPrefix(rhs, "T") || strings.HasPrefix(rhs, "E")
}
