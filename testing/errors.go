package testing

import (
	"github.com/PlayerR9/mygo-lib/common"
	"github.com/dustin/go-humanize"
)

/////////////////////////////////////////////////////////

type ErrTestFailed struct {
	Idx       int
	Operation string
	Reason    error
}

func (e ErrTestFailed) Error() string {
	var msg string

	if e.Reason == nil {
		msg = "no reason was provided"
	} else {
		msg = e.Reason.Error()
	}

	if e.Operation == "" {
		return "(" + humanize.Ordinal(e.Idx+1) + " test failed) " + msg
	} else {
		return "(" + humanize.Ordinal(e.Idx+1) + " test failed) " + e.Operation + ": " + msg
	}
}

func NewErrTestFailed(idx int, operation string, reason error) error {
	return &ErrTestFailed{
		Idx:       idx,
		Reason:    reason,
		Operation: operation,
	}
}

func AddErr(errs *[]error, idx int, operation string, reason error) error {
	if errs == nil {
		return common.NewErrNilParam("errs")
	}

	err := NewErrTestFailed(idx, operation, reason)
	*errs = append(*errs, err)

	return err
}
