package lexer

import (
	slgr "github.com/PlayerR9/SlParser/grammar"
	internal "github.com/PlayerR9/SlParser/lexer/internal"
	tr "github.com/PlayerR9/mygo-lib/CustomData/tree"
	"github.com/PlayerR9/mygo-lib/common"
)

/////////////////////////////////////////////////////////

// baseActive is a implementation of a lexer that uses a function for lexing.
type baseActive struct {
	// global is the global lexer.
	global *Lexer

	// tokens is the list of tokens.
	tokens []*tr.Node

	// pos is the current position in the stream.
	pos int

	// err is the last error.
	err error
}

// HasError implements the history.Subject interface.
func (ba baseActive) HasError() bool {
	return ba.err != nil
}

// GetError implements the history.Subject interface.
func (ba baseActive) GetError() error {
	return ba.err
}

// NextEvents implements the history.Subject interface.
func (ba *baseActive) NextEvents() ([]internal.Event, error) {
	if ba == nil {
		return nil, common.ErrNilReceiver
	}

	results, err := ba.global.Match(ba.pos)
	if err != nil {
		ba.err = NewErrLexing(ba.pos, err)

		return nil, nil
	}

	return results, nil
}

// ApplyEvent implements the history.Subject interface.
func (ba *baseActive) ApplyEvent(r internal.Event) error {
	if ba == nil {
		return common.ErrNilReceiver
	}

	if r.Type != slgr.EtToSkip {
		tk := slgr.NewToken(ba.pos, r.Type, string(r.Data), nil)
		ba.tokens = append(ba.tokens, tk)
	}

	ba.pos += len(r.Data)

	return nil
}

// Reset implements the Lexer interface.
func (ba *baseActive) Reset() {
	if ba == nil {
		return
	}

	if len(ba.tokens) > 0 {
		clear(ba.tokens)
		ba.tokens = nil
	}

	ba.pos = 0
	ba.err = nil
}

// Tokens returns the tokens generated from the input stream.
//
// Returns:
//   - []*tr.Node: The tokens generated from the input stream.
//
// Remember to append the EOF after calling this function
func (l baseActive) Tokens() []*tr.Node {
	tokens := make([]*tr.Node, len(l.tokens))
	copy(tokens, l.tokens)

	return tokens
}
