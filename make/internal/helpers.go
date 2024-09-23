package internal

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode"

	gcch "github.com/PlayerR9/go-commons/runes"
)

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

func unique(tokens []*Token) []*Token {
	for i := 0; i < len(tokens)-1; i++ {
		top := i + 1

		for j := i + 1; j < len(tokens); j++ {
			if tokens[j].Data != tokens[i].Data {
				tokens[top] = tokens[j]
				top++
			}
		}

		tokens = tokens[:top:top]
	}

	return tokens
}

func sort(tokens []*Token) error {
	buckets := make(map[TokenType][]*Token, 3)
	buckets[ExtraTk] = make([]*Token, 0)
	buckets[TerminalTk] = make([]*Token, 0)
	buckets[NonterminalTk] = make([]*Token, 0)

	for _, tk := range tokens {
		type_ := tk.Type

		prev, ok := buckets[type_]
		if !ok {
			return fmt.Errorf("bucket %q not found", type_.String())
		}

		buckets[type_] = append(prev, tk)
	}

	for type_, bucket := range buckets {
		slices.SortStableFunc(bucket, func(a, b *Token) int {
			return strings.Compare(a.Data, b.Data)
		})

		buckets[type_] = bucket
	}

	i := 0

	tks := buckets[ExtraTk]
	for _, tk := range tks {
		tokens[i] = tk
		i++
	}

	tks = buckets[TerminalTk]
	for _, tk := range tks {
		tokens[i] = tk
		i++
	}

	tks = buckets[NonterminalTk]
	for _, tk := range tks {
		tokens[i] = tk
		i++
	}

	return nil
}

func TokenSymbols(tokens []*Token) ([]*Token, error) {
	if len(tokens) == 0 {
		return nil, nil
	}

	tokens = unique(tokens)
	err := sort(tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func ExtractSymbols(tokens []*Token) []string {
	if len(tokens) == 0 {
		return nil
	}

	symbols := make([]string, 0, len(tokens))

	for _, tk := range tokens {
		s := tk.String()
		symbols = append(symbols, s)
	}

	return symbols
}

func FindLastTerminal(tokens []*Token) *Token {
	if len(tokens) == 0 {
		return nil
	}

	idx := -1

	for i := 0; i < len(tokens) && idx == -1; i++ {
		if tokens[i].Type == NonterminalTk {
			idx = i
		}
	}

	if idx == -1 {
		return tokens[len(tokens)-1]
	} else if idx == 0 {
		return nil
	}

	return tokens[idx-1]
}

func CheckEofExists(tokens []*Token) bool {
	for _, tk := range tokens {
		if tk.Type == ExtraTk && tk.Data == "EOF" {
			return true
		}
	}

	return false
}

func CandidatesForAst(tokens []*Token) []string {
	if len(tokens) == 0 {
		return nil
	}

	var candidates []string

	for _, tk := range tokens {
		if tk.IsCandidateForAst() {
			candidates = append(candidates, tk.Data)
		}
	}

	return candidates
}
