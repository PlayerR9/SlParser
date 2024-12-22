package common

import "errors"

var (
	// ErrNilReceiver occurs when a method is called on a receiver that was not
	// expected to be nil. This error can be checked with the == operator.
	//
	// Format:
	// 	"receiver must not be nil"
	ErrNilReceiver error
)

func init() {
	ErrNilReceiver = errors.New("receiver must not be nil")
}

// ErrBadParam occurs when a parameter is not valid.
type ErrBadParam struct {
	// ParamName is the name of the parameter that is invalid.
	ParamName string

	// Reason is the reason the parameter is invalid.
	Reason string
}

// Error implements error.
func (e ErrBadParam) Error() string {
	var reason string

	if e.Reason == "" {
		reason = "is not valid"
	} else {
		reason = e.Reason
	}

	if e.ParamName == "" {
		return "parameter " + reason
	} else {
		return "parameter (" + e.ParamName + ") " + reason
	}
}

// NewErrBadParam returns an error with the given parameter name and reason.
//
// Parameters:
//   - param_name: The name of the parameter that is invalid.
//   - reason: The reason the parameter is invalid.
//
// Returns:
//   - error: An instance of ErrBadParam. Never returns nil.
//
// Format:
//
//	"parameter (<param_name>) <reason>"
//
// Where:
//   - (<param_name>) is the name of the parameter that is invalid. If empty, it is ignored.
//   - <reason> is the reason the parameter is invalid. If empty, "is not valid" is used.
func NewErrBadParam(param_name, reason string) error {
	e := &ErrBadParam{
		ParamName: param_name,
		Reason:    reason,
	}

	return e
}

// NewErrNilParam returns an error with the given parameter name and reason "must not be nil".
//
// Parameters:
//   - param_name: The name of the parameter that is invalid.
//
// Returns:
//   - error: An instance of ErrBadParam. Never returns nil.
//
// Format:
//
//	"parameter (<param_name>) must not be nil"
//
// Where:
//   - (<param_name>) is the name of the parameter that is invalid. If empty, it is ignored.
func NewErrNilParam(param_name string) error {
	e := &ErrBadParam{
		ParamName: param_name,
		Reason:    "must not be nil",
	}
	return e
}
