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

	err = pkg.Lexer.Lex()

	tokens := pkg.Lexer.Tokens()

	fmt.Println("Here's the list of tokens:")
	for _, tk := range tokens {
		fmt.Println("\t", tk.String())
	}
	fmt.Println()

	if err != nil {
		fmt.Println(err.Error())

		lex_err, ok := err.(*lxr.Err)
		if ok {
			str := lxr.Display(data, lex_err.Pos)
			fmt.Println(string(str))

			fmt.Println("Hints:")
			for _, s := range lex_err.Suggestions {
				fmt.Println("\t", s)
			}

			os.Exit(2 + int(lex_err.Code))
		}
	}

	// [0, 1, 4, 9, 16, 25, 36, 49, 64, 81]
}
