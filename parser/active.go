package parser

import (
	"fmt"
	"slices"
	"strings"

	"github.com/PlayerR9/SlParser/grammar"
	slgr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/parser/internal"
	assert "github.com/PlayerR9/go-verify"

	ern "github.com/PlayerR9/go-evals/rank"
	// assert "github.com/PlayerR9/go-verify"
	ll "github.com/PlayerR9/mygo-lib/CustomData/listlike"
	tr "github.com/PlayerR9/mygo-lib/CustomData/tree"
	"github.com/PlayerR9/mygo-lib/common"
	gslc "github.com/PlayerR9/mygo-lib/slices"
)

/////////////////////////////////////////////////////////

// Active is the active parser.
type Active struct {
	// global is the parser to use. Must not be nil.
	global Parser

	// input_stream is the tokens to parse. Must not be nil.
	input_stream []*tr.Node

	// pos is the current position in the input stream. It indicates the position
	// of the first non-shifted token in the input stream.
	pos int

	// stack is the stack to use. Must not be nil.
	stack *ll.RefusableStack[*tr.Node]

	// is_done is a flag that indicates whether the parsing is done.
	is_done bool

	// state is the internal state of the parser.
	state *internalState
}

// ApplyEvent implements the history.Subject interface.
func (active *Active) ApplyEvent(event *internal.Event) error {
	if active == nil {
		return common.ErrNilReceiver
	} else if event == nil {
		return common.NewErrNilParam("event")
	}

	_ = active.state.SetEvent(event)

	en := event.ExpectedNext()
	if en == "" {
		switch act := event.Action(); act {
		case internal.ActShift:
			err := active.shift()
			if err != nil {
				return fmt.Errorf("error while shifting: %w", err)
			}
		case internal.ActReduce:
			err := active.reduce(event.Rule())
			if err != nil {
				return fmt.Errorf("error while reducing: %w", err)
			}
		case internal.ActAccept:
			err := active.reduce(event.Rule())
			if err != nil {
				return fmt.Errorf("error while reducing: %w", err)
			}

			active.is_done = true
		default:
			return NewUnsupportedValue("action type", act.String())
		}
	} else {
		is_terminal := slgr.IsTerminal(en)

		act := event.Action()

		switch act {
		case internal.ActShift:
			err := active.shift()
			if err != nil {
				return fmt.Errorf("error while shifting: %w", err)
			}

		case internal.ActReduce:
			err := active.reduce(event.Rule())
			if err != nil {
				return fmt.Errorf("error while reducing: %w", err)
			}

		case internal.ActAccept:
			err := active.reduce(event.Rule())
			if err != nil {
				return fmt.Errorf("error while reducing: %w", err)
			}

			active.is_done = true
		default:
			return NewUnsupportedValue("action type", act.String())
		}

		if !is_terminal {
			if act != internal.ActShift {
				_ = active.state.UpdatePhase(PhaseCheckBranch)

				top, err := active.stack.Peek()
				if err != nil {
					_ = active.state.SetGot("")
					_ = active.state.ToggleError()
				}

				topd := grammar.MustGet[*grammar.TokenData](top)

				if topd.Type != en {
					_ = active.state.SetGot(topd.Type)
					_ = active.state.ToggleError()
				}

				return nil
			}
		} else {
			if act == internal.ActShift {
				_ = active.state.UpdatePhase(PhaseCheckBranch)

				top, err := active.stack.Peek()
				if err != nil {
					_ = active.state.SetGot("")
					_ = active.state.ToggleError()
				}

				topd := grammar.MustGet[*grammar.TokenData](top)
				if topd.Type != en {
					_ = active.state.SetGot(topd.Type)
					_ = active.state.ToggleError()
				}

				return nil
			}
		}
	}

	return nil
}

// HasError implements the history.Subject interface.
func (active Active) HasError() bool {
	return active.state.HasError()
}

// GetError implements the history.Subject interface.
func (active Active) GetError() error {
	return active.state.makeError()
}

