package highlighter

import (
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"

	gr "github.com/PlayerR9/SlParser/grammar"
	sllx "github.com/PlayerR9/SlParser/lexer"
	tr "github.com/PlayerR9/mygo-lib/CustomData/tree"
	"github.com/PlayerR9/mygo-lib/colors"
	gch "github.com/PlayerR9/mygo-lib/runes"
)

type Highlighter struct {
	table        map[string]*colors.Style
	defaultStyle *colors.Style
}

func (h Highlighter) Highlight(w io.Writer, data []rune, tokens []*tr.Node) error {
	if len(tokens) == 0 {
		var b []byte

		_ = Encode(&b, data)

		n, err := w.Write(b)
		if err != nil {
			return err
		} else if n != len(b) {
			return io.ErrShortWrite
		}

		return nil
	}

	infos := make([]*gr.TokenData, 0, len(tokens))

	for i, token := range tokens {
		tkd, err := gr.Get[*gr.TokenData](token)
		if err != nil {
			return fmt.Errorf("at index %d: %w", i, err)
		}

		infos = append(infos, tkd)
	}

	var pos int

	for _, info := range infos {
		if info.Pos > len(data) {
			return fmt.Errorf("invalid position: %d", info.Pos)
		}

		if pos < info.Pos {
			var b []byte

			_ = Encode(&b, data[pos:info.Pos])

			n, err := w.Write(b)
			if err != nil {
				return err
			} else if n != len(b) {
				return io.ErrShortWrite
			}
		}

		pos = info.Pos + utf8.RuneCountInString(info.Data)

		type_ := info.Type

		style, ok := h.table[type_]
		if !ok {
			style = h.defaultStyle
		}

		b := []byte(style.String() + info.Data)

		n, err := w.Write(b)
		if err != nil {
			return err
		} else if n != len(b) {
			return io.ErrShortWrite
		}
	}

	if pos < len(data) {
		var b []byte

		_ = Encode(&b, data[pos:])

		n, err := w.Write(b)
		if err != nil {
			return err
		} else if n != len(b) {
			return io.ErrShortWrite
		}
	}

	return nil
}

func Highlight(h Highlighter, lexer sllx.Lexer, data []byte) ([]byte, error) {
	defer lexer.Reset()

	_, _ = lexer.Write(data)

	itr, err := lexer.Lex()
	if err != nil {
		return nil, err
	}
	defer itr.Stop()

	var tokens []*tr.Node

	for {
		result, err := itr.Next()
		if err == sllx.ErrExhausted {
			break
		} else if err != nil {
			return nil, err
		}

		err = result.GetError()
		if err == nil {
			tokens = result.Tokens()
			break
		}
	}

	var buff bytes.Buffer

	chars, err := gch.BytesToUtf8(data)
	if err != nil {
		return nil, err
	}

	err = h.Highlight(&buff, chars, tokens)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
