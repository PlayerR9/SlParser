package slices

import "github.com/PlayerR9/SlParser/mygo-lib/common"

// Builder is a builder for slices. It is only efficent for making many slices one after the other.
//
// An empty builder can be created with the `var b Builder[T]` syntax or with the
// `new(Builder[T])` constructor.
type Builder[E any] struct {
	// slice is the underlying slice being built.
	slice []E
}

// Reset implements common.Resetter.
func (b *Builder[E]) Reset() error {
	if b == nil {
		return common.ErrNilReceiver
	}

	if len(b.slice) == 0 {
		return nil
	}

	clear(b.slice)
	b.slice = nil

	return nil
}

// Append appends an element to the slice being built.
//
// Parameters:
//   - elem: The element to append.
//
// Returns:
//   - error: An error if the element could not be appended.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
func (b *Builder[E]) Append(elem E) error {
	if b == nil {
		return common.ErrNilReceiver
	}

	b.slice = append(b.slice, elem)

	return nil
}

// Build builds the slice being built.
//
// Returns:
//   - []E: The slice being built. Nil if no elements were appended.
func (b Builder[E]) Build() []E {
	if len(b.slice) == 0 {
		return nil
	}

	slice := make([]E, len(b.slice), len(b.slice))
	copy(slice, b.slice)

	return slice
}
