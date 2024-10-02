package kdd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/PlayerR9/go-errors/assert"
)

func TestLexer(t *testing.T) {
	assert.NotNil(Lexer, "Lexer")

	var buff bytes.Buffer

	buff.WriteString("source : rule EOF ;")

	Lexer.SetInputStream(&buff)
	err := Lexer.Lex()
	tokens := Lexer.Tokens()

	for _, tk := range tokens {
		fmt.Println(tk.String())
	}

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
