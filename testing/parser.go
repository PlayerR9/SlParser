package testing

import (
	"errors"
	"fmt"
	"slices"

	slgr "github.com/PlayerR9/SlParser/grammar"
	slpx "github.com/PlayerR9/SlParser/parser"
	"github.com/PlayerR9/go-evals/common"
)

func LinkTokens(tokens []*slgr.Token) []*slgr.Token {
	tokens = append(tokens, slgr.EOFToken)

	prev := tokens[0]

	for _, tk := range tokens[1:] {
		prev.Lookahead = tk
		prev = tk
	}

	return tokens
}

type ParserArgs struct {
	Tokens    []*slgr.Token
	Expecteds []string
}

func NewParserArgs(tokens []*slgr.Token, expecteds ...string) ParserArgs {
	tokens = LinkTokens(tokens)

	return ParserArgs{
		Tokens:    tokens,
		Expecteds: expecteds,
	}
}

func ParserTest(parser slpx.Parser, arg ParserArgs) error {
	if parser == nil {
		return common.NewErrNilParam("parser")
	}

	itr := parser.Parse(arg.Tokens)
	defer itr.Stop()

	var errs []error
	var solutions []*slpx.ParseTree

	for {
		pair, err := itr.Next()
		if err == slpx.ErrExhausted {
			break
		} else if err != nil {
			err := fmt.Errorf("failed to parse: %w", err)
			errs = append(errs, err)

			break
		}

		err = pair.GetError()
		if err != nil {
			err := fmt.Errorf("failed to parse: %w", err)
			errs = append(errs, err)

			continue
		}

		forest := pair.Forest()
		if len(forest) != 1 {
			err := fmt.Errorf("expected 1 parse tree, got %d", len(forest))
			errs = append(errs, err)

			continue
		}

		tree := forest[0]

		slice := tree.Slice()
		if len(slice) != len(arg.Expecteds) {
			err := fmt.Errorf("expected %d tokens, got %d", len(arg.Expecteds), len(slice))
			errs = append(errs, err)

			continue
		}

		for i, expected := range arg.Expecteds {
			type_ := slice[i].Type

			if type_ != expected {
				err := fmt.Errorf("expected token at index %d to be %s, got %s", i, expected, type_)
				errs = append(errs, err)
			}
		}

		ok := slices.ContainsFunc(solutions, tree.Equals)
		if !ok {
			solutions = append(solutions, tree)
		}
	}

	if len(solutions) > 1 {
		err := fmt.Errorf("expected 1 parse tree, got %d", len(solutions))
		errs = append(errs, err)
	}

	if len(solutions) == 1 {
		return nil
	}

	return errors.Join(errs...)
}
