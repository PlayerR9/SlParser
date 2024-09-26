package SlParser

import (
	"bytes"
	"index/suffixarray"
	"slices"
)

// TODO: Delete this once go-commons is updated.

// ReverseSearch searches for the last occurrence of a byte in a byte slice.
//
// Parameters:
//   - data: the byte slice to search in.
//   - from: the index to start the search from. If greater than or equal to the length of the byte slice,
//     it is treated as the length of the byte slice minus 1.
//   - sep: the byte to search for.
//
// Returns:
//   - int: the index of the last occurrence of the byte in the byte slice, or -1 if not found.
//
// WARNING: As a side effect, the sep parameter is reversed. Make sure to copy it before calling this function
// if you intend to reuse it later.
func ReverseSearch(data []byte, from int, sep []byte) int {
	if len(data) == 0 || from < 0 || len(sep) == 0 {
		return -1
	}

	len_data := len(data)

	if from >= len_data {
		from = len_data
	}

	sub_data := make([]byte, from+1)
	copy(sub_data, data)
	slices.Reverse(sub_data)

	rev_sep := make([]byte, len(sep))
	copy(rev_sep, sep)

	slices.Reverse(rev_sep)

	idx := suffixarray.New(sub_data)

	offsets := idx.Lookup(rev_sep, 1)

	if len(offsets) == 0 {
		return -1
	}

	return len(sub_data) - offsets[0]
}

// ForwardSearch searches for the first occurrence of a byte in a byte slice.
//
// Parameters:
//   - data: the byte slice to search in.
//   - from: the index to start the search from. If negative, it is treated as 0.
//   - sep: the byte to search for.
//
// Returns:
//   - int: the index of the first occurrence of the byte in the byte slice, or -1 if not found.
func ForwardSearch(data []byte, from int, sep []byte) int {
	if len(data) == 0 || len(sep) == 0 || from >= len(data) {
		return -1
	}

	if from < 0 {
		from = 0
	}

	offset := bytes.Index(data[from:], sep)
	if offset == -1 {
		return -1
	}

	return from + offset
}
