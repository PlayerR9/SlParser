package grammar

import (
	"iter"

	gcslc "github.com/PlayerR9/go-commons/slices"
	gcers "github.com/PlayerR9/go-errors"
)

// ParseTree is a generic data structure that represents a tree.
type ParseTree[T TokenTyper] struct {
	// root is the root of the tree.
	root *Token[T]

	// leaves is the leaves of the tree.
	leaves []*Token[T]

	// size is the number of nodes in the tree.
	size int
}

func (t ParseTree[T]) Pos() int {
	gcers.AssertNotNil(t.root, "t.root")

	return t.root.Pos
}

func (t ParseTree[T]) Type() T {
	gcers.AssertNotNil(t.root, "t.root")

	return t.root.Type
}

func (t ParseTree[T]) Lookahead() *Token[T] {
	gcers.AssertNotNil(t.root, "t.root")

	return t.root.Lookahead
}

func (t ParseTree[T]) Data() string {
	gcers.AssertNotNil(t.root, "t.root")

	return t.root.Data
}

// Cleanup is a method that cleans up the tree.
func (t *ParseTree[T]) Cleanup() {
	if t == nil {
		return
	}

	Cleanup(t.root)

	t.size = 1
	t.leaves = []*Token[T]{t.root}
}

// DeepCopy is a method that deeply copies the tree.
//
// Returns:
//   - *Tree: A copy of the tree. Never returns nil.
func (t ParseTree[T]) DeepCopy() *ParseTree[T] {
	var tree *ParseTree[T]

	root_copy := DeepCopy(t.root)

	tree = &ParseTree[T]{
		root:   root_copy,
		leaves: GetNodeLeaves(root_copy),
		size:   t.size,
	}

	return tree
}

// String implements the fmt.Stringer interface.
//
// Format:
//
//	root
//	├── node1
//	│   ├── node2
//	│   └── node3
//	└── node4
//	|   └── node5
//	|
//	| // ...
//	// ...
func (t ParseTree[T]) String() string {
	trav := PrintFn[T]()

	info, err := ApplyDFS(&t, trav)
	if err != nil {
		panic(err.Error())
	}

	return info.String()
}

// NewTree creates a new tree from the given root.
//
// Parameters:
//   - root: The root of the tree.
//
// Returns:
//   - *ParseTree[T]: A pointer to the newly created tree.
//   - error: An error of type *errors.ErrInvalidParameter if the root is nil.
func NewTree[T TokenTyper](root *Token[T]) (*ParseTree[T], error) {
	if root == nil {
		return nil, gcers.NewErrNilParameter("root")
	}

	stack := []*Token[T]{root}
	size := 1

	var leaves []*Token[T]

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		size++

		if top.IsLeaf() {
			leaves = append(leaves, top)
		} else {
			for child := range top.Child() {
				stack = append(stack, child)
			}
		}
	}

	return &ParseTree[T]{
		root:   root,
		leaves: leaves,
		size:   size,
	}, nil
}

// Root returns the root of the tree.
//
// Returns:
//   - *Token[T]: The root of the tree. Never returns nil.
func (t ParseTree[T]) Root() *Token[T] {
	gcers.AssertNotNil(t.root, "t.root")

	return t.root
}

// Leaves returns the leaves of the tree.
//
// Returns:
//   - []*Token[T]: The leaves of the tree. Never returns nil.
func (t ParseTree[T]) Leaves() []*Token[T] {
	gcers.AssertNotNil(t.root, "t.root")

	return t.leaves
}

// Size returns the number of nodes in the tree.
//
// Returns:
//   - int: The number of nodes in the tree.
func (t ParseTree[T]) Size() int {
	return t.size
}

// SetChildren sets the children of the root of the tree.
//
// Parameters:
//   - children: The children to set.
//
// Returns:
//   - error: An error of type *ErrMissingRoot if the tree does not have a root.
func (t *ParseTree[T]) SetChildren(children []*ParseTree[T]) error {
	if t == nil {
		return nil
	}

	children = gcslc.FilterZeroValues(children)
	if len(children) == 0 {
		return nil
	}

	root := t.root

	var leaves []*Token[T]
	var sub_children []*Token[T]

	t.size = 1

	for _, child := range children {
		leaves = append(leaves, child.leaves...)
		t.size += child.Size()

		sub_children = append(sub_children, child.root)
	}

	root.LinkChildren(sub_children)

	t.leaves = leaves

	return nil
}

