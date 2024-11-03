package lexer

import (
	"errors"
	"slices"
	"strconv"

	"github.com/PlayerR9/mygo-lib/common"
)

var (
	// ErrNotFound occurs when the input stream is not empty but the current rune
	// is not as expected.
	//
	// Format:
	//
	// 	"not found"
	ErrNotFound error

	// ErrUnexpectedSuccess occurs when the matcher succeeds with utf8.RuneError
	// as its character match.
	//
	// Format:
	//
	// 	"matcher must not succeeded with utf8.RuneError"
	ErrUnexpectedSuccess error
)

func init() {
	ErrNotFound = errors.New("not found")

	ErrUnexpectedSuccess = errors.New("matcher must not succeeded with utf8.RuneError")
}

// ErrLexing is an error that occurs during lexing.
type ErrLexing struct {
	// At is the position in the input stream where the error occurred.
	At int

	// Reason is the reason for the error.
	Reason error
}

// Error implements the error interface.
func (e ErrLexing) Error() string {
	var msg string

	if e.Reason == nil {
		msg = "something went wrong"
	} else {
		msg = e.Reason.Error()
	}

	return msg
}

// NewErrLexing returns an error that occurs during lexing.
//
// Parameters:
//   - at: The position in the input stream where the error occurred.
//   - reason: The reason for the error.
//
// Returns:
//   - error: An error that occurs during lexing. Never returns nil.
//
// Format:
//   - exactly as "reason"
func NewErrLexing(at int, reason error) error {
	return &ErrLexing{
		At:     at,
		Reason: reason,
	}
}

// ExtractAt extracts the position in the input stream where the error occurred.
//
// Parameters:
//   - err: The error to extract the position from.
//
// Returns:
//   - int: The position in the input stream where the error occurred.
//   - bool: A boolean indicating if the position was successfully extracted.
func ExtractAt(err error) (int, bool) {
	if err == nil {
		return 0, false
	}

	e, ok := err.(*ErrLexing)
	if !ok {
		return 0, false
	}

	return e.At, true
}

// NewErrNotAsExpected is a convenience function that creates a new ErrNotAsExpected error with
// the specified kind, got value, and expected values.
//
// See common.NewErrNotAsExpected for more information.
func NewErrNotAsExpected(quote bool, kind string, got *byte, expecteds ...byte) error {
	var got_str string

	if got != nil {
		got_str = string(*got)
	}

	unique := make([]string, 0, len(expecteds))

	for _, expected := range expecteds {
		str := string(expected)

		pos, ok := slices.BinarySearch(unique, str)
		if !ok {
			unique = slices.Insert(unique, pos, str)
		}
	}

	unique = unique[:len(unique):len(unique)]

	return common.NewErrNotAsExpected(quote, kind, got_str, unique...)
}

// ErrAfter is an error that occurs after another error.
type ErrAfter struct {
	// Quote is a flag that indicates that the error should be quoted.
	Quote bool

	// Previous is the previous value.
	Previous *byte

	// Inner is the inner error.
	Inner error
}

// Error implements the error interface.
func (e ErrAfter) Error() string {
	var previous string

	if e.Previous == nil {
		previous = "at the start"
	} else if e.Quote {
		previous = strconv.Quote(string(*e.Previous))
		previous = "after " + previous
	} else {
		previous = string(*e.Previous)
		previous = "after " + previous
	}

	var reason string

	if e.Inner == nil {
		reason = "something went wrong"
	} else {
		reason = e.Inner.Error()
	}

	return previous + ": " + reason
}

// NewErrAfter creates a new ErrAfter error.
//
// Parameters:
//   - quote: A flag indicating whether the previous value should be quoted.
//   - previous: The previous value associated with the error. If not provided, "at the start" is used.
//   - inner: The inner error that occurred. If not provided, "something went wrong" is used.
//
// Returns:
//   - error: The newly created ErrAfter error. Never returns nil.
//
// Format:
//
//	"after <previous>: <inner>"
func NewErrAfter(quote bool, previous *byte, inner error) error {
	return &ErrAfter{
		Quote:    quote,
		Previous: previous,
		Inner:    inner,
	}
}
