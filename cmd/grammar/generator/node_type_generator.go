package generator

import (
	"fmt"

	common "github.com/PlayerR9/mygo-lib/common"
	cgen "github.com/PlayerR9/mygo-lib/generator"
)

type NodeTypeData struct {
	PackageName string
	Generate    string
}

func (ntd *NodeTypeData) SetPkgName(pkg_name string) error {
	if ntd == nil {
		return common.ErrNilReceiver
	}

	ntd.PackageName = pkg_name

	return nil
}

func NewNodeTypeData() *NodeTypeData {
	return &NodeTypeData{
		Generate: fmt.Sprintf("//%s:%s %s %s %s", "go", "generate", "stringer", "-type=NodeType", "-linecomment"),
	}
}

var (
	NodeTypeGenerator *cgen.CodeGenerator[*NodeTypeData]
)

func init() {
	NodeTypeGenerator = cgen.Must(cgen.New[*NodeTypeData]("node_type", node_type_templ))
}

const node_type_templ string = `package {{ .PackageName }}

{{ .Generate }}

type NodeType int

const (
	/*SourceNode is the root of the AST.
	Node[SourceNode]
	*/
	SourceNode NodeType = iota // Source
)

// String implements the fmt.Stringer interface.
func (nt NodeType) String() string {
	return [...]string{
		"Source",
	}[nt]
}

// AreChildrenCritical checks whether the children of a node are critical depending on the type
// they have.
//
// A children is said to be critical if, when an error occurs, the entire processed should be stopped
// instead of continuing.
//
// Returns:
// 	- bool: True if they are, false otherwise.
func (nt NodeType) AreChildrenCritical() bool {
	return nt != SourceNode
}`
