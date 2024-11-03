package ast

import (
	"strconv"
	"strings"

	"github.com/PlayerR9/mygo-lib/common"
)

type NodeType interface {
	~int

	String() string
}

func Quote[N NodeType](b *strings.Builder, type_ N) error {
	if b == nil {
		return common.NewErrNilParam("b")
	}

	_, _ = b.WriteString(strconv.Quote(type_.String()))

	return nil
}

type Noder interface {
	IsNil() bool
	String() string

	// AreChildrenCritical checks whether or not the children of the node are considered critical
	// for the node's evaluation. If they are critical, any error is immediately returned
	// and the evaluation is aborted. If not, the errors are stored and the evaluation proceeds
	// to the next node.
	//
	// Returns:
	// 	- bool: True if the node's children are critical, false otherwise.
	AreChildrenCritical() bool
}