// NextEvents implements the history.Subject interface.
func (active *Active) NextEvents() ([]*internal.Event, error) {
	if active == nil {
		return nil, common.ErrNilReceiver
	}

	_ = active.state.UpdatePhase(PhasePrediction)

	// fmt.Println(active.DebugStackString())
	// fmt.Println()

	if active.is_done {
		return nil, nil
	}

	defer active.stack.Refuse()

	var item_copy []*internal.Item

	defer func() {
		_ = active.state.ChangeLastItems(item_copy)
	}()

	top1, err := active.stack.Pop()
	if err != nil {
		_ = active.state.SetCurrentRHS("")
		_ = active.state.ToggleError()

		return nil, nil
	}

	top1d := grammar.MustGet[*grammar.TokenData](top1)

	type_ := top1d.Type

	_ = active.state.SetCurrentRHS(type_)

	items, ok := active.global.ItemsOf(type_)
	if !ok || len(items) == 0 {
		_ = active.state.ToggleError()

		return nil, nil
	}

	if len(items) == 1 {
		item := items[0]

		next, ok := item.NextRhs()
		if !ok {
			next = ""
		}

		return []*internal.Event{
			internal.NewEvent(item, next),
		}, nil
	}

	item_copy = make([]*internal.Item, len(items))
	copy(item_copy, items)

	// la := top1.Lookahead

	// if la == nil {
	// 	fn := func(item *internal.Item[T]) bool {
	// 		return item.ExpectLookahead()
	// 	}

	// 	item_copy = internal.RejectSlice(item_copy, fn)
	// } else {
	// 	fn := func(item *internal.Item[T]) bool {
	// 		if item.ExpectLookahead() {
	// 			if item.HasLookahead(la.Type) {
	// 				return true
	// 			} else {
	// 				return false
	// 			}
	// 		} else {
	// 			return true
	// 		}
	// 	}

	// 	item_copy = internal.FilterSlice(item_copy, fn)
	// }

	var offset uint = 2

	sols := new(ern.Rank[*internal.Item])

	var early_exit bool

	for len(item_copy) > 0 && !early_exit {
		top, err := active.stack.Pop()
		if err != nil {
			early_exit = true
			break
		}

		topd := grammar.MustGet[*grammar.TokenData](top)

		type_ := topd.Type

		sameRhsFn := func(item *internal.Item) bool {
			pos := item.Pos()

			if pos < offset {
				err := sols.Add(int(offset), item)
				assert.Err(err, "sols.Add(%d, item)", int(offset))

				return false
			}

			rhs, ok := item.RhsAt(pos - offset)
			assert.True(ok, "item.RhsAt(%d - %d)", pos, offset)

			return rhs == type_
		}

		ok := gslc.FilterIfApplicable(&item_copy, sameRhsFn)
		if !ok {
			if sols.Size() != 0 {
				clear(item_copy)
				item_copy = nil
			} else {
				early_exit = true
			}
		}

		offset++
	}

	if early_exit {
		isNotDoneFn := func(item *internal.Item) bool {
			_, ok := item.RhsAt(item.Pos() - offset)
			return ok
		}

		ok := gslc.RejectIfApplicable(&item_copy, isNotDoneFn)

		if ok {
			for _, item := range item_copy {
				err := sols.Add(int(offset), item)
				assert.Err(err, "sols.Add(%d, item)", int(offset))
			}
		}
	}

	res := sols.Build()
	if len(res) > 0 {
		events := make([]*internal.Event, 0, len(res))

		for _, item := range res {
			next, ok := item.NextRhs()
			if !ok {
				next = ""
			}

			event := internal.NewEvent(item, next)
			events = append(events, event)
		}

		return events, nil
	} else {
		_ = active.state.ToggleError()
		return nil, nil
	}
}

// NewActive creates a new parser based on the given parser and tokens.
//
// Parameters:
//   - global: The parser to use. Must not be nil.
//   - tokens: The tokens to parse. Must not be nil.
//
// Returns:
//   - *Active[T]: The new parser.
//   - error: An error if the initial shift failed.
func NewActive(global *baseParser, tokens []*tr.Node) (*Active, error) {
	if global == nil {
		return nil, common.NewErrNilParam("global")
	}

	active := &Active{
		global:       global,
		input_stream: tokens,
		stack:        new(ll.RefusableStack[*tr.Node]),
		pos:          0,
		state:        newInternalState(),
	}

	err := active.shift() // Initial Shift
	if err != nil {
		return nil, fmt.Errorf("initial shift failed: %w", err)
	}

	return active, nil
}

// push is a helper method that pushes a token onto the parse stack.
//
// Parameters:
//   - tk: The token to push. Must not be nil.
//
// Returns:
//   - error: An error if the receiver is nil.
func (active *Active) push(tk *tr.Node) error {
	if active == nil {
		return common.ErrNilReceiver
	}

	_ = active.stack.Push(tk)

	return nil
}

// shift is a helper method that shifts a token from the input stream onto the
// parse stack.
//
// If there are no more tokens to shift, it sets the error field to a descriptive
// error.
//
// Returns:
//   - error: An error if the push failed.
func (active *Active) shift() error {
	if active == nil {
		return common.ErrNilReceiver
	}

	_ = active.state.UpdatePhase(PhaseShifting)

	if active.pos >= len(active.input_stream) {
		_ = active.state.ToggleError()

		return nil
	}

	tk := active.input_stream[active.pos]
	active.pos++

	err := active.push(tk)
	return err
}

