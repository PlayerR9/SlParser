// Code generated by SlParser.
package test

// token_type is the type of a token.
type token_type int

const (
	etk_EOF token_type = iota
	
	ttk_ClParen
	ttk_Dot
	ttk_Equal
	ttk_LowercaseId
	ttk_Newline
	ttk_OpParen
	ttk_Pipe
	ttk_UppercaseId
	
	ntk_Identifier
	ntk_OrExpr
	ntk_Rhs
	ntk_RhsCls
	ntk_Rule
	ntk_RuleLine
	ntk_Source
	ntk_Source1
)

// String implements the Grammar.TokenTyper interface.
func (t token_type) String() string {
	return [...]string{
		"End of File",
		// Add here your custom token names.
	}[t]
}

// GoString implements the Grammar.TokenTyper interface.
func (t token_type) GoString() string {
	return [...]string{
		"etk_EOF",
		
		"ttk_ClParen",
		"ttk_Dot",
		"ttk_Equal",
		"ttk_LowercaseId",
		"ttk_Newline",
		"ttk_OpParen",
		"ttk_Pipe",
		"ttk_UppercaseId",
		
		"ntk_Identifier",
		"ntk_OrExpr",
		"ntk_Rhs",
		"ntk_RhsCls",
		"ntk_Rule",
		"ntk_RuleLine",
		"ntk_Source",
		"ntk_Source1",
	}[t]
}