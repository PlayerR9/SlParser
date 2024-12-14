package internal

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/PlayerR9/SlParser/PlayerR9/mygo-lib/common"
	assert "github.com/PlayerR9/go-verify"
	"github.com/PlayerR9/mygo-data/sets"
)

// InternalState is the internal state of a parser.
//
// An empty internal state can be created with the `var is InternalState` syntax or with the
// `is := new(InternalState)` constructor.
type InternalState struct {
	// last_items is the last items selected by the decision table.
	last_items []*Item

	// has_error is a flag that indicates whether an error has occurred.
	has_error bool

	// current_rhs is the current RHS.
	current_rhs string

	// expecteds is the expected tokens.
	expecteds []string

	// got is the got token.
	got string

	// phase is the current phase.
	phase PhaseType

	// event is the current event.
	event *Event
}

// UpdatePhase updates the current phase of the internal state.
//
// Parameters:
//   - phase: The new phase to set.
//
// Returns:
//   - error: An error if the phase could not be updated.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
func (is *InternalState) UpdatePhase(phase PhaseType) error {
	if is == nil {
		return common.ErrNilReceiver
	}

	is.phase = phase

	return nil
}

// SetEvent sets the current event of the internal state.
//
// Parameters:
//   - event: The event to set.
//
// Returns:
//   - error: An error if the event could not be set.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
func (is *InternalState) SetEvent(event *Event) error {
	if is == nil {
		return common.ErrNilReceiver
	}

	is.event = event

	return nil
}

// ToggleError toggles the error status of the internal state.
//
// Returns:
//   - error: An error if the error status could not be toggled.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
func (is *InternalState) ToggleError() error {
	if is == nil {
		return common.ErrNilReceiver
	}

	is.has_error = !is.has_error

	return nil
}

// ChangeLastItems changes the last items of the internal state.
//
// Parameters:
//   - items: The new last items.
//
// Returns:
//   - error: An error if the last items could not be changed.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
func (is *InternalState) ChangeLastItems(items []*Item) error {
	if is == nil {
		return common.ErrNilReceiver
	}

	is.last_items = items

	return nil
}

// HasError returns the error status of the internal state.
//
// Returns:
//   - bool: True if the internal state has an error, false otherwise.
func (is InternalState) HasError() bool {
	return is.has_error
}

// SetCurrentRHS sets the current RHS of the internal state.
//
// Parameters:
//   - rhs: The RHS to set.
//
// Returns:
//   - error: An error if the RHS could not be set.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
func (is *InternalState) SetCurrentRHS(rhs string) error {
	if is == nil {
		return common.ErrNilReceiver
	}

	is.current_rhs = rhs

	return nil
}

// SetExpecteds sets the expected values of the internal state.
//
// Parameters:
//   - expecteds: The expected values to set.
//
// Returns:
//   - error: An error if the expected values could not be set.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
func (is *InternalState) SetExpecteds(expecteds []string) error {
	if is == nil {
		return common.ErrNilReceiver
	}

	is.expecteds = expecteds

	return nil
}

// SetGot sets the current got value of the internal state.
//
// Parameters:
//   - rhs: The got value to set.
//
// Returns:
//   - error: An error if the got value could not be set.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
func (is *InternalState) SetGot(rhs string) error {
	if is == nil {
		return common.ErrNilReceiver
	}

	is.got = rhs

	return nil
}

// MakeError creates a detailed error message based on the current phase and state of the parser.
//
// Returns:
//   - error: A formatted error message indicating the context and reason for the error.
//
// Errors:
//   - During PhaseReduction: Returns an error indicating a token mismatch with expected tokens.
//   - During PhaseShifting: Returns an error indicating EOF was reached unexpectedly.
//   - During PhaseCheckBranch: Returns an error indicating a token mismatch with the expected next token.
//   - Default case: Returns an error indicating a parsing issue, such as reaching EOF without an accept
//     state or encountering unsupported tokens.
func (is InternalState) MakeError() error {
	switch is.phase {
	case PhaseReduction:
		return fmt.Errorf("while reducing: %w", NewErrNotAsExpected(true, "token", is.got, is.expecteds...))
	case PhaseShifting:
		return fmt.Errorf("while shifting: %w", errors.New("EOF reached"))
	case PhaseCheckBranch:
		expected := is.event.ExpectedNext

		return fmt.Errorf("while checking branch: %w", NewErrNotAsExpected(true, "token", is.got, expected))
	default:
		items := is.last_items

		rhs := is.current_rhs
		if rhs == "" {
			return fmt.Errorf("while parsing: %w", errors.New("EOF reached yet no accept state was reached"))
		}

		if len(items) == 0 {
			return fmt.Errorf("while parsing: %w", NewUnsupportedValue("token", strconv.Quote(rhs)))
		}

		expecteds := new(sets.OrderedSet[string])

		for _, item := range items {
			indices := item.IndicesOf(rhs)

			for _, idx := range indices {
				next_rhs, ok := item.RhsAt(idx + 1)
				if ok {
					err := expecteds.Insert(next_rhs)
					assert.Err(err, "expecteds.Insert(%s)", strconv.Quote(next_rhs))
				}
			}
		}

		elems := expecteds.Slice()

		for i, elem := range elems {
			elems[i] = strconv.Quote(elem)
		}

		got := strconv.Quote(rhs)

		return fmt.Errorf("while parsing: %w", NewErrNotAsExpected(true, "token", got, elems...))
	}
}
