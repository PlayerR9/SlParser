package lexer

import (
	"bytes"
)

func Display(data []byte, pos int) []byte {
	data_before := make([]byte, pos+1)
	copy(data_before, data[:pos+1])

	data_after := make([]byte, len(data)-pos-2)
	copy(data_after, data[pos+2:])

	var builder bytes.Buffer

	builder.Write(data_before)
	builder.WriteRune('\n')
	builder.WriteRune('^')
	builder.WriteRune('\n')
	builder.Write(data_after)

	return builder.Bytes()
}
