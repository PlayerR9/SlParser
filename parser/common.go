package parser

import (
	gr "github.com/PlayerR9/SlParser/grammar"
)

// CheckTop is a function that checks if the top of the stack is in the allowed list.
//
// Parameters:
//   - parser: The parser to check.
//   - allowed: The list of allowed tokens.
//
// Returns:
//   - *gr.Token[T]: The top of the stack.
//   - bool: True if the top of the stack is in the allowed list, false otherwise.
//
// If the receiver is nil, then it returns nil and false.
//
// If no allowed tokens are provided, then it returns the top of the stack and false.
func CheckTop[T gr.TokenTyper](parser *Parser[T], allowed ...T) (*gr.Token[T], bool) {
	if parser == nil {
		return nil, false
	}

	top, ok := parser.Pop()
	if !ok || len(allowed) == 0 {
		return top, false
	}

	for _, a := range allowed {
		if top.Type == a {
			return top, true
		}
	}

	return top, false
}

// CheckLookahead is a function that checks if the lookahead is in the allowed list.
//
// Parameters:
//   - lookahead: The lookahead to check.
//   - allowed: The list of allowed tokens.
//
// Returns:
//   - bool: True if the lookahead is in the allowed list, false otherwise.
//
// If the receiver is nil or no allowed tokens are provided, then it returns false.
func CheckLookahead[T gr.TokenTyper](lookahead *gr.Token[T], allowed ...T) bool {
	if lookahead == nil || len(allowed) == 0 {
		return false
	}

	for _, a := range allowed {
		if lookahead.Type == a {
			return true
		}
	}

	return false
}
