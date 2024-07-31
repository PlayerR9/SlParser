package pkg

// TokenType is a type of token.
type TokenType int

const (
	TtkEOF TokenType = iota
	TtkDot
	TtkOpParen
	TtkClParen
	TtkPipe
	TtkEqualSign
	TtkNewline
	TtkUppercaseID
	TtkLowercaseID

	NtkSource
	NtkSource1
	NtkRule
	NtkRhsCls
	NtkRuleLine
	NtkRhs
	NtkIdentifier
	NtkOrExpr
)

func (t TokenType) String() string {
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
