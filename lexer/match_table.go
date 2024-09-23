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

		limit := 1

		for len(indices) > 1 {
			c, err := lexer.NextRune()
			if err == io.EOF {
				// TODO: Handle this case.
			} else if err != nil {
				// TODO: Handle this case.
			}

			var top int

			for i := 0; i < len(indices); i++ {
				idx := indices[i]

				word := mt.words[idx]

				if word[limit] == c {
					indices[top] = idx
					top++
				}
			}

			indices = indices[:top:top]

			limit++
		}

		return T(indices[0]), string(mt.words[indices[0]]), nil
	}

	return fn
}
