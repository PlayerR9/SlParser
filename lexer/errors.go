package lexer

import "errors"

var (
	// ErrCannotUnread occurs when the lexer cannot unread a character. This error
	// can be checked with the == operator.
	//
	// Format:
	// 	"nothing to unread"
	ErrCannotUnread error

	// ErrNotFound occurs when the lexer cannot find a token. This error can be
	// checked with the == operator.
	//
	// Format:
	// 	"token not found"
	ErrNotFound error
)

func init() {
	ErrCannotUnread = errors.New("nothing to unread")
	ErrNotFound = errors.New("token not found")
}
