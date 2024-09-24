package lexer

import "strings"

type LexerState struct {
	// last_char_read is the last character read.
	//
	// This is updated every time input stream is read.
	last_char_read *rune

	// last_char_err is the last error encountered.
	//
	// This is updated every time an error is encountered.
	last_err error

	// builder is the lexer builder.
	builder strings.Builder

	// char_buff is the character buffer.
	char_buff []rune
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

// SetChar sets the character buffer.
//
// Does nothing if the receiver is nil.
//
// Parameters:
//   - char: the character to set.
func (ls *LexerState) SetChar(char rune) {
	if ls == nil {
		return
	}

	ls.char_buff = append(ls.char_buff, char)
}

// RemoveChar removes the last character from the buffer.
//
// Does nothing if the receiver is nil or the buffer is empty.
func (ls *LexerState) RemoveChar() {
	if ls == nil || len(ls.char_buff) == 0 {
		return
	}

	ls.char_buff = ls.char_buff[:len(ls.char_buff)-1]
}

// GetData returns the data in the buffer.
//
// Does nothing if the receiver is nil.
//
// Returns:
//   - string: the data in the buffer.
func (ls *LexerState) GetData() string {
	if ls == nil {
		return ""
	}

	if len(ls.char_buff) > 0 {
		for i := 0; i < len(ls.char_buff); i++ {
			ls.builder.WriteRune(ls.char_buff[i])
		}

		ls.char_buff = ls.char_buff[:0]
	}

	str := ls.builder.String()
	ls.builder.Reset()

	return str
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
	ls.builder.Reset()

	if len(ls.char_buff) > 0 {
		ls.char_buff = ls.char_buff[:0]
	}
}
