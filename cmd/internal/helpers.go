package internal

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode"

	kdd "github.com/PlayerR9/SlParser/kdd"
	gcch "github.com/PlayerR9/go-commons/runes"
	gers "github.com/PlayerR9/go-errors"
)

// replace_underscore replaces underscores with underscores.
//
// This function skips underscores and uppercase the following letter. Furthermore
// any other letter is lowercased. Finally, the first letter is always uppercase.
//
// Parameters:
//   - chars: The characters to replace.
//
// Returns:
//   - string: The replaced string.
//
// Assertions:
//   - len(chars) > 0
func replace_underscore(chars []rune) string {
	gers.Assert(len(chars) > 0, "chars must not be empty")

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

// MakeLiteral makes a literal according to the type and its data.
//
// Parameters:
//   - type_: The type of the literal.
//   - data: The data of the literal.
//
// Returns:
//   - string: The literal.
//   - error: The error.
//
// Errors:
//   - errors.ErrInvalidUTF8Encoding: If the data is not properly encoded in UTF-8.
//   - error.Err (with code InvalidParameter): If the data is an empty string.
//   - any other error if the data does not starts with a letter nor the type_ is not one of the
//     supported types (that are, TerminalTk, NonterminalTk, ExtraTk).
func MakeLiteral(type_ TokenType, data string) (string, error) {
	if data == "" {
		return "", gers.NewErrInvalidParameter("data must not be an empty string")
	}

	chars, err := gcch.StringToUtf8(data)
	if err != nil {
		return "", err
	}

	if !unicode.IsLetter(chars[0]) {
		return "", errors.New("symbol must start with a letter")
	}

	if unicode.IsLower(chars[0]) {
		chars[0] = unicode.ToUpper(chars[0])
	}

	data = replace_underscore(chars)

	switch type_ {
	case TerminalTk:
		data = "Tt" + data
	case NonterminalTk:
		data = "Nt" + data
	case ExtraTk:
		data = "Et" + data
	default:
		return "", fmt.Errorf("type (%v) is not supported", type_)
	}

	return data, nil
}

// CandidatesForAst returns the candidates for the AST. (i.e., the list of all
// literals such that they are marked as candidates).
//
// Parameters:
//   - table: The info table.
//
// Returns:
//   - []string: The candidates.
//
// The returned slice is sorted in ascending order and contains no duplicates.
func CandidatesForAst(table map[*kdd.Node]*Info) []string {
	if len(table) == 0 {
		return nil
	}

	var candidates []string

	for _, info := range table {
		gers.AssertNotNil(info, "info")

		if !info.IsCandidate {
			continue
		}

		pos, ok := slices.BinarySearch(candidates, info.Literal)
		if !ok {
			candidates = slices.Insert(candidates, pos, info.Literal)
		}
	}

	return candidates
}

// sort sorts the given slice of Info objects in ascending order.
//
// The slice is first divided into three buckets: ExtraTk, TerminalTk, and
// NonterminalTk. Each bucket is then sorted in ascending order using the
// Literals of the Info objects. Finally, the buckets are concatenated in the
// order given above and the sorted slice is returned.
//
// Does nothing if the slice is empty. Plus, the slice is sorted in place.
//
// Parameters:
//   - infos: The slice of Info objects to be sorted.
//
// Returns:
//   - error: nil on success, otherwise an error.
func sort(infos []*Info) error {
	if len(infos) == 0 {
		return nil
	}

	// 1. Make buckets for each type.
	buckets := make(map[TokenType][]*Info, 3)
	buckets[ExtraTk] = make([]*Info, 0)
	buckets[TerminalTk] = make([]*Info, 0)
	buckets[NonterminalTk] = make([]*Info, 0)

	// 2. Divide the infos into buckets.
	for _, info := range infos {
		gers.AssertNotNil(info, "node")

		prev, ok := buckets[info.Type]
		if !ok {
			return fmt.Errorf("bucket %q not found", info.Type.String())
		}

		buckets[info.Type] = append(prev, info)
	}

	// 3. Sort the buckets.
	for type_, bucket := range buckets {
		slices.SortStableFunc(bucket, func(a, b *Info) int {
			return strings.Compare(a.Literal, b.Literal)
		})

		buckets[type_] = bucket
	}

	// 4. Concatenate the buckets.
	i := 0

	tks := buckets[ExtraTk]
	for _, tk := range tks {
		infos[i] = tk
		i++
	}

	tks = buckets[TerminalTk]
	for _, tk := range tks {
		infos[i] = tk
		i++
	}

	tks = buckets[NonterminalTk]
	for _, tk := range tks {
		infos[i] = tk
		i++
	}

	return nil
}

// LinearizeTable linearizes the given table into a list without duplicates.
//
// Does nothing if the table is empty.
//
// Parameters:
//   - table: The table to linearize.
//
// Returns:
//   - []*Info: The linearized list.
func LinearizeTable(table map[*kdd.Node]*Info) []*Info {
	if len(table) == 0 {
		return nil
	}

	// 1. Transform the map into a list without duplicates.
	list := make([]*Info, 0, len(table))

	for _, info := range table {
		gers.AssertNotNil(info, "info")

		ok := slices.ContainsFunc(list, info.Equals)
		if !ok {
			list = append(list, info)
		}
	}

	list = list[:len(list):len(list)]

	// 2. Sort the list in ascending order using bucket sort.
	err := sort(list)
	gers.AssertErr(err, "sort(list)")

	return list
}

// FindLastTerminal finds the last terminal in the given list.
//
// This function assumes that the list is sorted in ascending order.
// Make sure to call LinearizeTable first.
//
// Parameters:
//   - infos: The list to search.
//
// Returns:
//   - *Info: The last terminal node, or nil if not found.
func FindLastTerminal(infos []*Info) *Info {
	if len(infos) == 0 {
		return nil
	}

	idx := -1

	for i := 0; i < len(infos) && idx == -1; i++ {
		gers.AssertNotNil(infos[i], "tokens[i]")

		info := infos[i]
		if info.Type == NonterminalTk {
			idx = i
		}
	}

	if idx == -1 {
		return infos[len(infos)-1]
	} else if idx == 0 {
		return nil
	}

	return infos[idx-1]
}

// CheckEofExists checks if EOF exists in the given list.
//
// Parameters:
//   - infos: The list of info objects.
//
// Returns:
//   - bool: True if EOF exists in the list. False otherwise.
func CheckEofExists(infos []*Info) bool {
	if len(infos) == 0 {
		return false
	}

	for _, info := range infos {
		gers.AssertNotNil(info, "tk")

		node := info.Node()
		gers.AssertNotNil(node, "node")

		if node.Data == "EOF" {
			return true
		}
	}

	return false
}
