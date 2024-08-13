// Code generated by SlParser.
package test

import (
	"github.com/PlayerR9/grammar/ast"
	gr "github.com/PlayerR9/grammar/grammar"
)

// NodeType represents the type of a node in the AST tree.
type NodeType int

const (
	SourceNode NodeType = iota

	// Add here your custom node names...
	IdentifierNode
	OrExpr1Node
	OrExprNode
	RhsClsNode
	RhsNode
	Rule1Node
	RuleNode
	Source1Node

	// Add here your custom node types.
)

// String implements the NodeTyper interface.
func (t NodeType) String() string {
	return [...]string{
		"Source",
		// Add here your custom node names.
	}[t]
}

var (
	// ast_builder is the AST builder of the parser.
	ast_builder *ast.Make[*Node, token_type]
)

func init() {
	ast_builder = ast.NewMake[*Node, token_type]()

	parts := ast.NewPartsBuilder[*Node]()

	// Add here your custom AST builder rules...
		
	// ntk_Rule : ttk_UppercaseId ttk_Equal ntk_RhsCls ttk_Dot .
	// ntk_Rule : ttk_UppercaseId ttk_Equal ntk_RhsCls ttk_Rule1 .

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		root := prev.(*gr.Token[token_type])

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return nil, err
		}
		
		if len(children) != 4 {
			return nil, NewErrInvalidNumberOfChildren([]int{4}, len(children))
		}

		var sub_nodes []ast.Noder

		// Extract here any desired sub-node...

		n := NewNode(RuleNode, "", children[0].At)
		a.SetNode(&n)
		_ = a.AppendChildren(sub_nodes)

		return nil, nil
	})
		
	// ntk_Rhs : ntk_Identifier .
	// ntk_Rhs : ttk_OpParen ntk_OrExpr ttk_ClParen .

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		root := prev.(*gr.Token[token_type])

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return nil, err
		}

		switch len(children) {
				case 1:
			var sub_nodes []ast.Noder
		
			// Extract here any desired sub-node...
		
			n := NewNode(RhsNode, "", children[0].At)
			a.SetNode(&n)
			_ = a.AppendChildren(sub_nodes)
		case 3:
			var sub_nodes []ast.Noder
		
			// Extract here any desired sub-node...
		
			n := NewNode(RhsNode, "", children[0].At)
			a.SetNode(&n)
			_ = a.AppendChildren(sub_nodes)
		default:
			return nil, NewErrInvalidNumberOfChildren([]int{1, 3}, len(children))
		}

		return nil, nil
	})
	
	ast_builder.AddEntry(ntk_Rhs, parts.Build())
	parts.Reset()
		
	// ntk_OrExpr : ntk_Identifier ntk_OrExpr1 .

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		root := prev.(*gr.Token[token_type])

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return nil, err
		}
		
		if len(children) != 2 {
			return nil, NewErrInvalidNumberOfChildren([]int{2}, len(children))
		}

		var sub_nodes []ast.Noder

		// Extract here any desired sub-node...

		n := NewNode(OrExprNode, "", children[0].At)
		a.SetNode(&n)
		_ = a.AppendChildren(sub_nodes)

		return nil, nil
	})
		
	// ntk_OrExpr1 : ttk_Pipe ntk_Identifier .
	// ntk_OrExpr1 : ttk_Pipe ntk_Identifier ttk_Or ttk_Xpr1 .

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		root := prev.(*gr.Token[token_type])

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return nil, err
		}

		switch len(children) {
				case 2:
			var sub_nodes []ast.Noder
		
			// Extract here any desired sub-node...
		
			n := NewNode(OrExpr1Node, "", children[0].At)
			a.SetNode(&n)
			_ = a.AppendChildren(sub_nodes)
		case 4:
			var sub_nodes []ast.Noder
		
			// Extract here any desired sub-node...
		
			n := NewNode(OrExpr1Node, "", children[0].At)
			a.SetNode(&n)
			_ = a.AppendChildren(sub_nodes)
		default:
			return nil, NewErrInvalidNumberOfChildren([]int{2, 4}, len(children))
		}

		return nil, nil
	})
	
	ast_builder.AddEntry(ntk_OrExpr1, parts.Build())
	parts.Reset()
		
	// ntk_Source : ntk_Source1 etk_EOF .

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		root := prev.(*gr.Token[token_type])

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return nil, err
		}
		
		if len(children) != 2 {
			return nil, NewErrInvalidNumberOfChildren([]int{2}, len(children))
		}

		var sub_nodes []ast.Noder

		// Extract here any desired sub-node...

		n := NewNode(SourceNode, "", children[0].At)
		a.SetNode(&n)
		_ = a.AppendChildren(sub_nodes)

		return nil, nil
	})
		
	// ntk_Source1 : ntk_Rule .
	// ntk_Source1 : ntk_Rule ntk_Source1 .

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		root := prev.(*gr.Token[token_type])

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return nil, err
		}

		switch len(children) {
				case 1:
			var sub_nodes []ast.Noder
		
			// Extract here any desired sub-node...
		
			n := NewNode(Source1Node, "", children[0].At)
			a.SetNode(&n)
			_ = a.AppendChildren(sub_nodes)
		case 2:
			var sub_nodes []ast.Noder
		
			// Extract here any desired sub-node...
		
			n := NewNode(Source1Node, "", children[0].At)
			a.SetNode(&n)
			_ = a.AppendChildren(sub_nodes)
		default:
			return nil, NewErrInvalidNumberOfChildren([]int{1, 2}, len(children))
		}

		return nil, nil
	})
	
	ast_builder.AddEntry(ntk_Source1, parts.Build())
	parts.Reset()
		
	// ntk_Rule1 : ttk_Pipe ntk_RhsCls .
	// ntk_Rule1 : ttk_Pipe ntk_RhsCls ntk_Rule1 .

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		root := prev.(*gr.Token[token_type])

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return nil, err
		}

		switch len(children) {
				case 2:
			var sub_nodes []ast.Noder
		
			// Extract here any desired sub-node...
		
			n := NewNode(Rule1Node, "", children[0].At)
			a.SetNode(&n)
			_ = a.AppendChildren(sub_nodes)
		case 3:
			var sub_nodes []ast.Noder
		
			// Extract here any desired sub-node...
		
			n := NewNode(Rule1Node, "", children[0].At)
			a.SetNode(&n)
			_ = a.AppendChildren(sub_nodes)
		default:
			return nil, NewErrInvalidNumberOfChildren([]int{2, 3}, len(children))
		}

		return nil, nil
	})
	
	ast_builder.AddEntry(ntk_Rule1, parts.Build())
	parts.Reset()
		
	// ntk_RhsCls : ntk_Rhs .
	// ntk_RhsCls : ntk_Rhs ntk_RhsCls .

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		root := prev.(*gr.Token[token_type])

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return nil, err
		}

		switch len(children) {
				case 1:
			var sub_nodes []ast.Noder
		
			// Extract here any desired sub-node...
		
			n := NewNode(RhsClsNode, "", children[0].At)
			a.SetNode(&n)
			_ = a.AppendChildren(sub_nodes)
		case 2:
			var sub_nodes []ast.Noder
		
			// Extract here any desired sub-node...
		
			n := NewNode(RhsClsNode, "", children[0].At)
			a.SetNode(&n)
			_ = a.AppendChildren(sub_nodes)
		default:
			return nil, NewErrInvalidNumberOfChildren([]int{1, 2}, len(children))
		}

		return nil, nil
	})
	
	ast_builder.AddEntry(ntk_RhsCls, parts.Build())
	parts.Reset()
		
	// ntk_Identifier : ttk_UppercaseId .
	// ntk_Identifier : ttk_LowercaseId .

	parts.Add(func(a *ast.Result[*Node], prev any) (any, error) {
		root := prev.(*gr.Token[token_type])

		children, err := ast.ExtractChildren(root)
		if err != nil {
			return nil, err
		}
		
		if len(children) != 1 {
			return nil, NewErrInvalidNumberOfChildren([]int{1}, len(children))
		}

		var sub_nodes []ast.Noder

		// Extract here any desired sub-node...

		n := NewNode(IdentifierNode, "", children[0].At)
		a.SetNode(&n)
		_ = a.AppendChildren(sub_nodes)

		return nil, nil
	})
}