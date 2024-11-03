package generator

import (
	"github.com/PlayerR9/SlParser/cmd/grammar/pkg"
	common "github.com/PlayerR9/mygo-lib/common"
	cgen "github.com/PlayerR9/mygo-lib/generator"
	gslc "github.com/PlayerR9/mygo-lib/slices"
)

type ASTData struct {
	PackageName string
	Symbols     []string
}

func (ad *ASTData) SetPkgName(pkg_name string) error {
	if ad == nil {
		return common.ErrNilReceiver
	}

	ad.PackageName = pkg_name

	return nil
}

func NewASTData(rules []*pkg.Rule) *ASTData {
	var symbols []string

	for _, rule := range rules {
		_, _ = gslc.Merge(&symbols, rule.Symbols())
	}

	return &ASTData{
		Symbols: symbols,
	}
}

var (
	ASTGenerator *cgen.CodeGenerator[*ASTData]
)

func init() {
	ASTGenerator = cgen.Must(cgen.New[*ASTData]("ast", ast_templ))
}

const ast_templ string = `package {{ .PackageName }}

import (
	"errors"
	"github.com/PlayerR9/SlParser/ast"
	"github.com/PlayerR9/SlParser/grammar"
)

var (
	Ast ast.AST[*Node]
)

func init() {
	var builder ast.Builder[*Node]
	defer builder.Reset()
	{{ range $index, $symbol := .Symbols }}
	builder.Register({{ $symbol }}, func(token *grammar.Token) ([]*Node, error) {
		if token == nil {
			return nil, errors.New("token must not be nil")
		}

		// TODO: Write here the logic for turning the token into an AST node...

		panic("implement me")
	})
	{{ end }}
	Ast = builder.Build()
}`
