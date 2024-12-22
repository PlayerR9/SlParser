package grammar

import (
	"strconv"
	"strings"

	"github.com/PlayerR9/SlParser/mygo-lib/common"
	gslc "github.com/PlayerR9/SlParser/mygo-lib/slices"
)

// Token is a token in the grammar.
type Token struct {
	// Type is the type of the token.
	Type string

	// Data is the data of the token.
	Data string

	// Children is the children of the token.
	Children []*Token
}

// String implements TreeNode.
func (tk Token) String() string {
	var builder strings.Builder

	_, _ = builder.WriteString("Token[")
	_, _ = builder.WriteString(tk.Type)

	if tk.Data != "" {
		_, _ = builder.WriteString(" (")
		_, _ = builder.WriteString(strconv.Quote(tk.Data))
		_, _ = builder.WriteRune(')')
	}

	_, _ = builder.WriteRune(']')

	str := builder.String()
	return str
}

// NewToken creates a new Token with the specified type and data.
//
// Parameters:
//   - type_: The type of the token.
//   - data: The data associated with the token.
//
// Returns:
//   - *Token: A pointer to the newly created Token. Never returns nil.
func NewToken(type_, data string) *Token {
	tk := &Token{
		Type: type_,
		Data: data,
	}

	return tk
}

// PrependChildren prepends the given children to the token's children.
//
// Parameters:
//   - children: The children to be prepended.
//
// Returns:
//   - error: An error if the receiver is nil.
func (tk *Token) PrependChildren(children []*Token) error {
	if tk == nil {
		return common.ErrNilReceiver
	}

	children = gslc.RejectNils(children)
	if len(children) == 0 {
		return nil
	}

	tk.Children = append(children, tk.Children...)

	return nil
}

// AppendChildren appends the given children to the token's children.
//
// Parameters:
//   - children: The children to be appended.
//
// Returns:
//   - error: An error if the receiver is nil.
func (tk *Token) AppendChildren(children []*Token) error {
	if tk == nil {
		return common.ErrNilReceiver
	}

	children = gslc.RejectNils(children)
	if len(children) == 0 {
		return nil
	}

	tk.Children = append(tk.Children, children...)

	return nil
}

// GetChildren returns a copy of the children of the token.
//
// Returns:
//   - []*Token: A copy of the children of the token.
func (tk Token) GetChildren() []*Token {
	if len(tk.Children) == 0 {
		return nil
	}

	children := make([]*Token, len(tk.Children))
	copy(children, tk.Children)

	return children
}
