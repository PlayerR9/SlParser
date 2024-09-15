package lexer

import gr "github.com/PlayerR9/SlParser/grammar"

// LexFunc is a function that lexes a token.
//
// Parameters:
//   - lexer: the lexer. Assumed to be non-nil.
//   - char: the first character of the token.
//
// Returns:
//   - *grammar.Token[T]: the token.
//   - error: if an error occurred.
//
// If the returned token is nil, then the token is marked to be skipped.
type LexFunc[T gr.TokenTyper] func(lexer *Lexer[T], char rune) (*gr.Token[T], error)

// Builder is a lexer builder.
type Builder[T gr.TokenTyper] struct {
	// table is the lexer table.
	table map[rune]LexFunc[T]

	// def_fn is the default lexer function.
	def_fn LexFunc[T]
}

// NewBuilder creates a new lexer builder.
func NewBuilder[T gr.TokenTyper]() Builder[T] {
	return Builder[T]{
		table: make(map[rune]LexFunc[T]),
	}
}

// Register registers a new lexer function.
//
// Parameters:
//   - char: the first character of the token.
//   - fn: the lexer function.
//
// Behaviors:
//   - If the receiver or 'fn' are nil, then nothing is registered.
//   - If a 'char' is already registered, then the previous function is overwritten.
func (b *Builder[T]) Register(char rune, fn LexFunc[T]) {
	if b == nil || fn == nil {
		return
	}

	b.table[char] = fn
}

// Default sets the default lexer function.
//
// Parameters:
//   - fn: the default lexer function.
//
// Behaviors:
//   - If the receiver is nil, then nothing is set.
//   - If 'fn' is nil, then the previous function is removed.
//   - If a 'fn' is already registered, then the previous function is overwritten.
func (b *Builder[T]) Default(fn LexFunc[T]) {
	if b == nil {
		return
	}

	b.def_fn = fn
}

// Build builds the lexer.
//
// Returns:
//   - Lexer: the lexer.
func (b Builder[T]) Build() Lexer[T] {
	var table map[rune]LexFunc[T]

	if len(b.table) > 0 {
		table = make(map[rune]LexFunc[T], len(b.table))
		for char, fn := range b.table {
			table[char] = fn
		}
	}

	fn := b.def_fn

	return Lexer[T]{
		table:  table,
		def_fn: fn,
	}
}

// Reset resets the builder.
func (b *Builder[T]) Reset() {
	if b == nil {
		return
	}

	if len(b.table) > 0 {
		for key := range b.table {
			b.table[key] = nil
			delete(b.table, key)
		}

		b.table = make(map[rune]LexFunc[T])
	}

	b.def_fn = nil
}
