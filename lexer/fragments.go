package lexer

import (
	"slices"

	emtch "github.com/PlayerR9/go-evals/matcher"
	"github.com/PlayerR9/mygo-lib/common"
	gch "github.com/PlayerR9/mygo-lib/runes"
)

// matchNewline matches a newline.
type matchNewline struct {
	// matched are the characters that were matched.
	matched []rune

	// has_prev is true if the previous character was the '\r' character.
	has_prev bool

	// is_done is true if the match is done.
	is_done bool
}

// Match implements the Matcher interface.
func (m *matchNewline) Close() error {
	if m == nil {
		return common.ErrNilReceiver
	} else if m.is_done {
		return nil
	} else if m.has_prev {
		r := '\r'
		return gch.NewErrAfter(true, &r, gch.NewErrNotAsExpected(true, "", nil, '\n'))
	} else {
		return gch.NewErrNotAsExpected(true, "", nil, '\r', '\n')
	}
}

// Match implements the Matcher interface.
func (m *matchNewline) Match(char rune) error {
	if m == nil {
		return common.ErrNilReceiver
	} else if m.is_done {
		return emtch.ErrMatchDone
	}

	if m.has_prev {
		if char != '\n' {
			return gch.NewErrAfter(true, &char, gch.NewErrNotAsExpected(true, "", &char, '\n'))
		}

		m.is_done = true
	} else {
		switch char {
		case '\r':
			m.has_prev = true
		case '\n':
			m.is_done = true
		default:
			return gch.NewErrNotAsExpected(true, "", &char, '\r', '\n')
		}
	}

	m.matched = append(m.matched, char)

	return nil
}

// Reset implements the Matcher interface.
func (m *matchNewline) Reset() {
	if m == nil {
		return
	}

	m.is_done = false
	m.has_prev = false

	if len(m.matched) > 0 {
		clear(m.matched)
		m.matched = nil
	}
}

// Matched implements the Matcher interface.
func (m matchNewline) Matched() []rune {
	matched := make([]rune, len(m.matched))
	copy(matched, m.matched)

	return matched
}

// Newline returns a new matcher for newline sequences.
//
// Returns:
//   - Matcher: A matcher that detects newline sequences. Never returns nil.
func Newline() emtch.Matcher[rune] {
	return &matchNewline{
		is_done:  false,
		has_prev: false,
	}
}

// Literal returns a Matcher that matches a given literal string. If the string
// has only one character, then the character is matched directly.
//
// Parameters:
//   - word: The literal string to match.
//
// Returns:
//   - Matcher: The matcher. Nil if the word is empty.
//   - error: An error if converting the string to a UTF-8 array fails.
func Literal(word string) emtch.Matcher[rune] {
	if word == "" {
		return nil
	}

	chars, err := gch.StringToUtf8(word)
	if err != nil {
		panic(err)
	}

	return emtch.Slice(chars)
}

// One returns a new matcher that matches a single character.
//
// Parameters:
//   - char: The character to match.
//
// Returns:
//   - Matcher: The matcher.
func One(char rune) emtch.Matcher[rune] {
	return emtch.Single(char)
}

// Predicate returns a new matcher that matches according to a predicate.
//
// Parameters:
//   - group_name: The name of the group.
//   - predicate: The function to match.
//
// Returns:
//   - Matcher: The matcher. Nil if the predicate is nil.
func Predicate(group_name string, predicate func(char rune) bool) emtch.Matcher[rune] {
	return emtch.Fn(group_name, predicate)
}

// Group returns a new matcher that matches a group of characters.
//
// Parameters:
//   - chars: The characters to match.
//
// Returns:
//   - Matcher: The matcher. If no characters are provided, then nil is
//     returned.
func Group(chars ...rune) emtch.Matcher[rune] {
	unique := make([]rune, 0, len(chars))

	for _, char := range chars {
		pos, ok := slices.BinarySearch(unique, char)
		if !ok {
			unique = slices.Insert(unique, pos, char)
		}
	}

	return emtch.Group("["+string(unique)+"]", unique)
}

// Many returns a Matcher that matches a given inner Matcher as many times as
// possible.
//
// Parameters:
//   - inner: The inner Matcher.
//
// Returns:
//   - Matcher: The matcher. Nil if the inner matcher is nil.
func Many(inner emtch.Matcher[rune]) emtch.Matcher[rune] {
	return emtch.Greedy(inner)
}

// Sequence returns a Matcher that matches a sequence of provided Matchers
// in the order they are given. The sequence will be processed by iterating
// through each Matcher and attempting to match the input character.
//
// Parameters:
//   - seq: A variadic number of Matcher instances. Matchers in the sequence are
//     expected to be non-nil objects.
//
// Returns:
//   - Matcher: A Matcher that represents a sequence of Matchers. Returns nil
//     if no non-nil Matchers are provided.
func Sequence(seq ...emtch.Matcher[rune]) emtch.Matcher[rune] {
	return emtch.Sequence(seq...)
}

