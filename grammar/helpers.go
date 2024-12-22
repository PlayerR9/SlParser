package grammar

// CheckToken checks if the given token has the given type.
//
// Parameters:
//   - tk: The token to check.
//   - want: The type that the token should have.
//
// Returns:
//   - error: An error if the token does not have the given type.
//
// Errors:
//   - *ErrWant: If the token is nil or does not have the given type.
func CheckToken(tk *Token, want string) error {
	if tk == nil {
		err := NewErrWant(true, "token type", want, nil)
		return err
	}

	type_ := tk.Type

	if type_ != want {
		err := NewErrWant(true, "token type", want, &type_)
		return err
	}

	return nil
}
