package internal

import (
	"errors"
	"unicode"
	"unicode/utf8"

	"github.com/PlayerR9/SlParser/ast"
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

type Info struct {
	Type    TokenType
	Literal string

	*ast.Info[*kdd.Node]
}

func (info *Info) NextInfos() []*Info {
	if info == nil {
		return nil
	}

	new_frames := info.AppendFrame()

	var nexts []*Info

	for child := range info.Info.Child() {
		next := NewInfo()
		next.Init(child, new_frames)
	}

	return nexts
}

func NewInfo() *Info {
	return &Info{
		Type: InvalidTk,
		Info: ast.NewInfo[*kdd.Node](),
	}
}

var (
	InfoTableOf ast.InfoTableOfFn[*kdd.Node, *Info]
)

func init() {
	fn := func(node *kdd.Node) (*Info, error) {
		gers.AssertNotNil(node, "node")

		if node.Type != kdd.RhsNode {
			return nil, ast.IgnoreInfo
		}

		info := NewInfo()

		if node.Data == "EOF" {
			info.Type = ExtraTk
			info.Literal = "EtEOF"

			return info, nil
		}

		c, _ := utf8.DecodeRuneInString(node.Data)
		if c == utf8.RuneError {
			return nil, errors.New("found node with invalid utf8-encoded data")
		}

		if unicode.IsLower(c) {
			info.Type = NonterminalTk
			info.Literal = "Nt" + node.Data
		} else {
			info.Type = TerminalTk
			info.Literal = "Tt" + node.Data
		}

		return info, nil
	}

	InfoTableOf = ast.MakeInfoTable[*kdd.Node, *Info](fn)
}
