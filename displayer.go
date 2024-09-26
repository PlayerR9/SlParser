package SlParser

import (
	"bytes"
	"fmt"
	"io"

	lxr "github.com/PlayerR9/SlParser/lexer"
	fch "github.com/PlayerR9/go-commons/Formatting/runes"
	gcby "github.com/PlayerR9/go-commons/bytes"
	gers "github.com/PlayerR9/go-errors"
	gerr "github.com/PlayerR9/go-errors/error"
	"github.com/dustin/go-humanize"
)

type Displayer struct {
	err  *gerr.Err
	data []byte
	pos  int
	x    int
	y    int
}

func NewDisplayer(err *gerr.Err, data []byte, pos int) *Displayer {
	var x, y int

	for i := 0; i < pos-1; i++ {
		if data[i] == '\n' {
			x = 0
			y++
		} else {
			x++
		}
	}

	return &Displayer{
		err:  err,
		data: data,
		pos:  pos,
		x:    x,
		y:    y,
	}
}

func (d *Displayer) init_source() {
	if d == nil {
		return
	}

	var before, faulty_line, after []byte

	last_idx := ReverseSearch(d.data, d.pos, []byte{'\n'})
	if last_idx < 0 {
		faulty_line = make([]byte, len(d.data[:d.pos]))
		copy(faulty_line, d.data[:d.pos])

		last_idx = 0
	} else {
		before = make([]byte, last_idx-1)
		copy(before, d.data[:last_idx-1])

		faulty_line = make([]byte, d.pos-last_idx+1)
		copy(faulty_line, d.data[last_idx:d.pos+1])
	}

	first_idx := ForwardSearch(d.data, d.pos, []byte("\n"))
	if first_idx < 0 {
		faulty_line = append(faulty_line, d.data[d.pos:]...)
	} else {
		after = make([]byte, len(d.data)-first_idx)
		copy(after, d.data[first_idx:])

		faulty_line = append(faulty_line, d.data[d.pos:first_idx]...)
	}

	var builder bytes.Buffer

	if before != nil {
		builder.Write(before)
		builder.WriteRune('\n')
	}

	builder.Write(faulty_line)
	builder.WriteRune('\n')

	for i := 0; i < d.pos-last_idx-1; i++ {
		if faulty_line[i] == '\t' {
			builder.WriteRune('\t')
		} else {
			builder.WriteRune(' ')
		}
	}

	builder.WriteRune('^')

	if after != nil {
		if !bytes.HasPrefix(after, []byte("\n")) {
			builder.WriteRune('\n')
		}

		builder.Write(after)
	}

	d.data = builder.Bytes()
}

func (d *Displayer) write_source(w io.Writer) error {
	if d == nil || w == nil {
		return nil
	}

	d.init_source()

	style := fch.NewBoxStyle(fch.BtNormal, true, [4]int{0, 1, 0, 1})

	var table fch.RuneTable

	lines := bytes.Split(d.data, []byte("\n"))

	err := table.FromBytes(lines)
	gers.AssertErr(err, "fch.FromBytes(lines)")

	err = style.Apply(&table)
	if err != nil {
		return err
	}

	err = gcby.Write(w, table.Byte())
	if err != nil {
		return err
	}

	return nil
}

func (d Displayer) write_error(w io.Writer) error {
	var builder bytes.Buffer

	builder.WriteString("\n\nError at ")
	builder.WriteString(humanize.Ordinal(d.x + 1))
	builder.WriteString(" column of the ")
	builder.WriteString(humanize.Ordinal(d.y + 1))
	builder.WriteString(" line:\n\t")

	err := gers.DisplayError(&builder, d.err)
	if err != nil {
		return err
	}

	builder.WriteRune('\n')

	err = gcby.Write(w, builder.Bytes())
	if err != nil {
		return err
	}

	return err
}

///////////////////////////////////

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

	e, ok := err.(*gerr.Err)
	if !ok {
		err := gcby.Write(w, data)
		return 0, err
	}

	pos, err := gers.Value[lxr.ErrorCode, int](e, "pos")
	if err != nil {
		return 0, err
	}

	d := NewDisplayer(e, data, pos)
	err = d.write_source(w)
	if err != nil {
		return 0, fmt.Errorf("could not write source: %w", err)
	}

	err = d.write_error(w)
	if err != nil {
		return 0, fmt.Errorf("could not write error: %w", err)
	}

	return e.Code.Int(), nil
}
