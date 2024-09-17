package lexer

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
