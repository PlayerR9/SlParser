package ast

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
