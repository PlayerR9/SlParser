package parser

import (
	"iter"

	gr "github.com/PlayerR9/SlParser/grammar"
	bck "github.com/PlayerR9/go-commons/Evaluations/history"
	"github.com/PlayerR9/go-errors/assert"
)

// Parser is a parser.
type Parser[T gr.TokenTyper] struct {
	// tokens is the list of tokens.
	tokens []*gr.Token[T]

	// table is the parser table.
	table map[T]ParseFn[T]

	// next is the next function.
	next func() (*ActiveParser[T], bool)

	// stop is the stop function.
	stop func()
}

// SetTokens sets the list of tokens.
//
// Parameters:
//   - tokens: the list of tokens.
//
// Does nothing if the receiver is nil.
func (p *Parser[T]) SetTokens(tokens []*gr.Token[T]) {
	if p == nil {
		return
	}

	p.tokens = tokens

	if p.next != nil {
		p.stop()

		p.next = nil
		p.stop = nil
	}
}

// Parse parses the list of tokens. Successive calls to this function
// yields different parsing results.
//
// The results are valid and then, as soon as this returns an error,
// the next results are all invalid.
//
// Returns:
//   - *ActiveParser[T]: the active parser. Nil the parser has done.
//   - error: if an error occurred.
func (p *Parser[T]) Parse() ([]*gr.ParseTree[T], error) {
	if p == nil {
		return nil, nil
	}

	if p.next == nil {
		fn := func() *ActiveParser[T] {
			ap, err := NewActiveParser(p)
			assert.Err(err, "NewActiveParser(p)")

			return ap
		}

		seq := bck.Subject(fn)
		p.next, p.stop = iter.Pull(seq)
	}

	ap, ok := p.next()
	if !ok {
		p.stop()

		return nil, nil
	}

	if ap == nil {
		return nil, nil
	}

	return ap.Forest(), nil
}

// Reset resets the parser; allowing it to be reused.
func (p *Parser[T]) Reset() {
	if p == nil {
		return
	}

	if len(p.tokens) > 0 {
		for i := 0; i < len(p.tokens); i++ {
			p.tokens[i] = nil
		}

		p.tokens = p.tokens[:0]
	}

	if p.next != nil {
		p.next = nil
	}

	if p.stop != nil {
		p.stop()

		p.stop = nil
	}
}

// ParseFnOf returns the parse function for the given symbol.
//
// Parameters:
//   - symbol: the symbol.
//
// Returns:
//   - ParseFn[T]: the parse function.
//   - bool: true if the symbol was found, false otherwise.
func (p Parser[T]) ParseFnOf(symbol T) (ParseFn[T], bool) {
	if p.table == nil {
		return nil, false
	}

	fn, ok := p.table[symbol]
	if !ok {
		return nil, false
	}

	return fn, true
}
