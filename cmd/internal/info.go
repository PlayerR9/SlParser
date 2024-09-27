package internal

import (
	"errors"
	"fmt"
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

// IsNil checks whether the info is nil.
//
// Returns:
//   - bool: true if the info is nil, false otherwise.
func (info *Info) IsNil() bool {
	return info == nil
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
func (info *Info) NextInfos() ([]*Info, error) {
	if info == nil {
		return nil, errors.New("receiver is nil")
	}

	new_frames := info.AppendFrame()

	var nexts []*Info

	for child := range info.Info.Child() {
		if child.IsNil() {
			return nil, errors.New("found a nil child")
		}

		next, err := NewInfo(child, new_frames)
		if err != nil {
			return nil, err
		}

		nexts = append(nexts, next)
	}

	return nexts, nil
}

func NewInfo(node *kdd.Node, frames []string) (*Info, error) {
	if node == nil {
		return nil, gers.NewErrNilParameter("node")
	}

	info, err := ast.NewInfo(node, frames)
	if err != nil {
		return nil, err
	}

	next := &Info{
		Type: InvalidTk,
		Info: info,
	}

	return next, nil
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
	InfoTableOf *ast.InfoTableMaker[*kdd.Node, *Info]
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
		literal, err := make_literal(type_, node.Data)
		if err != nil {
			err.AddContext("node", node)
			err.AddFrame(fmt.Sprintf("make_literal(%s, %q)", type_.String(), node.Data))

			return nil, err
		}

		sub_info := gers.AssertNew(ast.NewInfo(node, []string{literal}))

		info := &Info{
			Type:        type_,
			Literal:     literal,
			IsCandidate: is_candidate,
			Info:        sub_info,
		}

		return info, nil
	}

	init_fn := func(node *kdd.Node, _ []string) (*Info, error) {
		if node == nil {
			return nil, gers.NewErrNilParameter("node")
		}

		sub_info, err := ast.NewInfo(node, nil)
		if err != nil {
			return nil, err
		}

		return &Info{
			Type: InvalidTk,
			Info: sub_info,
		}, nil
	}

	InfoTableOf = &ast.InfoTableMaker[*kdd.Node, *Info]{
		InitFn:     init_fn,
		MakeInfoFn: fn,
	}
}
