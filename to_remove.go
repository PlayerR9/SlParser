package SlParser

import (
	"fmt"

	gcerr "github.com/PlayerR9/errors/error"
	dba "github.com/PlayerR9/go-debug/assert"
)

// TODO: Remove this once errors is updated.

func Value[T gcerr.ErrorCoder](err gcerr.Err[T], key string) (any, bool) {
	if len(err.Context) == 0 {
		return nil, false
	}

	value, ok := err.Context[key]
	return value, ok
}

// TODO: Remove this once go-debug is updated.

// AssertConv tries to convert an element to the expected type and panics if it is not possible.
//
// Parameters:
//   - elem: the element to check.
//   - var_name: the name of the variable.
//
// Returns:
//   - T: the converted element.
func AssertConv[T any](elem any, target string) T {
	if elem == nil {
		msg := fmt.Sprintf("expected %q to be of type %T, got nil instead", target, *new(T))

		panic(dba.NewErrAssertFailed(msg))
	}

	res, ok := elem.(T)
	if !ok {
		msg := fmt.Sprintf("expected %q to be of type %T, got %T instead", target, *new(T), elem)

		panic(dba.NewErrAssertFailed(msg))
	}

	return res
}
