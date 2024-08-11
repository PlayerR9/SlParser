// Package PageInterval provides a data structure for managing page intervals.
package PageInterval

import (
	"io"
	"sort"
	"strings"

	gcers "github.com/PlayerR9/go-commons/errors"
)

// PageInterval represents a collection of page intervals, where each
// interval is represented by a pair of integers.
type PageInterval struct {
	// The 'intervals' field is a slice of integer pairs, where each pair
	// represents a start and end page number of an interval.
	intervals []*PageRange

	// The 'pageCount' field represents the total number of pages across all
	// intervals.
	pageCount int
}

/* // FString is a method of the PageInterval type that returns the formatted
// string representation of the PageInterval using the given traversor and
// options.
//
// Parameters:
//   - trav: The traversor to use for printing.
//   - opts: The options to use for formatting the string.
//
// Returns:
//   - error: An error if the traversor encounters an error while printing.
//
// Options:
//   - WithWS: Sets the whitespace to use between the intervals. By default, it
//     is a single space.
//   - WithSep: Sets the separator to use between the start and end page numbers
//     of an interval. By default, it is a colon. If the separator is an empty
//     string, it is set to a colon instead.
//
// Behaviors:
//   - If the traversor is empty, the function does nothing.
func (pi *PageInterval) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
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

	for i, interval := range pi.intervals {
		var str string

		if interval.First == interval.Second {
			str = strconv.Itoa(interval.First)
		} else {
			str = strings.Join([]string{
				strconv.Itoa(interval.First),
				settings.sep,
				strconv.Itoa(interval.Second),
			}, settings.ws)
		}

		err = trav.AppendString(str)
		if err != nil {
			return luint.NewErrAt(i+1, "interval", err)
		}

		trav.AcceptWord()
	}

	return nil
} */

// String is a method of the PageInterval type that returns a string
// representation of the PageInterval.
// Each interval is represented as "start : end" separated by a comma.
//
// Returns:
//   - string: A formatted string representation of the PageInterval.
func (pi *PageInterval) String() string {
	values := make([]string, 0, len(pi.intervals))

	for _, interval := range pi.intervals {
		values = append(values, interval.String())
	}

	return strings.Join(values, ",")
}

// Iterator is a method of the PageInterval type that returns an iterator for
// iterating over the pages in the PageInterval.
//
// Panics if an error occurs while creating the iterator.
//
// Returns:
//   - *PageIterator: An iterator for iterating over the pages in the PageInterval. Never returns nil.
func (pi *PageInterval) Iterator() *PageIterator {
	return &PageIterator{
		intervals: pi.intervals,
		current:   0,
	}
}

// PageCount is a method of the PageInterval type that returns the total number
// of pages across all intervals in the PageInterval.
//
// Returns:
//   - pageCount: The total number of pages across all intervals in the PageInterval.
func (pi *PageInterval) PageCount() int {
	return pi.pageCount
}

// Intervals is a method of the PageInterval type that returns the intervals
// stored in the PageInterval.
// Each interval is represented as a pair of integers, where the first integer
// is the start page number and the second integer is the end page number.
//
// Returns:
//   - []*PageRange: A slice of integer pairs representing the intervals in the
//     PageInterval.
func (pi *PageInterval) Intervals() []*PageRange {
	return pi.intervals
}

// NewPageInterval creates a new instance of PageInterval with
// empty intervals and a page count of 0.
//
// Returns:
//   - PageInterval: The new PageInterval.
//
// The PageInterval ensures non-overlapping, non-duplicate intervals and
// reduces the amount of intervals by merging two consecutive intervals
// into one.
//
// Example:
//
//	pi := NewPageInterval()
//	pi.AddPagesBetween(1, 5)
//	pi.AddPagesBetween(10, 15)
//
//	fmt.Println(pi.Intervals()) // Output: [[1 5] [10 15]]
//	fmt.Println(pi.PageCount()) // Output: 11
func NewPageInterval() *PageInterval {
	return &PageInterval{
		intervals: make([]*PageRange, 0),
		pageCount: 0,
	}
}

// HasPages is a method of the PageInterval type that checks if the PageInterval
// has any pages.
//
// Returns:
//   - bool: A boolean value that is true if the PageInterval has pages, and
//     false otherwise.
func (pi *PageInterval) HasPages() bool {
	return pi.pageCount > 0
}

// GetFirstPage is a method of the PageInterval type that returns the first
// page number in the PageInterval.
//
// Returns:
//   - int: The first page number in the PageInterval.
//   - error: An error of type *uc.ErrNoPagesInInterval if no pages have been set.
func (pi *PageInterval) GetFirstPage() (int, error) {
	if pi.pageCount <= 0 {
		return 0, NewErrNoPagesInInterval()
	}

	return pi.intervals[0].First, nil
}