// WithRightBound returns a Matcher that matches a given inner Matcher until a
// boundary character is encountered.
//
// Parameters:
//   - inner: The inner Matcher.
//   - bound: The boundary function. It takes a rune and returns true if it is a
//     boundary and false otherwise.
//
// Returns:
//   - Matcher: The matcher. Nil if the inner matcher is nil.
//
// If the bound function is nil, the inner Matcher is returned as is.
func WithRightBound(inner emtch.Matcher[rune], bound func(char rune) bool) emtch.Matcher[rune] {
	return emtch.WithBound(inner, bound)
}

// Greedy returns a Matcher that matches a given inner Matcher if and only if
// the next character does not satisfy the inner Matcher.
//
// Parameters:
//   - inner: The inner Matcher.
//
// Returns:
//   - Matcher: The matcher. Nil if the inner matcher is nil.
//
// It is equivalent to:
//
//	WithRightBound(inner, func(char rune) bool { err := inner.Match(char); return err != nil }).
func Greedy(inner emtch.Matcher[rune]) emtch.Matcher[rune] {
	return emtch.AutoBound(inner)
}

// Range returns a matcher that matches a group of characters between left and
// right.
//
// If left and right are equal, the returned matcher will match exactly one
// character. Otherwise, the returned matcher will match any character in the
// range [left, right].
//
// Parameters:
//   - left: The left boundary of the group.
//   - right: The right boundary of the group.
//
// Returns:
//   - Matcher: The matcher. Never returns nil.
func Range(left, right rune) emtch.Matcher[rune] {
	return emtch.Range(left, right)
}

/*
type matchOr struct {
	matchers []Matcher
	indices  []int
	eos      ernk.ErrRorSol[int]
	level    int
}

// Match implements the Matcher interface.
func (m *matchOr) Match(char rune) error {
	if m == nil {
		return common.ErrNilReceiver
	}

	if char == utf8.RuneError {

	} else {
		if len(m.indices) == 0 {
			if m.eos.HasError() {
				errs := m.eos.Errors()
				return errors.Join(errs...)
			}

			return ErrMatchDone
		}

		var top int

		for _, idx := range m.indices {
			match := m.matchers[idx]

			err := match.Match(char)
			if err == nil {
				m.indices[top] = idx
				top++
			} else if err == ErrMatchDone {
				_ = m.eos.AddSol(m.level, idx)
			} else {
				_ = m.eos.AddErr(m.level, err)
			}
		}

		m.indices = m.indices[:top:top]
		m.level++
	}

	/*
		if char == utf8.RuneError {
			if !m.match_inner {
				return ErrMatchDone
			}

			err := m.inner.Match(char)
			if err == nil {
				return fmt.Errorf("while matching with boundary: %w", ErrUnexpectedSuccess)
			} else if err != ErrMatchDone {
				return fmt.Errorf("while matching with boundary: %w", err)
			}

			m.match_inner = false

			m.matched = append(m.matched, m.inner.Matched()...)
			m.inner.Reset()

			return ErrMatchDone
		}

		if m.match_inner {
			err := m.inner.Match(char)
			if err == nil {
				return nil
			} else if err != ErrMatchDone {
				return fmt.Errorf("while matching with boundary: %w", err)
			}

			m.match_inner = false

			m.matched = append(m.matched, m.inner.Matched()...)

			m.inner.Reset()
		}

		err := m.inner.Match(char)
		if err != nil {
			return ErrMatchDone
		}

		return errors.New("boundary not satisfied")
}

// Reset implements the Matcher interface.
func (m *matchOr) Reset() {
	if m == nil {
		return
	}

	if len(m.indices) > 0 {
		clear(m.indices)
	}

	m.indices = make([]int, 0, len(m.matchers))
	for i := range m.matchers {
		m.indices = append(m.indices, i)
	}

	for _, match := range m.matchers {
		match.Reset()
	}

	m.eos.Reset()
	m.level = 0
}

// Matched implements the Matcher interface.
func (m matchOr) Matched() []rune {

	matched := make([]rune, len(m.matched))
	copy(matched, m.matched)

	return matched
}

func Or() Matcher {

}

func (m matchOr) MultiMatched() [][]rune {
	indices := m.eos.Sols()
	if len(indices) == 0 {
		return nil
	}

	table := make([][]rune, 0, len(indices))
	for _, idx := range indices {
		match := m.matchers[idx]

		matched := match.Matched()
		table = append(table, matched)
	}

	return table
}
*/
