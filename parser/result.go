package parser

import (
	"iter"

	"github.com/PlayerR9/SlParser/parser/internal"
	hst "github.com/PlayerR9/go-evals/history"
)

/////////////////////////////////////////////////////////

// Result is the result of an evaluation.
type Result struct {
	// pair is the result of the evaluation.
	pair hst.Pair[*internal.Event, *Active]
}

// Emulate emulates the parse process. It returns a sequence of all the active trees
// that were generated during the parse process.
//
// Returns:
//   - iter.Seq[*Active]: The sequence of active trees. Never returns nil.
func (r Result) Emulate() iter.Seq[*Active] {
	return hst.Emulate(r.pair)
}

// Forest returns the parse forest of the parser. The parse forest is a slice of
// parse trees, where each parse tree is a tree of tokens. The parse forest is
// constructed by popping all the tokens from the parse stack and constructing a
// parse tree for each one. The parse forest is useful for debugging and for
// visualizing the parse tree.
//
// Returns:
//   - []*ParseTree: A slice of parse trees. The slice is never empty, since the parser always has a parse stack.
func (r Result) Forest() []*ParseTree {
	return r.pair.Subject.Forest()
}

// GetError returns the error that occurred during the parse process.
//
// Returns:
//   - error: An error if the parse process failed. Otherwise, nil.
func (r Result) GetError() error {
	return r.pair.GetError()
}
