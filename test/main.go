package main

import (
	"github.com/PlayerR9/SLParser/pkg"
	"github.com/PlayerR9/grammar"
)

func main() {
	pkg.Parser.SetDebug(grammar.ShowNone)

	data := "Source = Source1 EOF ."

	_, err := pkg.Parser.Parse([]byte(data))
	if err != nil {
		panic(err)
	}
}
