package grammar

import (
	"errors"
	"fmt"

	tr "github.com/PlayerR9/mygo-lib/CustomData/tree"
	"github.com/PlayerR9/mygo-lib/common"
	gslc "github.com/PlayerR9/mygo-lib/slices"
)

/////////////////////////////////////////////////////////

// astT is for internal use.
type astT struct{}

var (
	// AST is the namespace that allows for the creation of AST nodes.
	AST astT
)

func init() {
	AST = astT{}
}

// ToASTFn is a function that converts a token to a node.
//
// Parameters:
//   - token: the token to convert. Assumed to not be nil.
//
// Returns:
//   - []*tr.Node: the nodes created by the token.
//   - error: an error if the token cannot be converted.
type ToASTFn func(token *tr.Node) ([]*tr.Node, error)

// Make makes an AST node from a token.
//
// Parameters:
//   - ast: The AST maker.
//   - token: The token to make an AST node from.
//
// Returns:
//   - []*tr.Node: The AST nodes.
//   - error: An error if the evaluation failed.
func (astT) Make(ast map[string]ToASTFn, token *tr.Node) ([]*tr.Node, error) {
	if ast == nil {
		return nil, common.NewErrNilParam("ast")
	} else if token == nil {
		return nil, common.NewErrNilParam("token")
	}

	tkd, err := Get[*TokenData](token)
	if err != nil {
		return nil, err
	}

	type_ := tkd.Type

	fn, ok := ast[type_]
	if !ok || fn == nil {
		fn = func(t *tr.Node) ([]*tr.Node, error) {
			common.TODO("Handle this case.")

			return nil, nil
		}
	}

	nodes, err := fn(token)
	return nodes, err
}

type Group struct {
	Tokens []string
	Nodes  []string
}

func NewGroup(tokens, nodes []string) Group {
	if len(tokens) > 1 {
		_ = gslc.Uniquefy(&tokens)
	}

	if len(nodes) > 1 {
		_ = gslc.Uniquefy(&nodes)
	}

	return Group{
		Tokens: tokens,
		Nodes:  nodes,
	}
}

func (g Group) Check(ast map[string]ToASTFn, idx int, child *tr.Node, nodes *[]*tr.Node) error {
	if nodes == nil {
		return common.NewErrNilParam("nodes")
	}

	if len(g.Tokens) > 0 {
		err := CheckNode(idx, "child", child, g.Tokens...)
		if err != nil {
			return err
		}
	}

	sub_nodes, err := AST.Make(ast, child)
	if err != nil {
		return err
	} else if len(sub_nodes) == 0 {
		return fmt.Errorf("while making %d child: expected at least 1 sub node, got %d instead", idx, len(sub_nodes))
	}

	if len(g.Nodes) == 0 {
		*nodes = append(*nodes, sub_nodes...)

		return nil
	}

	for j, sub_node := range sub_nodes {
		err := CheckNode(j, "sub node", sub_node, g.Nodes...)
		if err != nil {
			return err
		}
	}

	*nodes = append(*nodes, sub_nodes...)

	return nil
}

func Many(ast map[string]ToASTFn, lhs string, groups ...Group) error {
	if ast == nil {
		return common.NewErrNilParam("ast")
	} else if len(groups) == 0 {
		return common.NewErrBadParam("groups", "must contain at least 1 group")
	}

	fn := func(t *tr.Node) ([]*tr.Node, error) {
		if t.LastChild == nil {
			return nil, errors.New("missing last child")
		}

		is_base_case := MustGet[*TokenData](t.LastChild).Type != lhs

		var size int

		for n := t.LastChild; n != nil; n = n.PrevSibling {
			size++
		}

		if is_base_case {
			if size != len(groups) {
				return nil, fmt.Errorf("expected %d children, got %d instead", len(groups), size)
			}
		} else if size != len(groups)+1 {
			return nil, fmt.Errorf("expected %d children for last level, got %d instead", len(groups)+1, size)
		}

		var nodes []*tr.Node

		child := t.FirstChild

		for i, g := range groups {
			err := g.Check(ast, i, child, &nodes)
			if err != nil {
				return nodes, err
			}

			child = child.NextSibling
		}

		if is_base_case {
			return nodes, nil
		}

		err := CheckNode(len(groups), "child", t.LastChild, lhs)
		if err != nil {
			return nil, err
		}

		sub_nodes, err := AST.Make(ast, t.LastChild)
		if err != nil {
			return nil, err
		}

		return append(nodes, sub_nodes...), nil
	}

	ast[lhs] = fn

	return nil
}

// Transform transforms a given token into an AST node.
//
// Parameters:
//   - token: The token to transform.
//
// Returns:
//   - *tr.Node: The transformed node.
//   - error: An error if the transformation failed.
func (astT) Transform(token *tr.Node) (*tr.Node, error) {
	if token == nil {
		return nil, common.NewErrNilParam("token")
	}

	tkd, err := Get[*TokenData](token)
	if err != nil {
		return nil, err
	}

	n := NewNode(tkd.Pos, tkd.Type, tkd.Data)

	if token.FirstChild == nil {
		return n, nil
	}

	var children []*tr.Node

	for c := token.FirstChild; c != nil; c = c.NextSibling {
		child, err := AST.Transform(c)
		if err != nil {
			return nil, err
		}

		children = append(children, child)
	}

	_ = n.AppendChildren(children...)

	return n, nil
}
