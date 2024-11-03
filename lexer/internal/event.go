package internal

// Event is the result of a match.
type Event struct {
	// Type is the type of the match.
	Type string

	// Data is the data of the match.
	Data []rune
}
