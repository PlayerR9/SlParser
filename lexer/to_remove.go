package lexer

import gcers "github.com/PlayerR9/errors/error"

// TODO: Remove this once errors is updated.

// AddContext adds context to an error.
//
// Parameters:
//   - err: the error.
//   - key: the context key.
//   - value: the context value.
//
// Does nothing if 'err' is nil.
func AddContext(err *gcers.Err[ErrorCode], key string, value any) {
	if err == nil {
		return
	}

	if err.Context == nil {
		err.Context = make(map[string]any)
	}

	err.Context[key] = value
}
