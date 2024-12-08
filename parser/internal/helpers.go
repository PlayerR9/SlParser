package internal

import "strings"

// EitherOrString is a function that returns a string representation of a slice
// of strings. Empty strings are ignored.
//
// Parameters:
//   - values: The values to convert to a string.
//
// Returns:
//   - string: The string representation.
//
// Example:
//
//	EitherOrString([]string{"a", "b", "c"}) // "either a, b, or c"
func EitherOrString(elems []string) string {
	var str string

	switch len(elems) {
	case 0:
		// Do nothing
	case 1:
		str = elems[0]
	case 2:
		str = "either " + elems[0] + " or " + elems[1]
	default:
		str = "either " + strings.Join(elems[:len(elems)-1], ", ") + ", or " + elems[len(elems)-1]
	}

	return str
}
