package matcher

import "errors"

var (
	// ErrMatchDone occurs when a match is done. This error can be checked
	// with the == operator.
	//
	// Format:
	// 	"matcher is done"
	ErrMatchDone error
)

func init() {
	ErrMatchDone = errors.New("matcher is done")
}
