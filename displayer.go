package SlParser

import (
	"bytes"
	"errors"
	"io"

	lxr "github.com/PlayerR9/SlParser/lexer"
	gcers "github.com/PlayerR9/go-commons/errors"
)

// Display is a function that displays the given data.
//
// Parameters:
//   - data: The data.
//   - pos: The position.
//
// Returns:
//   - []byte: The displayed data.
func Display(data []byte, pos int) []byte {
	var before, faulty_line, after []byte

	last_idx := ReverseSearch(data, pos, []byte{'\n'})
	if last_idx < 0 {
		faulty_line = make([]byte, len(data[:pos]))
		copy(faulty_line, data[:pos])

		last_idx = 0
	} else {
		before = make([]byte, last_idx)
		copy(before, data[:last_idx])

		last_idx++

		faulty_line = make([]byte, pos-last_idx)
		copy(faulty_line, data[last_idx:pos])
	}

	first_idx := ForwardSearch(data, pos, []byte{'\n'})
	if first_idx < 0 {
		faulty_line = append(faulty_line, data[pos:]...)
	} else {
		after = make([]byte, len(data)-first_idx)
		copy(after, data[first_idx:])

		faulty_line = append(faulty_line, data[pos:first_idx]...)
	}

	var builder bytes.Buffer

	if before != nil {
		builder.Write(before)
		builder.WriteRune('\n')
	}

	builder.Write(faulty_line)
	builder.WriteRune('\n')

	for i := 0; i < pos-last_idx-1; i++ {
		if faulty_line[i] == '\t' {
			builder.WriteRune('\t')
		} else {
			builder.WriteRune(' ')
		}
	}

	builder.WriteRune('^')

	if after != nil {
		builder.WriteRune('\n')
		builder.Write(after)
	}

	return builder.Bytes()
}

func GetCoords(data []byte, pos int) (int, int) {
	var x, y int

	for i := 0; i < pos-1; i++ {
		if data[i] == '\n' {
			x = 0
			y++
		} else {
			x++
		}
	}

	return x + 1, y + 1
}

// write is a helper function that writes the data to the writer.
//
// Parameters:
//   - w: the writer.
//   - data: the data.
//
// Returns:
//   - error: if an error occurred.
func write(w io.Writer, data []byte) error {
	n, err := w.Write(data)
	if err != nil {
		return err
	} else if n != len(data) {
		return io.ErrShortWrite
	}

	return nil
}

// DisplayErr is a function that displays the error in the given writer
// if it is not nil.
//
// Parameters:
//   - w: the writer.
//   - data: the data.
//   - err: the error.
//
// Returns:
//   - int: the exit code of the error.
//   - error: if the writer failed to write the data.
//
// The exit code is 0 if an error occurred or 'err' is nil.
func DisplayErr(w io.Writer, data []byte, err error) (int, error) {
	if w == nil || err == nil {
		return 0, nil
	}

	data_err := []byte(err.Error())
	var lexing_err *lxr.Err

	ok := errors.As(err, &lexing_err)
	if !ok {
		err := write(w, data_err)
		return 0, err
	}

	display_data := Display(data, lexing_err.Pos)

	err = write(w, display_data)
	if err != nil {
		return 0, err
	}

	x_coord, y_coord := GetCoords(data, lexing_err.Pos)

	var builder bytes.Buffer

	builder.WriteString("\n\nError at ")
	builder.WriteString(gcers.GetOrdinalSuffix(x_coord))
	builder.WriteString(" column of the ")
	builder.WriteString(gcers.GetOrdinalSuffix(y_coord))
	builder.WriteString(" line:\n\t")
	builder.WriteString(lexing_err.Error())
	builder.WriteString(".\n\nHints:\n")

	for _, suggestion := range lexing_err.Suggestions {
		builder.WriteRune('\t')
		builder.WriteString(suggestion)
		builder.WriteRune('\n')
	}

	err = write(w, builder.Bytes())
	if err != nil {
		return 0, err
	}

	return int(lexing_err.Code), nil
}
