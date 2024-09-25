package parsing

import (
	"iter"
	"slices"
	"strconv"
	"strings"
)

// Node is a node in a tree.
type Node struct {
	Parent, FirstChild, NextSibling, LastChild, PrevSibling *Node
	Data                                                    string
	Type                                                    NodeType
	Pos                                                     int
}

// IsNil implements the ast.Noder interface.
func (n *Node) IsNil() bool {
	return n == nil
}

// IsLeaf implements the tree.Noder interface.
func (n Node) IsLeaf() bool {
	return n.FirstChild == nil
}

// IsSingleton implements the tree.Noder interface.
func (n Node) IsSingleton() bool {
	return n.FirstChild != nil && n.FirstChild == n.LastChild
}

// String implements the tree.Noder interface.
func (n Node) String() string {
	var builder strings.Builder

	builder.WriteString("Node[")
	builder.WriteString(strconv.Itoa(n.Pos))
	builder.WriteString(": ")
	builder.WriteString(n.Type.String())

	if n.Data != "" {
		builder.WriteString(" (")
		builder.WriteString(n.Data)
		builder.WriteRune(')')
	}

	builder.WriteRune(']')

	return builder.String()
}

// NewNode creates a new node with the given data.
//
// Parameters:
//   - pos: The position of the node.
//   - type_: The type of the node.
//   - data: The data of the node.
//
// Returns:
//   - *Node: A pointer to the newly created node. It is
//     never nil.
func NewNode(pos int, type_ NodeType, data string) *Node {
	return &Node{
		Pos:  pos,
		Data: data,
		Type: type_,
	}
}

// AddChild adds the target child to the node. Because this function clears the parent and sibling
// of the target, it does not add its relatives.
//
// Parameters:
//   - target: The child to add.
//
// If the receiver or the target are nil, it does nothing.
func (n *Node) AddChild(target *Node) {
	if n == nil || target == nil {
		return
	}

	target.NextSibling = nil
	target.PrevSibling = nil

	last_child := n.LastChild

	if last_child == nil {
		n.FirstChild = target
	} else {
		last_child.NextSibling = target
		target.PrevSibling = last_child
	}

	target.Parent = n
	n.LastChild = target
}

// BackwardChild scans the children of the node in reverse order (i.e., from the
// last child to the first one) and yields them one by one.
//
// Returns:
//   - iter.Seq[*Node]: A sequence of the children of the node.
func (n Node) BackwardChild() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for c := n.LastChild; c != nil; c = c.PrevSibling {
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
//   - iter.Seq[*Node]: A sequence of the children of the node.
func (n Node) Child() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
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
//   - []*Node: The children of the node.
func (n *Node) Cleanup() []*Node {
	if n == nil {
		return nil
	}

	var children []*Node

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}

	n.FirstChild = nil
	n.LastChild = nil
	n.Parent = nil

	prev := n.PrevSibling
	next := n.NextSibling

	if prev != nil {
		prev.NextSibling = next
	}

	if next != nil {
		next.PrevSibling = prev
	}

	n.PrevSibling = nil
	n.NextSibling = nil

	return children
}

// Copy creates a shally copy of the node.
//
// Although this function never returns nil, it does not copy any pointers.
func (n Node) Copy() *Node {
	return &Node{
		Pos:  n.Pos,
		Data: n.Data,
		Type: n.Type,
	}
}

// delete_child is a helper function to delete the child from the children of the node. No nil
// nodes are returned when this function is called. However, if target is nil, then nothing happens.
//
// Parameters:
//   - target: The child to remove.
//
// Returns:
//   - []Node: A slice of pointers to the children of the node.
func (n *Node) delete_child(target *Node) []*Node {
	if n == nil {
		return nil
	}

	ok := n.HasChild(target)
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

	if target == n.FirstChild {
		n.FirstChild = next

		if next == nil {
			n.LastChild = nil
		}
	} else if target == n.LastChild {
		n.LastChild = prev
	}

	target.Parent = nil
	target.PrevSibling = nil
	target.NextSibling = nil

	var children []*Node

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}

	return children
}

// DeleteChild deletes the child from the children of the node while
// returning the children of the target node.
//
// Parameters:
//   - target: The child to remove.
//
// Returns:
//   - []*Node: A slice of the children of the target node.
func (n *Node) DeleteChild(target *Node) []*Node {
	if n == nil || target == nil {
		return nil
	}

	children := n.delete_child(target)
	if len(children) == 0 {
		return nil
	}

	for _, child := range children {
		child.PrevSibling = nil
		child.NextSibling = nil
		child.Parent = nil
	}

	n.FirstChild = nil
	n.LastChild = nil

	return children
}

