package parser

import (
	"github.com/PlayerR9/SlParser/parser/internal"
	tr "github.com/PlayerR9/mygo-lib/CustomData/tree"
)

type Parser interface {
	Parse(tokens []*tr.Node) *Iterator
	ItemsOf(type_ string) ([]*internal.Item, bool)
}

func NewParser(table map[string][]*internal.Item) Parser {
	return &baseParser{
		decision_table: table,
	}
}
