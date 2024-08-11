package generation

import (
	"bytes"
	"strconv"
	"strings"
	"text/template"

	"github.com/PlayerR9/SLParser/cmd/pkg"
	upi "github.com/PlayerR9/SLParser/util/PageInterval"
)

var (
	t *template.Template
)

func init() {
	t = template.Must(template.New("ast_elem").Parse(ast_elem_templ))
}

type AstElemGen struct {
	Key   string
	Rules []string

	interval *upi.PageInterval
	IfCond   string
	Expected string
	Target   string
}

func (aeg *AstElemGen) String() string {
	var buff bytes.Buffer

	err := t.Execute(&buff, aeg)
	if err != nil {
		panic(err.Error())
	}

	return buff.String()
}

func NewAstElemGen(key string, rules []*pkg.Rule) *AstElemGen {
	aeg := &AstElemGen{
		Key:      key,
		interval: upi.NewPageInterval(),
	}

	ok := aeg.make_target()
	if !ok {
		return nil
	}

	for _, rule := range rules {
		_ = aeg.interval.AddPage(rule.Size())
		aeg.Rules = append(aeg.Rules, rule.StringOriginal())
	}

	aeg.if_cond("len(children)")

	return aeg
}

func (aeg *AstElemGen) make_target() bool {
	if !strings.HasPrefix(aeg.Key, "ntk_") {
		return false
	}

	aeg.Target = strings.TrimPrefix(aeg.Key, "ntk_")
	aeg.Target += "Node"

	return true
}

func (aeg *AstElemGen) if_cond(thing string) {
	intervals := aeg.interval.Intervals()
	if len(intervals) == 0 {
		return
	}

	elems := make([]string, 0, len(intervals))
	var builder strings.Builder

	for _, pr := range intervals {
		if pr.First == pr.Second {
			// <thing> != pr.First
			builder.WriteString(thing)
			builder.WriteString(" != ")
			builder.WriteString(strconv.Itoa(pr.First))

			elems = append(elems, builder.String())
		} else {
			// (<thing> < pr.First || <thing> > pr.Second)

			builder.WriteRune('(')
			builder.WriteString(thing)
			builder.WriteString(" < ")
			builder.WriteString(strconv.Itoa(pr.First))
			builder.WriteString(" || ")

			builder.WriteString(thing)
			builder.WriteString(" > ")
			builder.WriteString(strconv.Itoa(pr.Second))
			builder.WriteRune(')')

			elems = append(elems, builder.String())

		}

		builder.Reset()
	}

	aeg.IfCond = strings.Join(elems, " && ")

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

		if {{ .IfCond }} {
			return nil, NewErrInvalidNumberOfChildren({{ .Expected }}, len(children))
		}

		var sub_nodes []ast.Noder

		// Extract here any desired sub-node...

		a.SetNode(NewNode({{ .Target }}, "", children[0].At))
		a.AppendChildren(sub_nodes)

		return nil, nil
	})
		
	ast_builder.AddEntry({{ .Key }}, parts.Build())
	parts.Reset()`
