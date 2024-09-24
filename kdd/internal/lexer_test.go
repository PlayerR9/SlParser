package internal

import (
	"bytes"
	"fmt"
	"testing"

	gers "github.com/PlayerR9/go-errors"
)

func TestLexer(t *testing.T) {
	gers.AssertNotNil(Lexer, "Lexer")

	var buff bytes.Buffer

	buff.WriteString("source : rule EOF ;")

	Lexer.SetInputStream(&buff)
	err := Lexer.Lex()
	tokens := Lexer.Tokens()

	for _, tk := range tokens {
		fmt.Println(tk.String())
	}

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	t.Errorf("Not implemented yet")
}