// GetFirstChild returns the first child of the node.
//
// Returns:
//   - *Node: The first child of the node.
//   - bool: True if the node has a child, false otherwise.
func (n Node) GetFirstChild() (*Node, bool) {
	return n.FirstChild, n.FirstChild == nil
}

// GetParent returns the parent of the node.
//
// Returns:
//   - *Node: The parent of the node.
//   - bool: True if the node has a parent, false otherwise.
func (n Node) GetParent() (*Node, bool) {
	return n.Parent, n.Parent == nil
}

// LinkChildren is a method that links the children of the node.
//
// Parameters:
//   - children: The children to link.
func (n *Node) LinkChildren(children []*Node) {
	if n == nil {
		return
	}

	var valid_children []*Node

	for _, child := range children {
		if child == nil {
			continue
		}

		child.Parent = n

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

	n.FirstChild, n.LastChild = valid_children[0], valid_children[len(valid_children)-1]
}

// RemoveNode removes the node from the tree while shifting the children up one level to
// maintain the tree structure. The returned children can be used to create a forest of
// trees if the root node is removed.
//
// Returns:
//   - []*Node: A slice of pointers to the children of the node iff the node is the root.
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
func (n *Node) RemoveNode() []*Node {
	if n == nil {
		return nil
	}

	prev := n.PrevSibling
	next := n.NextSibling
	parent := n.Parent

	var sub_roots []*Node

	if parent == nil {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			sub_roots = append(sub_roots, c)
		}
	} else {
		children := parent.delete_child(n)

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

	n.Parent = nil
	n.PrevSibling = nil
	n.NextSibling = nil

	if len(sub_roots) == 0 {
		return nil
	}

	for _, child := range sub_roots {
		child.PrevSibling = nil
		child.NextSibling = nil
		child.Parent = nil
	}

	n.FirstChild = nil
	n.LastChild = nil

	return sub_roots
}

// AddChildren is a convenience function to add multiple children to the node at once.
// It is more efficient than adding them one by one. Therefore, the behaviors are the
// same as the behaviors of the Node.AddChild function.
//
// Parameters:
//   - children: The children to add.
func (n *Node) AddChildren(children []*Node) {
	if n == nil || len(children) == 0 {
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

	last_child := n.LastChild

	if last_child == nil {
		n.FirstChild = first_child
	} else {
		last_child.NextSibling = first_child
		first_child.PrevSibling = last_child
	}

	first_child.Parent = n
	n.LastChild = first_child

	// Deal with the rest of the children
	for i := 1; i < len(children); i++ {
		child := children[i]

		child.NextSibling = nil
		child.PrevSibling = nil

		last_child := n.LastChild
		last_child.NextSibling = child
		child.PrevSibling = last_child

		child.Parent = n
		n.LastChild = child
	}
}

// GetChildren returns the immediate children of the node.
//
// The returned nodes are never nil and are not copied. Thus, modifying the returned
// nodes will modify the tree.
//
// Returns:
//   - []*Node: A slice of pointers to the children of the node.
func (n Node) GetChildren() []*Node {
	var children []*Node

	for c := n.FirstChild; c != nil; c = c.NextSibling {
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
func (n Node) HasChild(target *Node) bool {
	if target == nil || n.FirstChild == nil {
		return false
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
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
func (n Node) IsChildOf(target *Node) bool {
	if target == nil {
		return false
	}

	var ancestors []*Node

	for target.Parent != nil {
		parent := target.Parent
		ancestors = append(ancestors, parent)
		target = parent
	}

	slices.Reverse(ancestors)

	for node := &n; node.Parent != nil; node = node.Parent {
		ok := slices.Contains(ancestors, node.Parent)
		if ok {
			return true
		}
	}

	return false
}

// Pair is a pair of a node and its info.
type Pair struct {
	// Node is the node of the pair.
	Node *Node

	// Info is the info of the pair.
	Info any
}

// NewPair creates a new pair of a node and its info.
//
// Parameters:
//   - node: The node of the pair.
//   - info: The info of the pair.
//
// Returns:
//   - Pair: The new pair.
func NewPair(node *Node, info any) Pair {
	return Pair{
		Node: node,
		Info: info,
	}
}

// Traverser is the traverser that holds the traversal logic.
type Traverser struct {
	// InitFn is the function that initializes the traversal info.
	//
	// Parameters:
	//   - root: The root node of the tree.
	//
	// Returns:
	//   - any: The initial traversal info.
	InitFn func(root *Node) any

	// DoFn is the function that performs the traversal logic.
	//
	// Parameters:
	//   - node: The current node of the tree.
	//   - info: The traversal info.
	//
	// Returns:
	//   - []Pair: The next traversal info.
	//   - error: The error that might occur during the traversal.
	DoFn func(node *Node, info any) ([]Pair, error)
}

// ApplyDFS applies the DFS traversal logic to the tree.
//
// Parameters:
//   - t: The tree to apply the traversal logic to.
//   - trav: The traverser that holds the traversal logic.
//
// Returns:
//   - any: The final traversal info.
//   - error: The error that might occur during the traversal.
func ApplyDFS(root *Node, trav Traverser) (any, error) {
	if root == nil {
		return nil, nil
	}

	info := trav.InitFn(root)

	stack := []Pair{NewPair(root, info)}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		nexts, err := trav.DoFn(top.Node, top.Info)
		if err != nil {
			return info, err
		}

		if len(nexts) > 0 {
			slices.Reverse(nexts)
			stack = append(stack, nexts...)
		}
	}

	return info, nil
}

// ApplyBFS applies the BFS traversal logic to the tree.
//
// Parameters:
//   - t: The tree to apply the traversal logic to.
//   - trav: The traverser that holds the traversal logic.
//
// Returns:
//   - any: The final traversal info.
//   - error: The error that might occur during the traversal.
func ApplyBFS(root *Node, trav Traverser) (any, error) {
	if root == nil {
		return nil, nil
	}

	info := trav.InitFn(root)

	queue := []Pair{NewPair(root, info)}

	for len(queue) > 0 {
		top := queue[0]
		queue = queue[1:]

		nexts, err := trav.DoFn(top.Node, top.Info)
		if err != nil {
			return info, err
		}

		if len(nexts) > 0 {
			queue = append(queue, nexts...)
		}
	}

	return info, nil
}

// _AstStringTrav is the stack element of the tree stringer.
type _AstStringTrav struct {
	// seen is the seen map of the tree stringer.
	seen map[*Node]struct{}

	// builder is the builder of the tree stringer.
	builder *strings.Builder

	// indent is the indentation string.
	indent string

	// is_last is the flag that indicates whether the node is the last node in the level.
	is_last bool

	// same_level is the flag that indicates whether the node is in the same level.
	same_level bool
}

// String implements the fmt.Stringer interface.
func (tse _AstStringTrav) String() string {
	str := tse.builder.String()
	if str != "" {
		str = strings.TrimSuffix(str, "\n")
	}

	return str
}

// set_is_last is a helper function that sets the is_last flag.
//
// Assume that the receiver is not nil.
func (tse *_AstStringTrav) set_is_last() {
	tse.is_last = true
}

// set_same_level is a helper function that sets the same_level flag.
//
// Assume that the receiver is not nil.
func (tse *_AstStringTrav) set_same_level() {
	tse.same_level = true
}

var (
	printer_trav Traverser
)

func init() {
	init_fn := func(root *Node) any {
		var builder strings.Builder

		return &_AstStringTrav{
			seen:       make(map[*Node]struct{}),
			builder:    &builder,
			indent:     "",
			is_last:    true,
			same_level: false,
		}
	}

	fn := func(node *Node, info any) ([]Pair, error) {
		inf := info.(*_AstStringTrav)

		if inf.indent != "" {
			inf.builder.WriteString(inf.indent)

			if !node.IsLeaf() || inf.is_last {
				inf.builder.WriteString("└── ")
			} else {
				inf.builder.WriteString("├── ")
			}
		}

		// Prevent cycles.
		_, ok := inf.seen[node]
		if ok {
			inf.builder.WriteString("... WARNING: Cycle detected!\n")

			return nil, nil
		}

		inf.builder.WriteString(node.String())
		inf.builder.WriteString("\n")

		inf.seen[node] = struct{}{}

		if node.IsLeaf() {
			return nil, nil
		}

		var indent strings.Builder

		indent.WriteString(inf.indent)

		if inf.same_level && !inf.is_last {
			indent.WriteString("│   ")
		} else {
			indent.WriteString("    ")
		}

		var elems []Pair

		for c := range node.Child() {
			se := &_AstStringTrav{
				seen:       inf.seen,
				builder:    inf.builder,
				indent:     indent.String(),
				is_last:    false,
				same_level: false,
			}

			elems = append(elems, NewPair(c, se))
		}

		if len(elems) >= 2 {
			for i := 0; i < len(elems); i++ {
				elems[i].Info.(*_AstStringTrav).set_same_level()
			}
		}

		elems[len(elems)-1].Info.(*_AstStringTrav).set_is_last()

		return elems, nil
	}

	printer_trav = Traverser{
		InitFn: init_fn,
		DoFn:   fn,
	}
}

// PrintAst returns the string representation of the AST.
//
// Parameters:
//   - root: The root node of the AST.
//
// Returns:
//   - string: The string representation of the AST.
func PrintAst(root *Node) string {
	if root == nil {
		return ""
	}

	info, err := ApplyDFS(root, printer_trav)
	if err != nil {
		panic(err.Error())
	}

	print_info := info.(*_AstStringTrav)
	return print_info.String()
}
