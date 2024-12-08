package lexer

import (
	emtch "github.com/PlayerR9/go-evals/matcher"
	"github.com/PlayerR9/mygo-lib/common"
)

/////////////////////////////////////////////////////////

// Builder is a builder for a Lexer. An empty builder can be created with
// the `var b Builder` syntax or with the `b := new(Builder)` constructor.
type Builder struct {
	// matchers is a list of matchers used by the lexer.
	matchers []emtch.Matcher[rune]

	// types is a mirror of the types associated with the matchers.
	types []string

	// size is the number of matchers in the builder.
	size int
}

// RegisterLiteral is a convenience method for registering a literal. It is the same
// as calling `Register(common.Must(Literal(literal)), type_)`.
//
// Parameters:
//   - literal: The literal to register.
//   - type_: The type to associate with the literal.
//
// Returns:
//   - error: An error if the receiver is nil.
//
// Does nothing is the literal is empty.
func (b *Builder) RegisterLiteral(literal, type_ string) error {
	if literal == "" {
		return nil
	} else if b == nil {
		return common.ErrNilReceiver
	}

	w := Literal(literal)

	b.matchers = append(b.matchers, w)
	b.types = append(b.types, type_)
	b.size++

	return nil
}

// Register associates the given matcher with the given type.
//
// Parameters:
//   - m: The matcher to associate with the given type.
//   - type_: The type to associate with the given matcher.
//
// Returns:
//   - error: An error if the receiver is nil.
func (b *Builder) Register(m emtch.Matcher[rune], type_ string) error {
	if m == nil {
		return nil
	} else if b == nil {
		return common.ErrNilReceiver
	}

	b.matchers = append(b.matchers, m)
	b.types = append(b.types, type_)
	b.size++

	return nil
}

// Build returns a Lexer with the given matchers and types. The returned Lexer's
// internal state is a copy of the Builder's internal state, so modifications to
// the Builder after calling Build will not affect the returned Lexer.
//
// Returns:
//   - Lexer: A Lexer with the given matchers and types. Never returns nil.
func (b Builder) Build() *Lexer {
	lexer := &Lexer{
		table:   make([]emtch.Matcher[rune], b.size),
		types_:  make([]string, b.size),
		indices: make([]int, 0, b.size),
	}

	copy(lexer.table, b.matchers)
	copy(lexer.types_, b.types)

	for idx := range b.matchers {
		lexer.indices = append(lexer.indices, idx)
	}

	return lexer
}

// Reset resets the builder's internal state for reuse. No-op if the receiver is nil.
func (b *Builder) Reset() {
	if b == nil {
		return
	}

	if len(b.matchers) > 0 {
		clear(b.matchers)
		b.matchers = nil
	}

	if len(b.types) > 0 {
		clear(b.types)
		b.types = nil
	}

	b.size = 0
}
