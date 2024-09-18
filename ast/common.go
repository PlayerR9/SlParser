package ast

import (
	gr "github.com/PlayerR9/SlParser/grammar"
	util "github.com/PlayerR9/SlParser/util"
)

func CheckType[T gr.TokenTyper](tk *gr.Token[T], type_ T) error {
	if tk == nil {
		return util.NewErrValue("first child", type_, nil, true)
	}

	if tk.Type != type_ {
		return util.NewErrValue("first child", type_, tk.Type, true)
	}

	return nil
}
