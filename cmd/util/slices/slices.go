package slices

import (
	"cmp"
	"slices"
)

// OrderedUniquefy returns a copy of elems without duplicates.
//
// Parameters:
//   - elems: The elements to uniquefy.
//
// Returns:
//   - []T: The unique elements.
//
// This function also sorts the elements.
func OrderedUniquefy[T cmp.Ordered](elems []T) []T {
	if len(elems) == 0 {
		return nil
	}

	var unique []T

	for _, elem := range elems {
		pos, ok := slices.BinarySearch(unique, elem)
		if !ok {
			unique = slices.Insert(unique, pos, elem)
		}
	}

	return unique
}
