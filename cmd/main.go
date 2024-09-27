package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/PlayerR9/SlParser/cmd/internal"
	kdd "github.com/PlayerR9/SlParser/kdd"
	"github.com/PlayerR9/go-generator"
)

var (
	Logger *log.Logger
)

func init() {
	Logger = log.New(os.Stdout, "[Sl Parser]: ", log.LstdFlags)
}

func main() {
	dir, err := internal.ParseFlags()
	if err != nil {
		generator.PrintFlags()

		Logger.Fatalf("Error parsing flags: %v", err)
	}

	input_loc := *internal.InputFileFlag

	data, err := os.ReadFile(input_loc)
	if err != nil {
		Logger.Fatalf("Error reading file: %v", err)
	}

	parser := kdd.NewParser()
	parser.SetDebugger(Logger)
	parser.SetMode(kdd.ShowAll)

	root, err := parser.Full(data)
	if err != nil {
		Logger.Fatalf("Error parsing file: %v", err)
	}

	table, err := internal.InfoTableOf.Apply(root)
	if err != nil {
		Logger.Fatalf("Error creating info table: %v", err)
	}

	rules, err := internal.ExtractRules(table, root)
	if err != nil {
		Logger.Fatalf("Error extracting rules: %v", err)
	}

	infos := internal.LinearizeTable(table)

	err = GenerateTokens(dir, infos, rules)
	if err != nil {
		Logger.Fatalf("Error generating tokens: %v", err)
	}

	err = GenerateLexer(dir)
	if err != nil {
		Logger.Fatalf("Error generating lexer: %v", err)
	}

	err = GenerateNode(dir)
	if err != nil {
		Logger.Fatalf("Error generating node: %v", err)
	}

	err = GenerateAst(dir, table)
	if err != nil {
		Logger.Fatalf("Error generating ast: %v", err)
	}

	err = GenerateParsing(dir)
	if err != nil {
		Logger.Fatalf("Error generating parsing: %v", err)
	}

	err = GenerateGen(dir)
	if err != nil {
		Logger.Fatalf("Error generating gen: %v", err)
	}

	// cmd := exec.Command("go", "generate", "./...")
	// err = cmd.Run()
	// if err != nil {
	// 	panic(err)
	// }

	Logger.Println("Successfully generated parser. Make sure to run go generate ./...")
}

// WARNING: The CheckEofExists() function makes a fundamentally wrong and/or
// too restrictive assumption about the EOF symbol.
//
// That's because a valid parsing may not contain an EOF symbol. Yet, this derives
// from how the github.com/PlayerR9/SlParser/parser package is implemented.
//
// If that ever changes, this function will need to be updated.
func GenerateTokens(dir string, tk_symbols []*internal.Info, rules []*internal.Rule) error {
	ok := internal.CheckEofExists(tk_symbols)
	if !ok {
		return errors.New("missing EOF")
	}

	gd, err := internal.NewTokenGen(tk_symbols)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		gd.Rules = append(gd.Rules, rule.String())
	}

	loc := filepath.Join(dir, "token.go")

	gen, err := internal.TokenGenerator.GenerateWithLoc(loc, gd)
	if err != nil {
		return err
	}

	err = gen.WriteFile()
	if err != nil {
		return err
	}

	return nil
}

func GenerateLexer(dir string) error {
	lexer_gen := internal.NewLexerGen()

	loc := filepath.Join(dir, "lexer.go")

	gen, err := internal.LexerGenerator.GenerateWithLoc(loc, lexer_gen)
	if err != nil {
		return err
	}

	err = gen.WriteFile()
	if err != nil {
		return err
	}

	return nil
}

func GenerateNode(dir string) error {
	node_gen := internal.NewNodeGen()

	loc := filepath.Join(dir, "node.go")

	gen, err := internal.NodeGenerator.GenerateWithLoc(loc, node_gen)
	if err != nil {
		return err
	}

	err = gen.WriteFile()
	if err != nil {
		return err
	}

	return nil
}

func GenerateAst(dir string, table map[*kdd.Node]*internal.Info) error {
	gen := internal.NewASTGen(table)

	loc := filepath.Join(dir, "ast.go")

	ast_gen, err := internal.ASTGenerator.GenerateWithLoc(loc, gen)
	if err != nil {
		return err
	}

	err = ast_gen.WriteFile()
	if err != nil {
		return err
	}

	return nil
}

func GenerateParsing(dir string) error {
	gen := internal.NewParsingGen()

	loc := filepath.Join(dir, "parsing.go")

	data, err := internal.ParsingGenerator.GenerateWithLoc(loc, gen)
	if err != nil {
		return err
	}

	err = data.WriteFile()
	return err
}

func GenerateGen(dir string) error {
	gen := internal.NewGenGen()

	loc := filepath.Join(dir, "generate.go")

	data, err := internal.GenGenerator.GenerateWithLoc(loc, gen)
	if err != nil {
		return err
	}

	err = data.WriteFile()
	return err
}
