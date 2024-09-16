package grammar

import (
	"errors"
	"fmt"
	"iter"
	"slices"

	"github.com/PlayerR9/tree/tree"
)

// TokenTyper is the interface that must be implemented by token types.
type TokenTyper interface {
	~int

	// String returns the string representation of the token type.
	//
	// Returns:
	//   - string: the string representation of the token type.
	String() string

	// IsTerminal checks whether the token type is terminal.
	//
	// Returns:
	//   - bool: true if the token type is terminal, false otherwise.
	IsTerminal() bool
}

// Token is a node in a tree.
type Token[T TokenTyper] struct {
	Parent, FirstChild, NextSibling, LastChild, PrevSibling *Token[T]
	Data                                                    string
	Lookahead                                               *Token[T]
	Type                                                    T
}

// IsLeaf implements the tree.Noder interface.
func (tn *Token[T]) IsLeaf() bool {
	return tn.FirstChild == nil
}

// IsSingleton implements the tree.Noder interface.
func (tn *Token[T]) IsSingleton() bool {
	return tn.FirstChild != nil && tn.FirstChild == tn.LastChild
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

	tk := &Token[T]{
		Type:      type_,
		Lookahead: last_tk.Lookahead,
	}

	tk.AddChildren(children)

	return tk, nil
}

// AddChild adds the target child to the node. Because this function clears the parent and sibling
// of the target, it does not add its relatives.
//
// Parameters:
//   - child: The child to add.
//
// If child is nil, it does nothing.
func (tn *Token[T]) AddChild(target *Token[T]) {
	if target == nil {
		return
	}

	target.NextSibling = nil
	target.PrevSibling = nil

	last_child := tn.LastChild

	if last_child == nil {
		tn.FirstChild = target
	} else {
		last_child.NextSibling = target
		target.PrevSibling = last_child
	}

	target.Parent = tn
	tn.LastChild = target
}

// BackwardChild scans the children of the node in reverse order (i.e., from the
// last child to the first one) and yields them one by one.
//
// Returns:
//   - iter.Seq[*Token[T]]: A sequence of the children of the node.
func (tn *Token[T]) BackwardChild() iter.Seq[*Token[T]] {
	return func(yield func(*Token[T]) bool) {
		for c := tn.LastChild; c != nil; c = c.PrevSibling {
			if !yield(c) {
				return
			}
		}
	}
}

// Child scans the children of the node in order (i.e., from the
// first child to the last one) and yields them one by one.
//
// Returns:
//   - iter.Seq[*Token[T]]: A sequence of the children of the node.
func (tn *Token[T]) Child() iter.Seq[*Token[T]] {
	return func(yield func(*Token[T]) bool) {
		for c := tn.FirstChild; c != nil; c = c.NextSibling {
			if !yield(c) {
				return
			}
		}
	}
}

// Cleanup cleans the node and returns its children.
// This function logically removes the node from the siblings and the parent.
//
// Finally, it is not safe to use in goroutines as pointers may be dereferenced while another
// goroutine is still using them.
//
// Returns:
//   - []*Token[T]: The children of the node.
func (tn *Token[T]) Cleanup() []*Token[T] {
	var children []*Token[T]

	for c := tn.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}

	tn.FirstChild = nil
	tn.LastChild = nil
	tn.Parent = nil

	prev := tn.PrevSibling
	next := tn.NextSibling

	if prev != nil {
		prev.NextSibling = next
	}

	if next != nil {
		next.PrevSibling = prev
	}

	tn.PrevSibling = nil
	tn.NextSibling = nil

	return children
}

// Copy creates a shally copy of the node.
//
// Although this function never returns nil, it does not copy any pointers.
func (tn *Token[T]) Copy() *Token[T] {
	return &Token[T]{
		Data:      tn.Data,
		Lookahead: tn.Lookahead,
		Type:      tn.Type,
	}
}

// delete_child is a helper function to delete the child from the children of the node. No nil
// nodes are returned when this function is called. However, if target is nil, then nothing happens.
//
// Parameters:
//   - target: The child to remove.
//
// Returns:
//   - []Token[T]: A slice of pointers to the children of the node.
func (tn *Token[T]) delete_child(target *Token[T]) []*Token[T] {
	ok := tn.HasChild(target)
	if !ok {
		return nil
	}

	prev := target.PrevSibling
	next := target.NextSibling

	if prev != nil {
		prev.NextSibling = next
	}

	if next != nil {
		next.PrevSibling = prev
	}

	if target == tn.FirstChild {
		tn.FirstChild = next

		if next == nil {
			tn.LastChild = nil
		}
	} else if target == tn.LastChild {
		tn.LastChild = prev
	}

	target.Parent = nil
	target.PrevSibling = nil
	target.NextSibling = nil

	children := target.GetChildren()

	return children
}

// DeleteChild deletes the child from the children of the node while
// returning the children of the target node.
//
// Parameters:
//   - target: The child to remove.
//
// Returns:
//   - []*Token[T]: A slice of the children of the target node.
func (tn *Token[T]) DeleteChild(target *Token[T]) []*Token[T] {
	if target == nil {
		return nil
	}

	children := tn.delete_child(target)
	if len(children) == 0 {
		return nil
	}

	for _, child := range children {
		child.PrevSibling = nil
		child.NextSibling = nil
		child.Parent = nil
	}

	tn.FirstChild = nil
	tn.LastChild = nil

	return children
}

// GetFirstChild returns the first child of the node.
//
// Returns:
//   - *Token[T]: The first child of the node.
//   - bool: True if the node has a child, false otherwise.
func (tn *Token[T]) GetFirstChild() (*Token[T], bool) {
	return tn.FirstChild, tn.FirstChild == nil
}

