package history

import (
	"fmt"

	assert "github.com/PlayerR9/go-verify"
)

// Subject is an interface representing an entity that can walk through a history.
type Subject[E Event] interface {
	// ApplyEvent applies an event to the subject.
	//
	// Parameters:
	//   - event: The event to apply.
	//
	// Returns:
	//   - error: An error if the event could not be applied.
	//
	// NOTES: Because the returned error causes the immediate stop of the history,
	// use it only for panic-level error handling. For any other error, use the
	// HasError method to signify something non-critical happened.
	ApplyEvent(event E) error

	// NextEvents returns the next events in the subject.
	//
	// Returns:
	//   - []E: The next events in the subject.
	//   - error: An error if the next events could not be returned.
	//
	// NOTES:
	// 	- Because the returned error causes the immediate stop of the history,
	// 	use it only for panic-level error handling. For any other error, use the
	// 	HasError method to signify something non-critical happened.
	NextEvents() ([]E, error)

	// HasError checks whether the subject has an error.
	//
	// Returns:
	//   - bool: True if the subject has an error, false otherwise.
	HasError() bool

	// GetError returns the error associated with the subject. However, this is mostly
	// used as a builder for the error and, as such, it always assume an error
	// has, indeed, occurred.
	//
	// Returns:
	//   - error: The error associated with the subject.
	GetError() error
}

// realign realigns the history with the subject by applying each event in the history
// to the subject using ApplyEvent method. It stops when an error occurs or when the
// history is fully walked. It asserts that the subject does not have an error after
// applying each event.
//
// Parameters:
//   - history: The history to realign.
//   - subject: The subject to apply the events to.
//
// Returns:
//   - error: An error if the history could not be realigned.
func realign[E Event](history *History[E], subject Subject[E]) error {
	assert.Cond(history != nil, "history != nil")
	assert.Cond(subject != nil, "subject != nil")

	for {
		event, err := history.Walk()
		if err != nil {
			break
		}

		err = subject.ApplyEvent(event)
		if err != nil {
			return fmt.Errorf("while applying event: %w", err)
		}

		ok := subject.HasError()
		if ok {
			err := subject.GetError()
			err = fmt.Errorf("subject has an error: %w", err)
			return err
		}
	}

	return nil
}
