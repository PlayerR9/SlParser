package grammar

import (
	"errors"
	"fmt"
	"iter"
)

// TokenTyper is the interface that must be implemented by token types.
type TokenTyper interface {
	~int

	// String returns the string representation of the token type.
	//
	// Returns:
	//   - string: the string representation of the token type.
	String() string
}

// Token is a token.
type Token[T TokenTyper] struct {
	// Type is the type of the token.
	Type T

	// Data is the data of the token.
	Data string

	// Children is the list of children of the token.
	Children []*Token[T]

	// Lookahead is the lookahead of the token.
	Lookahead *Token[T]
}

// IsLeaf implements tree.TreeNoder interface.
func (t Token[T]) IsLeaf() bool {
	return len(t.Children) == 0
}

// String implements tree.TreeNoder interface.
func (t Token[T]) String() string {
	if t.Data == "" {
		return fmt.Sprintf("Token[T][%s]", t.Type.String())
	} else {
		return fmt.Sprintf("Token[T][%s (%q)]", t.Type.String(), t.Data)
	}
}

// NewTerminalToken creates a new terminal token.
//
// Parameters:
//   - type_: the type of the token.
//   - data: the data of the token.
//
// Returns:
//   - *Token[T]: the new terminal token. Never returns nil.
func NewTerminalToken[T TokenTyper](type_ T, data string) *Token[T] {
	return &Token[T]{
		Type: type_,
		Data: data,
	}
}

// NewNonTerminalToken creates a new non-terminal token.
//
// Parameters:
//   - type_: the type of the token.
//   - children: the children of the token.
//
// Returns:
//   - *Token[T]: the new non-terminal token.
//   - error: an error if the children are empty.
func NewNonTerminalToken[T TokenTyper](type_ T, children []*Token[T]) (*Token[T], error) {
	if len(children) == 0 {
		return nil, errors.New("non-terminal token must have at least one child")
	}

	last_tk := children[len(children)-1]

	return &Token[T]{
		Type:      type_,
		Children:  children,
		Lookahead: last_tk.Lookahead,
	}, nil
}

// Child returns an iterator over the children of the token.
//
// Returns:
//   - iter.Seq[*Token[T]]: an iterator over the children of the token.
func (t Token[T]) Child() iter.Seq[*Token[T]] {
	return func(yield func(*Token[T]) bool) {
		for _, child := range t.Children {
			if !yield(child) {
				break
			}
		}
	}
}

// BackwardChild returns an iterator over the children of the token in reverse order.
//
// Returns:
//   - iter.Seq[*Token[T]]: an iterator over the children of the token in reverse order.
func (t Token[T]) BackwardChild() iter.Seq[*Token[T]] {
	return func(yield func(*Token[T]) bool) {
		for i := len(t.Children) - 1; i >= 0; i-- {
			if !yield(t.Children[i]) {
				break
			}
		}
	}
}
