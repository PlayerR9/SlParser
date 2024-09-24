package internal

import (
	"errors"

	kdd "github.com/PlayerR9/SlParser/kdd"
	"github.com/PlayerR9/go-generator"
)

type TokenGen struct {
	PackageName  string
	Symbols      []string
	LastTerminal string
	Rules        []string
}

func (gd *TokenGen) SetPackageName(pkg_name string) {
	if gd == nil {
		return
	}

	gd.PackageName = pkg_name
}

func NewTokenGen(tokens []*kdd.Node) (*TokenGen, error) {
	var symbols []string

	if len(tokens) > 0 {
		symbols = make([]string, 0, len(tokens))

		for _, tk := range tokens {
			if tk == nil {
				continue
			}

			symbols = append(symbols, tk.Data)
		}
	}

	lt, err := FindLastTerminal(tokens)
	if err != nil {
		return nil, err
	} else if lt == nil {
		return nil, errors.New("missing terminal")
	}

	gd := &TokenGen{
		Symbols:      symbols,
		LastTerminal: lt.Data,
	}

	return gd, nil
}

var (
	TokenGenerator *generator.CodeGenerator[*TokenGen]
)

func init() {
	var err error

	TokenGenerator, err = generator.NewCodeGeneratorFromTemplate[*TokenGen]("enum", token_templ)
	if err != nil {
		panic(err)
	}
}

const token_templ string = `// Code generated by SlParser. Do not edit.
package internal

import (
	"github.com/PlayerR9/SlParser/parser"
)

//go:generate stringer -type=TokenType

type TokenType int

const (
	EtInvalid TokenType = iota -1{{ range $index, $value := .Symbols }}
	{{ $value }}
	{{- end }}
)

func (t TokenType) IsTerminal() bool {
	return t <= {{ .LastTerminal }}
}
	
var (
	Parser *parser.Parser[TokenType]
)

func init() {
	is := parser.NewItemSet[TokenType]()
	{{ range $index, $value := .Rules }}
	{{ $value }}
	{{- end }}

	Parser = parser.Build(&is)
}`
