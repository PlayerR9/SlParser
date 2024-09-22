package lexer

import (
	"fmt"
	"io"

	gr "github.com/PlayerR9/SlParser/grammar"
)

type MatchTable[T gr.TokenTyper] struct {
	words [][]rune
}

func (mt MatchTable[T]) MakeFunc() LexFunc[T] {
	fn := func(lexer RuneStreamer, char rune) (T, string, error) {
		var indices []int

		for i, word := range mt.words {
			c := word[0]

			if c == char {
				indices = append(indices, i)
			}
		}

		if len(indices) == 0 {
			return T(-1), "", fmt.Errorf("no words start with %q", char)
		}

		for len(indices) > 1 {
			c, err := lexer.NextRune()
			if err == io.EOF {

			}
		}

		return T(indices[0]), string(mt.words[indices[0]]), nil
	}

	return fn
}
