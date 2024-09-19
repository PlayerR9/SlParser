package debug

import (
	"bytes"
	"io"
	"iter"
	"log"
)

// write is a helper function that writes a data to a writer.
//
// Does nothing if the data is empty.
//
// Parameters:
//   - w: the writer to write to.
//   - data: the data to write.
//
// Returns:
//   - error: the error that occurred.
//
// Errors:
//   - io.ErrShortWrite: if the writer is nil or the data could not be written fully.
//   - any other error returned by the writer.
func write(w io.Writer, data []byte) error {
	if len(data) == 0 {
		return nil
	} else if w == nil {
		return io.ErrShortWrite
	}

	n, err := w.Write(data)
	if err != nil {
		return err
	} else if n != len(data) {
		return io.ErrShortWrite
	}

	return nil
}

// DebugPrint is a helper function that prints a string to a writer.
//
// Parameters:
//   - w: the writer to write to.
//   - title: the title of the debug message. If empty, no title is printed.
//   - lines: the lines of the debug message. If nil, no lines are printed.
//
// Returns:
//   - error: the error that occurred.
//
// Errors:
//   - io.ErrShortWrite: if the writer is nil or the lines could not be written fully.
//   - any other error returned by the writer.
func DebugPrint(w io.Writer, title string, lines iter.Seq[string]) error {
	var buff bytes.Buffer

	if title != "" {
		buff.WriteString(title)
		buff.WriteRune('\n')
		buff.WriteRune('\n')
	}

	if lines != nil {
		for line := range lines {
			buff.WriteString(line)
			buff.WriteRune('\n')
		}

		buff.WriteRune('\n')
	}

	err := write(w, buff.Bytes())
	return err
}

// LogPrint acts as DebugPrint but with a logger instead.
//
// Parameters:
//   - l: the logger to write to.
//   - title: the title of the debug message. If empty, no title is printed.
//   - lines: the lines of the debug message. If nil, no lines are printed.
//
// Returns:
//   - error: the error that occurred.
//
// Errors:
//   - io.ErrShortWrite: if the writer is nil or the lines could not be written fully.
//   - any other error returned by the logger.
func LogPrint(l *log.Logger, title string, lines iter.Seq[string]) error {
	var buff bytes.Buffer

	if title != "" {
		l.Println(title)
		_, _ = buff.WriteRune('\n')
	}

	if lines != nil {
		for line := range lines {
			_, _ = buff.WriteString(line)
			_, _ = buff.WriteRune('\n')
		}

		_, _ = buff.WriteRune('\n')
	}

	w := l.Writer()
	err := write(w, buff.Bytes())
	return err
}
