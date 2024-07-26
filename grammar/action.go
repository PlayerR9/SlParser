package grammar

import "fmt"

// Actioner is an action of the grammar.
type Actioner interface {
	fmt.Stringer
}

// ActShift is a shift action.
type ActShift struct{}

// String implements the Actioner interface.
//
// Format:
//
//	"shift"
func (a *ActShift) String() string {
	return "shift"
}

// NewActShift is a constructor for an ActShift.
//
// Returns:
//   - *ActShift: The created ActShift. Never returns nil.
func NewActShift() *ActShift {
	return &ActShift{}
}

// ActReduce is a reduce action.
type ActReduce[T TokenTyper] struct {
	// rule is the rule of the action.
	rule *Rule[T]
}

// String implements the Actioner interface.
//
// Format:
//
//	"reduce"
func (a *ActReduce[T]) String() string {
	return "reduce"
}

// NewActReduce is a constructor for an ActReduce.
//
// Parameters:
//   - rule: The rule of the action.
//
// Returns:
//   - *ActReduce: The created ActReduce.
//
// Returns nil iff the rule is nil.
func NewActReduce[T TokenTyper](rule *Rule[T]) *ActReduce[T] {
	if rule == nil {
		return nil
	}

	return &ActReduce[T]{
		rule: rule,
	}
}

// ActAccept is an accept action.
type ActAccept[T TokenTyper] struct {
	// rule is the rule of the action.
	rule *Rule[T]
}

// String implements the Actioner interface.
//
// Format:
//
//	"accept"
func (a *ActAccept[T]) String() string {
	return "accept"
}

// NewActAccept is a constructor for an ActAccept.
//
// Parameters:
//   - rule: The rule of the action.
//
// Returns:
//   - *ActAccept: The created ActAccept.
//
// Returns nil iff the rule is nil.
func NewActAccept[T TokenTyper](rule *Rule[T]) *ActAccept[T] {
	if rule == nil {
		return nil
	}

	return &ActAccept[T]{
		rule: rule,
	}
}
