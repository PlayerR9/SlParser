package highlighter

import (
	"bytes"
	"io"
	"unicode/utf8"

	slgr "github.com/PlayerR9/SlParser/grammar"
	sllx "github.com/PlayerR9/SlParser/lexer"
	"github.com/PlayerR9/mygo-lib/colors"
	gch "github.com/PlayerR9/mygo-lib/runes"
)

type Highlighter struct {
	table        map[string]*colors.Style
	defaultStyle *colors.Style
}

func (h Highlighter) Highlight(w io.Writer, data []rune, tokens []*slgr.Token) error {
	var pos int

	for _, token := range tokens {
		if pos < token.Pos {
			var b []byte

			_ = Encode(&b, data[pos:token.Pos])

			n, err := w.Write(b)
			if err != nil {
				return err
			} else if n != len(b) {
				return io.ErrShortWrite
			}
		}

		pos = token.Pos + utf8.RuneCountInString(token.Data)

		type_ := token.Type

		style, ok := h.table[type_]
		if !ok {
			style = h.defaultStyle
		}

		b := []byte(style.String() + token.Data)

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

	var tokens []*slgr.Token

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
