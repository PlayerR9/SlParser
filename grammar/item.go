package grammar

import (
	"strings"

	uc "github.com/PlayerR9/lib_units/common"
)

// Item is an item of the grammar.
type Item[T TokenTyper] struct {
	// rule is the rule of the item.
	rule *Rule[T]

	// pos is the position of the item in the rule.
	pos int

	// action is the action of the item.
	action Actioner
}

// String implements the fmt.Stringer interface.
//
// Format:
//
//	RHS(n) RHS(n-1) ... RHS(1) -> LHS : action .
func (item *Item[T]) String() string {
	var values []string

	iter := item.rule.Iterator()
	var i int

	for {
		rhs, err := iter.Consume()
		if err != nil {
			break
		}

		if i == item.pos {
			values = append(values, "[")
			values = append(values, rhs.GoString())
			values = append(values, "]")
		} else {
			values = append(values, rhs.GoString())
		}

		i++
	}

	values = append(values, "->", item.rule.GetLhs().GoString(), ":", item.action.String(), ".")

	return strings.Join(values, " ")
}

// NewItem is a constructor for an Item.
//
// Parameters:
//   - rule: The rule of the item.
//   - pos: The position of the item in the rule.
//   - action: The action of the item.
//
// Returns:
//   - *Item: The created item.
//   - error: The error of type *common.ErrInvalidParameter if rule or action is nil or
//     if pos is out of bounds.
func NewItem[T TokenTyper](rule *Rule[T], pos int, action Actioner) (*Item[T], error) {
	if rule == nil {
		return nil, uc.NewErrNilParameter("rule")
	} else if action == nil {
		return nil, uc.NewErrNilParameter("action")
	}

	size := rule.Size()
	if pos < 0 || pos >= size {
		return nil, uc.NewErrInvalidParameter("pos", uc.NewErrOutOfBounds(pos, 0, size))
	}

	return &Item[T]{
		rule:   rule,
		pos:    pos,
		action: action,
	}, nil
}
