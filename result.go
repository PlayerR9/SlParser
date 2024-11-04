package SlParser

import (
	"errors"
	"fmt"
	"iter"

	"github.com/PlayerR9/SlParser/ast"
	slgr "github.com/PlayerR9/SlParser/grammar"
	sllx "github.com/PlayerR9/SlParser/lexer"
	slpx "github.com/PlayerR9/SlParser/parser"

	ernk "github.com/PlayerR9/go-evals/rank"
	"github.com/PlayerR9/mygo-lib/common"
)

var (
	// ErrMissingData is the error that is returned when data is missing.
	ErrMissingData error

	// ErrMissingTokens is the error that is returned when tokens are missing.
	ErrMissingTokens error

	// ErrMissingParseTree is the error that is returned when the parse tree is missing.
	ErrMissingParseTree error
)

func init() {
	ErrMissingData = errors.New("missing data")
	ErrMissingTokens = errors.New("missing tokens")
	ErrMissingParseTree = errors.New("missing parse tree")
}

// Result holds all the information regarding the parsing process.
type Result[N interface {
	Child() iter.Seq[N]

	ast.Noder
}] struct {
	// data is a pointer to a slice of bytes that is used to store the input data.
	data *[]byte

	// tokens is the list of tokens that were produced during the parsing process.
	tokens *[]*slgr.Token

	// lexer_err is the error that occurred during the lexing process.
	lexer_err *error

	// parse_tree is the parse tree that was produced during the parsing process.
	parse_tree **slpx.Result

	// node is the node that was produced during the parsing process.
	node *N

	// err is the error that occurred during the parsing process.
	err error
}

// HasError implements the Resulter interface.
func (r Result[N]) HasError() bool {
	return r.err != nil
}

// NewResult creates a new result.
//
// Parameters:
//   - data: The data to create the result from.
//
// Returns:
//   - Result[N]: The new result.
func NewResult[N interface {
	Child() iter.Seq[N]

	ast.Noder
}](data []byte) Result[N] {
	return Result[N]{
		data: &data,
	}
}

// SetError sets the error of the result.
//
// Parameters:
//   - err: The error to set.
//
// Returns:
//   - Result[T, N]: The result with the error set.
func (r Result[N]) SetError(err error) Result[N] {
	if err == nil {
		err = r.err
	}

	return Result[N]{
		data:       r.data,
		tokens:     r.tokens,
		lexer_err:  r.lexer_err,
		parse_tree: r.parse_tree,
		node:       r.node,
		err:        err,
	}
}

// Data returns the data of the result.
//
// Returns:
//   - []byte: The data of the result.
//   - error: An error if the data is not set.
//
// Errors:
//   - ErrMissingData: If the data is not set.
func (r Result[N]) Data() ([]byte, error) {
	if r.data == nil {
		return nil, ErrMissingData
	} else {
		return *r.data, nil
	}
}

// Tokens returns the tokens of the result.
//
// Returns:
//   - []*grammar.Token: The tokens of the result.
//   - error: An error if the tokens are not set.
func (r Result[N]) Tokens() ([]*slgr.Token, error) {
	if r.tokens == nil {
		return nil, ErrMissingTokens
	} else {
		return *r.tokens, nil
	}
}

// ParseTree returns the parse tree of the result.
//
// Returns:
//   - *slpx.Result: The parse tree of the result.
//   - error: An error if the parse tree is not set.
func (r Result[N]) ParseTree() (*slpx.Result, error) {
	if r.parse_tree == nil {
		return nil, ErrMissingParseTree
	} else {
		return *r.parse_tree, nil
	}
}

// Node returns the node of the result.
//
// Returns:
//   - N: The node of the result.
//   - error: An error if the node is not set.
func (r Result[N]) Node() (N, error) {
	if r.node == nil {
		return *new(N), errors.New("missing node")
	} else {
		return *r.node, nil
	}
}

// LexerErr returns the lexer error of the result.
//
// Returns:
//   - error: The lexer error of the result.
//   - error: An error if the lexer error is not set.
func (r Result[N]) LexerErr() (error, error) {
	if r.lexer_err == nil {
		return nil, errors.New("missing lexer error")
	} else {
		return *r.lexer_err, nil
	}
}

// Err returns the error of the result.
//
// Returns:
//   - error: The error of the result.
func (r Result[N]) Err() error {
	return r.err
}

// Lex processes the input data using the provided lexer and returns a slice of results.
//
// Parameters:
//   - lexer: The lexer to use for processing the input data.
//
// Returns:
//   - []*Result[T, N]: A slice containing the result of the lexing process. If successful,
//     it contains the tokens generated from the input data. Otherwise, it contains the
//     error that occurred during the lexing process.
//   - error: An error if the evaluation failed.
//
// Errors:
//   - ErrMissingData: If the Lex function is called before the data is set.
//   - any other error: When the lexer is nil or any other error occurs during the lexing process.
func (r Result[N]) Lex(lexer *sllx.Lexer) ([]Result[N], error) {
	if lexer == nil {
		return nil, common.NewErrNilParam("lexer")
	}

	defer lexer.Reset()

	var data []byte

	if r.data == nil {
		return nil, ErrMissingData
	} else {
		data = *r.data
	}

	_, err := lexer.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to set input stream: %w", err)
	}

	itr, _ := lexer.Lex()
	defer itr.Stop()

	eos := ernk.NewErrRorSol[Result[N]]()
	eos.ChangeOrder(true)

	for {
		pair, err := itr.Next()
		if err == sllx.ErrExhausted {
			break
		} else if err != nil {
			return nil, fmt.Errorf("failed to lex: %w", err)
		}

		err = pair.GetError()
		tokens := pair.Tokens()
		tokens = append(tokens, slgr.EOFToken)

		if err == nil {
			r := r.SetError(nil)
			r.tokens = &tokens

			_ = eos.AddSol(len(tokens), r)
		} else {
			_ = eos.AddErr(len(tokens), err)
		}
	}

	if !eos.HasError() {
		return eos.Sols(), nil
	}

	errs := eos.Errors()

	results := make([]Result[N], 0, len(errs))

	for _, err := range errs {
		r := r.SetError(err)
		r.lexer_err = &err

		results = append(results, r)
	}

	return results, nil
}