// GetLastPage is a method of the PageInterval type that returns the last
// page number in the PageInterval.
//
// Returns:
//   - int: The last page number in the PageInterval.
//   - error: An error of type *uc.ErrNoPagesInInterval if no pages have been set.
func (pi *PageInterval) GetLastPage() (int, error) {
	if pi.pageCount <= 0 {
		return 0, NewErrNoPagesInInterval()
	}

	return pi.intervals[len(pi.intervals)-1].Second, nil
}

// AddPage is a method of the PageInterval type that adds a page to the
// PageInterval, maintaining the non-overlapping, non-duplicate intervals.
//
// Parameters:
//   - page: The page number to add to the PageInterval.
//
// Returns:
//   - error: An error of type *uc.ErrInvalidParameter if the page number is less than 1.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: []PageRange{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	pi.AddPage(6)
//	fmt.Println(pi.intervals) // Output: [[1 6] [10 15]]
//	fmt.Println(pi.pageCount) // Output: 12
func (pi *PageInterval) AddPage(page int) error {
	if page < 1 {
		return gcers.NewErrInvalidParameter(
			"page",
			gcers.NewErrGT(0),
		)
	}

	criteriaPageGTE := func(i int) bool {
		return pi.intervals[i].First >= page
	}

	if len(pi.intervals) == 0 {
		pi.intervals = append(pi.intervals, newPageRange(page, page))
	} else {
		insertPos := sort.Search(len(pi.intervals), criteriaPageGTE)

		if insertPos > 0 && pi.intervals[insertPos-1].Second >= page-1 {
			insertPos--

			var max int

			if page > pi.intervals[insertPos].Second {
				max = page
			} else {
				max = pi.intervals[insertPos].Second
			}

			pi.intervals[insertPos].Second = max
		} else if insertPos < len(pi.intervals) && pi.intervals[insertPos].First <= page+1 {
			var min int

			if page < pi.intervals[insertPos].First {
				min = page
			} else {
				min = pi.intervals[insertPos].First
			}

			pi.intervals[insertPos].First = min
		} else {
			pi.intervals = append(pi.intervals[:insertPos],
				append([]*PageRange{newPageRange(page, page)}, pi.intervals[insertPos:]...)...,
			)
		}
	}

	pi.pageCount++
	reduce(pi)

	return nil
}

// RemovePage is a method of the PageInterval type that removes the specified
// page from the PageInterval.
// No changes are made if the page number is less than 1 or not found in the
// PageInterval.
//
// Parameters:
//   - page: The page number to remove from the PageInterval.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: []PageRange{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	pi.RemovePage(5)
//	fmt.Println(pi.intervals) // Output: [[1 4] [10 15]]
//	fmt.Println(pi.pageCount) // Output: 10
func (pi *PageInterval) RemovePage(page int) {
	if page < 1 {
		return // No-op
	}

	index := pi.findPageInterval(page)
	if index == -1 {
		return
	}

	if pi.intervals[index].First == pi.intervals[index].Second {
		pi.intervals = append(pi.intervals[:index], pi.intervals[index+1:]...)
	} else if pi.intervals[index].First == page {
		pi.intervals[index].First++
	} else if pi.intervals[index].Second == page {
		pi.intervals[index].Second--
	} else {
		newIntervals := make([]*PageRange, len(pi.intervals)+1)

		// Copy the intervals before the split
		copy(newIntervals, pi.intervals[:index+1])

		// Modify the interval at the split index
		newIntervals[index] = newPageRange(pi.intervals[index].First, page-1)

		// Add the new interval
		newIntervals[index+1] = newPageRange(page+1, pi.intervals[index].Second)

		// Copy the intervals after the split
		copy(newIntervals[index+2:], pi.intervals[index+1:])

		pi.intervals = newIntervals
	}

	pi.pageCount--

	reduce(pi)
}

// HasPage is a method of the PageInterval type that checks if the given page
// exists in the PageInterval.
//
// Parameters:
//   - page: The page number to check for in the PageInterval.
//
// Returns:
//   - bool: A boolean value that is true if the page exists in the PageInterval,
//     and false otherwise.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: []PageRange{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	hasPage := pi.HasPage(3)
//	fmt.Println(hasPage) // Output: true
func (pi *PageInterval) HasPage(page int) bool {
	return pi.findPageInterval(page) != -1
}

// AddPagesBetween is a method of the PageInterval type that adds pages between
// the first and last page numbers to the PageInterval.
//
// However, if the first page number is less than 1, it is set to 1 to remove
// invalid pages, same goes for the last page number.
// Finally, if the last page number is less than the first page number, the
// values are swapped.
//
// Parameters:
//   - first: The first page number to add to the PageInterval.
//   - last: The last page number to add to the PageInterval.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: []PageRange{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	pi.AddPagesBetween(6, 9)
//	fmt.Println(pi.intervals) // Output: [[1 15]]
//	fmt.Println(pi.pageCount) // Output: 15
func (pi *PageInterval) AddPagesBetween(first, last int) {
	if first < 1 {
		first = 1 // remove invalid pages
	}

	if last < 1 {
		last = 1 // remove invalid pages
	}

	if last < first {
		last, first = first, last // swap values
	}

	for i := first; i <= last; i++ {
		pi.AddPage(i)
	}
}

