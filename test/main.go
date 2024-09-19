package main

import (
	"fmt"
	"log"
	"os"

	sl "github.com/PlayerR9/SlParser"
	pkg "github.com/PlayerR9/SlParser/test/parsing"
	dbp "github.com/PlayerR9/SlParser/util/go-debug/debug"
	tr "github.com/PlayerR9/tree/tree"
)

var (
	Debugger *log.Logger
)

func init() {
	Debugger = log.New(os.Stdout, "[DEBUG]: ", log.LstdFlags)
}

type DebugMode int

const (
	ShowNone    DebugMode = 0
	ShowItemSet DebugMode = 1
	ShowTokens  DebugMode = 2
	ShowForest  DebugMode = 4
	ShowAST     DebugMode = 8
	ShowAll     DebugMode = ShowItemSet | ShowTokens | ShowForest | ShowAST
)

func main() {
	var debugmode DebugMode = ShowAll

	if debugmode&ShowItemSet != 0 {
		err := dbp.LogPrint(Debugger, "Here's the item set:", func(yield func(string) bool) {
			lines := pkg.PrintItemSet()

			for _, line := range lines {
				if !yield(line) {
					break
				}
			}
		})
		if err != nil {
			panic(err)
		}
	}

	data, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tokens, err := sl.Lex(pkg.Lexer, data)

	// DEBUG: Print the list of tokens.
	if debugmode&ShowTokens != 0 {
		err := dbp.LogPrint(Debugger, "Here's the list of tokens:", func(yield func(string) bool) {
			for _, tk := range tokens {
				if !yield(tk.String()) {
					return
				}
			}
		})
		if err != nil {
			panic(err)
		}
	}

	exit_code, err := sl.DisplayErr(os.Stdout, data, err)
	if err != nil {
		panic(err)
	} else if exit_code != 0 {
		os.Exit(exit_code + 3)
	}

	forest, err := sl.Parse(pkg.Parser, tokens)

	// DEBUG: Print the forest.
	if debugmode&ShowForest != 0 {
		err := dbp.LogPrint(Debugger, "Here's the forest:", func(yield func(string) bool) {
			for _, tree := range forest {
				if !yield(tree.String()) {
					return
				}
			}
		})
		if err != nil {
			panic(err)
		}
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

	if debugmode&ShowAST != 0 {
		err := dbp.LogPrint(Debugger, "Here's the AST:", func(yield func(string) bool) {
			_ = yield(tr.NewTree(node).String())
		})
		if err != nil {
			panic(err)
		}
	}

	// [0, 1, 4, 9, 16, 25, 36, 49, 64, 81]
}

func ParseCmd() {

}
