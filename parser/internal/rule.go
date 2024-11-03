package internal

import (
	"iter"

	gslc "github.com/PlayerR9/mygo-lib/slices"
)

type Rule struct {
	lhs  string
	rhss []string
}

func NewRule(lhs string, rhss []string) *Rule {
	return &Rule{
		lhs:  lhs,
		rhss: rhss,
	}
}

func (r Rule) Size() int {
	return len(r.rhss)
}

func (r Rule) RhsAt(idx int) (string, bool) {
	if idx < 0 || idx >= len(r.rhss) {
		return "", false
	}

	return r.rhss[idx], true
}

func (r Rule) Rhs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, rhs := range r.rhss {
			if !yield(rhs) {
				break
			}
		}
	}
}

func (r Rule) BackwardRhs() iter.Seq[string] {
	return func(yield func(string) bool) {
		for i := len(r.rhss) - 1; i >= 0; i-- {
			if !yield(r.rhss[i]) {
				break
			}
		}
	}
}

func (r Rule) Lhs() string {
	return r.lhs
}

func (r Rule) Symbols() []string {
	symbols := make([]string, len(r.rhss))
	copy(symbols, r.rhss)

	_ = gslc.Uniquefy(&symbols)

	_, _ = gslc.MayInsert(&symbols, r.lhs)

	return symbols
}

func (r Rule) IndicesOf(target string) []int {
	var count int

	for _, rhs := range r.rhss {
		if rhs == target {
			count++
		}
	}

	if count == 0 {
		return nil
	}

	slice := make([]int, 0, count)

	for idx, rhs := range r.rhss {
		if rhs == target {
			slice = append(slice, idx)
		}
	}

	return slice
}
