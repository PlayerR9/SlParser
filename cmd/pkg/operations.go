package pkg

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	ebnf "github.com/PlayerR9/SLParser/ebnf"
	gcch "github.com/PlayerR9/go-commons/runes"
	gcslc "github.com/PlayerR9/go-commons/slices"
	dbg "github.com/PlayerR9/go-debug/assert"
	tr "github.com/PlayerR9/grammar/traversing"
)

// EnumType represents the type of enum.
type EnumType int

const (
	LexerEnum EnumType = iota
	ParserEnum
	SpecialEnum
	NotEnum
)

// GetEnumType returns the type of enum.
//
// Parameters:
//   - data: The data to parse.
//
// Returns:
//   - EnumType: The type of enum.
func GetEnumType(data string) EnumType {
	if data == "" {
		return NotEnum
	}

	if data == "EOF" {
		return SpecialEnum
	}

	first_letter, _ := utf8.DecodeRuneInString(data)
	if first_letter == utf8.RuneError || !unicode.IsLetter(first_letter) {
		return NotEnum
	}

	if unicode.IsLower(first_letter) {
		return LexerEnum
	} else {
		return ParserEnum
	}
}

// ToEnum converts a string to an enum.
//
// Parameters:
//   - str: The string to convert.
//   - t_type: The type of enum.
//
// Returns:
//   - string: The enum string.
//   - error: An error if the conversion failed.
func ToEnum(str string, t_type EnumType) (string, error) {
	switch t_type {
	case SpecialEnum:
		var builder strings.Builder

		builder.WriteString("etk_")
		builder.WriteString(str)

		return builder.String(), nil
	case LexerEnum:
		chars, err := gcch.StringToUtf8(str)
		if err != nil {
			return "", err
		}

		var indices []int

		indices = append(indices, -1) // force uppercase the first letter

		for i := 0; i < len(chars); i++ {
			if chars[i] == '_' {
				indices = append(indices, i)
			}
		}

		for _, idx := range indices {
			chars[idx+1] = unicode.ToUpper(chars[idx+1])
		}

		chars = gcslc.SliceFilter(chars, func(r rune) bool {
			return r != '_'
		})

		var builder strings.Builder

		builder.WriteString("ttk_")

		for _, c := range chars {
			builder.WriteRune(c)
		}

		return builder.String(), nil
	case ParserEnum:
		var builder strings.Builder

		builder.WriteString("ntk_")
		builder.WriteString(str)

		return builder.String(), nil
	default:
		return "", fmt.Errorf("invalid enum type: %d", t_type)
	}
}

type EnumExtractor struct {
	special_enums []string
	lexer_enums   []string
	parser_enums  []string
}

// Reset implements the traverser.Traverser interface.
func (ee *EnumExtractor) Reset() {
	ee.lexer_enums = ee.lexer_enums[:0]
	ee.parser_enums = ee.parser_enums[:0]
	ee.special_enums = ee.special_enums[:0]
}

// Copy implements the traverser.Traverser interface.
//
// Returns a pointer to itself.
func (ee EnumExtractor) Copy() tr.Traverser {
	return &ee
}

// Apply implements the traverser.Traverser interface.
func (ee *EnumExtractor) Apply(node tr.TreeNoder) ([]tr.TravData, error) {
	if node == nil {
		return nil, nil
	}

	n := dbg.AssertConv[*ebnf.Node](node, "node")

	e_type := GetEnumType(n.Data)

	switch e_type {
	case LexerEnum:
		err := ee.AddLexerEnum(n.Data)
		if err != nil {
			return nil, fmt.Errorf("in node %q: %w", n.Data, err)
		}
	case ParserEnum:
		err := ee.AddParserEnum(n.Data)
		if err != nil {
			return nil, fmt.Errorf("in node %q: %w", n.Data, err)
		}
	case SpecialEnum:
		err := ee.AddSpecialEnum(n.Data)
		if err != nil {
			return nil, fmt.Errorf("in node %q: %w", n.Data, err)
		}
	}

	var data []tr.TravData

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		d := tr.TravData{
			Node: c,
			Data: ee.Copy(),
		}

		data = append(data, d)
	}

	return data, nil
}

func (ee *EnumExtractor) AddLexerEnum(enum string) error {
	new_enum, err := ToEnum(enum, LexerEnum)
	if err != nil {
		return err
	}

	ee.lexer_enums = gcslc.TryInsert(ee.lexer_enums, new_enum)

	return nil
}

func (ee *EnumExtractor) AddParserEnum(enum string) error {
	new_enum, err := ToEnum(enum, ParserEnum)
	if err != nil {
		return err
	}

	ee.parser_enums = gcslc.TryInsert(ee.parser_enums, new_enum)

	return nil
}

func (ee *EnumExtractor) AddSpecialEnum(enum string) error {
	new_enum, err := ToEnum(enum, SpecialEnum)
	if err != nil {
		return err
	}

	ee.special_enums = gcslc.TryInsert(ee.special_enums, new_enum)

	return nil
}

func (ee EnumExtractor) GetLexerEnums() []string {
	return ee.lexer_enums
}

func (ee EnumExtractor) GetParserEnums() []string {
	return ee.parser_enums
}

func (ee EnumExtractor) GetSpecialEnums() []string {
	return ee.special_enums
}

var (
	RenameNodes tr.SimpleDFS[*ebnf.Node]
)

func init() {
	f := func(node *ebnf.Node) error {
		e_type := GetEnumType(node.Data)

		switch e_type {
		case LexerEnum:
			new_name, err := ToEnum(node.Data, LexerEnum)
			if err != nil {
				return fmt.Errorf("in node %q: %w", node.Data, err)
			}

			node.Data = new_name
		case ParserEnum:
			new_name, err := ToEnum(node.Data, ParserEnum)
			if err != nil {
				return fmt.Errorf("in node %q: %w", node.Data, err)
			}

			node.Data = new_name
		case SpecialEnum:
			new_name, err := ToEnum(node.Data, SpecialEnum)
			if err != nil {
				return fmt.Errorf("in node %q: %w", node.Data, err)
			}

			node.Data = new_name
		}

		return nil
	}

	RenameNodes.SetDoFunc(f)
}
