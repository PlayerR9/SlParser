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

	defer pkg.Parser.Reset()

	pkg.Parser.SetTokens(tokens)

	var node *pkg.Node
	var last_error error

	for node == nil {
		ap, err := pkg.Parser.Parse()
		if err != nil {
			if last_error == nil {
				last_error = err
			}

			break
		} else if ap == nil {
			break
		}

		forest := ap.Forest()

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

		if len(forest) != 1 {
			last_error = fmt.Errorf("expected one forest, got %d instead", len(forest))

			continue
		}

		node, err = pkg.AstMaker.Convert(forest[0])
		if err != nil {
			last_error = err

			continue
		}
	}

	if node == nil {
		fmt.Println(last_error.Error())
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
