package main

import (
	"errors"
	"log"
	"os"

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
	err := internal.ParseFlags()
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

	tk_symbols := internal.ExtractSymbols(root)

	rules, err := internal.ExtractRules(root)
	if err != nil {
		Logger.Fatalf("Error extracting rules: %v", err)
	}

	err = GenerateTokens(tk_symbols, rules)
	if err != nil {
		Logger.Fatalf("Error generating tokens: %v", err)
	}

	err = GenerateLexer()
	if err != nil {
		Logger.Fatalf("Error generating lexer: %v", err)
	}

	err = GenerateNode()
	if err != nil {
		Logger.Fatalf("Error generating node: %v", err)
	}

	err = GenerateAst(tk_symbols)
	if err != nil {
		Logger.Fatalf("Error generating ast: %v", err)
	}

	err = GenerateError()
	if err != nil {
		Logger.Fatalf("Error generating errors: %v", err)
	}

	err = GenerateParsing()
	if err != nil {
		Logger.Fatalf("Error generating parsing: %v", err)
	}

	// cmd := exec.Command("go", "generate", "./...")
	// err = cmd.Run()
	// if err != nil {
	// 	panic(err)
	// }

	Logger.Println("Successfully generated parser. Make sure to run go generate ./...")
}

func GenerateTokens(tk_symbols []*kdd.Node, rules []*internal.Rule) error {
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

	gen, err := internal.TokenGenerator.Generate(internal.OutputLocFlag, "lexer", gd)
	if err != nil {
		return err
	}

	gen.ModifyPrefixPath("token_", "internal")

	err = gen.WriteFile()
	if err != nil {
		return err
	}

	return nil
}

func GenerateLexer() error {
	gen := internal.NewLexerGen()

	lexer_gen, err := internal.LexerGenerator.Generate(internal.OutputLocFlag, "lexer", gen)
	if err != nil {
		return err
	}

	lexer_gen.ModifyPrefixPath("lexer_", "internal")

	err = lexer_gen.WriteFile()
	if err != nil {
		return err
	}

	return nil
}

func GenerateNode() error {
	gen := internal.NewNodeGen()

	node_gen, err := internal.NodeGenerator.Generate(internal.OutputLocFlag, "node", gen)
	if err != nil {
		return err
	}

	node_gen.ModifyPrefixPath("node_")

	err = node_gen.WriteFile()
	if err != nil {
		return err
	}

	return nil
}

func GenerateAst(tk_symbols []*kdd.Node) error {
	gen, err := internal.NewASTGen(tk_symbols)
	if err != nil {
		return err
	}

	ast_gen, err := internal.ASTGenerator.Generate(internal.OutputLocFlag, "ast", gen)
	if err != nil {
		return err
	}

	ast_gen.ModifyPrefixPath("ast_")

	err = ast_gen.WriteFile()
	if err != nil {
		return err
	}

	return nil
}

func GenerateParsing() error {
	gen := internal.NewParsingGen()

	data, err := internal.ParsingGenerator.Generate(internal.OutputLocFlag, "parsing", gen)
	if err != nil {
		return err
	}

	data.ModifyPrefixPath("parsing_")

	err = data.WriteFile()
	return err
}

func GenerateError() error {
	gen := internal.NewErrorGen()

	data, err := internal.ErrorGenerator.Generate(internal.OutputLocFlag, "error", gen)
	if err != nil {
		return err
	}

	data.ModifyPrefixPath("errors_")

	err = data.WriteFile()
	return err
}
