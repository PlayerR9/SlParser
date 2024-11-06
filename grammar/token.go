package grammar

import (
	"fmt"
	"strconv"

	tr "github.com/PlayerR9/mygo-lib/CustomData/tree"
)

type TokenData struct {
	Type      string
	Data      string
	Pos       int
	Lookahead *tr.Node
}

func (d TokenData) String() string {
	if d.Data == "" {
		return fmt.Sprintf("%d:%s", d.Pos, d.Type)
	} else {
		return fmt.Sprintf("%d:%s (%s)", d.Pos, d.Type, strconv.Quote(d.Data))
	}
}

func (d TokenData) Equals(other tr.Infoer) bool {
	if other == nil {
		return false
	}

	v, ok := other.(*TokenData)
	if !ok {
		return false
	}

	return d.Type == v.Type && d.Data == v.Data
}

// NewToken is a convenience function to create a new token.
//
// Parameters:
//   - pos: The position of the token (in bytes).
//   - type_: The type of the token.
//   - data: The data of the token.
//
// Returns:
//   - *tree.Node: The new tree node. Never returns nil.
func NewToken(pos int, type_, data string, lookahead *tr.Node) *tr.Node {
	return tr.NewNode(&TokenData{
		Pos:       pos,
		Type:      type_,
		Data:      data,
		Lookahead: lookahead,
	})
}
