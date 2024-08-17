package cmp

import (
	"cmp"
	"slices"
)

// Min returns the minimum of a and b.
//
// Parameters:
//   - a: The first value.
//   - b: The second value.
//
// Returns:
//   - T: The minimum value.
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}

	return b
}

// Max returns the maximum of a and b.
//
// Parameters:
//   - a: The first value.
//   - b: The second value.
//
// Returns:
//   - T: The maximum value.
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}

	return b
}

// MinVar returns the minimum of the elements in the slice.
//
// Parameters:
//   - elems: The elements.
//
// Returns:
//   - T: The minimum value.
//   - bool: Whether the operation was successful.
func MinVar[T cmp.Ordered](elems []T) (T, bool) {
	if len(elems) == 0 {
		return *new(T), false
	}

	idx := 0
	min := elems[0]

	for i := 1; i < len(elems); i++ {
		if elems[i] < min {
			min = elems[i]
			idx = i
		}
	}

	return elems[idx], true
}

// MaxVar returns the maximum of the elements in the slice.
//
// Parameters:
//   - elems: The elements.
//
// Returns:
//   - T: The maximum value.
//   - bool: Whether the operation was successful.
func MaxVar[T cmp.Ordered](elems []T) (T, bool) {
	if len(elems) == 0 {
		return *new(T), false
	}

	idx := 0
	max := elems[0]

	for i := 1; i < len(elems); i++ {
		if elems[i] > max {
			max = elems[i]
			idx = i
		}
	}

	return elems[idx], true
}

// Mins returns the minimums of the elements in the slice.
//
// Parameters:
//   - elems: The elements.
//
// Returns:
//   - []int: The indices of the minimum values.
func Mins[T cmp.Ordered](elems []T) []int {
	if len(elems) == 0 {
		return nil
	} else if len(elems) == 1 {
		return []int{0}
	}

	min := elems[0]
	indices := []int{0}

	for i := 1; i < len(elems); i++ {
		if elems[i] < min {
			min = elems[i]
			indices = []int{i}
		} else if elems[i] == min {
			indices = append(indices, i)
		}
	}

	return indices
}

// Maxs returns the maximums of the elements in the slice.
//
// Parameters:
//   - elems: The elements.
//
// Returns:
//   - []int: The indices of the maximum values.
func Maxs[T cmp.Ordered](elems []T) []int {
	if len(elems) == 0 {
		return nil
	} else if len(elems) == 1 {
		return []int{0}
	}

	max := elems[0]
	res := []int{0}

	for i := 1; i < len(elems); i++ {
		if elems[i] > max {
			max = elems[i]
			res = []int{i}
		} else if elems[i] == max {
			res = append(res, i)
		}
	}

	return res
}

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