// GetParent returns the parent of the node.
//
// Returns:
//   - *Token[T]: The parent of the node.
//   - bool: True if the node has a parent, false otherwise.
func (tn *Token[T]) GetParent() (*Token[T], bool) {
	return tn.Parent, tn.Parent == nil
}

// LinkChildren is a method that links the children of the node.
//
// Parameters:
//   - children: The children to link.
func (tn *Token[T]) LinkChildren(children []*Token[T]) {
	var valid_children []*Token[T]

	for _, child := range children {
		if child == nil {
			continue
		}

		child.Parent = tn

		valid_children = append(valid_children, child)
	}
	if len(valid_children) == 0 {
		return
	}

	valid_children[0].PrevSibling = nil
	valid_children[len(valid_children)-1].NextSibling = nil

	if len(valid_children) == 1 {
		return
	}

	for i := 0; i < len(valid_children)-1; i++ {
		valid_children[i].NextSibling = valid_children[i+1]
	}

	for i := 1; i < len(valid_children); i++ {
		valid_children[i].PrevSibling = valid_children[i-1]
	}

	tn.FirstChild, tn.LastChild = valid_children[0], valid_children[len(valid_children)-1]
}

// RemoveNode removes the node from the tree while shifting the children up one level to
// maintain the tree structure. The returned children can be used to create a forest of
// trees if the root node is removed.
//
// Returns:
//   - []*Token[T]: A slice of pointers to the children of the node iff the node is the root.
//
// Example:
//
//	// Given the tree:
//	1
//	├── 2
//	├── 3
//	|	├── 4
//	|	└── 5
//	└── 6
//
//	// The tree after removing node 3:
//
//	1
//	├── 2
//	├── 4
//	├── 5
//	└── 6
func (tn *Token[T]) RemoveNode() []*Token[T] {
	prev := tn.PrevSibling
	next := tn.NextSibling
	parent := tn.Parent

	var sub_roots []*Token[T]

	if parent == nil {
		for c := tn.FirstChild; c != nil; c = c.NextSibling {
			sub_roots = append(sub_roots, c)
		}
	} else {
		children := parent.delete_child(tn)

		for _, child := range children {
			child.Parent = parent
		}
	}

	if prev != nil {
		prev.NextSibling = next
	} else {
		parent.FirstChild = next
	}

	if next != nil {
		next.PrevSibling = prev
	} else {
		parent.Parent.LastChild = prev
	}

	tn.Parent = nil
	tn.PrevSibling = nil
	tn.NextSibling = nil

	if len(sub_roots) == 0 {
		return nil
	}

	for _, child := range sub_roots {
		child.PrevSibling = nil
		child.NextSibling = nil
		child.Parent = nil
	}

	tn.FirstChild = nil
	tn.LastChild = nil

	return sub_roots
}

// AddChildren is a convenience function to add multiple children to the node at once.
// It is more efficient than adding them one by one. Therefore, the behaviors are the
// same as the behaviors of the Token.AddChild function.
//
// Parameters:
//   - children: The children to add.
func (tn *Token[T]) AddChildren(children []*Token[T]) {
	if len(children) == 0 {
		return
	}

	var top int

	for i := 0; i < len(children); i++ {
		child := children[i]

		if child != nil {
			children[top] = child
			top++
		}
	}

	children = children[:top]
	if len(children) == 0 {
		return
	}

	// Deal with the first child
	first_child := children[0]

	first_child.NextSibling = nil
	first_child.PrevSibling = nil

	last_child := tn.LastChild

	if last_child == nil {
		tn.FirstChild = first_child
	} else {
		last_child.NextSibling = first_child
		first_child.PrevSibling = last_child
	}

	first_child.Parent = tn
	tn.LastChild = first_child

	// Deal with the rest of the children
	for i := 1; i < len(children); i++ {
		child := children[i]

		child.NextSibling = nil
		child.PrevSibling = nil

		last_child := tn.LastChild
		last_child.NextSibling = child
		child.PrevSibling = last_child

		child.Parent = tn
		tn.LastChild = child
	}
}

// GetChildren returns the immediate children of the node.
//
// The returned nodes are never nil and are not copied. Thus, modifying the returned
// nodes will modify the tree.
//
// Returns:
//   - []*Token[T]: A slice of pointers to the children of the node.
func (tn *Token[T]) GetChildren() []*Token[T] {
	var children []*Token[T]

	for c := tn.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}

	return children
}

// HasChild returns true if the node has the given child.
//
// Because children of a node cannot be nil, a nil target will always return false.
//
// Parameters:
//   - target: The child to check for.
//
// Returns:
//   - bool: True if the node has the child, false otherwise.
func (tn *Token[T]) HasChild(target *Token[T]) bool {
	if target == nil || tn.FirstChild == nil {
		return false
	}

	for c := tn.FirstChild; c != nil; c = c.NextSibling {
		if c == target {
			return true
		}
	}

	return false
}

// IsChildOf returns true if the node is a child of the parent. If target is nil,
// it returns false.
//
// Parameters:
//   - target: The target parent to check for.
//
// Returns:
//   - bool: True if the node is a child of the parent, false otherwise.
func (tn *Token[T]) IsChildOf(target *Token[T]) bool {
	if target == nil {
		return false
	}

	parents := tree.GetNodeAncestors(target)

	for node := tn; node.Parent != nil; node = node.Parent {
		ok := slices.Contains(parents, node.Parent)
		if ok {
			return true
		}
	}

	return false
}
