package SlParser

import (
	gr "github.com/PlayerR9/SlParser/grammar"
	sllx "github.com/PlayerR9/SlParser/lexer"
	slpx "github.com/PlayerR9/SlParser/parser"
	evrsl "github.com/PlayerR9/go-evals/result"
	"github.com/PlayerR9/mygo-lib/common"
)

// MakeEvaluate makes an evaluator function that evaluates a sequence of SlParser
// results. The evaluator function takes a sequence of SlParser results and returns
// a new sequence of SlParser results. The new sequence of SlParser results is
// computed by first attempting to lex the input, then attempting to parse the
// lexer output, and finally attempting to convert the parse output to an AST.
//
// Parameters:
//   - lexer: The lexer to use for lexing.
//   - parser: The parser to use for parsing.
//   - ast: The AST maker to use for converting the parse output to an AST.
//
// Returns:
//   - evrsl.ApplyOnValidsFn[Result]: The evaluator function.
//   - error: An error if the operation fails.
func MakeEvaluate(lexer *sllx.Lexer, parser slpx.Parser, ast map[string]gr.ToASTFn) (evrsl.ApplyOnValidsFn[Result], error) {
	if lexer == nil {
		return nil, common.NewErrNilParam("lexer")
	} else if parser == nil {
		return nil, common.NewErrNilParam("parser")
	} else if ast == nil {
		return nil, common.NewErrNilParam("ast")
	}

	evaluateParseFn := func(elem Result) ([]Result, error) {
		return elem.Parse(parser)
	}

	evaluateASTFn := func(elem Result) ([]Result, error) {
		return elem.AST(ast)
	}

	evaluateLexerFn := func(elem Result) ([]Result, error) {
		return elem.Lex(lexer)
	}

	astRunFn, _ := evrsl.MakeRunFn(evaluateASTFn, nil)
	astApplyFn, _ := evrsl.MakeApplyFn(astRunFn)
	parseRunFn, _ := evrsl.MakeRunFn(evaluateParseFn, astApplyFn)
	parseApplyFn, _ := evrsl.MakeApplyFn(parseRunFn)
	lexerRunFn, _ := evrsl.MakeRunFn(evaluateLexerFn, parseApplyFn)

	applyFn, _ := evrsl.MakeApplyFn(lexerRunFn)
	return applyFn, nil
}
