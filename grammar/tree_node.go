package grammar

import (
	"slices"
)

// DeepCopy is a method that deep copies the node.
//
// Parameters:
//   - node: The node to copy.
//
// Returns:
//   - *Token[T]: The copied node.
func DeepCopy[T TokenTyper](node *Token[T]) *Token[T] {
	n := node.Copy()

	var children []*Token[T]

	for child := range node.Child() {
		child_copy := DeepCopy(child)
		children = append(children, child_copy)
	}

	n.LinkChildren(children)

	return n
}

// RootOf returns the root of the given node.
//
// Parameters:
//   - node: The node to get the root of.
//
// Returns:
//   - *Token[T]: The root of the given node.
func RootOf[T TokenTyper](node *Token[T]) *Token[T] {
	for {
		parent, ok := node.GetParent()
		if !ok {
			break
		}

		node = parent
	}

	return node
}

// GetNodeLeaves returns the leaves of the given node.
//
// This is expensive as leaves are not stored and so, every time this function is called,
// it has to do a DFS traversal to find the leaves. Thus, it is recommended to call
// this function once and then store the leaves somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, no nil nodes are returned.
func GetNodeLeaves[T TokenTyper](node *Token[T]) []*Token[T] {
	var leaves []*Token[T]

	stack := []*Token[T]{node}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if top.IsLeaf() {
			leaves = append(leaves, top)
		} else {
			for child := range top.Child() {
				stack = append(stack, child)
			}
		}
	}

	return leaves
}

// Size implements the *TreeNode[T] interface.
//
// This is expensive as it has to traverse the whole tree to find the size of the tree.
// Thus, it is recommended to call this function once and then store the size somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, the traversal is done in a depth-first manner.
func GetNodeSize[T TokenTyper](node *Token[T]) int {
	var size int

	stack := []*Token[T]{node}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		size++

		for child := range top.Child() {
			stack = append(stack, child)
		}
	}

	return size
}

// GetAncestors is used to get all the ancestors of the given node. This excludes
// the node itself.
//
// Parameters:
//   - node: The node to get the ancestors of.
//
// Returns:
//   - []T: The ancestors of the node.
//
// This is expensive since ancestors are not stored and so, every time this
// function is called, it has to traverse the tree to find the ancestors. Thus, it is
// recommended to call this function once and then store the ancestors somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, no nil nodes are returned.
func GetNodeAncestors[T TokenTyper](node *Token[T]) []*Token[T] {
	var ancestors []*Token[T]

	for {
		parent, ok := node.GetParent()
		if !ok {
			break
		}

		ancestors = append(ancestors, parent)

		node = parent
	}

	slices.Reverse(ancestors)

	return ancestors
}

// FindCommonAncestor returns the first common ancestor of the two nodes.
//
// This function is expensive as it calls GetNodeAncestors two times.
//
// Parameters:
//   - n1: The first node.
//   - n2: The second node.
//
// Returns:
//   - *Token[T]: The common ancestor.
//   - bool: True if the nodes have a common ancestor, false otherwise.
func FindCommonAncestor[T TokenTyper](n1, n2 *Token[T]) (*Token[T], bool) {
	if n1 == n2 {
		return n1, true
	}

	ancestors1 := GetNodeAncestors(n1)
	ancestors2 := GetNodeAncestors(n2)

	if len(ancestors1) > len(ancestors2) {
		ancestors1, ancestors2 = ancestors2, ancestors1
	}

	for _, node := range ancestors1 {
		if slices.Contains(ancestors2, node) {
			return node, true
		}
	}

	return nil, false
}

// Cleanup is used to delete all the children of the given node.
//
// Parameters:
//   - node: The node to delete the children of.
func Cleanup[T TokenTyper](node *Token[T]) {
	queue := node.Cleanup()

	for len(queue) > 0 {
		first := queue[0]
		queue = queue[1:]

		queue = append(queue, first.Cleanup()...)
	}
}
