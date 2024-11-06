package grammar

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/PlayerR9/go-evals/common"
	tr "github.com/PlayerR9/mygo-lib/CustomData/tree"
)

type NodeData struct {
	Type string
	Data string
	Pos  int
}

func (d NodeData) String() string {
	if d.Data == "" {
		return fmt.Sprintf("%d:%s", d.Pos, d.Type)
	} else {
		return fmt.Sprintf("%d:%s (%s)", d.Pos, d.Type, strconv.Quote(d.Data))
	}
}

func (d NodeData) Equals(other tr.Infoer) bool {
	if other == nil {
		return false
	}

	v, ok := other.(*NodeData)
	if !ok {
		return false
	}

	return d.Type == v.Type && d.Data == v.Data
}

// NewNode is a convenience function to create a new AST node.
//
// Parameters:
//   - pos: The position of the node (in bytes).
//   - type_: The type of the node.
//   - data: The data of the node.
//
// Returns:
//   - *tree.Node: The new tree node. Never returns nil.
func NewNode(pos int, type_, data string) *tr.Node {
	return tr.NewNode(&NodeData{
		Pos:  pos,
		Type: type_,
		Data: data,
	})
}

func GetNodeData(node *tr.Node) (*TokenData, error) {
	if node == nil {
		return nil, common.NewErrNilParam("node")
	}

	info := node.Info
	if info == nil {
		return nil, errors.New("node has no data")
	}

	v, ok := info.(*TokenData)
	if !ok {
		return nil, errors.New("node has the wrong data")
	}

	return v, nil
}
