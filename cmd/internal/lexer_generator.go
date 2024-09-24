package internal

import "github.com/PlayerR9/go-generator"

type LexerGen struct {
	PackageName string
}

func (gd *LexerGen) SetPackageName(pkg_name string) {
	if gd == nil {
		return
	}

	gd.PackageName = pkg_name
}

var (
	LexerGenerator *generator.CodeGenerator[*LexerGen]
)

func init() {
	var err error

	LexerGenerator, err = generator.NewCodeGeneratorFromTemplate[*LexerGen]("enum", templ)
	if err != nil {
		panic(err)
	}
}

var templ string = `
package internal

import (
	"github.com/PlayerR9/SlParser/lexer"
)

var (
	Lexer *lexer.Lexer[TokenType]
)

func init() {
	builder := lexer.NewBuilder[TokenType]()

	// TODO: Add here your own custom rules...
	
	Lexer = builder.Build()
}`
