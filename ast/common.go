package ast

import (
	gr "github.com/PlayerR9/SlParser/grammar"
	util "github.com/PlayerR9/SlParser/util"
	gcers "github.com/PlayerR9/go-commons/errors"
)

// CheckType is a helper function that checks the type of the token at the given
// position.
//
// Parameters:
//   - children: The list of children.
//   - at: The position of the token.
//   - type_: The type of the token.
//
// Returns:
//   - error: if an error occurred.
//
// Errors:
//   - *errors.ErrInvalidParameter: If 'at' is less than 0.
//   - *errors.ErrValue: If the token at the given position is nil or
//     'at' is out of range.
func CheckType[T gr.TokenTyper](children []*gr.Token[T], at int, type_ T) error {
	if at < 0 {
		return gcers.NewErrInvalidParameter("at", gcers.NewErrGTE(0))
	}

	pos_str := gcers.GetOrdinalSuffix(at+1) + " child"

	if at >= len(children) {
		return util.NewErrValue(pos_str, type_, nil, true)
	}

	tk := children[at]
	if tk == nil {
		return util.NewErrValue(pos_str, type_, nil, true)
	} else if tk.Type != type_ {
		return util.NewErrValue(pos_str, type_, tk.Type, true)
	}

	return nil
}
