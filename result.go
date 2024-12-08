package SlParser

import (
	"errors"
	"fmt"

	gr "github.com/PlayerR9/SlParser/grammar"
	sllx "github.com/PlayerR9/SlParser/lexer"
	slpx "github.com/PlayerR9/SlParser/parser"
	ernk "github.com/PlayerR9/go-evals/rank"
	tr "github.com/PlayerR9/mygo-lib/CustomData/tree"
	"github.com/PlayerR9/mygo-lib/common"
)

/////////////////////////////////////////////////////////

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
type Result struct {
	// data is a pointer to a slice of bytes that is used to store the input data.
	data *[]byte

	// tokens is the list of tokens that were produced during the parsing process.
	tokens *[]*tr.Node

	// lexer_err is the error that occurred during the lexing process.
	lexer_err *error

	// parse_tree is the parse tree that was produced during the parsing process.
	parse_tree **slpx.Result

	// node is the node that was produced during the parsing process.
	node **tr.Node

	// err is the error that occurred during the parsing process.
	err error
}

// HasError implements the Resulter interface.
func (r Result) HasError() bool {
	return r.err != nil
}

// NewResult creates a new result.
//
// Parameters:
//   - data: The data to create the result from.
//
// Returns:
//   - Result: The new result.
func NewResult(data []byte) Result {
	return Result{
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
func (r Result) SetError(err error) Result {
	if err == nil {
		err = r.err
	}

	return Result{
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
func (r Result) Data() ([]byte, error) {
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
func (r Result) Tokens() ([]*tr.Node, error) {
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
func (r Result) ParseTree() (*slpx.Result, error) {
	if r.parse_tree == nil {
		return nil, ErrMissingParseTree
	} else {
		return *r.parse_tree, nil
	}
}

// Node returns the node of the result.
//
// Returns:
//   - grammar.Node: The node of the result.
//   - error: An error if the node is not set.
func (r Result) Node() (*tr.Node, error) {
	if r.node == nil {
		return nil, errors.New("missing node")
	} else {
		return *r.node, nil
	}
}

// LexerErr returns the lexer error of the result.
//
// Returns:
//   - error: The lexer error of the result.
//   - error: An error if the lexer error is not set.
func (r Result) LexerErr() (error, error) {
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
func (r Result) Err() error {
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
func (r Result) Lex(lexer *sllx.Lexer) ([]Result, error) {
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

	eos := ernk.NewErrRorSol[Result]()
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
		tokens = append(tokens, gr.EOFToken)

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

	results := make([]Result, 0, len(errs))

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
//   - []*Result A slice containing the result of the parsing process. If successful,
//     it contains the parse trees generated from the input data. Otherwise, it contains the
//     error that occurred during the parsing process.
//   - error: An error if the evaluation failed.
//
// Errors:
//   - ErrMissingTokens: If the Parse function is called before the tokens are set.
//   - any other error: When the parser is nil or any other error occurs during the parsing process.
func (r Result) Parse(parser slpx.Parser) ([]Result, error) {
	if parser == nil {
		return nil, common.NewErrNilParam("parser")
	}

	var tokens []*tr.Node

	if r.tokens == nil {
		return nil, ErrMissingTokens
	} else {
		tokens = *r.tokens
	}

	eos := ernk.NewErrRorSol[Result]()

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
		results := make([]Result, 0, eos.Size())

		for _, r := range eos.Sols() {
			ok := HasTree(results, (*r.parse_tree).Forest()[0])
			if !ok {
				results = append(results, r)
			}
		}

		return results[:len(results):len(results)], nil
	}

	errs := eos.Errors()

	results := make([]Result, 0, len(errs))

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
//
// Parameters:
//   - table: The AST maker to use for transforming the parse tree into an abstract syntax tree.
//
// Returns:
//   - []Result: A slice containing the result of the transformation process. If
//     successful, it contains the abstract syntax tree nodes generated from the parse tree.
//     Otherwise, it contains the error that occurred during the transformation process.
//   - error: An error if the evaluation failed.
func (r Result) AST(table map[string]gr.ToASTFn) ([]Result, error) {
	if table == nil {
		return nil, common.NewErrNilParam("ast")
	} else if r.parse_tree == nil {
		return nil, ErrMissingParseTree
	}

	forest := (*r.parse_tree).Forest()

	if len(forest) != 1 {
		err := fmt.Errorf("unexpected number of parse trees: %d", len(forest))
		r = r.SetError(err)

		return []Result{r}, nil
	}

	tree := forest[0]

	root := tree.Root()

	node, _ := gr.AST.Transform(root)

	nodes, err := gr.AST.Make(table, node)
	if err == nil && len(nodes) != 1 {
		err = fmt.Errorf("expected one node, got %d instead", len(nodes))
	}

	if len(nodes) == 0 {
		r = r.SetError(err)

		return []Result{r}, nil
	}

	results := make([]Result, 0, len(nodes))

	for i := range nodes {
		result := r.SetError(err)
		result.node = &nodes[i]

		results = append(results, result)
	}

	return results, nil
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
