package generation

import (
	"bytes"
	"strconv"
	"strings"
	"text/template"

	"github.com/PlayerR9/SLParser/cmd/pkg"
	upi "github.com/PlayerR9/SLParser/util/PageInterval"
)

func aeg_make_target(key string) (string, bool) {
	if !strings.HasPrefix(key, "ntk_") {
		return "", false
	}

	target := strings.TrimPrefix(key, "ntk_")
	target += "Node"

	return target, true
}

type AstElemGen struct {
	Key   string
	Rules []string

	interval *upi.PageInterval
	Expected string
	Target   string
	Lengths  []string
}

func (aeg AstElemGen) String() string {
	t := template.Must(template.New("ast_elem").Parse(ast_elem_templ))

	var buff bytes.Buffer

	err := t.Execute(&buff, aeg)
	if err != nil {
		panic(err.Error())
	}

	return buff.String()
}

func NewAstElemGen(key string, rules []*pkg.Rule) *AstElemGen {
	target, ok := aeg_make_target(key)
	if !ok {
		return nil
	}

	aeg := &AstElemGen{
		Key:      key,
		interval: upi.NewPageInterval(),
		Target:   target,
	}

	for _, rule := range rules {
		_ = aeg.interval.AddPage(rule.Size())
		aeg.Rules = append(aeg.Rules, rule.StringOriginal())
	}

	aeg.if_cond()

	return aeg
}

func (aeg *AstElemGen) if_cond() {
	intervals := aeg.interval.Intervals()
	if len(intervals) == 0 {
		return
	}

	var builder strings.Builder

	expected := make([]string, 0, aeg.interval.PageCount())

	iter := aeg.interval.Iterator()

	for {
		value, err := iter.Consume()
		if err != nil {
			break
		}

		expected = append(expected, strconv.Itoa(value))
	}

	builder.WriteString("[]int{")
	builder.WriteString(strings.Join(expected, ", "))
	builder.WriteRune('}')

	aeg.Expected = builder.String()
	aeg.Lengths = expected
}

const ast_elem_templ string = `
{{- range $index, $rule := .Rules }}
	// {{ $rule }}
{{- end }}

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		root := prev.(*gr.Token[token_type])

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return nil, err
		}

		{{ if eq (len .Lengths) 1 }}
			if len(children) != {{ index .Lengths 0 }} {
				return nil, NewErrInvalidNumberOfChildren({{ .Expected }}, len(children))
			}

			var sub_nodes []ast.Noder

			// Extract here any desired sub-node...

			a.SetNode(NewNode({{ .Target }}, "", children[0].At))
			_ = a.AppendChildren(sub_nodes)
		{{ else if gt (len .Lengths) 1 }}
			switch len(children) {
			{{ range $index, $length := .Lengths }}
			case {{ $length }}:
				var sub_nodes []ast.Noder

				// Extract here any desired sub-node...

				a.SetNode(NewNode({{ .Target }}, "", children[0].At))
				_ = a.AppendChildren(sub_nodes)
			{{ end }}
			default:
				return nil, NewErrInvalidNumberOfChildren({{ .Expected }}, len(children))
			}
		{{ end }}

		return nil, nil
	})
		
	ast_builder.AddEntry({{ .Key }}, parts.Build())
	parts.Reset()`
