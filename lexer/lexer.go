package lexer

import (
	"fmt"

	internal "github.com/PlayerR9/SlParser/lexer/internal"
	emtch "github.com/PlayerR9/go-evals/matcher"
	"github.com/PlayerR9/mygo-lib/common"
	gch "github.com/PlayerR9/mygo-lib/runes"
)

// Lexer is a implementation of a lexer that uses a function for lexing.
type Lexer struct {
	// data is the input data.
	data []rune

	// table is the table of matchers.
	table []emtch.Matcher[rune]

	// types_ is the table of types.
	types_ []string

	// indices is the table of indices.
	indices []int
}

// Write adds the given data to the input stream.
//
// Parameters:
//   - data: The data to add to the input stream.
//
// Returns:
//   - int: The number of bytes written.
//   - error: An error if the operation fails.
//
// Always returns the length of data and a nil error, unless the receiver is nil; in which
// case it returns 0 and an error.
func (l *Lexer) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	} else if l == nil {
		return 0, common.ErrNilReceiver
	}

	chars, err := gch.BytesToUtf8(data)
	if err != nil {
		return 0, err
	}

	l.data = append(l.data, chars...)

	return len(data), nil
}

// Lex lexes the input stream, if it was set.
//
// Returns:
//   - error: An error if lexing failed.
func (l Lexer) Lex() (*Iterator, error) {
	initFn := func() (*baseActive, error) {
		return &baseActive{
			global: &l,
			pos:    0,
			err:    nil,
		}, nil
	}

	itr, _ := NewIterator(initFn)

	return itr, nil
}

// Reset resets the lexer's internal state for reuse. No-op if
// the receiver is nil.
func (l *Lexer) Reset() {
	if l == nil {
		return
	}

	if len(l.data) > 0 {
		clear(l.data)
		l.data = nil
	}
}

// Match uses the lexer's matcher to find sequences in the lexer's data.
//
// Parameters:
//   - from: The index to start matching from.
//
// Returns:
//   - []MatchResult: A slice of MatchResult containing the types of the matched sequences.
//   - error: An error if the matching process fails.
func (l Lexer) Match(from int) ([]internal.Event, error) {
	defer func() {
		for _, l := range l.table {
			l.Reset()
		}
	}()

	if len(l.table) == 0 || from == len(l.data) {
		return nil, nil
	} else if from < 0 || from > len(l.data) {
		return nil, common.NewErrBadParam("from", fmt.Sprintf("must be in [0, %d]", len(l.data)))
	}

	indices := make([]int, len(l.table))
	copy(indices, l.indices)

	pairs, err := emtch.Match(l.table, indices, l.data[from:])
	if err != nil {
		return nil, err
	}

	results := make([]internal.Event, 0, len(pairs))

	for _, p := range pairs {
		r := internal.Event{
			Type: l.types_[p.Idx],
			Data: p.Matched,
		}

		results = append(results, r)
	}

	return results, nil
}

// GetCharAt returns the character at index idx in the lexer's data or 0 if it does not exist.
//
// Parameters:
//   - idx: The index of the character to retrieve.
//
// Returns:
//   - rune: The character at index idx.
//   - bool: A boolean indicating whether the character exists or not.
func (l Lexer) GetCharAt(idx int) (rune, bool) {
	if idx < 0 || idx >= len(l.data) {
		return 0, false
	}

	return l.data[idx], true
}
