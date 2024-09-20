package parser

import (
	"iter"

	gr "github.com/PlayerR9/SlParser/grammar"
)

// Parser is a parser.
type Parser[T gr.TokenTyper] struct {
	// tokens is the list of tokens.
	tokens []*gr.Token[T]

	// table is the parser table.
	table map[T]ParseFn[T]

	// seq is the parser sequence.
	seq iter.Seq[*ActiveParser[T]]

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

	p.next, p.stop = iter.Pull(p.seq)
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
func (p *Parser[T]) Parse() (*ActiveParser[T], error) {
	if p == nil || p.next == nil {
		return nil, nil
	}

	ap, ok := p.next()
	if !ok {
		p.stop()

		return nil, nil
	}

	return ap, nil
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
