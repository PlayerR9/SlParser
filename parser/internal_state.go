package parser

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/PlayerR9/SlParser/parser/internal"
	"github.com/PlayerR9/mygo-lib/common"
)

/////////////////////////////////////////////////////////

type PhaseType int

const (
	PhaseReduction PhaseType = iota
	PhasePrediction
	PhaseShifting
	PhaseCheckBranch
)

type internalState struct {
	last_items  []*internal.Item
	has_error   bool
	current_rhs string
	expecteds   []string
	got         string
	phase       PhaseType
	event       *internal.Event
}

func newInternalState() *internalState {
	return &internalState{}
}

func (is *internalState) UpdatePhase(phase PhaseType) error {
	if is == nil {
		return common.ErrNilReceiver
	}

	is.phase = phase

	return nil
}

func (is *internalState) SetEvent(event *internal.Event) error {
	if is == nil {
		return common.ErrNilReceiver
	}

	is.event = event

	return nil
}

func (is *internalState) ToggleError() error {
	if is == nil {
		return common.ErrNilReceiver
	}

	is.has_error = !is.has_error

	return nil
}

func (is *internalState) ChangeLastItems(items []*internal.Item) error {
	if is == nil {
		return common.ErrNilReceiver
	}

	is.last_items = items

	return nil
}

func (is internalState) HasError() bool {
	return is.has_error
}

func (is *internalState) SetCurrentRHS(rhs string) error {
	if is == nil {
		return common.ErrNilReceiver
	}

	is.current_rhs = rhs

	return nil
}

func (is *internalState) SetExpecteds(expecteds []string) error {
	if is == nil {
		return common.ErrNilReceiver
	}

	is.expecteds = expecteds

	return nil
}

func (is *internalState) SetGot(rhs string) error {
	if is == nil {
		return common.ErrNilReceiver
	}

	is.got = rhs

	return nil
}

func (is internalState) makeError() error {
	switch is.phase {
	case PhaseReduction:
		return fmt.Errorf("while reducing: %w", common.NewErrNotAsExpected(true, "token", is.got, is.expecteds...))
	case PhaseShifting:
		return fmt.Errorf("while shifting: %w", errors.New("EOF reached"))
	case PhaseCheckBranch:
		expected := is.event.ExpectedNext()

		return fmt.Errorf("while checking branch: %w", common.NewErrNotAsExpected(true, "token", is.got, expected))
	default:
		items := is.last_items

		rhs := is.current_rhs
		if rhs == "" {
			return fmt.Errorf("while parsing: %w", errors.New("EOF reached yet no accept state was reached"))
		}

		if len(items) == 0 {
			return fmt.Errorf("while parsing: %w", NewUnsupportedValue("token", strconv.Quote(rhs)))
		}

		var expecteds []string

		for _, item := range items {
			indices := item.IndicesOf(rhs)

			for _, idx := range indices {
				next_rhs, ok := item.RhsAt(idx + 1)
				if !ok {
					continue
				}

				pos, ok := slices.BinarySearch(expecteds, next_rhs)
				if !ok {
					expecteds = slices.Insert(expecteds, pos, next_rhs)
				}
			}
		}

		for i, rhs := range expecteds {
			expecteds[i] = strconv.Quote(rhs)
		}

		got := strconv.Quote(rhs)

		return fmt.Errorf("while parsing: %w", common.NewErrNotAsExpected(true, "token", got, expecteds...))
	}
}
