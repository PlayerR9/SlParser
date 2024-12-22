package ast

import (
	"fmt"

	slgr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/mygo-lib/common"
)

// ApplyAST applies a function to the token token to convert it into a node of type N.
//
// The function attempts to construct a node from the token token using the provided
// conversion function. If the token token is missing or an error occurs during the
// conversion process, the result will contain an error.
//
// Parameters:
//   - tk: The token token of the tree.
//   - table: A table of functions that convert a token into a node of type N.
//
// Returns:
//   - N: The node resulting from the conversion of the token token.
//   - error: An error if the conversion process fails.
//
// Errors:
//   - common.ErrBadParam: If the token or table are nil.
//   - fmt.Errorf("unknown type: %s", strconv.Quote(token.Type)): If the token's type is not recognized.
//   - any other error: Function-specific.
func ApplyAST[N slgr.TreeNode](tk *slgr.Token, table ASTMaker[N]) (N, error) {
	if tk == nil {
		err := common.NewErrNilParam("tk")
		return *new(N), err
	}

	if table == nil {
		err := common.NewErrNilParam("table")
		return *new(N), err
	}

	fn, ok := table[tk.Type]
	if !ok || fn == nil {
		err := fmt.Errorf("unknown type: %s", tk.Type)
		return *new(N), err
	}

	n, err := fn(tk)
	if err != nil {
		return *new(N), err
	}

	return n, nil
}
