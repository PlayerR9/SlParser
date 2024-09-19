package parsing

import (
	gr "github.com/PlayerR9/SlParser/grammar"
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
	builder := prx.NewBuilder(&is)

	builder.Register(NttStatement, func(parser *prx.Parser[TokenType], top1 *gr.ParseTree[TokenType], lookahead *gr.Token[TokenType]) ([]*prx.Item[TokenType], error) {
		has_la := lookahead != nil && lookahead.Type == TttNewline

		var it1 *prx.Item[TokenType]

		if has_la {
			// source1 : statement # NEWLINE source1 ;

			it1 = prx.MustNewItem(rule2, 0)
		} else {
			// source1 : statement # ;

			it1 = prx.MustNewItem(rule1, 0)
		}

		return []*prx.Item[TokenType]{it1}, nil
	})

	/* builder.Register(NttSource1, func(parser *prx.Parser[TokenType], top1, lookahead *gr.Token[TokenType]) ([]*prx.Item[TokenType], error) {
		// source : NEWLINE source1 # EOF ;

		// source1 : statement NEWLINE source1 # ;
	}) */

	Parser = builder.Build()
}

// PrintItemSet prints the item set.
//
// Returns:
//   - []string: the lines of the item set.
func PrintItemSet() []string {
	return is.PrintTable()
}
