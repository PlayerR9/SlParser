package lexer

import (
	"bytes"
	"fmt"

	"github.com/PlayerR9/mygo-lib/common"
	gch "github.com/PlayerR9/mygo-lib/runes"
	"github.com/dustin/go-humanize"
)

// Prettify is a pretty-printer.
type Prettify struct {
	// data is the data to be pretty-printed.
	data []byte

	// tab_size is the size of the tab.
	tab_size int

	// at is the current position of the lexer. (-1 means the lexer has not
	// an error yet.)
	at int

	// err is the error that occurred during the lexer.
	err error
}

// NewPrettify returns a new instance of Prettify with the given tab size.
//
// Parameters:
//   - tab_size: The size of the tab.
//
// Returns:
//   - *Prettify: A new instance of Prettify. Never returns nil.
//   - error: An error if the tab size is invalid.
func NewPrettify(tab_size int) (*Prettify, error) {
	if tab_size <= 0 {
		return nil, common.NewErrBadParam("tab_size", "must be positive")
	}

	return &Prettify{
		tab_size: tab_size,
		at:       -1,
		err:      nil,
	}, nil
}

// SetData sets the data to be pretty-printed.
//
// Parameters:
//   - data: The data to be pretty-printed.
//
// Returns:
//   - error: An error if the receiver is nil.
func (p *Prettify) SetData(data []byte) error {
	if p == nil {
		return common.ErrNilReceiver
	}

	p.data = data

	return nil
}

// SetError sets the error of the lexer.
//
// Parameters:
//   - err: The error to be set.
//
// Returns:
//   - error: An error if the receiver is nil, or if the error cannot be set.
//
// The error is set by directly assigning the error to the receiver's err field.
// The receiver's at field is set to the value of the error if the error is
// an instance of ErrLexing, otherwise the receiver's at field is cleared.
func (p *Prettify) SetError(err error) error {
	if p == nil {
		return common.ErrNilReceiver
	}

	p.err = err

	at, ok := ExtractAt(err)
	if !ok {
		p.at = -1
	} else {
		if at < 0 || at >= len(p.data) {
			return fmt.Errorf("at %d is out of range [%d, %d)", at, 0, len(p.data))
		}

		p.at = at
	}

	return nil
}

// PrettifyData prettifies the data by highlighting the position of the lexer error.
//
// The prettified data is returned as a byte slice. The byte slice is created by
// wrapping the original data in a string and inserting a caret (^) at the
// position of the lexer error. The position of the caret is determined by the
// value of the at field, which is the index of the character in the original
// data where the lexer error occurred.
//
// The prettified data is formatted as follows:
//
//	...<before>...
//	<faulty_line>
//	<caret>
//	...<after>...
//
// Where:
//   - <before> is the part of the original data before the lexer error.
//   - <faulty_line> is the line of the original data containing the lexer error.
//   - <caret> is a caret (^) indicating the position of the lexer error.
//   - <after> is the part of the original data after the lexer error.
//
// If the lexer error is at the beginning of the data, the caret is placed at the
// beginning of the line. If the lexer error is at the end of the data, the caret
// is placed at the end of the line.
//
// Returns:
//   - []byte: The prettified data.
//   - error: An error if the receiver is nil, or if the data cannot be prettified.
func (p Prettify) PrettifyData() ([]byte, error) {
	if p.at == -1 {
		return p.data, nil
	}

	at := p.at

	before_idx := -1

	for i := 0; i < at; i++ {
		if p.data[i] == '\n' {
			before_idx = i
		}
	}

	after_idx := -1

	for i := at; i < len(p.data) && after_idx == -1; i++ {
		if p.data[i] == '\n' {
			after_idx = i
		}
	}

	var before, faulty_line, after []byte

	if before_idx == -1 || before_idx == len(p.data)-1 {
		before_idx = 0
	} else {
		before = p.data[:before_idx]
		before_idx++
	}

	if after_idx == -1 || after_idx == len(p.data)-1 {
		after_idx = len(p.data)
	} else {
		after = p.data[after_idx+1:]
	}

	faulty_line = p.data[before_idx:after_idx]

	var buff bytes.Buffer

	if len(before) > 0 {
		buff.Write(before)
		buff.Write([]byte{'\n'})
	}

	buff.Write(faulty_line)
	buff.Write([]byte{'\n'})

	var diff int

	if len(before) == 0 {
		diff = at
	} else {
		diff = at - len(before) - 1
	}

	if diff < 0 {
		diff = 0
	}

	buff.Write(bytes.Repeat([]byte{' '}, diff))
	buff.Write([]byte{'^'})

	if len(after) > 0 {
		buff.Write([]byte{'\n'})
		buff.Write(after)
	}

	return buff.Bytes(), nil
}

// GetCoords returns the coordinates of the given position in the given data.
//
// Parameters:
//   - data: The data to get the coordinates from.
//   - pos: The position in the data to get the coordinates for.
//
// Returns:
//   - [2]int: The coordinates of the given position in the given data. The
//     coordinates are [x, y] where x is the column and y is the line.
//   - error: An error if the position is invalid, or if the data cannot be
//     processed.
func (p Prettify) GetCoords(data []byte, pos int) ([2]int, error) {
	if pos < 0 || pos > len(data) {
		return [2]int{}, common.NewErrBadParam("pos", fmt.Sprintf("must be in range [%d, %d), got %d", 0, len(data), pos))
	}

	data = data[:pos:pos]

	chars, err := gch.BytesToUtf8(data)
	if err != nil {
		return [2]int{}, err
	}

	err = gch.Normalize(&chars, p.tab_size)
	if err != nil {
		return [2]int{}, err
	}

	x, y := 0, 0

	for _, char := range chars {
		switch char {
		case '\n':
			y++
			x = 0
		default:
			x++
		}
	}

	return [2]int{x, y}, nil
}

// PrettifyLexError returns a prettified lexer error.
//
// Returns:
//   - []byte: The prettified lexer error.
//   - error: An error if the receiver is nil, or if the lexer error cannot be
//     prettified.
func (p Prettify) PrettifyLexError() ([]byte, error) {
	if p.at == -1 {
		return []byte(p.err.Error()), nil
	}

	coords, err := p.GetCoords(p.data, p.at)
	if err != nil {
		return nil, fmt.Errorf("failed to get coordinates: %w", err)
	}

	var buff bytes.Buffer

	fmt.Fprintf(&buff, "Lexer error at %s column of %s line:\n\t%v.\n",
		humanize.Ordinal(coords[0]+1), humanize.Ordinal(coords[1]+1), p.err,
	)

	return buff.Bytes(), nil
}
