// Package PageInterval provides a data structure for managing page intervals.
package page_interval

import (
	"slices"
	"sort"
	"strings"

	gcers "github.com/PlayerR9/go-commons/errors"
	itr "github.com/PlayerR9/go-commons/iterator"
)

// PageIterator represents an iterator that iterates over a collection of
// page ranges.
type PageIterator struct {
	// intervals is the collection of page ranges.
	intervals []PageRange

	// sub_iter is the sub iterator of the current page range.
	sub_iter itr.Iterable

	// current is the current page number of the iterator.
	current int
}

// Apply implements the iterator.Iterable interface.
//
// The argument passed to the function is the current page number of the
// iterator and it is of type int.
func (it *PageIterator) Apply(fn itr.IteratorFunc) error {
	if it.sub_iter == nil {
		if it.current >= len(it.intervals) {
			return itr.ErrExausted
		}

		it.sub_iter = it.intervals[it.current].Iterator()
	}

	for {
		err := it.sub_iter.Apply(fn)
		if err == nil {
			return nil
		} else if err != itr.ErrExausted {
			return err
		}

		it.current++

		if it.current >= len(it.intervals) {
			break
		}

		it.sub_iter = it.intervals[it.current].Iterator()
	}

	return itr.ErrExausted
}

// Reset implements the iterator.Iterable interface.
func (it *PageIterator) Reset() {
	it.current = 0
	it.sub_iter = nil
}

// PageReverseIterator represents an iterator that iterates over a
// collection of page ranges in reverse order.
type PageReverseIterator struct {
	// intervals is the collection of page ranges.
	intervals []PageRange

	// sub_iter is the sub iterator of the current page range.
	sub_iter itr.Iterable

	// current is the current page number of the iterator.
	current int
}

// Apply implements the iterator.Iterable interface.
//
// The argument passed to the function is the current page number of the
// iterator and it is of type int.
func (it *PageReverseIterator) Apply(fn itr.IteratorFunc) error {
	if it.sub_iter == nil {
		if it.current < 0 {
			return itr.ErrExausted
		}

		it.sub_iter = it.intervals[it.current].ReverseIterator()
	}

	for {
		err := it.sub_iter.Apply(fn)
		if err == nil {
			return nil
		} else if err != itr.ErrExausted {
			return err
		}

		if it.current < 0 {
			return itr.ErrExausted
		}

		it.sub_iter = it.intervals[it.current].ReverseIterator()
		it.current--
	}
}

// Reset implements the iterator.Iterable interface.
func (it *PageReverseIterator) Reset() {
	it.current = len(it.intervals) - 1
	it.sub_iter = nil
}

