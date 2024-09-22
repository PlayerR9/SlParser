package main

import (
	"fmt"
	"log"
	"os"

	sl "github.com/PlayerR9/SlParser"
	ast "github.com/PlayerR9/SlParser/ast"
	lxr "github.com/PlayerR9/SlParser/lexer"
	prx "github.com/PlayerR9/SlParser/parser"
	pkg "github.com/PlayerR9/SlParser/test/parsing"
	dbp "github.com/PlayerR9/go-debug/debug"
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
	err := ParseCmd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// [0, 1, 4, 9, 16, 25, 36, 49, 64, 81]
}

func ParseCmd() error {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		return err
	}

	p := &Parsing{
		debug_mode: ShowAll,
		lexer:      pkg.Lexer,
		parser:     pkg.Parser,
		ast:        pkg.AstMaker,
	}

	_, err = p.Full(data)
	if err != nil {
		return err
	}

	return nil
}

type Parsing struct {
	debug_mode DebugMode

	lexer  *lxr.Lexer[pkg.TokenType]
	parser *prx.Parser[pkg.TokenType]
	ast    *ast.AstMaker[*pkg.Node, pkg.TokenType]
}

func (p *Parsing) SetMode(mode DebugMode) {
	if p == nil {
		return
	}

	p.debug_mode = mode
}

func (p Parsing) Full(data []byte) (*pkg.Node, error) {
	if p.debug_mode&ShowItemSet != 0 {
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

	tokens, err := sl.Lex(pkg.Lexer, data)

	// DEBUG: Print the list of tokens.
	if p.debug_mode&ShowTokens != 0 {
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
		return nil, err
	}

	defer pkg.Parser.Reset()

	pkg.Parser.SetTokens(tokens)

	var node *pkg.Node
	var last_error error

	for node == nil {
		forest, err := pkg.Parser.Parse()
		if err != nil {
			if last_error == nil {
				last_error = err
			}

			break
		} else if len(forest) == 0 {
			break
		}

		// DEBUG: Print the forest.
		if p.debug_mode&ShowForest != 0 {
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
		return nil, last_error
	}

	if p.debug_mode&ShowAST != 0 {
		err := dbp.LogPrint(Debugger, "Here's the AST:", func(yield func(string) bool) {
			_ = yield(tr.NewTree(node).String())
		})
		if err != nil {
			panic(err)
		}
	}

	return node, nil
}
