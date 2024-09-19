package parser

import "io"

// TODO: Remove this once go-commons is updated.

// SeenMap is a map that keeps track of seen values.
type SeenMap[T comparable] struct {
	// table is the map that keeps track of seen values.
	table map[T]struct{}
}

// NewSeenMap creates a new SeenMap.
//
// Returns:
//   - *SeenMap[T]: The new SeenMap. Never returns nil.
func NewSeenMap[T comparable]() *SeenMap[T] {
	return &SeenMap[T]{
		table: make(map[T]struct{}),
	}
}

// SetSeen sets the value as seen.
//
// Parameters:
//   - v: The value to set as seen.
//
// Returns:
//   - bool: True if the value was set as seen, false otherwise.
func (sm *SeenMap[T]) SetSeen(v T) bool {
	if sm == nil {
		return false
	}

	_, ok := sm.table[v]
	if ok {
		return false
	}

	sm.table[v] = struct{}{}

	return true
}

// IsSeen checks whether the value is seen.
//
// Parameters:
//   - v: The value to check.
//
// Returns:
//   - bool: True if the value is seen, false otherwise.
func (sm SeenMap[T]) IsSeen(v T) bool {
	_, ok := sm.table[v]
	return ok
}

// FilterNotSeen returns the elements that are not seen.
//
// Parameters:
//   - elems: The elements to filter.
//
// Returns:
//   - []T: The elements that are not seen.
func (sm SeenMap[T]) FilterNotSeen(elems []T) []T {
	var filtered []T

	for _, e := range elems {
		_, ok := sm.table[e]
		if !ok {
			filtered = append(filtered, e)
		}
	}

	return filtered
}

// FilterSeen returns the elements that are seen.
//
// Parameters:
//   - elems: The elements to filter.
//
// Returns:
//   - []T: The elements that are seen.
func (sm SeenMap[T]) FilterSeen(elems []T) []T {
	var filtered []T

	for _, e := range elems {
		_, ok := sm.table[e]
		if ok {
			filtered = append(filtered, e)
		}
	}

	return filtered
}

// Reset resets the SeenMap.
func (sm *SeenMap[T]) Reset() {
	if sm == nil {
		return
	}

	if len(sm.table) > 0 {
		for k := range sm.table {
			delete(sm.table, k)
		}

		sm.table = make(map[T]struct{})
	}
}

// Size gets the size of the SeenMap.
//
// Returns:
//   - int: The size of the SeenMap. Never returns a negative number.
func (sm SeenMap[T]) Size() int {
	return len(sm.table)
}

// IsEmpty checks whether the SeenMap is empty.
//
// Returns:
//   - bool: True if the SeenMap is empty, false otherwise.
func (sm SeenMap[T]) IsEmpty() bool {
	return len(sm.table) == 0
}

// PredicateFilter is a type that defines a slice filter function.
//
// Parameters:
//   - T: The type of the elements in the slice.
//
// Returns:
//   - bool: True if the element satisfies the filter function, otherwise false.
type PredicateFilter[T any] func(T) bool

// SliceFilter is the same as SliceFilter, but without any side-effects.
func SliceFilter[T any](slice []T, filter PredicateFilter[T]) []T {
	if len(slice) == 0 {
		return nil
	} else if filter == nil {
		return slice
	}

	var top int

	for i := 0; i < len(slice); i++ {
		elem := slice[i]

		ok := filter(elem)
		if ok {
			slice[top] = elem
			top++
		}
	}

	return slice[:top:top]
}

// PureSliceFilter is the same as SliceFilter, but without any side-effects.
func PureSliceFilter[T any](slice []T, filter PredicateFilter[T]) []T {
	if len(slice) == 0 {
		return nil
	} else if filter == nil {
		return slice
	}

	slice_copy := make([]T, 0, len(slice))

	for i := 0; i < len(slice); i++ {
		ok := filter(slice[i])
		if ok {
			slice_copy = append(slice_copy, slice[i])
		}
	}

	return slice_copy[:len(slice_copy):len(slice_copy)]
}

// Write writes the string to the writer.
//
// Parameters:
//   - w: the writer.
//   - str: the string.
//
// Returns:
//   - error: if an error occurred.
//
// Errors:
//   - io.ErrShortWrite if the string is not fully written or the writer is nil.
//   - any other error returned by the writer.
func Write(w io.Writer, str string) error {
	if str == "" {
		return nil
	} else if w == nil {
		return io.ErrShortWrite
	}

	data := []byte(str)

	n, err := w.Write(data)
	if err != nil {
		return err
	} else if n != len(data) {
		return io.ErrShortWrite
	}

	return nil
}
