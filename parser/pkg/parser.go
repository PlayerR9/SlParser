package pkg

import (
	"fmt"
	"slices"

	gr "github.com/PlayerR9/grammar/grammar"
	uprx "github.com/PlayerR9/grammar/parser"
	luc "github.com/PlayerR9/lib_units/common"
	llq "github.com/PlayerR9/listlike/queue"
	lls "github.com/PlayerR9/listlike/stack"
)

type Parser struct {
	tokens *llq.ArrayQueue[*gr.Token[TokenType]]
	stack  *lls.ArrayStack[*gr.Token[TokenType]]
	popped *lls.ArrayStack[*gr.Token[TokenType]]
}

// SetInputStream implements the Grammar.Parser interface.
func (p *Parser) SetInputStream(tokens []*gr.Token[TokenType]) {
	luc.AssertNil(p.stack, "p.stack")
	luc.AssertNil(p.popped, "p.popped")

	p.tokens = llq.NewArrayQueue[*gr.Token[TokenType]]()
	p.tokens.EnqueueMany(tokens)

	p.stack.Clear()
	p.popped.Clear()
}

// Pop implements the Grammar.Parser interface.
func (p *Parser) Pop() (*gr.Token[TokenType], bool) {
	luc.AssertNil(p.stack, "p.stack")

	top, ok := p.stack.Pop()
	if !ok {
		return nil, false
	}

	luc.AssertNil(p.popped, "p.popped")

	p.popped.Push(top)

	return top, true
}

// Peek implements the Grammar.Parser interface.
func (p *Parser) Peek() (*gr.Token[TokenType], bool) {
	luc.AssertNil(p.stack, "p.stack")

	top, ok := p.stack.Peek()
	if !ok {
		return nil, false
	}

	return top, true
}

