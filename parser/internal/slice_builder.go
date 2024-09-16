package internal

// TODO: Remove this once go-commons is updated.

// SliceBuilder is a slice builder.
type SliceBuilder[T any] struct {
	// slice is the slice.
	slice []T
}

// Append appends a value to the slice. Does nothing if the receiver is nil.
//
// Parameters:
//   - v: the value to append.
func (sb *SliceBuilder[T]) Append(v T) {
	if sb == nil {
		return
	}

	sb.slice = append(sb.slice, v)
}

// Build builds the slice.
//
// Returns:
//   - []T: the slice.
func (sb SliceBuilder[T]) Build() []T {
	slice := make([]T, len(sb.slice))
	copy(slice, sb.slice)

	return slice
}

// Reset resets the slice.
func (sb *SliceBuilder[T]) Reset() {
	if sb == nil {
		return
	}

	if len(sb.slice) > 0 {
		for i := range sb.slice {
			sb.slice[i] = *new(T)
		}

		sb.slice = sb.slice[:0]
	}
}
