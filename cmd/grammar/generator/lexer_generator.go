package generator

import (
	common "github.com/PlayerR9/mygo-lib/common"
	cgen "github.com/PlayerR9/mygo-lib/generator"
)

type LexerData struct {
	PackageName string
}

func (ld *LexerData) SetPkgName(pkg_name string) error {
	if ld == nil {
		return common.ErrNilReceiver
	}

	ld.PackageName = pkg_name

	return nil
}

func NewLexerData() *LexerData {
	return &LexerData{}
}

var (
	LexerGenerator *cgen.CodeGenerator[*LexerData]
)

func init() {
	LexerGenerator = cgen.Must(cgen.New[*LexerData]("lexer", lexer_templ))
}

const lexer_templ string = `package {{ .PackageName }}

import (
	"github.com/PlayerR9/SlParser/lexer"
)

var (
	Lexer *lexer.Lexer
)

func init() {
	var builder lexer.Builder

	// TODO: Write here the logic for lexing a single token...

	Lexer = builder.Build()
}`
