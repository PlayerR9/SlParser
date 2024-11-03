package grammar

import (
	"fmt"
	"iter"
	"strconv"

	"github.com/PlayerR9/mygo-lib/common"
	gcslc "github.com/PlayerR9/mygo-lib/slices"
)

// Token is a token.
type Token struct {
	// Parent, FirstChild, LastChild, NextSibling, PrevSibling are all the
	// pointers required to traverse the tree.
	Parent, FirstChild, LastChild, NextSibling, PrevSibling *Token

	// Type is the type of the token.
	Type string

	// Data is the data of the token.
	Data string

	// Pos is the position of the token (in bytes).
	Pos int

	// Lookahead is the lookahead token.
	Lookahead *Token
}

// IsNil implements the tree.TreeNode interface.
func (t *Token) IsNil() bool {
	return t == nil
}

// IsLeaf implements the tree.TreeNode interface.
func (t Token) IsLeaf() bool {
	return t.FirstChild == nil
}

// String implements the tree.TreeNode interface.
func (t Token) String() string {
	var data string

	if t.Data != "" {
		data = " (" + strconv.Quote(t.Data) + ")"
	}

	return fmt.Sprintf("Token[%d:%s%s]", t.Pos, t.Type, data)
}

// NewToken is a convenience function to create a new token.
//
// Parameters:
//   - pos: The position of the token (in bytes).
//   - type_: The type of the token.
//   - data: The data of the token.
//
// Returns:
//   - T: The new token. Never returns nil.
func NewToken(pos int, type_, data string) *Token {
	return &Token{
		Pos:  pos,
		Type: type_,
		Data: data,
	}
}

// link_children is a helper function that links the children of the token.
//
// Parameters:
//   - parent: The parent of the children.
//   - children: The children of the token.
//
// Returns:
//   - []*Token[T]: The linked children without nil tokens.
//
// WARNING: This function modifies the input slice.
func link_children(parent *Token, children []*Token) []*Token {
	for _, tk := range children {
		tk.Parent = parent
	}

	prev := children[0]

	for _, tk := range children[1:] {
		tk.PrevSibling = prev
		prev.NextSibling = tk
		prev = tk
	}

	return children
}

// AppendChildren adds the given children tokens to the current token.
//
// Parameters:
//   - children: Variadic parameter of type Token representing the children to be added.
//
// Returns:
//   - error: Returns an error if the receiver is nil.
func (t *Token) AppendChildren(children ...*Token) error {
	gcslc.RejectNils(&children)
	if len(children) == 0 {
		return nil
	} else if t == nil {
		return common.ErrNilReceiver
	}

	children = link_children(t, children)

	if t.FirstChild == nil {
		t.FirstChild = children[0]
	} else {
		t.LastChild.NextSibling = children[0]
		children[0].PrevSibling = t.LastChild
	}

	t.LastChild = children[len(children)-1]

	return nil
}

// Child is an iterator over the children of the token that goes from the
// first child to the last child.
//
// Returns:
//   - iter.Seq[*Token[T]]: An iterator over the children of the token. Never
//     returns nil.
func (t Token) Child() iter.Seq[*Token] {
	return func(yield func(*Token) bool) {
		for child := t.FirstChild; child != nil; child = child.NextSibling {
			if !yield(child) {
				break
			}
		}
	}
}

// BackwardChild is an iterator over the children of the token that goes from
// the last child to the first child.
//
// Returns:
//   - iter.Seq[*Token[T]]: An iterator over the children of the token. Never
//     returns nil.
func (t Token) BackwardChild() iter.Seq[*Token] {
	return func(yield func(*Token) bool) {
		for child := t.LastChild; child != nil; child = child.PrevSibling {
			if !yield(child) {
				break
			}
		}
	}
}

// Equals checks whether the given token is equal to the current token.
//
// Parameters:
//   - other: The token to compare with.
//
// Returns:
//   - bool: True if the tokens are equal, false otherwise.
func (t Token) Equals(other *Token) bool {
	return other != nil && t.Data == other.Data && t.Type == other.Type
}
