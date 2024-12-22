package runes

import "errors"

var (
	// ErrInvalidUtf8 occurs when a certain byte sequence is not valid utf-8. This
	// error can be checked with the == operator.
	//
	// Format:
	// 	"invalid utf-8"
	ErrInvalidUtf8 error
)

func init() {
	ErrInvalidUtf8 = errors.New("invalid utf-8")
}
