package pkg

import (
	"fmt"

	gr "github.com/PlayerR9/grammar/grammar"
	uprx "github.com/PlayerR9/grammar/parsing"
)

var (
	internal_parser *uprx.Parser[token_type]
)

func init() {
	f := func(parser *uprx.Parser[token_type], lookahead *gr.Token[token_type]) (uprx.Actioner, error) {
		top1, ok := parser.Pop()
		if !ok {
			return nil, fmt.Errorf("p.stack is empty")
		}

		var act uprx.Actioner

		switch top1.Type {
		case ttk_EOF:
			// [ EOF ] Source1 -> Source : accept .
			tmp, _ := uprx.NewAcceptAction(uprx.NewRule(ntk_Source, []token_type{ttk_EOF, ntk_Source1}))
			// luc.AssertErr(err, "NewAcceptAction(TkSource, []TokenType{TkEOF, TkSource1})")

			act = tmp
		case ntk_Source1:
			top2, ok := parser.Pop()
			if !ok || top2.Type != ttk_Newline {
				// EOF [ Source1 ] -> Source : shift .

				act = uprx.NewShiftAction()
			} else {
				// [ Source1 ] newline Rule -> Source1 : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(ntk_Source1, []token_type{ntk_Source1, ttk_Newline, ntk_Rule}))
				// luc.AssertErr(err, "NewReduceAction(TkSource1, []TokenType{TkSource1, TkNewline, TkRule})")

				act = tmp
			}
		case ntk_Rule:
			if lookahead == nil || lookahead.Type != ttk_Newline {
				// [ Rule ] -> Source1 : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(ntk_Source1, []token_type{ntk_Rule}))
				// luc.AssertErr(err, "NewReduceAction(TkSource1, []TokenType{TkRule})")

				act = tmp
			} else {
				// Source1 newline [ Rule ] -> Source1 : shift .

				act = uprx.NewShiftAction()
			}
		case ttk_Newline:
			// Source1 [ newline ] Rule -> Source1 : shift .
			// RuleLine RhsCls equal [ newline ] uppercase_id -> Rule : shift .
			// RuleLine RhsCls pipe [ newline ] -> RuleLine : shift .
			// dot [ newline ] -> RuleLine : shift .

			act = uprx.NewShiftAction()
		case ttk_Dot:
			top2, ok := parser.Pop()
			if !ok {
				return nil, uprx.NewErrUnexpectedToken(&top1.Type, nil, ntk_RhsCls, ttk_Newline)
			}

			switch top2.Type {
			case ntk_RhsCls:
				// [ dot ] RhsCls equal uppercase_id -> Rule : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(ntk_Rule, []token_type{ttk_Dot, ntk_RhsCls, ttk_EqualSign, ttk_UppercaseID}))
				// luc.AssertErr(err, "NewReduceAction(TkRule, []TokenType{TkDot, TkRhsCls, TkEqualSign, TkUppercaseID})")

				act = tmp
			case ttk_Newline:
				// [ dot ] newline -> RuleLine : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(ntk_RuleLine, []token_type{ttk_Dot, ttk_Newline}))
				// luc.AssertErr(err, "NewReduceAction(TkRuleLine, []TokenType{TkDot, TkNewline})")

				act = tmp
			default:
				return nil, uprx.NewErrUnexpectedToken(&top1.Type, &top2.Type, ntk_RhsCls, ttk_Newline)
			}
		case ntk_RhsCls:
			top2, ok := parser.Pop()
			if !ok {
				return nil, uprx.NewErrUnexpectedToken(&top1.Type, nil, ntk_Rhs, ttk_EqualSign, ttk_Pipe)
			}

			if top2.Type == ntk_Rhs {
				// [ RhsCls ] Rhs -> RhsCls : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(ntk_RhsCls, []token_type{ntk_RhsCls, ntk_Rhs}))
				// luc.AssertErr(err, "NewReduceAction(TkRhsCls, []TokenType{TkRhsCls, TkRhs})")

				act = tmp
			} else {
				// dot [ RhsCls ] equal uppercase_id -> Rule : shift .

				// RuleLine [ RhsCls ] equal newline uppercase_id -> Rule : shift .
				// RuleLine [ RhsCls ] pipe newline -> RuleLine : shift .
				// -- RuleLine RhsCls pipe newline -> RuleLine .
				// -- dot newline -> RuleLine .

				act = uprx.NewShiftAction()
			}
		case ttk_EqualSign:
			// dot RhsCls [ equal ] uppercase_id -> Rule : shift .
			// RuleLine RhsCls [ equal ] newline uppercase_id -> Rule : shift .

			act = uprx.NewShiftAction()
		case ttk_UppercaseID:
			if lookahead == nil {
				// [ uppercase_id ] -> Identifier : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(ntk_Identifier, []token_type{ttk_UppercaseID}))
				// luc.AssertErr(err, "NewReduceAction(TkIdentifier, []TokenType{TkUppercaseID})")

				act = tmp
			} else {
				switch lookahead.Type {
				case ttk_EqualSign:
					// dot RhsCls equal [ uppercase_id ] -> Rule : shift .

					act = uprx.NewShiftAction()
				case ttk_Newline:
					// RuleLine RhsCls equal newline [ uppercase_id ] -> Rule : shift .

					act = uprx.NewShiftAction()
				default:
					// [ uppercase_id ] -> Identifier : reduce .

					tmp, _ := uprx.NewReduceAction(uprx.NewRule(ntk_Identifier, []token_type{ttk_UppercaseID}))
					// luc.AssertErr(err, "NewReduceAction(TkIdentifier, []TokenType{TkUppercaseID})")

					act = tmp
				}
			}
		case ntk_RuleLine:
			top2, ok := parser.Pop()
			if !ok {
				return nil, uprx.NewErrUnexpectedToken(&top1.Type, nil, ntk_RhsCls)
			} else if top2.Type != ntk_RhsCls {
				return nil, uprx.NewErrUnexpectedToken(&top1.Type, &top2.Type, ntk_RhsCls)
			}

			top3, ok := parser.Pop()
			if !ok {
				return nil, uprx.NewErrUnexpectedToken(&top2.Type, nil, ttk_EqualSign, ttk_Pipe)
			}

			switch top3.Type {
			case ttk_EqualSign:
				// [ RuleLine ] RhsCls equal newline uppercase_id -> Rule : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(ntk_Rule, []token_type{ntk_RuleLine, ntk_RhsCls, ttk_EqualSign, ttk_Newline, ttk_UppercaseID}))
				// luc.AssertErr(err, "NewReduceAction(TkRule, []TokenType{TkRuleLine, TkRhsCls, TkEqualSign, TkNewline, TkUppercaseID})")

				act = tmp
			case ttk_Pipe:
				// [ RuleLine ] RhsCls pipe newline -> RuleLine : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(ntk_RuleLine, []token_type{ntk_RuleLine, ntk_RhsCls, ttk_Pipe, ttk_Newline}))
				// luc.AssertErr(err, "NewReduceAction(TkRuleLine, []TokenType{TkRuleLine, TkRhsCls, TkPipe, TkTab, TkNewline})")

				act = tmp
			default:
				return nil, uprx.NewErrUnexpectedToken(&top2.Type, &top3.Type, ttk_EqualSign, ttk_Pipe)
			}
		case ttk_Pipe:
			// RuleLine RhsCls [ pipe ] newline -> RuleLine : shift .
			// Identifier [ pipe ] Identifier -> OrExpr : shift .
			// OrExpr [ pipe ] Identifier -> OrExpr : shift .

			act = uprx.NewShiftAction()
		case ntk_Rhs:
			if lookahead == nil || (lookahead.Type != ttk_UppercaseID && lookahead.Type != ttk_LowercaseID && lookahead.Type != ttk_OpParen) {
				// [ Rhs ] -> RhsCls : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(ntk_RhsCls, []token_type{ntk_Rhs}))
				// luc.AssertErr(err, "NewReduceAction(TkRhsCls, []TokenType{TkRhs})")

				act = tmp
			} else {
				// RhsCls [ Rhs ] -> RhsCls : shift .
				// -- Rhs -> RhsCls .
				// -- -- Identifier -> Rhs .
				// -- -- -- uppercase_id -> Identifier .
				// -- -- -- lowercase_id -> Identifier .
				// -- -- cl_paren OrExpr op_paren -> Rhs .
				// -- RhsCls Rhs -> RhsCls .
				// -- -- Identifier -> Rhs .
				// -- -- -- uppercase_id -> Identifier .
				// -- -- -- lowercase_id -> Identifier .
				// -- -- cl_paren OrExpr op_paren -> Rhs .

				act = uprx.NewShiftAction()
			}
		case ntk_Identifier:
			if lookahead == nil || lookahead.Type != ttk_Pipe {
				top2, ok := parser.Pop()
				if !ok || top2.Type != ttk_Pipe {
					// [ Identifier ] -> Rhs : reduce .

					tmp, _ := uprx.NewReduceAction(uprx.NewRule(ntk_Rhs, []token_type{ntk_Identifier}))
					// luc.AssertErr(err, "NewReduceAction(TkRhs, []TokenType{TkIdentifier})")

					act = tmp
				} else {
					// [ Identifier ] pipe Identifier -> OrExpr : reduce .

					tmp, _ := uprx.NewReduceAction(uprx.NewRule(ntk_OrExpr, []token_type{ntk_Identifier, ttk_Pipe, ntk_Identifier}))
					// luc.AssertErr(err, "NewReduceAction(TkOrExpr, []TokenType{TkIdentifier, TkPipe, TkIdentifier})")

					act = tmp
				}
			} else {
				// OrExpr pipe [ Identifier ] -> OrExpr : shift .

				act = uprx.NewShiftAction()
			}
		case ttk_ClParen:
			// [ cl_paren ] OrExpr op_paren -> Rhs : reduce .

			tmp, _ := uprx.NewReduceAction(uprx.NewRule(ntk_Rhs, []token_type{ttk_ClParen, ntk_OrExpr, ttk_OpParen}))
			// luc.AssertErr(err, "NewReduceAction(TkRhs, []TokenType{TkClParen, TkOrExpr, TkOpParen})")

			act = tmp
		case ntk_OrExpr:
			if lookahead == nil || lookahead.Type != ttk_ClParen {
				// [ OrExpr ] pipe Identifier -> OrExpr : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(ntk_OrExpr, []token_type{ntk_OrExpr, ttk_Pipe, ntk_Identifier}))
				// luc.AssertErr(err, "NewReduceAction(TkOrExpr, []TokenType{TkOrExpr, TkPipe, TkIdentifier})")

				act = tmp
			} else {
				// cl_paren [ OrExpr ] op_paren -> Rhs : shift .

				act = uprx.NewShiftAction()
			}
		case ttk_OpParen:
			// cl_paren OrExpr [ op_paren ] -> Rhs : shift .

			act = uprx.NewShiftAction()
		case ttk_LowercaseID:
			// [ lowercase_id ] -> Identifier : reduce .

			tmp, _ := uprx.NewReduceAction(uprx.NewRule(ntk_Identifier, []token_type{ttk_LowercaseID}))
			// luc.AssertErr(err, "NewReduceAction(TkIdentifier, []TokenType{TkLowercaseID})")

			act = tmp
		default:
			return nil, fmt.Errorf("unexpected token: %s", top1.String())
		}

		return act, nil
	}

	internal_parser = uprx.NewParser(f)
}
