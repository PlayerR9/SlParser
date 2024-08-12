// Code generated by SlParser.
package test

import (
	"fmt"

	"github.com/PlayerR9/grammar/grammar"
	"github.com/PlayerR9/grammar/parsing"
)

var (
	// internal_parser is the parser of the grammar.
	internal_parser *parsing.Parser[token_type]
)

func init() {
	decision_func := func(p *parsing.Parser[token_type], lookahead *grammar.Token[token_type]) (parsing.Actioner, error) {
		top1, ok := p.Pop()
		if !ok {
			return nil, fmt.Errorf("p.stack is empty")
		}

		var act parsing.Actioner

		switch top1.Type {
		case etk_EOF:
			// [ etk_EOF ] ntk_Source1 -> ntk_Source : ACCEPT .

			act, _ = parsing.NewAcceptAction(parsing.NewRule(ntk_Source, []token_type{etk_EOF, ntk_Source1}))
		case ntk_Identifier:
			// [ ntk_Identifier ] ttk_Pipe ntk_Identifier -> ntk_OrExpr : REDUCE .
			// ntk_Identifier ttk_Pipe [ ntk_Identifier ] -> ntk_OrExpr : SHIFT .
			// ntk_OrExpr ttk_Pipe [ ntk_Identifier ] -> ntk_OrExpr : SHIFT .
			// [ ntk_Identifier ] -> ntk_Rhs : REDUCE .

			panic("not implemented")
		case ntk_OrExpr:
			// ttk_ClParen [ ntk_OrExpr ] ttk_OpParen -> ntk_Rhs : SHIFT .
			// [ ntk_OrExpr ] ttk_Pipe ntk_Identifier -> ntk_OrExpr : REDUCE .

			panic("not implemented")
		case ntk_Rhs:
			// [ ntk_Rhs ] -> ntk_RhsCls : REDUCE .
			// ntk_RhsCls [ ntk_Rhs ] -> ntk_RhsCls : SHIFT .

			panic("not implemented")
		case ntk_RhsCls:
			// ttk_Dot [ ntk_RhsCls ] ttk_Equal ttk_UppercaseId -> ntk_Rule : SHIFT .
			// ntk_RuleLine [ ntk_RhsCls ] ttk_Equal ttk_Newline ttk_UppercaseId -> ntk_Rule : SHIFT .
			// ntk_RuleLine [ ntk_RhsCls ] ttk_Pipe ttk_Newline -> ntk_RuleLine : SHIFT .
			// [ ntk_RhsCls ] ntk_Rhs -> ntk_RhsCls : REDUCE .

			panic("not implemented")
		case ntk_Rule:
			// [ ntk_Rule ] -> ntk_Source1 : REDUCE .
			// ntk_Source1 ttk_Newline [ ntk_Rule ] -> ntk_Source1 : SHIFT .

			panic("not implemented")
		case ntk_RuleLine:
			// [ ntk_RuleLine ] ntk_RhsCls ttk_Equal ttk_Newline ttk_UppercaseId -> ntk_Rule : REDUCE .
			// [ ntk_RuleLine ] ntk_RhsCls ttk_Pipe ttk_Newline -> ntk_RuleLine : REDUCE .

			panic("not implemented")
		case ntk_Source1:
			// etk_EOF [ ntk_Source1 ] -> ntk_Source : SHIFT .
			// [ ntk_Source1 ] ttk_Newline ntk_Rule -> ntk_Source1 : REDUCE .

			panic("not implemented")
		case ttk_ClParen:
			// [ ttk_ClParen ] ntk_OrExpr ttk_OpParen -> ntk_Rhs : REDUCE .

			act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Rhs, []token_type{ttk_ClParen, ntk_OrExpr, ttk_OpParen}))
		case ttk_Dot:
			// [ ttk_Dot ] ntk_RhsCls ttk_Equal ttk_UppercaseId -> ntk_Rule : REDUCE .
			// [ ttk_Dot ] ttk_Newline -> ntk_RuleLine : REDUCE .

			panic("not implemented")
		case ttk_Equal:
			// ttk_Dot ntk_RhsCls [ ttk_Equal ] ttk_UppercaseId -> ntk_Rule : SHIFT .
			// ntk_RuleLine ntk_RhsCls [ ttk_Equal ] ttk_Newline ttk_UppercaseId -> ntk_Rule : SHIFT .

			act = parsing.NewShiftAction()
		case ttk_LowercaseId:
			// [ ttk_LowercaseId ] -> ntk_Identifier : REDUCE .

			act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Identifier, []token_type{ttk_LowercaseId}))
		case ttk_Newline:
			// ntk_Source1 [ ttk_Newline ] ntk_Rule -> ntk_Source1 : SHIFT .
			// ntk_RuleLine ntk_RhsCls ttk_Equal [ ttk_Newline ] ttk_UppercaseId -> ntk_Rule : SHIFT .
			// ntk_RuleLine ntk_RhsCls ttk_Pipe [ ttk_Newline ] -> ntk_RuleLine : SHIFT .
			// ttk_Dot [ ttk_Newline ] -> ntk_RuleLine : SHIFT .

			act = parsing.NewShiftAction()
		case ttk_OpParen:
			// ttk_ClParen ntk_OrExpr [ ttk_OpParen ] -> ntk_Rhs : SHIFT .

			act = parsing.NewShiftAction()
		case ttk_Pipe:
			// ntk_RuleLine ntk_RhsCls [ ttk_Pipe ] ttk_Newline -> ntk_RuleLine : SHIFT .
			// ntk_Identifier [ ttk_Pipe ] ntk_Identifier -> ntk_OrExpr : SHIFT .
			// ntk_OrExpr [ ttk_Pipe ] ntk_Identifier -> ntk_OrExpr : SHIFT .

			act = parsing.NewShiftAction()
		case ttk_UppercaseId:
			// ntk_RuleLine ntk_RhsCls ttk_Equal ttk_Newline [ ttk_UppercaseId ] -> ntk_Rule : SHIFT .
			// [ ttk_UppercaseId ] -> ntk_Identifier : REDUCE .
			// ttk_Dot ntk_RhsCls ttk_Equal [ ttk_UppercaseId ] -> ntk_Rule : SHIFT .

			panic("not implemented")
		default:
			return nil, fmt.Errorf("unexpected token: %s", top1.String())
		}

		return act, nil
	}

	internal_parser = parsing.NewParser(decision_func)
}