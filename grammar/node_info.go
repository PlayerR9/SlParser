package grammar

import (
	"fmt"
	"strconv"
)

type NodeInfo struct {
	Pos  int
	Type string
	Data string
}

func (n NodeInfo) String() string {
	if n.Data == "" {
		return fmt.Sprintf("%d:%s", n.Pos, n.Type)
	} else {
		return fmt.Sprintf("%d:%s (%s)", n.Pos, n.Type, strconv.Quote(n.Data))
	}
}
