package lexer

type LexerState struct {
	// last_char_read is the last character read.
	//
	// This is updated every time input stream is read.
	last_char_read *rune

	// last_char_err is the last error encountered.
	//
	// This is updated every time an error is encountered.
	last_err error
}

// UpdateLastCharRead updates the last character read.
//
// Does nothing if the receiver is nil.
//
// Parameters:
//   - char: the character to update.
func (ls *LexerState) UpdateLastCharRead(char *rune) {
	if ls == nil {
		return
	}

	ls.last_char_read = char
}

// UpdateLastErr updates the last error encountered.
//
// Does nothing if the receiver is nil.
//
// Parameters:
//   - err: the error to update.
func (ls *LexerState) UpdateLastErr(err error) {
	if ls == nil {
		return
	}

	ls.last_err = err
}

// Reset resets the state.
//
// Does nothing if the receiver is nil.
func (ls *LexerState) Reset() {
	if ls == nil {
		return
	}

	ls.last_char_read = nil
	ls.last_err = nil
}
