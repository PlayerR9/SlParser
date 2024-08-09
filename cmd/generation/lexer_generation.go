package generation

import (
	ggen "github.com/PlayerR9/go-generator/generator"
)

type LexerGen struct {
	PackageName string
}

func (g *LexerGen) SetPackageName(pkg_name string) {
	g.PackageName = pkg_name
}

var (
	LexerGenerator *ggen.CodeGenerator[*LexerGen]
)

func init() {
	tmp, err := ggen.NewCodeGeneratorFromTemplate[*LexerGen]("", lexer_templ)
	if err != nil {
		Logger.Fatalf("Error creating code generator: %s", err.Error())
	}

	LexerGenerator = tmp
}

const lexer_templ string = `// Code generated by SlParser.
package {{ .PackageName }}

import (
	"github.com/PlayerR9/grammar/grammar"
	"github.com/PlayerR9/grammar/lexing"
)

var (
	// matcher is the matcher of the grammar.
	matcher *lexing.Matcher[token_type]
)

func init() {
	matcher = lexing.NewMatcher[token_type]()

	// Add here your custom matcher rules.
}

var (
	// internal_lexer is the lexer of the grammar.
	internal_lexer *lexing.Lexer[token_type]
)

func init() {
	lex_one := func(l *lexing.Lexer[token_type]) (*grammar.Token[token_type], error) {
		at := l.Pos()

		match, _ := matcher.Match(l)

		if match.IsValidMatch() {
			symbol, data := match.GetMatch()

			return grammar.NewToken(symbol, data, at, nil), nil
		}

		// Lex here...
	
		panic("Implement me!")
	}

	internal_lexer = lexing.NewLexer(lex_one)
}`
