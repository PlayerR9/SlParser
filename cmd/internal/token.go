package internal

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	kdd "github.com/PlayerR9/SlParser/kdd"
	gers "github.com/PlayerR9/go-errors"
)

// TokenType is the type of a token.
type TokenType int

const (
	// InvalidTk is the invalid token type.
	InvalidTk TokenType = iota - 1 // Invalid

	// ExtraTk is the token type for extra symbols such as the EOF symbol.
	ExtraTk // Et

	// TerminalTk is the token type for terminal symbols.
	TerminalTk // Tt

	// NonterminalTk is the token type for nonterminal symbols.
	NonterminalTk // Nt
)

type NodeInfo struct {
	Type TokenType
}

func NewNodeInfo(type_ TokenType) *NodeInfo {
	return &NodeInfo{
		Type: type_,
	}
}

type InfoTable struct {
	table map[*kdd.Node]*NodeInfo
}

func NewInfoTable() *InfoTable {
	return &InfoTable{
		table: make(map[*kdd.Node]*NodeInfo),
	}
}

func (table *InfoTable) AddInfo(node *kdd.Node, info *NodeInfo) {
	if table == nil || node == nil {
		return
	}

	gers.AssertNotNil(table.table, "info.table")

	if info == nil {
		delete(table.table, node)
	} else {
		table.table[node] = info
	}
}

func make_info_rec(root *kdd.Node, info *InfoTable) {
	gers.AssertNotNil(root, "root")
	gers.AssertNotNil(info, "info")

	children := root.GetChildren()

	for _, child := range children {
		make_info_rec(child, info)
	}

	switch root.Type {
	case kdd.SourceNode:
		// Do nothing
	case kdd.RuleNode:
		// Do nothing
	case kdd.RhsNode:
		var type_ TokenType

		if root.Data == "EOF" {
			type_ = ExtraTk
		} else {
			c, _ := utf8.DecodeRuneInString(root.Data)
			gers.Assert(c != utf8.RuneError, "root.Data is not a valid UTF-8 string")

			if unicode.IsUpper(c) {
				type_ = TerminalTk
			} else {
				type_ = NonterminalTk
			}
		}

		node_info := NewNodeInfo(type_)

		info.AddInfo(root, node_info)
	default:
		panic(fmt.Sprintf("unknown node type: %v", root.Type))
	}
}

func MakeInfo(root *kdd.Node) (*InfoTable, error) {
	err := kdd.CheckAST(root, -1)
	if err != nil {
		return nil, err
	}

	info := NewInfoTable()

	make_info_rec(root, info)

	return info, nil
}

// Token is a token.
type Token struct {
	// Type is the type of the token.
	Type TokenType

	// Data is the data of the token.
	Data string
}

func (t Token) String() string {
	return t.Type.String() + t.Data
}

func NewToken(type_ TokenType, data string) *Token {
	return &Token{
		Type: type_,
		Data: data,
	}
}

func IsCandidateForAst(type_ TokenType, data string) bool {
	if type_ != NonterminalTk || data == "" {
		return false
	}

	r, _ := utf8.DecodeLastRuneInString(data)
	return !unicode.IsDigit(r) && unicode.IsLetter(r)
}
