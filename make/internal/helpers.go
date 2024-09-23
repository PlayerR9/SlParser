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

func Sort(tokens []*Token) error {
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

	// StableSort each buckets

	for key, bucket := range buckets {
		slices.Sort(bucket)

		buckets[key] = bucket
	}

	// Concatenate each buckets
	elems := buckets['E']

	i := 0

	for _, elem := range elems {
		slice[i] = elem
		i++
	}

	elems = buckets['T']

	for _, elem := range elems {
		slice[i] = elem
		i++
	}

	elems = buckets['N']

	for _, elem := range elems {
		slice[i] = elem
		i++
	}
}

func TokenSymbols(tokens []*Token) []*Token {
	if len(tokens) == 0 {
		return nil
	}

	tokens = unique(tokens)

	return tokens
}

func ExtractSymbols(tokens []*Token) []string {
	var symbols []string

	for _, tk := range tokens {
		s := tk.String()

		pos, ok := slices.BinarySearch(symbols, s)
		if !ok {
			symbols = slices.Insert(symbols, pos, s)
		}
	}

	return symbols
}

func FindLastTerminal(symbols []string) (string, bool) {
	if len(symbols) == 0 {
		return "", false
	}

	idx := -1

	for i := 0; i < len(symbols) && idx == -1; i++ {
		if strings.HasPrefix(symbols[i], "Ntt") {
			idx = i
		}
	}

	if idx == -1 {
		return symbols[len(symbols)-1], true
	}

	if idx == 0 {
		return "", false
	}

	return symbols[idx-1], true
}
