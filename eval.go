package SlParser

import (
	"iter"

	slast "github.com/PlayerR9/SlParser/ast"
	sllx "github.com/PlayerR9/SlParser/lexer"
	slpx "github.com/PlayerR9/SlParser/parser"
	evrsl "github.com/PlayerR9/go-evals/result"
	"github.com/PlayerR9/mygo-lib/common"
)

func MakeEvaluate[N interface {
	Child() iter.Seq[N]

	slast.Noder
}](lexer *sllx.Lexer, parser slpx.Parser, ast *slast.ASTMaker[N]) (evrsl.ApplyOnValidsFn[Result[N]], error) {
	if lexer == nil {
		return nil, common.NewErrNilParam("lexer")
	} else if parser == nil {
		return nil, common.NewErrNilParam("parser")
	} else if ast == nil {
		return nil, common.NewErrNilParam("ast")
	}

	evaluateParseFn := func(elem Result[N]) ([]Result[N], error) {
		return elem.Parse(parser)
	}

	evaluateASTFn := func(elem Result[N]) ([]Result[N], error) {
		return elem.AST(ast)
	}

	evaluateLexerFn := func(elem Result[N]) ([]Result[N], error) {
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
