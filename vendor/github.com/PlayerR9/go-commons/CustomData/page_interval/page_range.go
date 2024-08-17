package page_interval

import (
	"iter"
	"strconv"
	"strings"
)

// PageRange represents a pair of integers that represent the start and end
// page numbers of an interval.
// The first integer is the start page number and the second integer is the
// end page number of the interval. (both inclusive)
//
// For instance, the PageRange [1, 5] represents the interval from page 1 to
// page 5.
type PageRange struct {
	// first is the start page number of the interval.
	first int

	// second is the end page number of the interval.
	second int
}

/* // FString returns the string representation of the PageRange using the given
// traversor and options.
//
// Parameters:
//   - trav: The traversor to use for printing.
//   - ws: The whitespace to use between the elements. By default, it is a single space.
//   - sep: The separator to use between the key and value. By default, it is a colon.
//
// Behaviors:
//   - If sep is an empty string, it is set to a colon.
//   - ws can be empty. The default value is a single space.
//   - The default call for AString is: AString(trav, " ", "").
//   - If trav is empty, the function does nothing.
func (pr *PageRange) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	if trav == nil {
		return nil
	}

	settings := &settingsTable{
		ws:  " ",
		sep: ":",
	}

	for _, opt := range opts {
		opt(settings)
	}

	var err error

	if pr.First == pr.Second {
		err = trav.AppendString(strconv.Itoa(pr.First))
	} else {
		err = trav.AppendJoinedString(settings.ws, strconv.Itoa(pr.First), settings.sep, strconv.Itoa(pr.Second))
	}
	if err != nil {
		return err
	}

	trav.AcceptWord()

	return nil
} */

// String implements the fmt.Stringer interface.
//
// Format:
//
//	<from>:<to> // for example, 1:5
//	<from> // for example, 1 if first is equal to second
func (pr PageRange) String() string {
	if pr.first == pr.second {
		return strconv.Itoa(pr.first)
	}

	var builder strings.Builder

	builder.WriteString(strconv.Itoa(pr.first))
	builder.WriteRune(':')
	builder.WriteString(strconv.Itoa(pr.second))

	return builder.String()
}

// NewPageRange creates a new instance of PageRange with the given start and
// end page numbers.
//
// Parameters:
//   - start: The start page number of the interval.
//   - end: The end page number of the interval.
//
// Returns:
//   - PageRange: The new PageRange.
//
// If start is greater than end, the start and end are swapped. Negative numbers
// are treated as positive.
func NewPageRange(start, end int) PageRange {
	if start < 1 {
		start *= -1
	}

	if end < 1 {
		end *= -1
	}

	if start > end {
		start, end = end, start
	}

	return PageRange{start, end}
}

// All returns an iterator that iterates over the pages in the interval from
// the first page number to the second page number.
//
// Returns:
//   - iter.Seq[int]: The iterator. Never returns nil.
func (pr PageRange) All() iter.Seq[int] {
	return func(yield func(page int) bool) {
		for i := pr.first; i <= pr.second; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// Backward returns an iterator that iterates over the pages in the interval
// from the second page number to the first page number.
//
// Returns:
//   - iter.Seq[int]: The iterator. Never returns nil.
func (pr PageRange) Backward() iter.Seq[int] {
	return func(yield func(page int) bool) {
		for i := pr.second; i >= pr.first; i-- {
			if !yield(i) {
				return
			}
		}
	}
}
