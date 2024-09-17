package main

import (
	"fmt"
	"os"

	lxr "github.com/PlayerR9/SlParser/lexer"
	pkg "github.com/PlayerR9/SlParser/test/parsing"
)

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	is := lxr.NewStream().FromBytes(data)

	pkg.Lexer.SetInputStream(is)

	err = pkg.Lexer.Lex()
	tokens := pkg.Lexer.Tokens()
	fmt.Println(tokens)

	if err != nil {
		panic(err)
	}

	// [0, 1, 4, 9, 16, 25, 36, 49, 64, 81]
}
