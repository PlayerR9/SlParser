package internal

import assert "github.com/PlayerR9/go-verify"

// Event is the result of a match.
type Event struct {
	// Item is the item that was matched.
	Item *Item

	// ExpectedNext is the rhs symbol that is expected to follow the matched item.
	// This is used for optimization.
	ExpectedNext string
}

// NewEvent creates a new Event.
//
// Parameters:
//   - item: The item that was matched. If nil, the function returns nil.
//   - expected_next: The rhs symbol that is expected to follow the matched item.
//     This is used for optimization.
//
// Returns:
//   - *Event: The newly created Event. Nil if the item was nil.
func NewEvent(item *Item, expected_next string) *Event {
	if item == nil {
		return nil
	}

	event := &Event{
		Item:         item,
		ExpectedNext: expected_next,
	}

	return event
}

// Action returns the action type associated with the event.
//
// Returns:
//   - ActionType: The action type of the event.
func (e Event) Action() ActionType {
	act := e.Item.Action()
	return act
}

// Rule returns the rule associated with the event.
//
// Returns:
//   - *Rule: The rule of the event. Never returns nil.
func (e Event) Rule() *Rule {
	rule := e.Item.Rule()
	assert.Cond(rule != nil, "rule != nil")

	return rule
}
