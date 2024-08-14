package iterator

import (
	"errors"
)

var (
	// ErrExausted is the error that is used when the iterator is done. Readers must return
	// this error as is and not wrap it as callers will test this error using ==.
	ErrExausted error
)

func init() {
	ErrExausted = errors.New("iterator is exausted")
}

// IteratorFunc is the iterator function. The error ErrExausted signals the end of the iteration.
//
// Parameters:
//   - elem: The element to apply the function to.
//
// Returns:
//   - error: An error if the function failed.
type IteratorFunc func(elem any) error

// Iterable is an interface that defines the behavior of an iterator.
type Iterable interface {
	// Apply applies the iterator function on the current element.
	// The error io.EOF signals the successful end of the iteration.
	//
	// Parameters:
	//   - fn: The function to apply. Assumed to be non-nil.
	//
	// Returns:
	//   - error: An error if the function failed.
	//
	// Successful calls to Apply will also scan the next element.
	Apply(fn IteratorFunc) error

	// Reset resets the iterator. Used for initialization.
	Reset()
}

// Iterate applies the iterator function on the iterator.
// The error ErrExausted signals the end of the iteration.
//
// Parameters:
//   - it: The iterator. Assumed to be non-nil.
//   - fn: The function to apply. Assumed to be non-nil.
//
// Returns:
//   - error: An error of type *ErrIteration if the iteration failed.
func Iterate(it any, fn IteratorFunc) error {
	if fn == nil || it == nil {
		return nil
	}

	var iter Iterable

	switch it := it.(type) {
	case Iterable:
		iter = it
	case Iterater:
		iter = it.Iterator()
	default:
		return nil
	}

	iter.Reset()

	var err error
	idx := -1

	for err == nil {
		err = iter.Apply(fn)
		idx++
	}

	if err == ErrExausted {
		return nil
	}

	return NewErrIteration(idx, err)
}

// Iterater is an interface that defines the behavior of an iterator.
type Iterater interface {
	// Iterator returns an iterator over the collection.
	//
	// Returns:
	//   - Iterable: The iterator. It should never return nil.
	Iterator() Iterable
}
