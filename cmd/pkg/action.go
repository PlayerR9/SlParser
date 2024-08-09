package pkg

type ActionType int

const (
	Shift ActionType = iota
	Reduce
	Accept
)

func (a ActionType) String() string {
	return [...]string{
		"SHIFT",
		"REDUCE",
		"ACCEPT",
	}[a]
}

// DetermineAction determines the action of the item given the position and the symbol.
//
// Parameters:
//   - pos: The position of the item in the rule.
//   - symbol: The symbol of the item in the rule.
//
// Returns:
//   - ActionType: The action of the item.
func DetermineAction(pos int, symbol string) ActionType {
	if pos != 0 {
		return Shift
	}

	if symbol == "etk_EOF" {
		return Accept
	} else {
		return Reduce
	}
}
