package slp

import (
	"errors"

	slpast "github.com/PlayerR9/SlParser/ast"
	slgr "github.com/PlayerR9/SlParser/grammar"
	sllx "github.com/PlayerR9/SlParser/lexer"
	"github.com/PlayerR9/SlParser/mygo-lib/common"
	slpx "github.com/PlayerR9/SlParser/parser"
)

// Result is a struct that holds the result of a parsing operation.
type Result[N slgr.TreeNode] struct {
	// data is a pointer to the input data.
	data *[]byte

	// tokens is a pointer to the list of tokens.
	tokens *[]*slgr.Token

	// root is a pointer to the root token.
	root **slgr.Token

	// node is a pointer to the root node.
	node *N

	// err is an error that occurred during the parsing process.
	err error
}

// HasError checks if the result has an error.
//
// Returns:
//   - bool: True if the result has an error, false otherwise.
func (r Result[N]) HasError() bool {
	return r.err != nil
}

// GetError returns the error associated with the result.
//
// Returns:
//   - error: The error if available, or nil if not.
func (r Result[N]) GetError() error {
	return r.err
}

// GetData returns the input data associated with the result.
//
// Returns:
//   - []byte: The input data if available, or nil if not.
//   - error: An error if the data is missing, otherwise nil.
func (r Result[N]) GetData() ([]byte, error) {
	if r.data == nil {
		return nil, errors.New("missing data")
	} else {
		return *r.data, nil
	}
}

// GetTokens returns the list of tokens associated with the result.
//
// Returns:
//   - []*slgr.Token: The list of tokens if available, or nil if not.
//   - error: An error if the tokens are missing, otherwise nil.
func (r Result[N]) GetTokens() ([]*slgr.Token, error) {
	if r.tokens == nil {
		return nil, errors.New("missing tokens")
	} else {
		return *r.tokens, nil
	}
}

// GetRoot returns the root token associated with the result.
//
// Returns:
//   - *slgr.Token: The root token if available, or nil if not.
//   - error: An error if the root token is missing, otherwise nil.
func (r Result[N]) GetRoot() (*slgr.Token, error) {
	if r.root == nil {
		return nil, errors.New("missing root")
	} else {
		return *r.root, nil
	}
}

// GetNode returns the root node associated with the result.
//
// Returns:
//   - N: The root node if available, or a newly allocated node of type N if not.
//   - error: An error if the root node is missing, otherwise nil.
func (r Result[N]) GetNode() (N, error) {
	if r.node == nil {
		return *new(N), errors.New("missing node")
	} else {
		return *r.node, nil
	}
}

// NewResult creates a new Result with the specified input data.
//
// Parameters:
//   - data: The input data to be processed.
//
// Returns:
//   - Result[N]: A new result with the specified input data.
func NewResult[N slgr.TreeNode](data []byte) Result[N] {
	r := Result[N]{
		data: &data,
	}

	return r
}

// Copy returns a deep copy of the result.
//
// Returns:
//   - Result[N]: A deep copy of the result.
func (r Result[N]) Copy() Result[N] {
	var r_copy Result[N]

	if r.data == nil {
		r_copy.data = nil
	} else {
		data := make([]byte, len(*r.data))
		copy(data, *r.data)

		r_copy.data = &data
	}

	if r.tokens == nil {
		r_copy.tokens = nil
	} else {
		tokens := make([]*slgr.Token, len(*r.tokens))
		copy(tokens, *r.tokens)

		r_copy.tokens = &tokens
	}

	if r.root == nil {
		r_copy.root = nil
	} else {
		r_copy.root = r.root
	}

	if r.node == nil {
		r_copy.node = nil
	} else {
		r_copy.node = r.node
	}

	r_copy.err = r.err

	return r_copy
}

// Lex tokenizes the input data using the provided lexing function.
//
// The function attempts to lex the input data into a list of tokens. If the
// input data is missing or an error occurs during lexing, the result will
// contain an error.
//
// Parameters:
//   - lexer: A lexer that can be used to lex the input data.
//
// Returns:
//   - Result[N]: A new result containing the list of lexed tokens or an error
//     if the lexing process fails.
func (r Result[N]) Lex(lexer *sllx.Lexer) Result[N] {
	if lexer == nil {
		r := r.Copy()
		r.err = common.NewErrNilParam("lexer")

		return r
	}

	if r.data == nil {
		r := r.Copy()
		r.err = errors.New("missing data")

		return r
	}

	tokens, err := sllx.Lex(lexer, *r.data)
	if err != nil {
		r := r.Copy()
		r.err = err
		return r
	} else {
		r := r.Copy()
		r.tokens = &tokens
		return r
	}
}

// Parse parses the list of tokens using the provided parsing function.
//
// The function attempts to parse the list of tokens into a root token. If the
// list of tokens is missing or an error occurs during parsing, the result will
// contain an error.
//
// Parameters:
//   - parse_one_fn: A function that parses one token from the list of tokens.
//
// Returns:
//   - Result[N]: A new result containing the root token or an error if the
//     parsing process fails.
func (r Result[N]) Parse(parser *slpx.Parser) Result[N] {
	if r.tokens == nil {
		r := r.Copy()
		r.err = errors.New("missing tokens")

		return r
	}

	forest, err := slpx.Parse(parser, *r.tokens)
	if err != nil {
		r := r.Copy()
		r.err = err
		return r
	}

	if len(forest) != 1 {
		r := r.Copy()
		r.err = errors.New("expected one root token")
		return r
	}

	r = r.Copy()
	r.root = &forest[0]
	return r
}

// Ast applies a function to the root token to convert it into a node of type N.
//
// The function attempts to construct a node from the root token using the provided
// conversion function. If the root token is missing or an error occurs during the
// conversion process, the result will contain an error.
//
// Parameters:
//   - fn: A function that converts a token into a node of type N.
//
// Returns:
//   - Result[N]: A new result containing the node or an error if the conversion
//     process fails.
func (r Result[N]) Ast(table slpast.ASTMaker[N]) Result[N] {
	if r.root == nil {
		r := r.Copy()
		r.err = errors.New("missing root")

		return r
	}

	node, err := slpast.ApplyAST(*r.root, table)
	if err != nil {
		r := r.Copy()
		r.err = err
		return r
	} else {
		r := r.Copy()
		r.node = &node
		return r
	}
}