// PageInterval represents a collection of page intervals, where each
// interval is represented by a pair of integers.
type PageInterval struct {
	// The 'intervals' field is a slice of integer pairs, where each pair
	// represents a start and end page number of an interval.
	intervals []PageRange

	// The 'page_count' field represents the total number of pages across all
	// intervals.
	page_count int
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

// String implements the fmt.Stringer interface.
//
// Format:
//
//	<from>:<to>, <from>:<to>, ... // for example, 1:5, 6:10
func (pi PageInterval) String() string {
	values := make([]string, 0, len(pi.intervals))

	for _, interval := range pi.intervals {
		values = append(values, interval.String())
	}

	return strings.Join(values, ",")
}

// Iterator implements the iterator.Iterater interface.
func (pi PageInterval) Iterator() itr.Iterable {
	return &PageIterator{
		intervals: pi.intervals,
		current:   0,
	}
}

// NewPageInterval creates a new instance of PageInterval with
// empty intervals and a page count of 0. This is not necessary as
// var pi PageInterval is equivalent to NewPageInterval().
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
func NewPageInterval() PageInterval {
	return PageInterval{
		intervals:  make([]PageRange, 0),
		page_count: 0,
	}
}

// PageCount is a method of the PageInterval type that returns the total number
// of pages across all intervals in the PageInterval.
//
// Returns:
//   - int: The total number of pages across all intervals in the PageInterval.
func (pi PageInterval) PageCount() int {
	return pi.page_count
}

// Intervals is a method of the PageInterval type that returns the intervals
// stored in the PageInterval.
// Each interval is represented as a pair of integers, where the first integer
// is the start page number and the second integer is the end page number.
//
// Returns:
//   - []PageRange: A slice of integer pairs representing the intervals in the
//     PageInterval.
func (pi PageInterval) Intervals() []PageRange {
	return pi.intervals
}

// HasPages is a method of the PageInterval type that checks if the PageInterval
// has any pages.
//
// Returns:
//   - bool: A boolean value that is true if the PageInterval has pages, and
//     false otherwise.
func (pi PageInterval) HasPages() bool {
	return pi.page_count > 0
}

// GetFirstPage is a method of the PageInterval type that returns the first
// page number in the PageInterval.
//
// Returns:
//   - int: The first page number in the PageInterval.
//   - bool: True if the PageInterval has pages, and false otherwise.
func (pi PageInterval) GetFirstPage() (int, bool) {
	if pi.page_count <= 0 {
		return 0, false
	}

	return pi.intervals[0].first, true
}

// GetLastPage is a method of the PageInterval type that returns the last
// page number in the PageInterval.
//
// Returns:
//   - int: The last page number in the PageInterval.
//   - bool: True if the PageInterval has pages, and false otherwise.
func (pi PageInterval) GetLastPage() (int, bool) {
	if pi.page_count <= 0 {
		return 0, false
	}

	return pi.intervals[len(pi.intervals)-1].second, true
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

	criteria_page_gte := func(i int) bool {
		return pi.intervals[i].first >= page
	}

	if len(pi.intervals) == 0 {
		pi.intervals = append(pi.intervals, NewPageRange(page, page))
	} else {
		insert_pos := sort.Search(len(pi.intervals), criteria_page_gte)

		if insert_pos > 0 && pi.intervals[insert_pos-1].second >= page-1 {
			insert_pos--

			var max int

			if page > pi.intervals[insert_pos].second {
				max = page
			} else {
				max = pi.intervals[insert_pos].second
			}

			pi.intervals[insert_pos].second = max
		} else if insert_pos < len(pi.intervals) && pi.intervals[insert_pos].first <= page+1 {
			var min int

			if page < pi.intervals[insert_pos].first {
				min = page
			} else {
				min = pi.intervals[insert_pos].first
			}

			pi.intervals[insert_pos].first = min
		} else {
			pi.intervals = append(pi.intervals[:insert_pos],
				append([]PageRange{NewPageRange(page, page)}, pi.intervals[insert_pos:]...)...,
			)
		}
	}

	pi.page_count++
	pi.reduce()

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

	index := pi.find_page_interval(page)
	if index == -1 {
		return
	}

	if pi.intervals[index].first == pi.intervals[index].second {
		pi.intervals = append(pi.intervals[:index], pi.intervals[index+1:]...)
	} else if pi.intervals[index].first == page {
		pi.intervals[index].first++
	} else if pi.intervals[index].second == page {
		pi.intervals[index].second--
	} else {
		new_intervals := make([]PageRange, len(pi.intervals)+1)

		// Copy the intervals before the split
		copy(new_intervals, pi.intervals[:index+1])

		// Modify the interval at the split index
		new_intervals[index] = NewPageRange(pi.intervals[index].first, page-1)

		// Add the new interval
		new_intervals[index+1] = NewPageRange(page+1, pi.intervals[index].second)

		// Copy the intervals after the split
		copy(new_intervals[index+2:], pi.intervals[index+1:])

		pi.intervals = new_intervals
	}

	pi.page_count--

	pi.reduce()
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
func (pi PageInterval) HasPage(page int) bool {
	return pi.find_page_interval(page) != -1
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
// Returns:
//   - itr.Iterable: An iterator for iterating over the intervals in the
//     PageInterval in reverse order. Never returns nil.
func (pi PageInterval) ReverseIterator() itr.Iterable {
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
func (pi *PageInterval) reduce() {
	if len(pi.intervals) < 2 {
		return
	}

	criteria_sort := func(i, j int) bool {
		return pi.intervals[i].first < pi.intervals[j].first
	}

	sort.SliceStable(pi.intervals, criteria_sort)

	merged_intervals := make([]PageRange, 0, len(pi.intervals))
	current_interval := pi.intervals[0]

	for i := 1; i < len(pi.intervals); i++ {
		nextInterval := pi.intervals[i]
		if current_interval.second >= nextInterval.first-1 {
			if nextInterval.second > current_interval.second {
				current_interval.second = nextInterval.second
			}
		} else {
			merged_intervals = append(merged_intervals, current_interval)
			current_interval = nextInterval
		}
	}

	merged_intervals = append(merged_intervals, current_interval)
	pi.intervals = merged_intervals
}

// find_page_interval searches for the interval that contains the given page
// number in the PageInterval.
//
// Parameters:
//   - page: The page number to search for in the PageInterval.
//
// Returns:
//   - int: The index of the interval in the intervals slice if found, otherwise -1.
func (pi PageInterval) find_page_interval(page int) int {
	if page < 1 || pi.page_count == 0 {
		return -1
	}

	is_page_between := func(pr PageRange) bool {
		return pr.first <= page && page <= pr.second
	}

	return slices.IndexFunc(pi.intervals, is_page_between)
}
