package parser

import (
	"errors"
	"iter"

	gcers "github.com/PlayerR9/go-commons/errors"
)

// History is a history of items.
type History[T any] struct {
	// timeline is the timeline of the history.
	timeline []T

	// current is the current index in the timeline.
	current int
}

// Copy creates a copy of the history.
//
// Returns:
//   - *History[T]: The copy. Never returns nil.
func (h History[T]) Copy() *History[T] {
	timeline := make([]T, len(h.timeline))
	copy(timeline, h.timeline)

	return &History[T]{
		timeline: timeline,
		current:  h.current,
	}
}

func (h *History[T]) Restart() {
	if h == nil {
		return
	}

	h.current = 0
}

// AddEvent adds an event to the history. Does nothing if the receiver
// is nil.
//
// Parameters:
//   - event: The event to add to the timeline.
func (h *History[T]) AddEvent(event T) {
	if h == nil {
		return
	}

	h.timeline = append(h.timeline, event)
}

func (h *History[T]) Event() iter.Seq[T] {
	if h == nil {
		return func(yield func(T) bool) {}
	}

	return func(yield func(T) bool) {
		for i := h.current; i < len(h.timeline); i++ {
			if !yield(h.timeline[i]) {
				h.current = i

				return
			}
		}

		h.current = len(h.timeline)
	}
}

func Advance[T any, S interface {
	ApplyEvent(event T) (bool, error)
}](history *History[T], subject S) (bool, error) {
	if history == nil {
		return true, gcers.NewErrNilParameter("history")
	}

	if history.current >= len(history.timeline) {
		return true, nil
	}

	event := history.timeline[history.current]
	history.current++

	ok, err := subject.ApplyEvent(event)
	return ok, err
}

func Nexts[T any, S interface {
	DetermineNextEvents() ([]T, error)
}](history *History[T], subject S) ([]*History[T], error) {
	events, err := subject.DetermineNextEvents()
	if err != nil {
		return nil, err
	} else if len(events) == 0 {
		return nil, nil
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

	return new_histories, nil
}

func Execute[T any, S interface {
	ApplyEvent(event T) (bool, error)
	DetermineNextEvents() ([]T, error)
}](history *History[T], subject S) ([]*History[T], error) {
	var possible []*History[T]

	is_done := false

	for !is_done {
		tmp, err := Nexts(history, subject)
		if err != nil {
			return possible, err
		}

		if len(tmp) > 1 {
			possible = append(possible, tmp[1:]...)
		}

		history = tmp[0]

		is_done, err = Advance(history, subject)
		if err != nil {
			return possible, err
		}
	}

	return possible, nil
}

func Align[T any, S interface {
	ApplyEvent(event T) (bool, error)
}](history *History[T], subject S) error {
	for event := range history.Event() {
		done, err := subject.ApplyEvent(event)
		if err != nil {
			return err
		} else if done {
			return errors.New("reached done before the end of the history")
		}
	}

	return nil
}
