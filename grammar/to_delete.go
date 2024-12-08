package grammar

import (
	"fmt"

	tr "github.com/PlayerR9/mygo-lib/CustomData/tree"
	common "github.com/PlayerR9/mygo-lib/common"
)

/////////////////////////////////////////////////////////

// TODO: Delete this once mygo-lib is updated.

// ErrInvalidType occurs when a type is not as expected.
type ErrInvalidType struct {
	// Types are the expected types.
	Types []any

	// Got is the actual type.
	Got any
}

// Error implements the error interface.
func (e ErrInvalidType) Error() string {
	var expected string

	switch len(e.Types) {
	case 0:
		expected = "no types"
	case 1:
		expected = fmt.Sprintf("%T", e.Types[0])
	default:
		elems := make([]string, 0, len(e.Types))

		for _, elem := range e.Types {
			elems = append(elems, fmt.Sprintf("%T", elem))
		}

		expected = common.EitherOrString(elems)
	}

	var got string

	if e.Got == nil {
		got = "<nil>"
	} else {
		got = fmt.Sprintf("%T", e.Got)
	}

	return "want " + expected + ", got " + got
}

// NewErrInvalidType creates a new ErrInvalidType error with the specified expected types and actual type.
//
// Parameters:
//   - wants: The expected types.
//   - got: The actual type.
//
// Returns:
//   - error: The new ErrInvalidType error. Never returns nil.
//
// Format:
//
//	"want <want>, got <got>"
//
// Where:
//   - <want>: The expected type.
//   - <got>: The actual type. If nil, "<nil>" is used instead.
func NewErrInvalidType(got any, wants ...any) error {
	return &ErrInvalidType{
		Types: wants,
		Got:   got,
	}
}

// Get returns the value of the node if it is of type T, or an error if the node is nil or the
// information of the node is not of type T.
//
// Parameters:
//   - node: The node to get the value of.
//
// Returns:
//   - T: The value of the node if the node is not nil and the information of the node is of type T.
//   - error: An error if the node is nil, or the information of the node is not of type T.
//
// Errors:
//   - common.ErrBadParam: If the node is nil.
//   - common.ErrInvalidType: If the information of the node is not of type T, including if node.Info is nil.
func Get[T tr.Infoer](node *tr.Node) (T, error) {
	if node == nil {
		return *new(T), common.NewErrNilParam("node")
	}

	info := node.Info
	if info == nil {
		return *new(T), NewErrInvalidType(nil, *new(T))
	}

	v, ok := info.(T)
	if !ok {
		return *new(T), NewErrInvalidType(info, v)
	}

	return v, nil
}

// MustGet returns the value of the node if it is of type T, or panics if the node is nil or the
// information of the node is not of type T.
//
// Parameters:
//   - node: The node to get the value of.
//
// Panics:
//   - common.ErrNilParam: If the node is nil.
//   - common.ErrInvalidType: If the information of the node is not of type T, including if node.Info is nil.
//
// Returns:
//   - T: The value of the node if the node is not nil and the information of the node is of type T.
func MustGet[T tr.Infoer](node *tr.Node) T {
	if node == nil {
		panic(common.NewErrNilParam("node"))
	}

	info := node.Info
	if info == nil {
		panic(NewErrInvalidType(nil, *new(T)))
	}

	v, ok := info.(T)
	if !ok {
		panic(NewErrInvalidType(info, v))
	}

	return v
}
