package internal

import (
	gers "github.com/PlayerR9/go-errors"
	"github.com/PlayerR9/go-generator"
)

type GenGen struct {
	PackageName string
	Args        []string
}

func (g *GenGen) SetPackageName(pkg_name string) {
	if g == nil {
		return
	}

	g.PackageName = pkg_name
}

func NewGenGen() *GenGen {
	args := []string{
		"generate stringer -type=NodeType -linecomment",
		"generate stringer -type=TokenType",
	}

	return &GenGen{
		Args: args,
	}
}

var (
	GenGenerator *generator.CodeGenerator[*GenGen]
)

func init() {
	var err error

	GenGenerator, err = generator.NewCodeGeneratorFromTemplate[*GenGen]("gen", gen_templ)
	gers.AssertErr(err, "generator.NewCodeGeneratorFromTemplate[*GenGen](%q, gen_templ)", "gen")
}

const gen_templ string = `// Code generated by SlParser. Do not edit.
package {{ .PackageName }}
{{ range $index, $value := .Args }}
//go:{{$value}}
{{- end}}
`