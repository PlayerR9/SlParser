package ast

import (
	"fmt"
	"strconv"
)

// TODO: Remove this once go-commons is updated.

// Quote returns a string representation of the given element, enclosed in double
// quotes.
//
// Parameters:
//   - elem: The element to get the string representation of.
//
// Returns:
//   - string: The string representation of the element. If the element is nil, the
//     function returns an empty string enclosed in double quotes.
func Quote(elem any) string {
	if elem == nil {
		return "\"\""
	}

	str := fmt.Sprintf("%v", elem)
	str = strconv.Quote(str)

	return str
}
