package main

import (
	"github.com/PlayerR9/SlParser/ebnf"
	"github.com/PlayerR9/grammar"
)

func main() {
	ebnf.Parser.SetDebug(grammar.ShowNone)

	data := "Source = Source1 EOF ."

	_, err := ebnf.Parser.Parse([]byte(data))
	if err != nil {
		panic(err)
	}
}
