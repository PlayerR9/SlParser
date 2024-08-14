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
			var has_la bool

			if lookahead == nil || lookahead.Type != ttk_Pipe {
				has_la = false
			} else {
				has_la = true
			}

			top2, ok := p.Pop()
			if !ok || top2.Type != ttk_Pipe {
				if has_la {
					// ntk_OrExpr1 [ ntk_Identifier ] -> ntk_OrExpr : SHIFT .

					act = parsing.NewShiftAction()
				} else {
					// [ ntk_Identifier ] -> ntk_Rhs : REDUCE .

					act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Rhs, []token_type{ntk_Identifier}))
				}
			} else {
				if has_la {
					// ntk_OrExpr1 [ ntk_Identifier ] ttk_Pipe -> ntk_OrExpr1 : SHIFT .

					act = parsing.NewShiftAction()
				} else {
					// [ ntk_Identifier ] ttk_Pipe -> ntk_OrExpr1 : REDUCE .

					act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_OrExpr1, []token_type{ntk_Identifier, ttk_Pipe}))
				}
			}
		case ntk_OrExpr:
			// ttk_ClParen [ ntk_OrExpr ] ttk_OpParen -> ntk_Rhs : SHIFT .

			act = parsing.NewShiftAction()
		case ntk_OrExpr1:
			top2, ok := p.Pop()
			if !ok {
				return nil, parsing.NewErrUnexpectedToken(&top1.Type, nil, ntk_Identifier)
			} else if top2.Type != ntk_Identifier {
				return nil, parsing.NewErrUnexpectedToken(&top1.Type, &top2.Type, ntk_Identifier)
			}

			top3, ok := p.Pop()
			if !ok || top3.Type != ttk_Pipe {
				// [ ntk_OrExpr1 ] ntk_Identifier -> ntk_OrExpr : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_OrExpr, []token_type{ntk_OrExpr1, ntk_Identifier}))
			} else {
				// [ ntk_OrExpr1 ] ntk_Identifier ttk_Pipe -> ntk_OrExpr1 : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_OrExpr1, []token_type{ntk_OrExpr1, ntk_Identifier, ttk_Pipe}))
			}
		case ntk_Rhs:
			if lookahead == nil || (lookahead.Type != ttk_OpParen && lookahead.Type != ttk_LowercaseId && lookahead.Type != ttk_UppercaseId) {
				// [ ntk_Rhs ] -> ntk_Rhs1 : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Rhs1, []token_type{ntk_Rhs}))
			} else {
				// ntk_Rhs1 [ ntk_Rhs ] -> ntk_Rhs1 : SHIFT .

				act = parsing.NewShiftAction()
			}
		case ntk_Rhs1:
			top2, ok := p.Pop()
			if !ok {
				return nil, parsing.NewErrUnexpectedToken(&top1.Type, nil, ttk_Equal, ttk_Pipe, ntk_Rhs)
			}

			switch top2.Type {
			case ttk_Equal:
				// ttk_Semicolon [ ntk_Rhs1 ] ttk_Equal ttk_LowercaseId -> ntk_Rule : SHIFT .
				// ntk_Rule1 [ ntk_Rhs1 ] ttk_Equal ttk_LowercaseId -> ntk_Rule : SHIFT .

				act = parsing.NewShiftAction()
			case ttk_Pipe:
				if lookahead == nil || lookahead.Type != ttk_Pipe {
					// [ ntk_Rhs1 ] ttk_Pipe -> ntk_Rule1 : REDUCE .

					act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Rule1, []token_type{ntk_Rhs1, ttk_Pipe}))
				} else {
					// ntk_Rule1 [ ntk_Rhs1 ] ttk_Pipe -> ntk_Rule1 : SHIFT .

					act = parsing.NewShiftAction()
				}
			case ntk_Rhs:
				// [ ntk_Rhs1 ] ntk_Rhs -> ntk_Rhs1 : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Rhs1, []token_type{ntk_Rhs1, ntk_Rhs}))
			default:
				return nil, parsing.NewErrUnexpectedToken(&top1.Type, &top2.Type, ttk_Equal, ttk_Pipe, ntk_Rhs)
			}
		case ntk_Rule:
			if lookahead == nil || lookahead.Type != ttk_LowercaseId {
				// [ ntk_Rule ] -> ntk_Source1 : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Source1, []token_type{ntk_Rule}))
			} else {
				// ntk_Source1 [ ntk_Rule ] -> ntk_Source1 : SHIFT .

				act = parsing.NewShiftAction()
			}
		case ntk_Rule1:
			top2, ok := p.Pop()
			if !ok {
				return nil, parsing.NewErrUnexpectedToken(&top1.Type, nil, ntk_Rhs1)
			} else if top2.Type != ntk_Rhs1 {
				return nil, parsing.NewErrUnexpectedToken(&top1.Type, &top2.Type, ntk_Rhs1)
			}

			top3, ok := p.Pop()
			if !ok {
				return nil, parsing.NewErrUnexpectedToken(&top2.Type, nil, ttk_Pipe, ttk_Equal)
			}

			switch top3.Type {
			case ttk_Pipe:
				// [ ntk_Rule1 ] ntk_Rhs1 ttk_Pipe -> ntk_Rule1 : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Rule1, []token_type{ntk_Rule1, ntk_Rhs1, ttk_Pipe}))
			case ttk_Equal:
				// [ ntk_Rule1 ] ntk_Rhs1 ttk_Equal ttk_LowercaseId -> ntk_Rule : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Rule, []token_type{ntk_Rule1, ntk_Rhs1, ttk_Equal, ttk_LowercaseId}))
			default:
				return nil, parsing.NewErrUnexpectedToken(&top2.Type, &top3.Type, ttk_Pipe, ttk_Equal)
			}
		case ntk_Source1:
			top2, ok := p.Pop()
			if !ok || top2.Type != ntk_Rule {
				// etk_EOF [ ntk_Source1 ] -> ntk_Source : SHIFT .

				act = parsing.NewShiftAction()
			} else {
				// [ ntk_Source1 ] ntk_Rule -> ntk_Source1 : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Source1, []token_type{ntk_Source1, ntk_Rule}))
			}
		case ttk_ClParen:
			// [ ttk_ClParen ] ntk_OrExpr ttk_OpParen -> ntk_Rhs : REDUCE .

			act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Rhs, []token_type{ttk_ClParen, ntk_OrExpr, ttk_OpParen}))
		case ttk_Equal:
			// ttk_Semicolon ntk_Rhs1 [ ttk_Equal ] ttk_LowercaseId -> ntk_Rule : SHIFT .
			// ntk_Rule1 ntk_Rhs1 [ ttk_Equal ] ttk_LowercaseId -> ntk_Rule : SHIFT .

			act = parsing.NewShiftAction()
		case ttk_LowercaseId:
			if lookahead == nil || lookahead.Type != ttk_Equal {
				// [ ttk_LowercaseId ] -> ntk_Identifier : REDUCE .

				act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Identifier, []token_type{ttk_LowercaseId}))
			} else {
				// ttk_Semicolon ntk_Rhs1 ttk_Equal [ ttk_LowercaseId ] -> ntk_Rule : SHIFT .
				// ntk_Rule1 ntk_Rhs1 ttk_Equal [ ttk_LowercaseId ] -> ntk_Rule : SHIFT .

				act = parsing.NewShiftAction()
			}
		case ttk_OpParen:
			// ttk_ClParen ntk_OrExpr [ ttk_OpParen ] -> ntk_Rhs : SHIFT .

			act = parsing.NewShiftAction()
		case ttk_Pipe:
			// ntk_Identifier [ ttk_Pipe ] -> ntk_OrExpr1 : SHIFT .
			// ntk_OrExpr1 ntk_Identifier [ ttk_Pipe ] -> ntk_OrExpr1 : SHIFT .
			// ntk_Rhs1 [ ttk_Pipe ] -> ntk_Rule1 : SHIFT .
			// ntk_Rule1 ntk_Rhs1 [ ttk_Pipe ] -> ntk_Rule1 : SHIFT .

			act = parsing.NewShiftAction()
		case ttk_Semicolon:
			// [ ttk_Semicolon ] ntk_Rhs1 ttk_Equal ttk_LowercaseId -> ntk_Rule : REDUCE .

			act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Rule, []token_type{ttk_Semicolon, ntk_Rhs1, ttk_Equal, ttk_LowercaseId}))
		case ttk_UppercaseId:
			// [ ttk_UppercaseId ] -> ntk_Identifier : REDUCE .

			act, _ = parsing.NewReduceAction(parsing.NewRule(ntk_Identifier, []token_type{ttk_UppercaseId}))
		default:
			return nil, fmt.Errorf("unexpected token: %s", top1.String())
		}

		return act, nil
	}

	internal_parser = parsing.NewParser(decision_func)
}
