package parser

import (
	slgr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/parser/internal"
)

type Parser interface {
	Parse(tokens []*slgr.Token) *Iterator
	ItemsOf(type_ string) ([]*internal.Item, bool)
}

func NewParser(table map[string][]*internal.Item) Parser {
	return &baseParser{
		decision_table: table,
	}
}
