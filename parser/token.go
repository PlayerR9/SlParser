package parser

// token_type is a type of token.
type token_type int

const (
	ttk_EOF token_type = iota
	ttk_Dot
	ttk_OpParen
	ttk_ClParen
	ttk_Pipe
	ttk_EqualSign
	ttk_Newline
	ttk_UppercaseID
	ttk_LowercaseID

	ntk_Source
	ntk_Source1
	ntk_Rule
	ntk_RhsCls
	ntk_RuleLine
	ntk_Rhs
	ntk_Identifier
	ntk_OrExpr
)

func (t token_type) String() string {
	return [...]string{
		"End of File",
		"dot",
		"open parenthesis",
		"close parenthesis",
		"pipe",
		"equal sign",
		"newline",
		"uppercase identifier",
		"lowercase identifier",

		"Source",
		"Source (I)",
		"Rule",
		"Rule clause",
		"Rule line",
		"Right-hand side",
		"Identifier",
		"OR expression",
	}[t]
}

func (t token_type) GoString() string {
	return [...]string{
		"TtkEOF",

		"TtkDot",
		"TtkOpParen",
		"TtkClParen",
		"TtkPipe",
		"TtkEqualSign",
		"TtkNewline",
		"TtkUppercaseID",
		"TtkLowercaseID",

		"NtkSource",
		"NtkSource1",
		"NtkRule",
		"NtkRhsCls",
		"NtkRuleLine",
		"NtkRhs",
		"NtkIdentifier",
		"NtkOrExpr",
	}[t]
}
