package main

import (
	"fmt"
	"os"

	dspl "github.com/PlayerR9/SlParser/display"
	gr "github.com/PlayerR9/SlParser/grammar"
	lxr "github.com/PlayerR9/SlParser/lexer"
	pkg "github.com/PlayerR9/SlParser/test/parsing"
	tr "github.com/PlayerR9/tree/tree"
)

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tokens := FullLex(data)

	FullParse(tokens)

	// [0, 1, 4, 9, 16, 25, 36, 49, 64, 81]
}

func FullLex(data []byte) []*gr.Token[pkg.TokenType] {
	is := lxr.NewStream().FromBytes(data)

	pkg.Lexer.SetInputStream(is)

	err := pkg.Lexer.Lex()

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
			x, y := dspl.GetCoords(data, lex_err.Pos)

			fmt.Printf("Error at (%d column, %d line)\n", x, y)

			str := dspl.Display(data, lex_err.Pos)
			fmt.Println(string(str))

			fmt.Println("Hints:")
			for _, s := range lex_err.Suggestions {
				fmt.Println("\t", s)
			}

			os.Exit(2 + int(lex_err.Code))
		}
	}

	return tokens
}

func FullParse(tokens []*gr.Token[pkg.TokenType]) *gr.Token[pkg.TokenType] {
	pkg.Parser.SetTokens(tokens)

	err := pkg.Parser.Parse()

	forest := pkg.Parser.Forest()

	fmt.Println("Here's the list of forest:")

	for _, f := range forest {
		tree := tr.NewTree(f)

		fmt.Println(tree.String())
		fmt.Println()
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
	}

	if len(forest) != 1 {
		fmt.Println(fmt.Errorf("expected one forest, got %d instead", len(forest)))
	}

	return forest[0]
}
