package internal

import "github.com/PlayerR9/SlParser/mygo-lib/common"

// Predicate is a function that returns true if the element is valid.
//
// Parameters:
//   - e: The element to check.
//
// Returns:
//   - bool: True if the element is valid.
type Predicate[E any] func(e E) bool

// Filter filters a slice of E's according to a predicate.
//
// Parameters:
//   - s: The slice of E's to filter.
//   - predicate: The predicate to use for filtering.
//
// Returns:
//   - error: An error if the predicate is nil.
//
// Errors:
//   - common.ErrBadParam: If the receiver is nil.
func Filter[S ~[]E, E any](s *S, predicate Predicate[E]) error {
	if predicate == nil {
		err := common.NewErrNilParam("predicate")
		return err
	}

	if s == nil || len(*s) == 0 {
		return nil
	}

	var end uint

	for _, v := range *s {
		ok := predicate(v)
		if ok {
			(*s)[end] = v
			end++
		}
	}

	if end == 0 {
		clear(*s)
		*s = nil
	} else {
		clear((*s)[end:])
		*s = (*s)[:end]
	}

	return nil
}

// Reject rejects a slice of E's according to a predicate.
//
// Parameters:
//   - s: The slice of E's to reject.
//   - predicate: The predicate to use for rejecting.
//
// Returns:
//   - error: An error if the predicate is nil.
//
// Errors:
//   - common.ErrBadParam: If the receiver is nil.
func Reject[S ~[]E, E any](s *S, predicate Predicate[E]) error {
	if predicate == nil {
		err := common.NewErrNilParam("predicate")
		return err
	}

	if s == nil || len(*s) == 0 {
		return nil
	}

	var end uint

	for _, v := range *s {
		ok := predicate(v)
		if !ok {
			(*s)[end] = v
			end++
		}
	}

	if end == 0 {
		clear(*s)
		*s = nil
	} else {
		clear((*s)[end:])
		*s = (*s)[:end]
	}

	return nil
}

// RejectNils removes all nil values from a slice of pointers.
//
// Parameters:
//   - s: The slice of pointers to E to filter.
//
// Returns:
//   - S: The filtered slice of pointers to E.
func RejectNils[S ~[]*E, E any](s S) S {
	if len(s) == 0 {
		return *new(S)
	}

	count := 0

	for _, v := range s {
		if v != nil {
			count++
		}
	}

	if count == 0 {
		return *new(S)
	}

	result := make(S, 0, count)

	for _, v := range s {
		if v != nil {
			result = append(result, v)
		}
	}

	return result
}
