package parser

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	slgr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/mygo-lib/common"
	"github.com/PlayerR9/SlParser/parser/internal"
	lls "github.com/PlayerR9/mygo-data/stack"

	assert "github.com/PlayerR9/go-verify"
)

// Parser is a parser that can be used to parse input data into a list of tokens.
type Parser struct {
	// tokens is the list of tokens to be parsed.
	tokens []*slgr.Token

	// parse_one_fn is the function used to parse one token from the list of tokens.
	parse_one_fn ParseOneFn

	// stack is the stack of tokens that are currently being parsed.
	stack *lls.RefusableStack[*slgr.Token]
}

// Reset implements common.Resetter.
func (p *Parser) Reset() error {
	if p == nil {
		return common.ErrNilReceiver
	}

	if len(p.tokens) > 0 {
		clear(p.tokens)
		p.tokens = nil
	}

	err := p.stack.Reset()
	if err != nil {
		err := fmt.Errorf("while resetting stack: %w", err)
		return err
	}

	return nil
}

// SetInputStream sets the input stream of tokens to be parsed.
//
// Parameters:
//   - tokens: The list of tokens to be used as the input stream.
//
// Returns:
//   - error: An error if the receiver is nil.
func (p *Parser) SetInputStream(tokens []*slgr.Token) error {
	if p == nil {
		return common.ErrNilReceiver
	}

	p.tokens = tokens

	return nil
}

// Push pushes a token onto the stack.
//
// Parameters:
//   - tk: The token to be pushed onto the stack.
//
// Returns:
//   - error: An error if the receiver is nil.
func (p *Parser) Push(tk *slgr.Token) error {
	if p == nil {
		return common.ErrNilReceiver
	}

	if tk == nil {
		err := common.NewErrNilParam("tk")
		return err
	}

	err := p.stack.Push(tk)
	return err
}

// Pop pops the top token off the stack.
//
// Returns:
//   - *slgr.Token: The popped token, or nil if the stack is empty.
//   - error: An error if the receiver is nil or if the stack is empty.
func (p *Parser) Pop() (*slgr.Token, error) {
	if p == nil {
		return nil, common.ErrNilReceiver
	}

	tk, err := p.stack.Pop()
	if err != nil {
		return nil, err
	}

	return tk, nil
}

// shift shifts the next token from the input stream onto the stack.
//
// The function attempts to shift the next token from the input stream onto the stack.
// If the input stream is empty, the function returns an error.
//
// Returns:
//   - error: An error if the receiver is nil or if the input stream is empty.
func (p *Parser) shift() error {
	assert.Cond(p != nil, "p != nil")

	if len(p.tokens) == 0 {
		return errors.New("no tokens to shift")
	}

	tk := p.tokens[0]
	p.tokens = p.tokens[1:]

	assert.Cond(tk != nil, "tk != nil")

	err := p.Push(tk)
	assert.Err(err, "p.Push(tk)")

	return nil
}

// reduce reduces the top of the stack into a single token.
//
// The function attempts to reduce the top of the stack by popping tokens
// according to the provided right-hand side (rhss) symbols and combines them
// into a new token with the specified left-hand side (lhs) symbol. If the stack
// does not match the expected right-hand side symbols, or if the stack is empty,
// the function returns an error.
//
// Parameters:
//   - lhs: The left-hand side symbol for the new token.
//   - rhss: The right-hand side symbols to match against the stack.
//
// Returns:
//   - error: An error if the receiver is nil, if the stack is empty, or if the
//     symbols do not match.
func (p *Parser) reduce(rule *internal.Rule) error {
	assert.Cond(p != nil, "p != nil")
	assert.Cond(rule != nil, "rule != nil")

	rhss := rule.Rhss()
	slices.Reverse(rhss)

	for _, rhs := range rhss {
		tk, err := p.Pop()
		if err != nil {
			if err == lls.ErrEmptyStack {
				err = fmt.Errorf("want %s, got nothing", strconv.Quote(rhs))
			}

			return err
		}

		if tk.Type != rhs {
			err := fmt.Errorf("want %s, got %s", strconv.Quote(rhs), strconv.Quote(tk.Type))
			return err
		}
	}

	children := p.stack.Popped()

	err := p.stack.Accept()
	assert.Err(err, "p.stack.Accept()")

	// slices.Reverse(children)

	lhs := rule.Lhs()

	tk := slgr.NewToken(lhs, "")

	err = tk.AppendChildren(children)
	assert.Err(err, "tk.AppendChildren(children)")

	err = p.Push(tk)
	assert.Err(err, "p.Push(tk)")

	return nil
}

// Parse parses the input stream of tokens into a single token.
//
// Returns:
//   - *slgr.Token: The parsed token, or nil if the parsing process fails.
//   - error: An error if the receiver is nil or if the parsing process fails.
func (p *Parser) Parse() error {
	if p == nil {
		return common.ErrNilReceiver
	}

	err := p.shift() // Initial shift.
	if err != nil {
		return errors.New("initial shift failed")
	}

	is_done := false

	for !is_done {
		act, err := p.parse_one_fn(p)
		assert.Err(p.stack.Refuse(), "p.stack.Refuse()")

		if err != nil {
			return err
		}

		switch act := act.(type) {
		case *ShiftAction:
			err := p.shift()
			if err != nil {
				return fmt.Errorf("while shifting: %w", err)
			}
		case *ReduceAction:
			err := p.reduce(act.rule)
			if err != nil {
				return fmt.Errorf("while reducing: %w", err)
			}
		case *AcceptAction:
			err := p.reduce(act.rule)
			if err != nil {
				return fmt.Errorf("while reducing: %w", err)
			}

			is_done = true
		default:
			return fmt.Errorf("unknown action type: %T", act)
		}
	}

	return nil
}

// GetForest returns a slice of all the tokens in the parse forest.
//
// The function returns a slice of all the tokens in the parse forest. The
// tokens are ordered by their position in the input stream. The function does
// not modify the internal state of the receiver.
//
// Returns:
//   - []*slgr.Token: A slice of tokens in the parse forest.
func (p Parser) GetForest() []*slgr.Token {
	forest := p.stack.Slice()
	return forest
}

// Parse parses the input stream of tokens using the provided parser.
//
// Parameters:
//   - parser: The parser to be used to parse the input stream.
//   - tokens: The list of tokens to be used as the input stream.
//
// Returns:
//   - []*slgr.Token: The parsed token, or nil if the parsing process fails.
//   - error: An error if the receiver is nil or if the parsing process fails.
func Parse(parser *Parser, tokens []*slgr.Token) ([]*slgr.Token, error) {
	if parser == nil {
		err := common.NewErrNilParam("parser")
		return nil, err
	}

	defer parser.Reset()

	eof_tk := slgr.NewToken(EtEOF, "")
	tokens = append(tokens, eof_tk)

	err := parser.SetInputStream(tokens)
	assert.Err(err, "parser.SetInputStream(tokens)")

	err = parser.Parse()
	if err != nil {
		return nil, err
	}

	forest := parser.GetForest()
	return forest, nil
}
