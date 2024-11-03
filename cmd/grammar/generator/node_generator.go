package generator

import (
	common "github.com/PlayerR9/mygo-lib/common"
	cgen "github.com/PlayerR9/mygo-lib/generator"
)

type NodeData struct {
	PackageName string
}

func (nd *NodeData) SetPkgName(pkg string) error {
	if nd == nil {
		return common.ErrNilReceiver
	}

	nd.PackageName = pkg

	return nil
}

func NewNodeData() *NodeData {
	return &NodeData{}
}

var (
	NodeGenerator *cgen.CodeGenerator[*NodeData]
)

func init() {
	NodeGenerator = cgen.Must(cgen.New[*NodeData]("node", node_templ))
}

const node_templ string = `package {{ .PackageName }}

import (
	"fmt"
	"io"
	"iter"

	"github.com/PlayerR9/mygo-lib/common"
	"github.com/PlayerR9/mygo-lib/slices"
	"github.com/PlayerR9/mygo-lib/trees"
)

// Node is the node in the tree.
type Node struct {
	// Parent, FirstChild, LastChild, NextSibling, and PrevSibling are pointers of
	// the node.
	Parent, FirstChild, LastChild, NextSibling, PrevSibling *Node

	// Pos is the position in the source code.
	Pos int

	// Type is the type of the node.
	Type NodeType

	// Data is the data of the node.
	Data string
}

// IsNil implements the TreeNoder interface.
func (n *Node) IsNil() bool {
	return n == nil
}

// IsLeaf implements the TreeNoder interface.
func (n Node) IsLeaf() bool {
	return n.FirstChild == nil
}

// String implements the TreeNoder interface.
func (n Node) String() string {
	if n.Data != "" {
		return fmt.Sprintf("Node[%d:%s (%q)]", n.Pos, n.Type.String(), n.Data)
	} else {
		return fmt.Sprintf("Node[%d:%s]", n.Pos, n.Type.String())
	}
}

// NewNode creates a new node with the given position, type, and data.
//
// Parameters:
//   - pos: The position of the node.
//   - t: The type of the node.
//   - data: The data of the node.
//
// Returns:
//   - *Node: The new node. Never returns nil.
func NewNode(pos int, t NodeType, data string) *Node {
	return &Node{
		Pos:  pos,
		Type: t,
		Data: data,
	}
}

// AreChildrenCritical checks whether the children of the node are critical or not.
//
// A node is said to be critical if, when an error occurs, the entire process
// is immediately returned instead of continuing.
//
// Returns:
//		- bool: True if the children are critical, false otherwise.
func (n Node) AreChildrenCritical() bool {
	return n.Type.AreChildrenCritical()
}

// GetType returns the type of the node.
//
// Returns:
//   - NodeType: The type of the node.
func (n Node) GetType() NodeType {
	return n.Type
}

// Child iterates over the children of the node from the first child to the last child.
//
// Returns:
//   - iter.Seq[*Node]: An iterator over the children of the node. Never returns nil.
func (n Node) Child() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			if !yield(child) {
				break
			}
		}
	}
}

// BackwardChild is an iterator over the children of the node that goes from the
// last child to the first child.
//
// Returns:
//   - iter.Seq[*Node]: An iterator over the children of the node. Never returns nil.
func (n Node) BackwardChild() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for child := n.LastChild; child != nil; child = child.PrevSibling {
			if !yield(child) {
				break
			}
		}
	}
}

// link_nodes links the given children nodes to the specified parent node,
// setting up the parent, next sibling, and previous sibling pointers.
//
// Parameters:
//   - parent: The parent node to which the children will be linked.
//   - children: A slice of nodes to be linked as children.
//
// Returns:
//   - []*Node: The linked children nodes, excluding any nil nodes.
func link_nodes(parent *Node, children []*Node) []*Node {
	slices.RejectNils(&children)
	if len(children) == 0 {
		return nil
	}

	for _, c := range children {
		c.Parent = parent
	}

	prev := children[0]

	for _, c := range children[1:] {
		prev.NextSibling = c
		c.PrevSibling = prev
		prev = c
	}

	return children
}

// PrependChildren adds the given children nodes to the beginning of the current node's children list.
//
// Parameters:
//   - children: Variadic parameter of type Node representing the children to be added.
//
// Returns:
//   - error: Returns an error if the operation fails or if the receiver is nil, otherwise returns nil.
func (n *Node) PrependChildren(children ...*Node) error {
	children = link_nodes(n, children)
	if len(children) == 0 {
		return nil
	} else if n == nil {
		return common.ErrNilReceiver
	}

	if n.FirstChild == nil {
		n.LastChild = children[len(children)-1]
	} else {
		n.FirstChild.PrevSibling = children[len(children)-1]
		children[len(children)-1].NextSibling = n.FirstChild
	}

	n.FirstChild = children[0]
	n.Pos = children[0].Pos

	return nil
}

// AppendChildren adds the given children nodes to the end of the current node's children list.
//
// Parameters:
//   - children: Variadic parameter of type Node representing the children to be appended.
//
// Returns:
//   - error: Returns an error if the operation fails or if the receiver is nil, otherwise returns nil.
func (n *Node) AppendChildren(children ...*Node) error {
	children = link_nodes(n, children)
	if len(children) == 0 {
		return nil
	} else if n == nil {
		return common.ErrNilReceiver
	}

	if n.LastChild == nil {
		n.FirstChild = children[0]
		n.Pos = children[0].Pos
	} else {
		n.LastChild.NextSibling = children[0]
		children[0].PrevSibling = n.LastChild
	}

	n.LastChild = children[len(children)-1]

	return nil
}

var (
	// WriteTree is a function that writes a tree to the given writer.
	//
	// Parameters:
	//   - w: The writer to write the tree to.
	//   - node: The root node of the tree to write.
	//
	// Returns:
	//   - int: The number of bytes written to the writer.
	//   - error: Returns an error if the operation fails.
	WriteTree func(w io.Writer, node *Node) (int, error)
)

func init() {
	WriteTree = trees.MakeWriteTree[*Node]()
}`
