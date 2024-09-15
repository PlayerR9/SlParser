package pkg

import (
	"strings"

	gcers "github.com/PlayerR9/go-commons/errors"
	gcint "github.com/PlayerR9/go-commons/ints"
)

// Item is an item of the grammar.
type Item struct {
	// rule is the rule of the item.
	rule *Rule

	// pos is the position of the item in the rule.
	pos int

	// Action is the action of the item.
	Action ActionType
}

// String implements the fmt.Stringer interface.
//
// Format:
//
//	RHS(n) RHS(n-1) ... RHS(1) -> LHS : action .
func (item *Item) String() string {
	var values []string

	for i, rhs := range item.rule.GetRhss() {
		if i == item.pos {
			values = append(values, "[")
			values = append(values, rhs)
			values = append(values, "]")
		} else {
			values = append(values, rhs)
		}
	}

	values = append(values, "->", item.rule.GetLhs(), ":", item.Action.String(), ".")

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
func NewItem(rule *Rule, pos int, action ActionType) (*Item, error) {
	if rule == nil {
		return nil, gcers.NewErrNilParameter("rule")
	} else if action < 0 || action >= 3 {
		return nil, gcers.NewErrInvalidParameter("action", gcint.NewErrOutOfBounds(int(action), 0, 2))
	}

	size := rule.Size()
	if pos < 0 || pos >= size {
		return nil, gcers.NewErrInvalidParameter("pos", gcint.NewErrOutOfBounds(pos, 0, size))
	}

	return &Item{
		rule:   rule,
		pos:    pos,
		Action: action,
	}, nil
}

// GetItemTempl returns the template of the item.
//
// Parameters:
//   - pkg_name: The name of the package.
//   - tt_name: The name of the token type.
//
// Returns:
//   - string: The template of the item.
func (item *Item) GetItemTempl(pkg_name, tt_name string) string {
	var builder strings.Builder

	switch item.Action {
	case Shift:
		builder.WriteString("act = ")
		builder.WriteString(pkg_name)
		builder.WriteString(".NewShiftAction()")
	case Reduce:
		builder.WriteString("act, _ = ")
		builder.WriteString(pkg_name)
		builder.WriteString(".NewReduceAction(")
		builder.WriteString(item.rule.GetRuleTempl(pkg_name, tt_name))
		builder.WriteRune(')')
	case Accept:
		builder.WriteString("act, _ = ")
		builder.WriteString(pkg_name)
		builder.WriteString(".NewAcceptAction(")
		builder.WriteString(item.rule.GetRuleTempl(pkg_name, tt_name))
		builder.WriteRune(')')
	}

	return builder.String()
}
