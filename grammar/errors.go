package grammar

import (
	"strconv"
	"strings"
)

// ErrWant is an error that indicates that the want and got values are not equal.
type ErrWant struct {
	// Quote indicates if the want and got values should be quoted.
	Quote bool

	// Kind is the type of the want and got values.
	Kind string

	// Want is the expected value.
	Want string

	// Got is the actual value.
	Got *string
}

// Error implements error.
func (e ErrWant) Error() string {
	var want, got string

	if e.Want == "" {
		want = "something"
	} else if e.Quote {
		want = strconv.Quote(e.Want)
	} else {
		want = e.Want
	}

	if e.Got == nil {
		got = "nothing"
	} else if e.Quote {
		got = strconv.Quote(*e.Got)
	} else {
		got = *e.Got
	}

	var builder strings.Builder

	_, _ = builder.WriteString("want ")

	if e.Kind != "" {
		_, _ = builder.WriteString(e.Kind)
		_, _ = builder.WriteString(" to be ")
	}

	_, _ = builder.WriteString(want)
	_, _ = builder.WriteString(", got ")
	_, _ = builder.WriteString(got)

	str := builder.String()
	return str
}

// NewErrWant returns an error with a "want <want>, got <got>" message.
//
// Parameters:
//   - quote: Whether to quote the want and got strings.
//   - kind: The kind of value that is being compared, if any.
//   - want: The value that was expected.
//   - got: The value that was actually received.
//
// Returns:
//   - error: An instance of ErrWant. Never returns nil.
//
// Format:
//
//	"want <kind> to be <want>, got <got>"
//
// Where:
//   - kind: The kind of value that is being compared, if any.
//   - want: The value that was expected.
//   - got: The value that was actually received.
//
// If the quote parameter is true, the want and got values will be quoted when a value is provided.
func NewErrWant(quote bool, kind, want string, got *string) error {
	e := &ErrWant{
		Quote: quote,
		Kind:  kind,
		Want:  want,
		Got:   got,
	}
	return e
}
