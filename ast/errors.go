package ast

import (
	"strconv"
)

// ErrUnsupportedType is an error that occurs when a token type is not supported.
type ErrUnsupportedType struct {
	// Quote is a flag that indicates whether the type is quoted or not.
	Quote bool

	// Type is the type that is not supported.
	Type string
}

// Error implements the error interface.
func (e ErrUnsupportedType) Error() string {
	if e.Quote {
		return strconv.Quote(e.Type) + " is not a supported token type"
	} else {
		return e.Type + " is not a supported token type"
	}
}

// NewErrUnsupportedType creates a new instance of ErrUnsupportedType.
//
// Parameters:
//   - quote: A flag that indicates whether the type is quoted or not.
//   - type_: The type that is not supported.
//
// Returns:
//   - error: The new instance of ErrUnsupportedType. Never returns nil.
func NewErrUnsupportedType(quote bool, type_ string) error {
	return &ErrUnsupportedType{
		Quote: quote,
		Type:  type_,
	}
}

/*
type ErrBadData struct {
	Reason error
}

func (e ErrBadData) Error() string {
	var msg string

	if e.Reason == nil {
		msg = "something went wrong"
	} else {
		msg = e.Reason.Error()
	}

	return "(ErrBadData) " + msg
}

func NewErrBadData(reason error) error {
	return &ErrBadData{
		Reason: reason,
	}
}

type ErrBadChildren struct {
	Reason error
}

func (e ErrBadChildren) Error() string {
	var msg string

	if e.Reason == nil {
		msg = "something went wrong"
	} else {
		msg = e.Reason.Error()
	}

	return "(ErrBadChildren) " + msg
}

func NewErrBadChildren(reason error) error {
	return &ErrBadChildren{
		Reason: reason,
	}
}

var (
	ErrNilNode error
)

func init() {
	ErrNilNode = errors.New("nil nodes are not allowed")
}
*/
