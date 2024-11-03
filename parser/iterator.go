package parser

import (
	"errors"
	"iter"

	"github.com/PlayerR9/SlParser/parser/internal"
	hst "github.com/PlayerR9/go-evals/history"
	"github.com/PlayerR9/mygo-lib/common"
)

var (
	// ErrExhausted occurs when an iterator is exhausted.
	// This can be checked with the == operator.
	//
	// Format:
	//
	// 	"iterator is exhausted"
	ErrExhausted error
)

func init() {
	ErrExhausted = errors.New("iterator is exhausted")
}

// Iterator is an iterator over the results of an active parser.
type Iterator struct {
	next func() (hst.Pair[*internal.Event, *Active], error, bool)
	stop func()
}

// NewIterator creates a new iterator over the results of an active parser.
//
// Parameters:
//   - init_fn: A function that creates the initial active parser.
//
// Returns:
//   - *Iterator: An iterator over the results of the active parser.
//   - error: An error if the init_fn returns an error.
func NewIterator(init_fn func() (*Active, error)) (*Iterator, error) {
	if init_fn == nil {
		return nil, common.NewErrNilParam("init_fn")
	}

	itr := hst.Evaluate(init_fn)
	next, stop := iter.Pull2(itr)

	return &Iterator{
		next: next,
		stop: stop,
	}, nil
}

// Stop stops the iterator.
//
// After calling Stop, the iterator is no longer valid.
func (itr Iterator) Stop() {
	itr.stop()
}

// Next returns the next result of the iterator.
//
// Returns:
//   - *Result: The next result.
//   - error: An error if the iterator is exhausted or if the evaluation fails.
//
// Errors:
//   - ErrExhausted: The iterator is exhausted.
func (itr Iterator) Next() (*Result, error) {
	pair, err, ok := itr.next()
	if !ok {
		return nil, ErrExhausted
	}

	result := &Result{
		pair: pair,
	}

	return result, err
}