// RemovePagesBetween is a method of the PageInterval type that removes pages
// between the specified first and last page numbers from the PageInterval.
//
// However, if the first page number is less than 1, it is set to 1 to remove
// invalid pages, same goes for the last page number.
// Finally, if the last page number is less than the first page number, the
// values are swapped.
//
// Parameters:
//   - first, last: The first and last page numbers to remove from the PageInterval,
//     respectively.
//
// Example:
//
//	pi := PageInterval{
//	    intervals: []PageRange{{1, 5}, {10, 15}},
//	    pageCount: 11,
//	}
//
//	pi.RemovePagesBetween(3, 4)
//	fmt.Println(pi.intervals) // Output: [[1 2] [5 5] [10 15]]
//	fmt.Println(pi.pageCount) // Output: 9
func (pi *PageInterval) RemovePagesBetween(first, last int) {
	if first < 1 {
		first = 1 // remove invalid pages
	}

	if last < 1 {
		last = 1 // remove invalid pages
	}

	if last < first {
		last, first = first, last // swap values
	}

	for i := first; i <= last; i++ {
		pi.RemovePage(i)
	}
}

// ReverseIterator is a method of the PageInterval type that returns a
// PageIntervalReverseIterator for iterating over the intervals in the
// PageInterval in reverse order.
//
// Panics if an error occurs while creating the iterator.
//
// Returns:
//   - uc.Iterater[int]: An iterator for iterating over the intervals in the
//     PageInterval in reverse order.
func (pi *PageInterval) ReverseIterator() *PageReverseIterator {
	return &PageReverseIterator{
		intervals: pi.intervals,
		current:   len(pi.intervals) - 1,
	}
}

// reduce merges overlapping intervals in the PageInterval.
// It sorts the intervals based on the start value and then merges any
// overlapping intervals.
// The merged intervals are stored in the intervals field of the PageInterval.
// If the PageInterval contains less than two intervals, no operation is
// performed.
//
// Parameters:
//   - pi: A pointer to the PageInterval to reduce.
func reduce(pi *PageInterval) {
	if len(pi.intervals) < 2 {
		return
	}

	criteriaSort := func(i, j int) bool {
		return pi.intervals[i].First < pi.intervals[j].First
	}

	sort.Slice(pi.intervals, criteriaSort)

	mergedIntervals := make([]*PageRange, 0, len(pi.intervals))
	currentInterval := pi.intervals[0]

	for i := 1; i < len(pi.intervals); i++ {
		nextInterval := pi.intervals[i]
		if currentInterval.Second >= nextInterval.First-1 {
			if nextInterval.Second > currentInterval.Second {
				currentInterval.Second = nextInterval.Second
			}
		} else {
			mergedIntervals = append(mergedIntervals, currentInterval)
			currentInterval = nextInterval
		}
	}

	mergedIntervals = append(mergedIntervals, currentInterval)
	pi.intervals = mergedIntervals
}

type PageIterator struct {
	intervals []*PageRange

	sub_iter *PageRangeIterator

	current int
}

func (it *PageIterator) Consume() (int, error) {
	if it.sub_iter == nil {
		if it.current >= len(it.intervals) {
			return -1, io.EOF
		}

		it.sub_iter = it.intervals[it.current].Iterator()
	}

	for {
		value, err := it.sub_iter.Consume()
		if err == nil {
			return value, nil
		}

		if it.current >= len(it.intervals) {
			return -1, io.EOF
		}

		it.sub_iter = it.intervals[it.current].Iterator()
		it.current++
	}
}

func (it *PageIterator) Reset() {
	it.current = 0
	it.sub_iter = nil
}

type PageReverseIterator struct {
	intervals []*PageRange

	sub_iter *PageRangeReverseIterator

	current int
}

func (it *PageReverseIterator) Consume() (int, error) {
	if it.sub_iter == nil {
		if it.current < 0 {
			return -1, io.EOF
		}

		it.sub_iter = it.intervals[it.current].ReverseIterator()
	}

	for {
		value, err := it.sub_iter.Consume()
		if err == nil {
			return value, nil
		}

		if it.current < 0 {
			return -1, io.EOF
		}

		it.sub_iter = it.intervals[it.current].ReverseIterator()
		it.current--
	}
}

func (it *PageReverseIterator) Reset() {
	it.current = len(it.intervals) - 1
	it.sub_iter = nil
}
