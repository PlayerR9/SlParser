package grammar

import "fmt"

// TokenTyper is a token type.
type TokenTyper interface {
	~int

	// IsAccept returns whether the token is an accept token.
	//
	// Returns:
	// 	- bool: True if the token is an accept token, false otherwise.
	IsAccept() bool

	fmt.GoStringer
	fmt.Stringer
}

// Token is a token.
type Token struct {
}
