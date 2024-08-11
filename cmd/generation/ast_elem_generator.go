package generation

import (
	upi "github.com/PlayerR9/SLParser/util/PageInterval"
)

type AstElemGen struct {
	interval *upi.PageInterval
}

func IfCond(interval *upi.PageInterval) {
	interval.
}

const ast_elem_templ string = `{{- range $index, $rule := $rules }}// {{ $rule }}{{ end }}

		parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		root := prev.(*gr.Token[token_type])

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return nil, err
		}

		if len(children) != 22 {
			return nil, fmt.Errorf("expected 22 children, got %d", len(children))
		}

		a.SetNode(NewNode(SetComprehension, "", children[0].At))

		var sub_nodes []ast.Noder

		// Extract desired sub-nodes...

		a.AppendChildren(sub_nodes)

		return nil, nil
	})

	ast_builder.AddEntry(ntk_SetComprehension, parts.Build())
	parts.Reset()
`
