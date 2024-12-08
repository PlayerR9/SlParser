package grammar

import (
	"strconv"

	faults "github.com/PlayerR9/go-fault"
)

/////////////////////////////////////////////////////////

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

var (
	ErrNotAsExpected faults.Descriptor
)

func init() {
	ErrNotAsExpected = faults.New(NotAsExpected, "check failed")
}
