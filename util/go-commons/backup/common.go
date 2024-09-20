package backup

import (
	"iter"
	"slices"
)

type pairing[T any, S any] struct {
	History *History[T]
	Subject S
}

func new_pairing[T any, S any](history *History[T], subject S) *pairing[T, S] {
	if history == nil {
		history = &History[T]{}
	}

	return &pairing[T, S]{
		History: history,
		Subject: subject,
	}
}

func advance[T any, S interface {
	ApplyEvent(event T) bool
}](history *History[T], subject S) bool {
	if history == nil || history.current >= len(history.timeline) {
		return true
	}

	event := history.timeline[history.current]
	history.current++

	ok := subject.ApplyEvent(event)
	return ok
}

func nexts[T any, S interface {
	DetermineNextEvents() []T
	HasError() bool
}](history *History[T], subject S) ([]*History[T], bool) {
	events := subject.DetermineNextEvents()
	if subject.HasError() {
		return nil, false
	} else if len(events) == 0 {
		return nil, true
	}

	new_histories := make([]*History[T], 0, len(events))

	if history == nil {
		history = &History[T]{}
	}

	for _, event := range events {
		h := history.Copy()
		h.AddEvent(event)

		new_histories = append(new_histories, h)
	}

	for i := 1; i < len(new_histories); i++ {
		new_histories[i].Restart()
	}

	return new_histories, true
}

func execute_one[T any, S interface {
	ApplyEvent(event T) bool
	DetermineNextEvents() []T
	HasError() bool
}](history *History[T], subject S) ([]*History[T], bool) {
	var possible []*History[T]

	is_done := false

	for !is_done {
		tmp, ok := nexts(history, subject)
		if !ok || len(tmp) == 0 {
			break
		}

		if len(tmp) > 1 {
			possible = append(possible, tmp[1:]...)
		}

		history = tmp[0]

		is_done = advance(history, subject)
		if subject.HasError() {
			break
		}
	}

	return possible, is_done
}

// Subject returns a sequence of all possible states of a subject that can be
// reached by applying all possible events to the subject in all possible
// orders.
//
// The subject must implement the following interface:
//
//	interface {
//	    Align(history *History[T]) bool
//	    ApplyEvent(event T) (bool, error)
//	    DetermineNextEvents() ([]T, error)
//	}
//
// The function init_fn must return a new instance of the subject.
//
// The function returns a sequence of all possible states of the subject. The
// sequence is lazy, i.e. it is generated on-the-fly as the caller iterates
// over it.
//
// Parameters:
//   - init_fn: A function that returns a new instance of the subject.
//
// Returns:
//   - iter.Seq[S]: A sequence of all possible states of the subject.
//
// If the function 'init_fn' is nil, it defaults to 'var subject S'.
func Subject[T any, S interface {
	Align(history *History[T]) bool
	ApplyEvent(event T) bool
	DetermineNextEvents() []T
	HasError() bool
}](init_fn func() S) iter.Seq[S] {
	if init_fn == nil {
		init_fn = func() S {
			return *new(S)
		}
	}

	fn := func(yield func(S) bool) {
		sbj := init_fn()
		pair := new_pairing[T](nil, sbj)

		var invalid_subjects []S

		stack := []*pairing[T, S]{pair}

		for len(stack) > 0 {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			ok := top.Subject.Align(top.History)
			if !ok {
				invalid_subjects = append(invalid_subjects, top.Subject)

				continue
			}

			possible, ok := execute_one(top.History, top.Subject)

			if len(possible) > 0 {
				pairs := make([]*pairing[T, S], 0, len(possible))

				for _, path := range possible {
					sbj := init_fn()

					pair := new_pairing(path, sbj)
					pairs = append(pairs, pair)
				}

				slices.Reverse(pairs)

				stack = append(stack, pairs...)
			}

			if !ok {
				invalid_subjects = append(invalid_subjects, top.Subject)
			} else if !yield(top.Subject) {
				return
			}
		}

		for _, sbj := range invalid_subjects {
			if !yield(sbj) {
				return
			}
		}
	}

	return fn
}
