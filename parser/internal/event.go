package internal

/////////////////////////////////////////////////////////

type Event struct {
	item          *Item
	expected_next string
}

func NewEvent(item *Item, expected_next string) *Event {
	if item == nil {
		return nil
	}

	return &Event{
		item:          item,
		expected_next: expected_next,
	}
}

func (e Event) Action() ActionType {
	return e.item.Action()
}

func (e Event) Rule() *Rule {
	return e.item.rule
}

func (e Event) ExpectedNext() string {
	return e.expected_next
}
