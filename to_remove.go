package SlParser

import "bytes"

// TODO: Remove this once go-commons is updated.

// FindByte searches for the first occurrence of a byte in a byte slice starting from a given index.
//
// Parameters:
//   - data: the byte slice to search in.
//   - from: the index to start the search from. If negative, it is treated as 0.
//   - sep: the byte to search for.
//
// Returns:
//   - int: the index of the first occurrence of the byte in the byte slice, or -1 if not found.
func FindByte(data []byte, from int, sep byte) int {
	if len(data) == 0 || from >= len(data) {
		return -1
	}

	len_data := len(data)

	if from < 0 {
		from = 0
	}

	for i := from; i < len_data; i++ {
		if data[i] == sep {
			return i
		}
	}

	return -1
}

// FindByteReversed searches for the first occurrence of a byte in a byte slice starting from a given index in reverse order.
//
// Parameters:
//   - data: the byte slice to search in.
//   - from: the index to start the search from. If greater than or equal to the length of the byte slice,
//     it is treated as the length of the byte slice minus 1.
//   - sep: the byte to search for.
//
// Returns:
//   - int: the index of the first occurrence of the byte in the byte slice in reverse order, or -1 if not found.
func FindByteReversed(data []byte, from int, sep byte) int {
	if len(data) == 0 || from < 0 {
		return -1
	}

	len_data := len(data)

	if from >= len_data {
		from = len_data - 1
	}

	for i := from; i >= 0; i-- {
		if data[i] == sep {
			return i
		}
	}

	return -1
}

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
func ReverseSearch(data []byte, from int, sep []byte) int {
	if from < 0 || len(sep) == 0 || len(data) == 0 {
		return -1
	}

	sep_len := len(sep)

	if from+sep_len >= len(data) {
		from = len(data) - sep_len
	}

	if sep_len == 1 {
		return FindByteReversed(data, from, sep[0])
	}

	for {
		idx := FindByteReversed(data, from, sep[0])
		if idx == -1 {
			return -1
		}

		if bytes.Equal(data[idx:idx+sep_len], sep) {
			return idx
		}

		from = idx
	}
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
	if len(sep) == 0 || len(data) == 0 || from+len(sep) >= len(data) {
		return -1
	}

	sep_len := len(sep)

	if from < 0 {
		from = 0
	}

	if sep_len == 1 {
		return FindByte(data, from, sep[0])
	}

	for {
		idx := FindByte(data, from, sep[0])
		if idx == -1 {
			return -1
		}

		if bytes.Equal(data[idx:idx+sep_len], sep) {
			return idx
		}

		from = idx
	}
}
