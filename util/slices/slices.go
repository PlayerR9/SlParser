package slices

import (
	"cmp"
	"slices"
)

// DeleteElem deletes an element from a slice and returns the new slice.
// If the element is not found, the slice is returned unchanged.
//
// Parameters:
//   - slice: The slice to delete the element from.
//   - elem: The element to delete.
//
// Returns:
//   - []T: The new slice with the element deleted.
func DeleteElem[T cmp.Ordered](slice []T, elem T) []T {
	pos, ok := slices.BinarySearch(slice, elem)
	if !ok {
		return slice
	}

	return slices.Delete(slice, pos, pos+1)
}