// GetChildren returns the direct children of the root of the tree.
//
// Children are never nil.
//
// Returns:
//   - []*ParseTree[T]: A slice of the direct children of the root. Nil if the tree does not have a root.
func (t ParseTree[T]) GetChildren() []*ParseTree[T] {
	var children []*ParseTree[T]

	for child := range t.root.Child() {
		tree, err := NewTree(child)
		gcers.AssertErr(err, "NewTree(child)")

		children = append(children, tree)
	}

	return children
}

// DFS applies the DFS traversal logic to the tree.
//
// Returns:
//   - iter.Seq[*Token[T]]: The traversal sequence.
func (t ParseTree[T]) DFS() iter.Seq[*Token[T]] {
	fn := func(yield func(*Token[T]) bool) {
		stack := []*Token[T]{t.root}

		for len(stack) > 0 {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if !yield(top) {
				break
			}

			for child := range top.BackwardChild() {
				stack = append(stack, child)
			}
		}
	}

	return fn
}

// BFS applies the BFS traversal logic to the tree.
//
// Returns:
//   - iter.Seq[*Token[T]]: The traversal sequence.
func (t ParseTree[T]) BFS() iter.Seq[*Token[T]] {
	fn := func(yield func(*Token[T]) bool) {
		queue := []*Token[T]{t.root}

		for len(queue) > 0 {
			top := queue[0]
			queue = queue[1:]

			if !yield(top) {
				break
			}

			for child := range top.Child() {
				queue = append(queue, child)
			}
		}
	}

	return fn
}

// RegenerateLeaves regenerates the leaves of the tree. No op if the tree is nil.
//
// Behaviors:
//   - The leaves are updated in a DFS order.
//   - Expensive operation; use it only when necessary (i.e., leaves changed unexpectedly.)
//   - This also updates the size of the tree.
func (tree *ParseTree[T]) RegenerateLeaves() {
	if tree == nil {
		return
	}

	var leaves []*Token[T]
	var size int

	for node := range tree.DFS() {
		size++

		if node.IsLeaf() {
			leaves = append(leaves, node)
		}
	}

	tree.leaves = leaves
	tree.size = size
}

// UpdateLeaves updates the leaves of the tree. No op if the tree is nil.
//
// Behaviors:
//   - The leaves are updated in a DFS order.
//   - Less expensive than RegenerateLeaves. However, if nodes has been deleted
//     from the tree, this may give unexpected results.
//   - This also updates the size of the tree.
func (tree *ParseTree[T]) UpdateLeaves() {
	if tree == nil {
		return
	}

	if len(tree.leaves) == 0 {
		tree.leaves = []*Token[T]{tree.root}
		tree.size = 1

		return
	}

	var new_leaves []*Token[T]
	size := tree.size - len(tree.leaves)

	stack := tree.leaves

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		size++

		ok := top.IsLeaf()
		if ok {
			new_leaves = append(new_leaves, top)
		}
	}

	tree.leaves = new_leaves
	tree.size = size
}

// HasChild returns true if the tree has the given child in any of its nodes
// in a BFS order.
//
// Parameters:
//   - filter: The filter to apply. Must return true iff the node is the one we are looking for.
//     This function must assume node is never nil.
//
// Returns:
//   - bool: True if the tree has the child, false otherwise.
//
// If either tree or filter is nil, false is returned.
func (tree ParseTree[T]) HasChild(filter func(node *Token[T]) bool) bool {
	if filter == nil {
		return false
	}

	for node := range tree.BFS() {
		if filter(node) {
			return true
		}
	}

	return false
}

// FilterChildren returns all the children of the tree that satisfy the given filter
// in a BFS order.
//
// Parameters:
//   - filter: The filter to apply. Must return true iff the node is the one we want to keep.
//     This function must assume node is never nil.
//
// Returns:
//   - []*Token[T]: A slice of the children that satisfy the filter.
//
// If either tree or filter is nil, an empty slice and false are returned.
func (tree ParseTree[T]) FilterChildren(filter func(node *Token[T]) bool) []*Token[T] {
	if filter == nil {
		return nil
	}

	var children []*Token[T]

	for node := range tree.DFS() {
		if filter(node) {
			children = append(children, node)
		}
	}

	return children
}

