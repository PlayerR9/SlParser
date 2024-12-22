package matcher

import (
	"fmt"
	"strconv"

	"github.com/PlayerR9/SlParser/mygo-lib/common"
)

// matchSingle is a matcher that matches a single character.
type matchSingle struct {
	// target is the character to match. (Never nil.)
	target *rune

	// is_done is a flag that indicates whether the matcher is done.
	is_done bool
}

// Match implements Matcher.
func (ms *matchSingle) Match(char rune) error {
	if ms == nil {
		return common.ErrNilReceiver
	}

	if ms.is_done {
		return ErrMatchDone
	}

	if char == *ms.target {
		ms.is_done = true

		return nil
	}

	err := fmt.Errorf("want %s, got %s", strconv.QuoteRune(*ms.target), strconv.QuoteRune(char))
	return err
}

// Close implements Matcher.
func (ms *matchSingle) Close() error {
	if ms == nil {
		return common.ErrNilReceiver
	}

	if !ms.is_done {
		err := fmt.Errorf("want %s, got nothing", strconv.QuoteRune(*ms.target))
		return err
	}

	return nil
}

// Matched implements Matcher.
func (ms *matchSingle) Matched() []rune {
	if ms == nil {
		return nil
	}

	if !ms.is_done {
		return nil
	}

	matched := []rune{*ms.target}

	return matched
}

// Single creates a matcher that matches a single specified character.
//
// The function initializes a new matchSingle instance with the provided target
// character and returns it as a Matcher. The matcher will successfully match
// if the input character is the same as the target character.
//
// Parameters:
//   - target: The rune to be matched.
//
// Returns:
//   - Matcher: A matcher that matches the specified target character. Never nil.
func Single(target rune) Matcher {
	ms := &matchSingle{
		target: &target,
	}

	return ms
}

// matchAny is a matcher that matches any character.
type matchAny struct {
	// matched is the rune that has been matched.
	matched rune

	// is_done is a flag that indicates whether the matcher is done.
	is_done bool
}

// Match implements Matcher.
func (ma *matchAny) Match(char rune) error {
	if ma == nil {
		return common.ErrNilReceiver
	}

	if ma.is_done {
		return ErrMatchDone
	}

	ma.matched = char
	ma.is_done = true

	return nil
}

// Close implements Matcher.
func (ma *matchAny) Close() error {
	if ma == nil {
		return common.ErrNilReceiver
	}

	if !ma.is_done {
		err := fmt.Errorf("want something, got nothing")
		return err
	}

	return nil
}

// Matched implements Matcher.
func (ma *matchAny) Matched() []rune {
	if ma == nil {
		return nil
	}

	if !ma.is_done {
		return nil
	}

	matched := []rune{ma.matched}

	return matched
}

// Any creates a matcher that matches any character.
//
// The function initializes a new matchAny instance and returns it as a Matcher.
//
// Returns:
//   - Matcher: A matcher that matches any character. Never nil.
func Any() Matcher {
	ma := &matchAny{}

	return ma
}

// matchSlice is a matcher that matches a slice of characters.
type matchSlice struct {
	// target is the slice of characters to match. (Never nil.)
	target *[]rune

	// targetLen is the length of the target slice. (Never nil.)
	targetLen *uint

	// idx is the index of the next character to match.
	idx uint
}

// Match implements Matcher.
func (ms *matchSlice) Match(char rune) error {
	if ms == nil {
		return common.ErrNilReceiver
	}

	if ms.idx >= *ms.targetLen {
		return ErrMatchDone
	}

	target := (*ms.target)[ms.idx]

	if char == target {
		ms.idx++

		return nil
	}

	err := fmt.Errorf("want %s, got %s", strconv.QuoteRune(target), strconv.QuoteRune(char))
	return err
}

// Close implements Matcher.
func (ms *matchSlice) Close() error {
	if ms == nil {
		return common.ErrNilReceiver
	}

	if ms.idx >= *ms.targetLen {
		return nil
	}

	target := (*ms.target)[ms.idx]

	err := fmt.Errorf("want %s, got nothing", strconv.QuoteRune(target))
	return err
}

// Matched implements Matcher.
func (ms *matchSlice) Matched() []rune {
	if ms == nil {
		return nil
	}

	if ms.idx < *ms.targetLen {
		return nil
	}

	matched := make([]rune, *ms.targetLen)
	copy(matched, *ms.target)

	return matched
}

// Slice creates a matcher that matches a slice of characters.
//
// The function initializes a new matchSlice instance with the provided target
// slice and returns it as a Matcher. The matcher will successfully match if the
// input character is the same as the next character in the target slice.
//
// Parameters:
//   - target: The slice of characters to be matched. (Never nil.)
//
// Returns:
//   - Matcher: A matcher that matches the specified target slice.
//
// Returns nil if the target slice is empty.
func Slice(target []rune) Matcher {
	sliceLen := uint(len(target))

	switch sliceLen {
	case 0:
		return nil
	case 1:
		m := &matchSingle{
			target: &target[0],
		}

		return m
	default:
		slice := make([]rune, len(target))
		copy(slice, target)

		ms := &matchSlice{
			target:    &slice,
			targetLen: &sliceLen,
		}

		return ms
	}
}
