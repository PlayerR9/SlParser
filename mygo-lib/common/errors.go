package common

import "errors"

var (
	// ErrNilReceiver occurs when a method is called on a receiver who was
	// expected to be non-nil. This error can be checked with the == operator.
	//
	// Format:
	// 	"receiver must not be nil"
	ErrNilReceiver error
)

func init() {
	ErrNilReceiver = errors.New("receiver must not be nil")
}
