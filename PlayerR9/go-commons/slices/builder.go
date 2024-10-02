package slices

// Builder is a slice builder.
type Builder[T any] struct {
	// elems is the slice of elements.
	elems []T
}

// NewBuilder creates a new builder. This is not necessary since
// it is equivalent to doing the following:
//
//	var builder Builder[T]
//
// Returns:
//   - *Builder: The new builder. Never returns nil.
func NewBuilder[T any]() *Builder[T] {
	return &Builder[T]{
		elems: make([]T, 0),
	}
}

// Prepend adds the given elements to the beginning of the slice.
//
// Parameters:
//   - elems: The elements to add to the beginning of the slice.
//
// Does nothing if the builder is nil or the slice is empty.
func (b *Builder[T]) Prepend(elems ...T) {
	if b == nil || len(elems) == 0 {
		return
	}

	b.elems = append(elems, b.elems...)
}

// Append adds the given elements to the end of the slice.
//
// Parameters:
//   - elems: The elements to add to the end of the slice.
//
// Does nothing if the builder is nil or the slice is empty.
func (b *Builder[T]) Append(elems ...T) {
	if b == nil || len(elems) == 0 {
		return
	}

	b.elems = append(b.elems, elems...)
}

// Build creates a new slice from the elements added to the builder.
//
// Returns:
//   - []T: The new slice. Nil if there are no elements in the builder.
func (b Builder[T]) Build() []T {
	if len(b.elems) == 0 {
		return nil
	}

	slice := make([]T, len(b.elems))
	copy(slice, b.elems)

	return slice
}

// Reset resets the builder for reuse.
func (b *Builder[T]) Reset() {
	if b == nil {
		return
	}

	if len(b.elems) > 0 {
		b.elems = b.elems[:0]
	}
}
