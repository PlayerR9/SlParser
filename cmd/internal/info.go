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
	InvalidTk TokenType = iota - 1

	// ExtraTk is the token type for extra symbols such as the EOF symbol.
	ExtraTk

	// TerminalTk is the token type for terminal symbols.
	TerminalTk

	// NonterminalTk is the token type for nonterminal symbols.
	NonterminalTk
)

// Info is the info of a kdd node.
type Info struct {
	// Type indicates the type of the node when the lexer/parser is generated.
	Type TokenType

	// Literal is the literal of the node.
	// This is the enumerated value of the node type.
	Literal string

	// IsCandidate indicates whether the node is a candidate for the AST.
	IsCandidate bool

	*ast.Info[*kdd.Node]
}

/* // NewInfo creates a new info.
//
// Returns:
//   - *Info: The new info. Never returns nil.
//
// The info is initialized with the invalid token type. Make sure
// to change the type before using the info.
func NewInfo() *Info {
	return &Info{
		Type: InvalidTk,
		Info: ast.NewInfo[*kdd.Node](),
	}
}
*/

// Equals checks whether the info is equal to another info.
//
// Two infos are said to be equal if they have the same literal. Also, if other is
// nil, then false is returned.
//
// Parameters:
//   - other: The other info.
//
// Returns:
//   - bool: True if the infos are equal. False otherwise.
func (info Info) Equals(other *Info) bool {
	return other != nil && info.Literal == other.Literal
}

// NextInfos returns the information of the next nodes.
//
// Returns:
//   - []*Info: The information of the next nodes. No nil nodes are returned.
//
// As with NewInfo, the info is initialized with the invalid token type.
func (info *Info) NextInfos() []*Info {
	if info == nil {
		return nil
	}

	new_frames := info.AppendFrame()

	var nexts []*Info

	for child := range info.Info.Child() {
		next := &Info{
			Type: InvalidTk,
			Info: ast.NewInfo[*kdd.Node](),
		}

		next.Init(child, new_frames)

		nexts = append(nexts, next)
	}

	return nexts
}

var (
	// InfoTableOf is a function that creates an info table given the root node.
	//
	// Parameters:
	//   - root: The root node.
	//
	// Returns:
	//   - map[*kdd.Node]*Info: The info table.
	//   - error: An error if the info table could not be created.
	InfoTableOf ast.InfoTableOfFn[*kdd.Node, *Info]
)

func init() {
	fn := func(node *kdd.Node) (*Info, error) {
		gers.AssertNotNil(node, "node")

		if node.Type != kdd.RhsNode {
			return nil, ast.IgnoreInfo
		}

		// 1. Determine the type of the node.
		var type_ TokenType

		if node.Data == "EOF" {
			type_ = ExtraTk
		} else {
			c, _ := utf8.DecodeRuneInString(node.Data)
			if c == utf8.RuneError {
				return nil, errors.New("found node with invalid utf8-encoded data")
			}

			if unicode.IsLower(c) {
				type_ = NonterminalTk
			} else {
				type_ = TerminalTk
			}
		}

		// 2. Determine whether the node is a candidate for the AST.
		var is_candidate bool

		if type_ != NonterminalTk {
			is_candidate = false
		} else {
			r, _ := utf8.DecodeLastRuneInString(node.Data)

			is_candidate = !unicode.IsDigit(r) && unicode.IsLetter(r)
		}

		// 3. Determine the literal of the node.
		literal, err := MakeLiteral(type_, node.Data)
		gers.AssertErr(err, "MakeLiteral(%s, %q)", type_.String(), node.Data)

		info := &Info{
			Type:        type_,
			Literal:     literal,
			IsCandidate: is_candidate,
			Info:        ast.NewInfo[*kdd.Node](),
		}

		return info, nil
	}

	InfoTableOf = ast.MakeInfoTable(fn)
}
