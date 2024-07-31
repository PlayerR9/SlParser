package Parser

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	uast "github.com/PlayerR9/EbnfParser/cmd/util/ast"
	ast "github.com/PlayerR9/grammar/ast"
	luch "github.com/PlayerR9/lib_units/runes"
	lus "github.com/PlayerR9/lib_units/slices"
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
//   - root: The root node of the AST tree.
//
// Returns:
//   - EnumType: The type of enum.
func GetEnumType(root *ast.Node[NodeType]) EnumType {
	if root == nil || root.Data == "" {
		return NotEnum
	}

	if root.Data == "EOF" {
		return SpecialEnum
	}

	first_letter, _ := utf8.DecodeRuneInString(root.Data)
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

		builder.WriteString("Etk")
		builder.WriteString(str)

		return builder.String(), nil
	case LexerEnum:
		chars, err := luch.StringToUtf8(str)
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

		chars = lus.SliceFilter(chars, func(r rune) bool {
			return r != '_'
		})

		var builder strings.Builder

		builder.WriteString("Ttk")

		for _, c := range chars {
			builder.WriteRune(c)
		}

		return builder.String(), nil
	case ParserEnum:
		var builder strings.Builder

		builder.WriteString("Ntk")
		builder.WriteString(str)

		return builder.String(), nil
	default:
		return "", fmt.Errorf("invalid enum type: %d", t_type)
	}
}

type ExtractEnumsData struct {
	special_enums []string
	lexer_enums   []string
	parser_enums  []string
}

func (d *ExtractEnumsData) AddLexerEnum(enum string) error {
	new_enum, err := ToEnum(enum, LexerEnum)
	if err != nil {
		return err
	}

	// Replace this with slices.TryInsert when lib_units is updated:
	// lexer_enums = lus.TryInsert(lexer_enums, top.Data)

	pos, ok := slices.BinarySearch(d.lexer_enums, new_enum)
	if !ok {
		d.lexer_enums = slices.Insert(d.lexer_enums, pos, new_enum)
	}

	return nil
}

func (d *ExtractEnumsData) AddParserEnum(enum string) error {
	new_enum, err := ToEnum(enum, ParserEnum)
	if err != nil {
		return err
	}

	// Replace this with slices.TryInsert when lib_units is updated:
	// parser_enums = lus.TryInsert(parser_enums, top.Data)

	pos, ok := slices.BinarySearch(d.parser_enums, new_enum)
	if !ok {
		d.parser_enums = slices.Insert(d.parser_enums, pos, new_enum)
	}

	return nil
}

func (d *ExtractEnumsData) AddSpecialEnum(enum string) error {
	new_enum, err := ToEnum(enum, SpecialEnum)
	if err != nil {
		return err
	}

	// Replace this with slices.TryInsert when lib_units is updated:
	// special_enums = lus.TryInsert(special_enums, top.Data)

	pos, ok := slices.BinarySearch(d.special_enums, new_enum)
	if !ok {
		d.special_enums = slices.Insert(d.special_enums, pos, new_enum)
	}

	return nil
}

func (d *ExtractEnumsData) GetLexerEnums() []string {
	return d.lexer_enums
}

func (d *ExtractEnumsData) GetParserEnums() []string {
	return d.parser_enums
}

func (d *ExtractEnumsData) GetSpecialEnums() []string {
	return d.special_enums
}

var (
	ExtractEnums *uast.SimpleDFS[*ast.Node[NodeType], *ExtractEnumsData]
)

func init() {
	ee_do := func(node *ast.Node[NodeType], data *ExtractEnumsData) error {
		e_type := GetEnumType(node)

		switch e_type {
		case LexerEnum:
			err := data.AddLexerEnum(node.Data)
			if err != nil {
				return fmt.Errorf("in node %q: %w", node.Data, err)
			}
		case ParserEnum:
			err := data.AddParserEnum(node.Data)
			if err != nil {
				return fmt.Errorf("in node %q: %w", node.Data, err)
			}
		case SpecialEnum:
			err := data.AddSpecialEnum(node.Data)
			if err != nil {
				return fmt.Errorf("in node %q: %w", node.Data, err)
			}
		}

		return nil
	}

	ee_init := func() *ExtractEnumsData {
		return &ExtractEnumsData{
			lexer_enums:  make([]string, 0),
			parser_enums: make([]string, 0),
		}
	}

	ExtractEnums = uast.NewSimpleDFS(ee_do, ee_init)
}

var (
	RenameNodes *uast.SimpleDFS[*ast.Node[NodeType], any]
)

func init() {
	f := func(node *ast.Node[NodeType], data any) error {
		e_type := GetEnumType(node)

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

	RenameNodes = uast.NewSimpleDFS(f, nil)
}
