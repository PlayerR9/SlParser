package lexer

import (
	gcers "github.com/PlayerR9/go-errors"
)

// FragWithOptions lexes a fragment with options.
//
// Parameters:
//   - frag_fn: the fragment function.
//   - options: the lexer options.
//
// Returns:
//   - LexFragment: the lexer fragment.
//
// By default, the lexer does allow optional fragments and only lexes once.
//   - Use WithLexMany(true) to enable one or more fragments.
//
// If 'frag_fn' is nil, then a function that returns an error is returned.
func ApplyMany(stream RuneStreamer, frag LexFragment) error {
	if frag == nil {
		return gcers.NewErrNilParameter("lexer.ApplyMany()", "frag")
	}

	if stream == nil {
		return NotFound
	}

	err := frag(stream)
	if err != nil {
		return err
	}

	for {
		err := frag(stream)
		if err == NotFound {
			break
		} else if err != nil {
			return err
		}
	}

	return nil
}
