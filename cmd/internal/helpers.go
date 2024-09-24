package internal

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode"

	kdd "github.com/PlayerR9/SlParser/kdd"
	gcch "github.com/PlayerR9/go-commons/runes"
	gers "github.com/PlayerR9/go-errors"
)

// TypeOf returns the type of a kdd.Node.
//
// Parameters:
//   - n: The node.
//
// Returns:
//   - TokenType: The type of the node.
//   - error: An error if the node is nil or not a RHS node.
func TypeOf(n *kdd.Node) (TokenType, error) {
	if n == nil {
		return ExtraTk, errors.New("node must not be nil")
	}

	if n.Type != kdd.RhsNode {
		return ExtraTk, fmt.Errorf("node must be a RHS node, got %s instead", n.Type.String())
	}

	if !n.IsTerminal {
		return NonterminalTk, nil
	}

	if n.Data == "EOF" {
		return ExtraTk, nil
	}

	return TerminalTk, nil
}

func CheckEofExists(tokens []*kdd.Node) bool {
	if len(tokens) == 0 {
		return false
	}

	for _, tk := range tokens {
		gers.AssertNotNil(tk, "tk")

		if tk.Type == kdd.RhsNode && tk.Data == "EOF" {
			return true
		}
	}

	return false
}

func FindLastTerminal(tokens []*kdd.Node) (*kdd.Node, error) {
	if len(tokens) == 0 {
		return nil, nil
	}

	idx := -1

	for i := 0; i < len(tokens) && idx == -1; i++ {
		gers.AssertNotNil(tokens[i], "tokens[i]")

		type_, err := TypeOf(tokens[i])
		if err != nil {
			return nil, fmt.Errorf("at index %d: %w", i, err)
		}

		if type_ == NonterminalTk {
			idx = i
		}
	}

	if idx == -1 {
		return tokens[len(tokens)-1], nil
	} else if idx == 0 {
		return nil, nil
	}

	return tokens[idx-1], nil
}

func CandidatesForAst(tokens []*kdd.Node) ([]string, error) {
	if len(tokens) == 0 {
		return nil, nil
	}

	var candidates []string

	for i, tk := range tokens {
		if tk == nil {
			continue
		}

		type_, err := TypeOf(tk)
		if err != nil {
			return nil, fmt.Errorf("at index %d: %w", i, err)
		}

		ok := IsCandidateForAst(type_, tk.Data)
		if !ok {
			continue
		}

		pos, ok := slices.BinarySearch(candidates, tk.Data)
		if !ok {
			candidates = slices.Insert(candidates, pos, tk.Data)
		}
	}

	return candidates, nil
}

/////////////////////////

func replace_underscore(chars []rune) string {
	var builder strings.Builder

	capitalize_next := true

	for i := 0; i < len(chars); i++ {
		c := chars[i]

		if c == '_' {
			capitalize_next = true
		} else if capitalize_next {
			builder.WriteRune(unicode.ToUpper(c))
			capitalize_next = false
		} else {
			builder.WriteRune(unicode.ToLower(c))
		}
	}

	return builder.String()
}

func MakeToken(symbol []byte) (*Token, error) {
	if len(symbol) == 0 {
		return nil, errors.New("symbol must not be empty")
	}

	if bytes.Equal(symbol, []byte("EOF")) {
		tk := NewToken(ExtraTk, "EOF")
		return tk, nil
	}

	chars, err := gcch.BytesToUtf8(symbol)
	if err != nil {
		return nil, err
	}

	if !unicode.IsLetter(chars[0]) {
		return nil, errors.New("symbol must start with a letter")
	}

	var type_ TokenType

	if unicode.IsUpper(chars[0]) {
		type_ = TerminalTk
	} else {
		type_ = NonterminalTk
		chars[0] = unicode.ToUpper(chars[0])
	}

	tk := NewToken(type_, replace_underscore(chars))
	return tk, nil
}

/* func unique(tokens []string) []string {
	for i := 0; i < len(tokens)-1; i++ {
		top := i + 1

		for j := i + 1; j < len(tokens); j++ {
			if tokens[j] != tokens[i] {
				tokens[top] = tokens[j]
				top++
			}
		}

		tokens = tokens[:top:top]
	}

	return tokens
} */
