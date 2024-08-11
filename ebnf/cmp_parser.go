// Code generated by SlParser.
package ebnf

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
			top2, ok := p.Pop()
			if ok && top2.Type == ttk_Pipe {
				// [ ntk_Identifier ] ttk_Pipe ntk_Identifier -> ntk_OrExpr : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_OrExpr, []token_type{ntk_Identifier, ttk_Pipe, ntk_Identifier}))
			} else {
				if lookahead == nil || lookahead.Type != ttk_Pipe {
					// [ ntk_Identifier ] -> ntk_Rhs : REDUCE .

					act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Rhs, []token_type{ntk_Identifier}))
				} else {
					// ntk_Identifier ttk_Pipe [ ntk_Identifier ] -> ntk_OrExpr : SHIFT .
					// ntk_OrExpr ttk_Pipe [ ntk_Identifier ] -> ntk_OrExpr : SHIFT .

					act = parsing.NewShiftAction()
				}
			}
		case ntk_OrExpr:
			top2, ok := p.Pop()
			if ok && top2.Type == ttk_Pipe {
				// [ ntk_OrExpr ] ttk_Pipe ntk_Identifier -> ntk_OrExpr : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_OrExpr, []token_type{ntk_OrExpr, ttk_Pipe, ntk_Identifier}))
			} else {
				// ttk_ClParen [ ntk_OrExpr ] ttk_OpParen -> ntk_Rhs : SHIFT .

				act = parsing.NewShiftAction()
			}
		case ntk_Rhs:
			if lookahead == nil || (lookahead.Type != ttk_LowercaseId && lookahead.Type != ttk_UppercaseId && lookahead.Type != ttk_OpParen) {
				// [ ntk_Rhs ] -> ntk_RhsCls : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_RhsCls, []token_type{ntk_Rhs}))
			} else {
				// ntk_RhsCls [ ntk_Rhs ] -> ntk_RhsCls : SHIFT .
				// -- [ ntk_Identifier ] -> ntk_Rhs : REDUCE .
				// -- -- [ ttk_LowercaseId ] -> ntk_Identifier : REDUCE .
				// -- -- [ ttk_UppercaseId ] -> ntk_Identifier : REDUCE .
				// -- ttk_ClParen ntk_OrExpr ttk_OpParen -> ntk_Rhs : SHIFT .

				act = parsing.NewShiftAction()
			}
		case ntk_RhsCls:
			top2, ok := p.Pop()
			if ok && top2.Type == ntk_Rhs {
				// [ ntk_RhsCls ] ntk_Rhs -> ntk_RhsCls : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_RhsCls, []token_type{ntk_RhsCls, ntk_Rhs}))
			} else {
				// ttk_Dot [ ntk_RhsCls ] ttk_Equal ttk_UppercaseId -> ntk_Rule : SHIFT .
				// ntk_RuleLine [ ntk_RhsCls ] ttk_Equal ttk_Newline ttk_UppercaseId -> ntk_Rule : SHIFT .
				// ntk_RuleLine [ ntk_RhsCls ] ttk_Pipe ttk_Newline -> ntk_RuleLine : SHIFT .

				act = parsing.NewShiftAction()
			}
		case ntk_Rule:
			if lookahead == nil || lookahead.Type != ttk_Newline {
				// [ ntk_Rule ] -> ntk_Source1 : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Source1, []token_type{ntk_Rule}))
			} else {
				// ntk_Source1 ttk_Newline [ ntk_Rule ] -> ntk_Source1 : SHIFT .

				act = parsing.NewShiftAction()
			}
		case ntk_RuleLine:
			top2, ok := p.Pop()
			if !ok {
				return nil, parsing.NewErrUnexpectedToken(&top1.Type, nil, ntk_RhsCls)
			} else if top2.Type != ntk_RhsCls {
				return nil, parsing.NewErrUnexpectedToken(&top1.Type, &top2.Type, ntk_RhsCls)
			}

			top3, ok := p.Pop()
			if !ok {
				return nil, parsing.NewErrUnexpectedToken(&top2.Type, nil, ttk_Equal, ttk_Pipe)
			}

			switch top3.Type {
			case ttk_Equal:
				// [ ntk_RuleLine ] ntk_RhsCls ttk_Equal ttk_Newline ttk_UppercaseId -> ntk_Rule : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Rule, []token_type{ntk_RuleLine, ntk_RhsCls, ttk_Equal, ttk_Newline, ttk_UppercaseId}))
			case ttk_Pipe:
				// [ ntk_RuleLine ] ntk_RhsCls ttk_Pipe ttk_Newline -> ntk_RuleLine : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_RuleLine, []token_type{ntk_RuleLine, ntk_RhsCls, ttk_Pipe, ttk_Newline}))
			default:
				return nil, parsing.NewErrUnexpectedToken(&top3.Type, &top2.Type, ttk_Equal, ttk_Pipe)
			}
		case ntk_Source1:
			top2, ok := p.Pop()
			if ok && top2.Type == ttk_Newline {
				// [ ntk_Source1 ] ttk_Newline ntk_Rule -> ntk_Source1 : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Source1, []token_type{ntk_Source1, ttk_Newline, ntk_Rule}))
			} else {
				// etk_EOF [ ntk_Source1 ] -> ntk_Source : SHIFT .

				act = parsing.NewShiftAction()
			}
		case ttk_ClParen:
			// [ ttk_ClParen ] ntk_OrExpr ttk_OpParen -> ntk_Rhs : REDUCE .

			act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Rhs, []token_type{ttk_ClParen, ntk_OrExpr, ttk_OpParen}))
		case ttk_Dot:
			top2, ok := p.Pop()
			if !ok {
				return nil, parsing.NewErrUnexpectedToken(&top1.Type, nil, ntk_RhsCls, ttk_Newline)
			}

			switch top2.Type {
			case ntk_RhsCls:
				// [ ttk_Dot ] ntk_RhsCls ttk_Equal ttk_UppercaseId -> ntk_Rule : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Rule, []token_type{ttk_Dot, ntk_RhsCls, ttk_Equal, ttk_UppercaseId}))
			case ttk_Newline:
				// [ ttk_Dot ] ttk_Newline -> ntk_RuleLine : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_RuleLine, []token_type{ttk_Dot, ttk_Newline}))
			default:
				return nil, parsing.NewErrUnexpectedToken(&top1.Type, &top2.Type, ntk_RhsCls, ttk_Newline)
			}
		case ttk_Equal:
			// ttk_Dot ntk_RhsCls [ ttk_Equal ] ttk_UppercaseId -> ntk_Rule : SHIFT .
			// ntk_RuleLine ntk_RhsCls [ ttk_Equal ] ttk_Newline ttk_UppercaseId -> ntk_Rule : SHIFT .

			act = parsing.NewShiftAction()
		case ttk_LowercaseId:
			// [ ttk_LowercaseId ] -> ntk_Identifier : REDUCE .

			act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Identifier, []token_type{ttk_LowercaseId}))
		case ttk_Newline:
			// ttk_Dot [ ttk_Newline ] -> ntk_RuleLine : SHIFT .
			// ntk_Source1 [ ttk_Newline ] ntk_Rule -> ntk_Source1 : SHIFT .
			// ntk_RuleLine ntk_RhsCls ttk_Equal [ ttk_Newline ] ttk_UppercaseId -> ntk_Rule : SHIFT .
			// ntk_RuleLine ntk_RhsCls ttk_Pipe [ ttk_Newline ] -> ntk_RuleLine : SHIFT .

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
			if lookahead == nil || (lookahead.Type != ttk_Newline && lookahead.Type != ttk_Equal) {
				// [ ttk_UppercaseId ] -> ntk_Identifier : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Identifier, []token_type{ttk_UppercaseId}))
			} else {
				// ttk_Dot ntk_RhsCls ttk_Equal [ ttk_UppercaseId ] -> ntk_Rule : SHIFT .
				// ntk_RuleLine ntk_RhsCls ttk_Equal ttk_Newline [ ttk_UppercaseId ] -> ntk_Rule : SHIFT .

				act = parsing.NewShiftAction()
			}
		default:
			return nil, fmt.Errorf("unexpected token: %s", top1.String())
		}

		return act, nil
	}

	internal_parser = parsing.NewParser(decision_func)
}