// reduce reduces the top of the parse stack by the given rule.
//
// It checks that the top of the parse stack matches the rhs of the rule and
// creates a new token with the lhs of the rule and the input stream of the top
// of the parse stack. It sets the lookahead of the new token to the lookahead of
// the top of the parse stack.
//
// Finally, it pushes the new token onto the parse stack.
//
// Parameters:
//   - rule: The rule to use for reduction. Must not be nil.
//
// Returns:
//   - error: An error if reduction failed.
func (active *Active) reduce(rule *internal.Rule) error {
	if active == nil {
		return common.ErrNilReceiver
	}

	_ = active.state.UpdatePhase(PhaseReduction)

	rhss := rule.Rhss()
	slices.Reverse(rhss)

	for _, rhs := range rhss {
		_ = active.state.SetExpecteds([]string{rhs})

		top, err := active.stack.Pop()
		if err != nil {
			_ = active.state.SetGot("")
			_ = active.state.ToggleError()

			return nil
		}

		topd := grammar.MustGet[*grammar.TokenData](top)

		type_ := topd.Type

		if type_ != rhs {
			_ = active.state.SetGot(type_)
			_ = active.state.ToggleError()

			return nil
		}
	}

	popped := active.stack.Popped()
	active.stack.Accept()

	tk := slgr.NewToken(
		grammar.MustGet[*grammar.TokenData](popped[0]).Pos,
		rule.Lhs(),
		"",
		grammar.MustGet[*grammar.TokenData](popped[len(popped)-1]).Lookahead,
	)

	_ = tk.AppendChildren(popped...)

	_ = active.push(tk)

	return nil
}

// Forest returns the parse forest of the parser. The parse forest is a slice of
// parse trees, where each parse tree is a tree of tokens. The parse forest is
// constructed by popping all the tokens from the parse stack and constructing a
// parse tree for each one. The parse forest is useful for debugging and for
// visualizing the parse tree.
//
// The function never returns an error. If the parse stack is empty, the function
// returns a slice with a single element, which is a parse tree with a single
// token, which is the end of file token.
//
// Returns:
//   - []*ParseTree[T]: A slice of parse trees. The slice is never empty, since the parser always has a parse stack.
func (active *Active) Forest() []*ParseTree {
	if active == nil {
		return nil
	}

	defer active.stack.Refuse()

	forest := make([]*ParseTree, 0, active.stack.Size())

	for {
		tk, err := active.stack.Pop()
		if err != nil {
			break
		}

		tree, _ := NewParseTree(tk)

		forest = append(forest, tree)
	}

	return forest
}

// Shadow creates a shadow of the active parser. The shadow is a new parser that
// is a copy of the original parser. The shadow is useful for debugging and for
// testing the parser. The function never returns an error. If the function fails,
// it panics.
//
// The function returns a pointer to the shadow parser. The pointer is never nil.
//
// Returns:
//   - *Active[T]: A pointer to the shadow parser.
//   - error: An error if the function failed.
func (active Active) Shadow() (*Active, error) {
	shadow := &Active{
		global:       active.global,
		input_stream: active.input_stream,
		stack:        new(ll.RefusableStack[*tr.Node]),
		state:        newInternalState(),
		pos:          0,
	}

	err := shadow.shift() // Initial Shift
	if err != nil {
		return nil, fmt.Errorf("initial shift failed: %w", err)
	}

	return shadow, nil
}

// DebugStackString returns a string representation of the parse stack in a
// human-readable format. It constructs a parse forest by popping all tokens
// from the stack, creating a parse tree for each token, and then converting
// each parse tree to a string. The resulting strings are reversed to reflect
// the original order of the parse stack and concatenated with newline
// separators.
//
// Returns:
//   - string: A human-readable string representation of the parse stack.
func (active Active) DebugStackString() string {
	defer active.stack.Refuse()

	forest := make([]*ParseTree, 0, active.stack.Size())

	for {
		tk, err := active.stack.Pop()
		if err != nil {
			break
		}

		tree, _ := NewParseTree(tk)

		forest = append(forest, tree)
	}

	var lines []string

	for _, tree := range forest {
		lines = append(lines, tree.String())
	}

	slices.Reverse(lines)

	return strings.Join(lines, "\n")
}

// DebugInputStreamString returns a string representation of the input stream in a human-readable format. It
// constructs a slice of strings by iterating over the input stream and calling the String() method on each
// token. The resulting strings are reversed to reflect the original order of the input stream and
// concatenated with " <- " separators.
//
// Returns:
//   - string: A human-readable string representation of the input stream.
func (active Active) DebugInputStreamString() string {
	elems := make([]string, 0, len(active.input_stream[active.pos:]))

	for _, tk := range active.input_stream[active.pos:] {
		elems = append(elems, tk.String())
	}

	slices.Reverse(elems)

	return strings.Join(elems, " <- ")
}
