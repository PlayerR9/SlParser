package parsing

import (
	prx "github.com/PlayerR9/SlParser/parser"
	dba "github.com/PlayerR9/go-debug/assert"
)

var (
	is prx.ItemSet[TokenType]

	// source1 : statement ;
	rule1 *prx.Rule[TokenType]

	// source1 : statement NEWLINE source1 ;
	rule2 *prx.Rule[TokenType]
)

func init() {
	is = prx.NewItemSet[TokenType]()

	// source : NEWLINE source1 EOF ;
	_, err := is.AddRule(NttSource, TttNewline, NttSource1, EttEOF)
	dba.AssertErr(err, "is.AddRule(NttSource, TttNewline, NttSource1, EttEOF)")

	// source1 : statement ;
	rule1, err = is.AddRule(NttSource1, NttStatement)
	dba.AssertErr(err, "is.AddRule(NttSource1, NttStatement)")

	// source1 : statement NEWLINE source1 ;
	rule2, err = is.AddRule(NttSource1, NttStatement, TttNewline, NttSource1)
	dba.AssertErr(err, "is.AddRule(NttSource1, NttStatement, TttNewline, NttSource1)")

	// statement : LIST_COMPREHENSION ;
	_, err = is.AddRule(NttStatement, TttListComprehension)
	dba.AssertErr(err, "is.AddRule(NttStatement, TttListComprehension)")

	// statement : PRINT_STMT ;
	_, err = is.AddRule(NttStatement, TttPrintStmt)
	dba.AssertErr(err, "is.AddRule(NttStatement, TttPrintStmt)")
}

var (
	Parser *prx.Parser[TokenType]
)

func init() {
	Parser = prx.Build(&is)
}

// PrintItemSet prints the item set.
//
// Returns:
//   - []string: the lines of the item set.
func PrintItemSet() []string {
	return is.PrintTable()
}
