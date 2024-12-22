package parser

import (
	"errors"

	slgr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/mygo-lib/common"
	assert "github.com/PlayerR9/go-verify"
	lls "github.com/PlayerR9/mygo-data/stack"
)

// ParseOneFn is the function used to parse one token from the input data.
//
// Parameters:
//   - parser: The input data to be parsed.
//
// Returns:
//   - Action: The parsed action, or nil if the parsing process fails.
//   - error: An error if the parsing process fails.
type ParseOneFn func(parser *Parser) (Action, error)

// Builder is a builder for parsers.
type Builder struct {
	// parse_one_fn is the function used to parse the input tokens.
	parse_one_fn ParseOneFn
}

// Reset implements common.Resetter.
func (b *Builder) Reset() error {
	if b == nil {
		return common.ErrNilReceiver
	}

	b.parse_one_fn = nil

	return nil
}

// SetParseOneFn sets the parsing function used by the parser.
//
// Parameters:
//   - fn: The new parsing function. Must not be nil.
//
// Returns:
//   - error: An error if the receiver is nil or if the parameter is nil.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
//   - common.ErrBadParam: If the parameter is nil.
func (b *Builder) SetParseOneFn(fn ParseOneFn) error {
	if b == nil {
		return common.ErrNilReceiver
	}

	if fn == nil {
		err := common.NewErrNilParam("fn")
		return err
	}

	b.parse_one_fn = fn

	return nil
}

// Build creates a new parser using the values set on the builder.
//
// Returns:
//   - *Parser: The newly created parser. Never returns nil.
func (b Builder) Build() *Parser {
	var fn ParseOneFn

	if b.parse_one_fn == nil {
		fn = func(_ *Parser) (Action, error) {
			err := errors.New("no parsing function provided")
			return nil, err
		}
	} else {
		fn = b.parse_one_fn
	}

	stack, err := lls.RefusableOf(new(lls.ArrayStack[*slgr.Token]))
	assert.Err(err, "lls.RefusableOf(new(lls.ArrayStack[*slgr.Token]))")

	parser := &Parser{
		parse_one_fn: fn,
		stack:        stack,
	}

	return parser
}