// Parse processes the input data using the provided parser and returns a slice of results.
//
// Parameters:
//   - parser: The parser to use for processing the input data.
//
// Returns:
//   - []*Result[T, N]: A slice containing the result of the parsing process. If successful,
//     it contains the parse trees generated from the input data. Otherwise, it contains the
//     error that occurred during the parsing process.
//   - error: An error if the evaluation failed.
//
// Errors:
//   - ErrMissingTokens: If the Parse function is called before the tokens are set.
//   - any other error: When the parser is nil or any other error occurs during the parsing process.
func (r Result[N]) Parse(parser slpx.Parser) ([]Result[N], error) {
	if parser == nil {
		return nil, common.NewErrNilParam("parser")
	}

	var tokens []*slgr.Token

	if r.tokens == nil {
		return nil, ErrMissingTokens
	} else {
		tokens = *r.tokens
	}

	eos := ernk.NewErrRorSol[Result[N]]()

	itr := parser.Parse(tokens)
	defer itr.Stop()

	for {
		pair, err := itr.Next()
		if err == slpx.ErrExhausted {
			break
		} else if err != nil {
			return nil, fmt.Errorf("failed to parse: %w", err)
		}

		forest := pair.Forest()
		err = pair.GetError()

		if err != nil {
			_ = eos.AddErr(0, err)
		} else if len(forest) != 1 {
			_ = eos.AddErr(1, fmt.Errorf("unexpected number of parse trees: %d", len(forest)))
		} else {
			r := r.SetError(nil)
			r.parse_tree = &pair

			_ = eos.AddSol(0, r)
		}
	}

	if !eos.HasError() {
		results := make([]Result[N], 0, eos.Size())

		for _, r := range eos.Sols() {
			ok := HasTree(results, (*r.parse_tree).Forest()[0])
			if !ok {
				results = append(results, r)
			}
		}

		return results[:len(results):len(results)], nil
	}

	errs := eos.Errors()

	results := make([]Result[N], 0, len(errs))

	for _, err := range errs {
		r := r.SetError(err)
		results = append(results, r)
	}

	return results, nil
}

// AST transforms the parse tree into an abstract syntax tree.
//
// The `ast` function must return one abstract syntax tree node for each
// parse tree root node. If the `ast` function returns an error, the entire
// result is marked as invalid.
//
// If the `ast` function is nil, an error is returned.
//
// If the parse tree is missing, an error is returned.
//
// If the `ast` function returns more or less than one abstract syntax tree
// node, an error is returned.
func (r Result[N]) AST(ast_fn *ast.ASTMaker[N]) ([]Result[N], error) {
	if ast_fn == nil {
		return nil, common.NewErrNilParam("ast")
	}

	var tree *slpx.ParseTree

	if r.parse_tree == nil {
		return nil, ErrMissingParseTree
	} else {
		tree = (*r.parse_tree).Forest()[0]
	}

	root := tree.Root()

	errs := make([]error, 0, 2)

	nodes, err := ast_fn.Make(root)
	if err != nil {
		errs = append(errs, err)
	}

	if len(nodes) != 1 {
		errs = append(errs, fmt.Errorf("expected one node, got %d instead", len(nodes)))
	}

	err = errors.Join(errs...)

	if len(nodes) == 0 {
		r.SetError(err)

		return []Result[N]{r}, nil
	} else {
		results := make([]Result[N], 0, len(nodes))

		for i := range nodes {
			result := r.SetError(err)
			result.node = &nodes[i]

			results = append(results, result)
		}

		return results, nil
	}
}

type ModifyFn[T interface {
	SetError(err error) T
}, R any] func(result *T, elem R)

// ApplyResults applies a function to each element in the slice and returns a new slice
// of results.
//
// Parameters:
//   - parent: The parent result.
//   - elems: The slice of elements to apply the function to.
//   - err: The error to set on each result.
//   - fn: The function to apply to each element.
//
// Returns:
//   - []Result[T, N]: A slice of results with the function applied to each element.
//   - error: An error if the function is nil.
func ApplyResults[T interface {
	SetError(err error) T
}, R any](parent T, elems []R, err error, fn ModifyFn[T, R]) ([]T, error) {
	if fn == nil {
		return nil, common.NewErrNilParam("fn")
	} else if len(elems) == 0 {
		return nil, nil
	}

	results := make([]T, 0, len(elems))

	for i := range elems {
		result := parent.SetError(err)
		fn(&result, elems[i])

		results = append(results, result)
	}

	return results, nil
}
