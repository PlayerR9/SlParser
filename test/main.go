package main

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"os"

	sl "github.com/PlayerR9/SlParser"
	pkg "github.com/PlayerR9/SlParser/test/parsing"
	tr "github.com/PlayerR9/tree/tree"
)

func main() {
	err := pkg.PrintItemSet(os.Stdout)
	if err != nil {
		fmt.Println(err)
	}

	data, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tokens, err := sl.Lex(pkg.Lexer, data)

	// DEBUG: Print the list of tokens.
	fmt.Println("[DEBUG]: Here's the list of tokens:")
	for _, tk := range tokens {
		fmt.Println("\t", tk.String())
	}
	fmt.Println()

	exit_code, err := sl.DisplayErr(os.Stdout, data, err)
	if err != nil {
		panic(err)
	} else if exit_code != 0 {
		os.Exit(exit_code + 3)
	}

	forest, err := sl.Parse(pkg.Parser, tokens)

	// DEBUG: Print the forest.
	fmt.Println("[DEBUG]: Here is the forest:")

	for _, tree := range forest {
		fmt.Println(tree.String())
		fmt.Println()
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	} else if len(forest) != 1 {
		fmt.Println(fmt.Errorf("expected one forest, got %d instead", len(forest)))
		os.Exit(2)
	}

	node, err := pkg.AstMaker.Convert(forest[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	source_tree := tr.NewTree(node)
	fmt.Println(source_tree.String())

	// [0, 1, 4, 9, 16, 25, 36, 49, 64, 81]
}

func DebugPrint(w io.Writer, title string, elems iter.Seq[string]) error {
	var buff bytes.Buffer

	if title != "" {
		buff.WriteString("[DEBUG]: ")
		buff.WriteString(title)
	}

	if elems != nil {
		for elem := range elems {
			buff.WriteString(elem)
			buff.WriteRune('\n')
		}
	}

	data := buff.Bytes()

	if len(data) == 0 {
		return nil
	}

	if w == nil {
		return io.ErrShortWrite
	}

	err := sl.Write(w, data)
	return err
}
