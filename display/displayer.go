package display

import (
	"bytes"
)

func write_arrow(faulty_line []byte, pos int) []byte {

	var builder bytes.Buffer

	for i := 0; i < pos; i++ {
		if faulty_line[i] == '\t' {
			builder.WriteRune('\t')
		} else {
			builder.WriteRune(' ')
		}
	}

	builder.WriteRune('^')

	return builder.Bytes()
}

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
	builder.Write(write_arrow(faulty_line, pos-last_idx-1))

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
