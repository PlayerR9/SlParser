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
		fmt.Println(err)
		os.Exit(1)
	}

	is := lxr.NewStream().FromBytes(data)

	pkg.Lexer.SetInputStream(is)
	pkg.Lexer.Lex()
	tokens := pkg.Lexer.Tokens()
	fmt.Println(tokens)

	if err := pkg.Lexer.Error(); err != nil {
		fmt.Println(err.Error())

		fmt.Println("Hints:")
		for _, s := range err.Suggestions {
			fmt.Println("\t", s)
		}

		os.Exit(2 + int(err.Code))
	}

	// [0, 1, 4, 9, 16, 25, 36, 49, 64, 81]
}