// SearchNodes searches for the first node that satisfies the given filter in a BFS order.
//
// Parameters:
//   - tree: The tree to search.
//   - filter: The filter to apply. Must return true iff the node is the one we are looking for.
//     This function must assume node is never nil.
//
// Returns:
//   - *Token[T]: The node that satisfies the filter.
//   - bool: True if the node was found, false otherwise.
//
// Nodes that are not of type T will be ignored. If either tree or filter is nil, false is returned.
func (tree ParseTree[T]) SearchNodes(filter func(node *Token[T]) bool) (*Token[T], bool) {
	if filter == nil {
		return nil, false
	}

	for node := range tree.BFS() {
		if filter(node) {
			return node, true
		}
	}

	return nil, false
}

// rec_snake_traversal is an helper function that returns all the paths
// from n to the leaves of the tree rooted at n.
//
// Returns:
//   - [][]*Token[T]: A slice of slices of elements.
//
// Behaviors:
//   - The paths are returned in the order of a BFS traversal.
//   - It is a recursive function.
func rec_snake_traversal[T TokenTyper](n *Token[T]) [][]*Token[T] {
	ok := n.IsLeaf()
	if ok {
		return [][]*Token[T]{
			{n},
		}
	}

	var result [][]*Token[T]

	for child := range n.Child() {
		subResults := rec_snake_traversal(child)

		for _, tmp := range subResults {
			tmp = append([]*Token[T]{n}, tmp...)
			result = append(result, tmp)
		}
	}

	return result
}

// SnakeTraversal returns all the paths from the root to the leaves of the tree.
//
// Returns:
//   - [][]*Token[T]: A slice of slices of elements. Nil if the tree is empty.
//
// Behaviors:
//   - The paths are returned in the order of a BFS traversal.
func (tree ParseTree[T]) SnakeTraversal() [][]*Token[T] {
	sol := rec_snake_traversal(tree.root)
	return sol
}

// replaceLeafWithTree is a helper function that replaces a leaf with a tree.
//
// Parameters:
//   - tree: The tree to replace.
//   - at: The index of the leaf to replace.
//   - children: The children of the leaf.
//
// Behaviors:
//   - The leaf is replaced with the children.
//   - The size of the tree is updated.
func (tree *ParseTree[T]) replaceLeafWithTree(at int, values []*Token[T]) {
	// uc.AssertParam("at", at >= 0 && at < len(tree.leaves), uc.NewErrOutOfBounds(at, 0, len(tree.leaves)))

	leaf := tree.leaves[at]

	// Make the subtree
	leaf.LinkChildren(values)

	// Update the size of the tree
	tree.size += len(values) - 1

	// Replace the current leaf with the leaf's children
	sub_leaves := GetNodeLeaves(leaf)

	if at == len(tree.leaves)-1 {
		tree.leaves = append(tree.leaves[:at], sub_leaves...)
	} else if at == 0 {
		tree.leaves = append(sub_leaves, tree.leaves[at+1:]...)
	} else {
		tree.leaves = append(tree.leaves[:at], append(sub_leaves, tree.leaves[at+1:]...)...)
	}
}

// ProcessLeaves applies the given function to the leaves of the tree and replaces
// the leaves with the children returned by the function.
//
// Parameters:
//   - f: The function to apply to the leaves.
//
// Returns:
//   - error: An error returned by the function.
//
// Behaviors:
//   - The function is applied to the leaves in order.
//   - The function must return a slice of values of type T.
//   - If the function returns an error, the process stops and the error is returned.
//   - The leaves are replaced with the children returned by the function.
func (tree *ParseTree[T]) ProcessLeaves(f func(node *Token[T]) ([]*Token[T], error)) error {
	if tree == nil || f == nil {
		return nil
	}

	for i, leaf := range tree.leaves {
		children, err := f(leaf)
		if err != nil {
			return err
		}

		if len(children) != 0 {
			conv := make([]*Token[T], 0, len(children))

			for _, child := range children {
				conv = append(conv, child)
			}

			tree.replaceLeafWithTree(i, conv)
		}
	}

	return nil
}
