// Code generated by SlParser.
package ebnf

// token_type is the type of a token.
type token_type int

const (
	etk_EOF token_type = iota

	ttk_ClParen
	ttk_Equal
	ttk_LowercaseId
	ttk_OpParen
	ttk_Pipe
	ttk_Semicolon
	ttk_UppercaseId

	ntk_Identifier
	ntk_OrExpr
	ntk_OrExpr1
	ntk_Rhs
	ntk_Rhs1
	ntk_Rule
	ntk_Rule1
	ntk_Source
	ntk_Source1
)

// String implements the Grammar.TokenTyper interface.
func (t token_type) String() string {
	return [...]string{
		"End of File",
		// Add here your custom token names.

		"close parenthesis",
		"equal sign",
		"lowercase identifier",
		"open parenthesis",
		"pipe",
		"semicolon",
		"uppercase identifier",

		"Identifier",
		"OR expression",
		"OR expression (I)",
		"Right-hand side",
		"Right-hand side (I)",
		"Rule",
		"Rule (I)",
		"Source",
		"Source (I)",
	}[t]
}

// GoString implements the Grammar.TokenTyper interface.
func (t token_type) GoString() string {
	return [...]string{
		"etk_EOF",

		"ttk_ClParen",
		"ttk_Equal",
		"ttk_LowercaseId",
		"ttk_OpParen",
		"ttk_Pipe",
		"ttk_Rule1",
		"ttk_Semicolon",
		"ttk_UppercaseId",

		"ntk_Identifier",
		"ntk_OrExpr",
		"ntk_OrExpr1",
		"ntk_Rhs",
		"ntk_Rhs1",
		"ntk_Rule",
		"ntk_Rule1",
		"ntk_Source",
		"ntk_Source1",
	}[t]
}
