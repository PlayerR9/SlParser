package internal

import (
	gers "github.com/PlayerR9/go-errors"
	"github.com/PlayerR9/go-generator"
)

// LexerGen is a generator for the lexer.
type LexerGen struct{}

// SetPackageName implements the generator.PackageNameSetter interface.
//
// This function is a no-op.
func (gd *LexerGen) SetPackageName(pkg_name string) {}

// NewLexerGen creates a new lexer generator.
//
// Returns:
//   - *LexerGen: The new lexer generator. Never returns nil.
func NewLexerGen() *LexerGen {
	return &LexerGen{}
}

var (
	// LexerGenerator is the lexer generator.
	LexerGenerator *generator.CodeGenerator[*LexerGen]
)

func init() {
	var err error

	LexerGenerator, err = generator.NewCodeGeneratorFromTemplate[*LexerGen]("enum", lexer_templ)
	gers.AssertErr(err, "generator.NewCodeGeneratorFromTemplate[*LexerGen](%q, templ)", "lexer_templ")
}

// lexer_templ is the template for the lexer.
var lexer_templ string = `
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
