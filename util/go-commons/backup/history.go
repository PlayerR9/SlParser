package backup

import (
	"iter"

	dba "github.com/PlayerR9/go-debug/assert"
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

func Align[T any, S interface {
	ApplyEvent(event T) bool
	HasError() bool
}](history *History[T], subject S) {
	for event := range history.Event() {
		done := subject.ApplyEvent(event)
		if subject.HasError() {
			break
		}

		dba.Assert(!done, "somehow the subject was done before the end of the history")
	}
}
