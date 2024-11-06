package testing

import (
	"errors"
	"fmt"

	"github.com/PlayerR9/SlParser/grammar"
	slgr "github.com/PlayerR9/SlParser/grammar"
	sllx "github.com/PlayerR9/SlParser/lexer"
	tr "github.com/PlayerR9/mygo-lib/CustomData/tree"
	"github.com/PlayerR9/mygo-lib/common"
)

func CheckTokens(tokens []*tr.Node, expecteds ...string) error {
	if len(tokens) != len(expecteds) {
		return fmt.Errorf("expected %d tokens, got %d", len(expecteds), len(tokens))
	}

	for i, expected := range expecteds {
		tk := tokens[i]

		tkd, err := grammar.Get[*grammar.TokenData](tk)
		if err != nil {
			return fmt.Errorf("at index %d: %w", i, err)
		}

		type_ := tkd.Type

		if type_ != expected {
			return fmt.Errorf("expected token at index %d to be %s, got %s", i, expected, type_)
		}
	}

	return nil
}

type LexerArg struct {
	InputStr  string
	Expecteds []string
}

func NewLexerArg(input_str string, expecteds ...string) LexerArg {
	return LexerArg{
		InputStr:  input_str,
		Expecteds: expecteds,
	}
}

func LexerTest(lexer *sllx.Lexer, idx int, arg LexerArg) error {
	if lexer == nil {
		return common.NewErrNilParam("lexer")
	}

	_, err := lexer.Write([]byte(arg.InputStr))
	if err != nil {
		return NewErrTestFailed(idx, "failed to set input stream", err)
	}

	itr, _ := lexer.Lex()
	defer itr.Stop()

	var errs []error

	for {
		pair, err := itr.Next()
		if err == sllx.ErrExhausted {
			break
		} else if err != nil {
			err := NewErrTestFailed(idx, "failed to lex", err)
			errs = append(errs, err)

			continue
		}

		tokens := pair.Tokens()
		tokens = append(tokens, slgr.EOFToken)

		err = pair.GetError()
		if err != nil {
			err := NewErrTestFailed(idx, "failed to lex", err)
			errs = append(errs, err)

			continue
		}

		err = CheckTokens(tokens, arg.Expecteds...)
		if err != nil {
			err := NewErrTestFailed(idx, "failed to check tokens", err)
			errs = append(errs, err)

			continue
		}

		return nil
	}

	return errors.Join(errs...)
}
