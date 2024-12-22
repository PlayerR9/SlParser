package parser

import (
	"github.com/PlayerR9/SlParser/parser/internal"
)

type Action interface {
}

type ShiftAction struct {
}

func NewShiftAction() Action {
	act := &ShiftAction{}
	return act
}

type ReduceAction struct {
	rule *internal.Rule
}

func NewReduceAction(lhs string, rhss ...string) Action {
	rule := internal.NewRule(lhs, rhss)

	act := &ReduceAction{
		rule: rule,
	}
	return act
}

type AcceptAction struct {
	rule *internal.Rule
}

func NewAcceptAction(lhs string, rhss ...string) Action {
	rule := internal.NewRule(lhs, rhss)

	act := &AcceptAction{
		rule: rule,
	}
	return act
}
