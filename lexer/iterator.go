package lexer

import (
	"errors"
	"iter"

	"github.com/PlayerR9/SlParser/PlayerR9/mygo-lib/common"
	ehst "github.com/PlayerR9/SlParser/go-evals/history"
	internal "github.com/PlayerR9/SlParser/lexer/internal"
	tr "github.com/PlayerR9/mygo-lib/CustomData/tree"
)

/////////////////////////////////////////////////////////

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

// Result is the result of an active lexer.
type Result struct {
	// active is the active lexer.
	active *baseActive
}

// Tokens returns the tokens generated from the MatchResult of the active lexer.
//
// Returns:
//   - []*tr.Node: The tokens generated from the MatchResult.
func (r Result) Tokens() []*tr.Node {
	return r.active.Tokens()
}

// GetError returns the error of the active lexer.
//
// Returns:
//   - error: The error of the active lexer.
func (r Result) GetError() error {
	return r.active.GetError()
}

// Iterator is an iterator over the results of an active lexer.
type Iterator struct {
	// next is the next result of the iterator.
	next func() (ehst.Pair[internal.Event, *baseActive], error, bool)

	// stop is the stop function of the iterator.
	stop func()
}

// NewIterator creates a new iterator over the results of an active lexer.
//
// Parameters:
//   - initFn: A function that initializes and returns an active lexer.
//
// Returns:
//   - *Iterator: An iterator over the results of the active lexer.
//   - error: An error if the initFn returns an error.
func NewIterator(initFn func() (*baseActive, error)) (*Iterator, error) {
	if initFn == nil {
		return nil, common.NewErrNilParam("initFn")
	}

	itr := ehst.Evaluate(initFn)

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
	p, err, ok := itr.next()
	if !ok {
		return nil, ErrExhausted
	}

	return &Result{
		active: p.Subject,
	}, err
}
