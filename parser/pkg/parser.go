package pkg

import (
	"fmt"

	gr "github.com/PlayerR9/grammar/grammar"
	uprx "github.com/PlayerR9/grammar/parser"
)

var (
	Parser *uprx.Parser[TokenType]
)

func init() {
	f := func(parser *uprx.Parser[TokenType], lookahead *gr.Token[TokenType]) (uprx.Actioner, error) {
		top1, ok := parser.Pop()
		if !ok {
			return nil, fmt.Errorf("p.stack is empty")
		}

		var act uprx.Actioner

		switch top1.Type {
		case TtkEOF:
			// [ EOF ] Source1 -> Source : accept .
			tmp, _ := uprx.NewAcceptAction(uprx.NewRule(NtkSource, []TokenType{TtkEOF, NtkSource1}))
			// luc.AssertErr(err, "NewAcceptAction(TkSource, []TokenType{TkEOF, TkSource1})")

			act = tmp
		case NtkSource1:
			top2, ok := parser.Pop()
			if !ok || top2.Type != TtkNewline {
				// EOF [ Source1 ] -> Source : shift .

				act = uprx.NewShiftAction()
			} else {
				// [ Source1 ] newline Rule -> Source1 : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(NtkSource1, []TokenType{NtkSource1, TtkNewline, NtkRule}))
				// luc.AssertErr(err, "NewReduceAction(TkSource1, []TokenType{TkSource1, TkNewline, TkRule})")

				act = tmp
			}
		case NtkRule:
			if lookahead == nil || lookahead.Type != TtkNewline {
				// [ Rule ] -> Source1 : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(NtkSource1, []TokenType{NtkRule}))
				// luc.AssertErr(err, "NewReduceAction(TkSource1, []TokenType{TkRule})")

				act = tmp
			} else {
				// Source1 newline [ Rule ] -> Source1 : shift .

				act = uprx.NewShiftAction()
			}
		case TtkNewline:
			// Source1 [ newline ] Rule -> Source1 : shift .
			// RuleLine RhsCls equal [ newline ] uppercase_id -> Rule : shift .
			// RuleLine RhsCls pipe [ newline ] -> RuleLine : shift .
			// dot [ newline ] -> RuleLine : shift .

			act = uprx.NewShiftAction()
		case TtkDot:
			top2, ok := parser.Pop()
			if !ok {
				return nil, uprx.NewErrUnexpectedToken(&top1.Type, nil, NtkRhsCls, TtkNewline)
			}

			switch top2.Type {
			case NtkRhsCls:
				// [ dot ] RhsCls equal uppercase_id -> Rule : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(NtkRule, []TokenType{TtkDot, NtkRhsCls, TtkEqualSign, TtkUppercaseID}))
				// luc.AssertErr(err, "NewReduceAction(TkRule, []TokenType{TkDot, TkRhsCls, TkEqualSign, TkUppercaseID})")

				act = tmp
			case TtkNewline:
				// [ dot ] newline -> RuleLine : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(NtkRuleLine, []TokenType{TtkDot, TtkNewline}))
				// luc.AssertErr(err, "NewReduceAction(TkRuleLine, []TokenType{TkDot, TkNewline})")

				act = tmp
			default:
				return nil, uprx.NewErrUnexpectedToken(&top1.Type, &top2.Type, NtkRhsCls, TtkNewline)
			}
		case NtkRhsCls:
			top2, ok := parser.Pop()
			if !ok {
				return nil, uprx.NewErrUnexpectedToken(&top1.Type, nil, NtkRhs, TtkEqualSign, TtkPipe)
			}

			if top2.Type == NtkRhs {
				// [ RhsCls ] Rhs -> RhsCls : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(NtkRhsCls, []TokenType{NtkRhsCls, NtkRhs}))
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
		case TtkEqualSign:
			// dot RhsCls [ equal ] uppercase_id -> Rule : shift .
			// RuleLine RhsCls [ equal ] newline uppercase_id -> Rule : shift .

			act = uprx.NewShiftAction()
		case TtkUppercaseID:
			if lookahead == nil {
				// [ uppercase_id ] -> Identifier : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(NtkIdentifier, []TokenType{TtkUppercaseID}))
				// luc.AssertErr(err, "NewReduceAction(TkIdentifier, []TokenType{TkUppercaseID})")

				act = tmp
			} else {
				switch lookahead.Type {
				case TtkEqualSign:
					// dot RhsCls equal [ uppercase_id ] -> Rule : shift .

					act = uprx.NewShiftAction()
				case TtkNewline:
					// RuleLine RhsCls equal newline [ uppercase_id ] -> Rule : shift .

					act = uprx.NewShiftAction()
				default:
					// [ uppercase_id ] -> Identifier : reduce .

					tmp, _ := uprx.NewReduceAction(uprx.NewRule(NtkIdentifier, []TokenType{TtkUppercaseID}))
					// luc.AssertErr(err, "NewReduceAction(TkIdentifier, []TokenType{TkUppercaseID})")

					act = tmp
				}
			}
		case NtkRuleLine:
			top2, ok := parser.Pop()
			if !ok {
				return nil, uprx.NewErrUnexpectedToken(&top1.Type, nil, NtkRhsCls)
			} else if top2.Type != NtkRhsCls {
				return nil, uprx.NewErrUnexpectedToken(&top1.Type, &top2.Type, NtkRhsCls)
			}

			top3, ok := parser.Pop()
			if !ok {
				return nil, uprx.NewErrUnexpectedToken(&top2.Type, nil, TtkEqualSign, TtkPipe)
			}

			switch top3.Type {
			case TtkEqualSign:
				// [ RuleLine ] RhsCls equal newline uppercase_id -> Rule : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(NtkRule, []TokenType{NtkRuleLine, NtkRhsCls, TtkEqualSign, TtkNewline, TtkUppercaseID}))
				// luc.AssertErr(err, "NewReduceAction(TkRule, []TokenType{TkRuleLine, TkRhsCls, TkEqualSign, TkNewline, TkUppercaseID})")

				act = tmp
			case TtkPipe:
				// [ RuleLine ] RhsCls pipe newline -> RuleLine : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(NtkRuleLine, []TokenType{NtkRuleLine, NtkRhsCls, TtkPipe, TtkNewline}))
				// luc.AssertErr(err, "NewReduceAction(TkRuleLine, []TokenType{TkRuleLine, TkRhsCls, TkPipe, TkTab, TkNewline})")

				act = tmp
			default:
				return nil, uprx.NewErrUnexpectedToken(&top2.Type, &top3.Type, TtkEqualSign, TtkPipe)
			}
		case TtkPipe:
			// RuleLine RhsCls [ pipe ] newline -> RuleLine : shift .
			// Identifier [ pipe ] Identifier -> OrExpr : shift .
			// OrExpr [ pipe ] Identifier -> OrExpr : shift .

			act = uprx.NewShiftAction()
		case NtkRhs:
			if lookahead == nil || (lookahead.Type != TtkUppercaseID && lookahead.Type != TtkLowercaseID && lookahead.Type != TtkOpParen) {
				// [ Rhs ] -> RhsCls : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(NtkRhsCls, []TokenType{NtkRhs}))
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
		case NtkIdentifier:
			if lookahead == nil || lookahead.Type != TtkPipe {
				top2, ok := parser.Pop()
				if !ok || top2.Type != TtkPipe {
					// [ Identifier ] -> Rhs : reduce .

					tmp, _ := uprx.NewReduceAction(uprx.NewRule(NtkRhs, []TokenType{NtkIdentifier}))
					// luc.AssertErr(err, "NewReduceAction(TkRhs, []TokenType{TkIdentifier})")

					act = tmp
				} else {
					// [ Identifier ] pipe Identifier -> OrExpr : reduce .

					tmp, _ := uprx.NewReduceAction(uprx.NewRule(NtkOrExpr, []TokenType{NtkIdentifier, TtkPipe, NtkIdentifier}))
					// luc.AssertErr(err, "NewReduceAction(TkOrExpr, []TokenType{TkIdentifier, TkPipe, TkIdentifier})")

					act = tmp
				}
			} else {
				// OrExpr pipe [ Identifier ] -> OrExpr : shift .

				act = uprx.NewShiftAction()
			}
		case TtkClParen:
			// [ cl_paren ] OrExpr op_paren -> Rhs : reduce .

			tmp, _ := uprx.NewReduceAction(uprx.NewRule(NtkRhs, []TokenType{TtkClParen, NtkOrExpr, TtkOpParen}))
			// luc.AssertErr(err, "NewReduceAction(TkRhs, []TokenType{TkClParen, TkOrExpr, TkOpParen})")

			act = tmp
		case NtkOrExpr:
			if lookahead == nil || lookahead.Type != TtkClParen {
				// [ OrExpr ] pipe Identifier -> OrExpr : reduce .

				tmp, _ := uprx.NewReduceAction(uprx.NewRule(NtkOrExpr, []TokenType{NtkOrExpr, TtkPipe, NtkIdentifier}))
				// luc.AssertErr(err, "NewReduceAction(TkOrExpr, []TokenType{TkOrExpr, TkPipe, TkIdentifier})")

				act = tmp
			} else {
				// cl_paren [ OrExpr ] op_paren -> Rhs : shift .

				act = uprx.NewShiftAction()
			}
		case TtkOpParen:
			// cl_paren OrExpr [ op_paren ] -> Rhs : shift .

			act = uprx.NewShiftAction()
		case TtkLowercaseID:
			// [ lowercase_id ] -> Identifier : reduce .

			tmp, _ := uprx.NewReduceAction(uprx.NewRule(NtkIdentifier, []TokenType{TtkLowercaseID}))
			// luc.AssertErr(err, "NewReduceAction(TkIdentifier, []TokenType{TkLowercaseID})")

			act = tmp
		default:
			return nil, fmt.Errorf("unexpected token: %s", top1.String())
		}

		return act, nil
	}

	Parser = uprx.NewParser[TokenType](f)
}