// GetDecision implements the Grammar.Parser interface.
func (p *Parser) GetDecision(lookahead *gr.Token[TokenType]) (uprx.Actioner, error) {
	defer p.Refuse()

	top1, ok := p.Pop()
	if !ok {
		return nil, fmt.Errorf("p.stack is empty")
	}

	var act uprx.Actioner

	switch top1.Type {
	case TtkEOF:
		// [ EOF ] Source1 -> Source : accept .
		tmp, err := uprx.NewAcceptAction(uprx.NewRule(NtkSource, []TokenType{TtkEOF, NtkSource1}))
		luc.AssertErr(err, "NewAcceptAction(TkSource, []TokenType{TkEOF, TkSource1})")

		act = tmp
	case NtkSource1:
		top2, ok := p.Pop()
		if !ok || top2.Type != TtkNewline {
			// EOF [ Source1 ] -> Source : shift .

			act = uprx.NewShiftAction()
		} else {
			// [ Source1 ] newline Rule -> Source1 : reduce .

			tmp, err := uprx.NewReduceAction(uprx.NewRule(NtkSource1, []TokenType{NtkSource1, TtkNewline, NtkRule}))
			luc.AssertErr(err, "NewReduceAction(TkSource1, []TokenType{TkSource1, TkNewline, TkRule})")

			act = tmp
		}
	case NtkRule:
		if lookahead == nil || lookahead.Type != TtkNewline {
			// [ Rule ] -> Source1 : reduce .

			tmp, err := uprx.NewReduceAction(uprx.NewRule(NtkSource1, []TokenType{NtkRule}))
			luc.AssertErr(err, "NewReduceAction(TkSource1, []TokenType{TkRule})")

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
		top2, ok := p.Pop()
		if !ok {
			return nil, uprx.NewErrUnexpectedToken(&top1.Type, nil, NtkRhsCls, TtkNewline)
		}

		switch top2.Type {
		case NtkRhsCls:
			// [ dot ] RhsCls equal uppercase_id -> Rule : reduce .

			tmp, err := uprx.NewReduceAction(uprx.NewRule(NtkRule, []TokenType{TtkDot, NtkRhsCls, TtkEqualSign, TtkUppercaseID}))
			luc.AssertErr(err, "NewReduceAction(TkRule, []TokenType{TkDot, TkRhsCls, TkEqualSign, TkUppercaseID})")

			act = tmp
		case TtkNewline:
			// [ dot ] newline -> RuleLine : reduce .

			tmp, err := uprx.NewReduceAction(uprx.NewRule(NtkRuleLine, []TokenType{TtkDot, TtkNewline}))
			luc.AssertErr(err, "NewReduceAction(TkRuleLine, []TokenType{TkDot, TkNewline})")

			act = tmp
		default:
			return nil, uprx.NewErrUnexpectedToken(&top1.Type, &top2.Type, NtkRhsCls, TtkNewline)
		}
	case NtkRhsCls:
		top2, ok := p.Pop()
		if !ok {
			return nil, uprx.NewErrUnexpectedToken(&top1.Type, nil, NtkRhs, TtkEqualSign, TtkPipe)
		}

		if top2.Type == NtkRhs {
			// [ RhsCls ] Rhs -> RhsCls : reduce .

			tmp, err := uprx.NewReduceAction(uprx.NewRule(NtkRhsCls, []TokenType{NtkRhsCls, NtkRhs}))
			luc.AssertErr(err, "NewReduceAction(TkRhsCls, []TokenType{TkRhsCls, TkRhs})")

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

			tmp, err := uprx.NewReduceAction(uprx.NewRule(NtkIdentifier, []TokenType{TtkUppercaseID}))
			luc.AssertErr(err, "NewReduceAction(TkIdentifier, []TokenType{TkUppercaseID})")

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

				tmp, err := uprx.NewReduceAction(uprx.NewRule(NtkIdentifier, []TokenType{TtkUppercaseID}))
				luc.AssertErr(err, "NewReduceAction(TkIdentifier, []TokenType{TkUppercaseID})")

				act = tmp
			}
		}
	case NtkRuleLine:
		top2, ok := p.Pop()
		if !ok {
			return nil, uprx.NewErrUnexpectedToken(&top1.Type, nil, NtkRhsCls)
		} else if top2.Type != NtkRhsCls {
			return nil, uprx.NewErrUnexpectedToken(&top1.Type, &top2.Type, NtkRhsCls)
		}

		top3, ok := p.Pop()
		if !ok {
			return nil, uprx.NewErrUnexpectedToken(&top2.Type, nil, TtkEqualSign, TtkPipe)
		}

		switch top3.Type {
		case TtkEqualSign:
			// [ RuleLine ] RhsCls equal newline uppercase_id -> Rule : reduce .

			tmp, err := uprx.NewReduceAction(uprx.NewRule(NtkRule, []TokenType{NtkRuleLine, NtkRhsCls, TtkEqualSign, TtkNewline, TtkUppercaseID}))
			luc.AssertErr(err, "NewReduceAction(TkRule, []TokenType{TkRuleLine, TkRhsCls, TkEqualSign, TkNewline, TkUppercaseID})")

			act = tmp
		case TtkPipe:
			// [ RuleLine ] RhsCls pipe newline -> RuleLine : reduce .

			tmp, err := uprx.NewReduceAction(uprx.NewRule(NtkRuleLine, []TokenType{NtkRuleLine, NtkRhsCls, TtkPipe, TtkNewline}))
			luc.AssertErr(err, "NewReduceAction(TkRuleLine, []TokenType{TkRuleLine, TkRhsCls, TkPipe, TkTab, TkNewline})")

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

			tmp, err := uprx.NewReduceAction(uprx.NewRule(NtkRhsCls, []TokenType{NtkRhs}))
			luc.AssertErr(err, "NewReduceAction(TkRhsCls, []TokenType{TkRhs})")

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
			top2, ok := p.Pop()
			if !ok || top2.Type != TtkPipe {
				// [ Identifier ] -> Rhs : reduce .

				tmp, err := uprx.NewReduceAction(uprx.NewRule(NtkRhs, []TokenType{NtkIdentifier}))
				luc.AssertErr(err, "NewReduceAction(TkRhs, []TokenType{TkIdentifier})")

				act = tmp
			} else {
				// [ Identifier ] pipe Identifier -> OrExpr : reduce .

				tmp, err := uprx.NewReduceAction(uprx.NewRule(NtkOrExpr, []TokenType{NtkIdentifier, TtkPipe, NtkIdentifier}))
				luc.AssertErr(err, "NewReduceAction(TkOrExpr, []TokenType{TkIdentifier, TkPipe, TkIdentifier})")

				act = tmp
			}
		} else {
			// OrExpr pipe [ Identifier ] -> OrExpr : shift .

			act = uprx.NewShiftAction()
		}
	case TtkClParen:
		// [ cl_paren ] OrExpr op_paren -> Rhs : reduce .

		tmp, err := uprx.NewReduceAction(uprx.NewRule(NtkRhs, []TokenType{TtkClParen, NtkOrExpr, TtkOpParen}))
		luc.AssertErr(err, "NewReduceAction(TkRhs, []TokenType{TkClParen, TkOrExpr, TkOpParen})")

		act = tmp
	case NtkOrExpr:
		if lookahead == nil || lookahead.Type != TtkClParen {
			// [ OrExpr ] pipe Identifier -> OrExpr : reduce .

			tmp, err := uprx.NewReduceAction(uprx.NewRule(NtkOrExpr, []TokenType{NtkOrExpr, TtkPipe, NtkIdentifier}))
			luc.AssertErr(err, "NewReduceAction(TkOrExpr, []TokenType{TkOrExpr, TkPipe, TkIdentifier})")

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

		tmp, err := uprx.NewReduceAction(uprx.NewRule(NtkIdentifier, []TokenType{TtkLowercaseID}))
		luc.AssertErr(err, "NewReduceAction(TkIdentifier, []TokenType{TkLowercaseID})")

		act = tmp
	default:
		return nil, fmt.Errorf("unexpected token: %s", top1.String())
	}

	return act, nil
}

// Shift implements the Grammar.Parser interface.
func (p *Parser) Shift() bool {
	luc.AssertNil(p.tokens, "p.tokens")

	first, ok := p.tokens.Dequeue()
	if !ok {
		return false
	}

	luc.AssertNil(p.stack, "p.stack")

	p.stack.Push(first)

	return true
}

// GetPopped implements the Grammar.Parser interface.
func (p *Parser) GetPopped() []*gr.Token[TokenType] {
	popped := p.popped.Slice()

	slices.Reverse(popped)

	return popped
}

// Push implements the Grammar.Parser interface.
func (p *Parser) Push(token *gr.Token[TokenType]) {
	if token == nil {
		return
	}

	luc.AssertNil(p.stack, "p.stack")

	p.stack.Push(token)
}

// Refuse implements the Grammar.Parser interface.
func (p *Parser) Refuse() {
	for {
		top, ok := p.popped.Pop()
		if !ok {
			break
		}

		p.stack.Push(top)
	}
}

// Accept implements the Grammar.Parser interface.
func (p *Parser) Accept() {
	p.popped.Clear()
}

func NewParser() *Parser {
	return &Parser{
		stack:  lls.NewArrayStack[*gr.Token[TokenType]](),
		popped: lls.NewArrayStack[*gr.Token[TokenType]](),
	}
}

// FullParse is just a wrapper around the Grammar.FullParse function.
//
// Parameters:
//   - tokens: The input stream of the parser.
//
// Returns:
//   - []*gr.TokenTree[TokenType]: The syntax forest of the input stream.
//   - error: An error if the parser encounters an error while parsing the input stream.
func FullParse(tokens []*gr.Token[TokenType]) ([]*gr.TokenTree[TokenType], error) {
	parser := NewParser()

	forest, err := uprx.FullParse(parser, tokens)
	if err != nil {
		return forest, err
	}

	return forest, nil
}
