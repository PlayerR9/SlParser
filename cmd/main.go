package main

import (
	"bytes"
	"errors"
	"log"
	"os"

	"github.com/PlayerR9/SlParser/cmd/internal"
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

	lines := bytes.Split(data, []byte{'\n'})

	var tokens []*internal.Token

	var rules []*internal.Rule

	for i, line := range lines {
		if len(line) == 0 {
			continue
		}

		fields := bytes.Fields(line)
		if len(fields) == 0 {
			continue
		}

		if len(fields) <= 3 {
			Logger.Fatalf("invalid line at %d: %q", i, line)
		}

		if !bytes.Equal(fields[1], []byte(":")) {
			Logger.Fatalf("missing colon at %d: %q", i, line)
		}

		if !bytes.Equal(fields[len(fields)-1], []byte(";")) {
			Logger.Fatalf("missing semicolon at %d: %q", i, line)
		}

		lhs, err := internal.MakeToken(fields[0])
		if err != nil {
			Logger.Fatalf("invalid lhs at %d: %v", i, err)
		}

		tokens = append(tokens, lhs)

		var rhss []*internal.Token

		for j := 2; j < len(fields)-1; j++ {
			tk, err := internal.MakeToken(fields[j])
			if err != nil {
				Logger.Fatalf("invalid token at %d of line %d: %v", j, i, err)
			}

			rhss = append(rhss, tk)
			tokens = append(tokens, tk)
		}

		rule := internal.NewRule(lhs, rhss)
		rules = append(rules, rule)
	}

	tk_symbols, err := internal.TokenSymbols(tokens)
	if err != nil {
		Logger.Fatalf("Error parsing tokens: %v", err)
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

	// cmd := exec.Command("go", "generate", "./...")
	// err = cmd.Run()
	// if err != nil {
	// 	panic(err)
	// }

	Logger.Println("Successfully generated parser. Make sure to run go generate ./...")
}

func GenerateTokens(tk_symbols []*internal.Token, rules []*internal.Rule) error {
	ok := internal.CheckEofExists(tk_symbols)
	if !ok {
		return errors.New("missing EOF")
	}

	all_symbols := internal.ExtractSymbols(tk_symbols)

	last_terminal := internal.FindLastTerminal(tk_symbols)
	if last_terminal == nil {
		return errors.New("missing terminal")
	}

	gd := &internal.TokenGen{
		Symbols:      all_symbols,
		LastTerminal: last_terminal.String(),
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
	gd := &internal.LexerGen{}

	gen, err := internal.LexerGenerator.Generate(internal.OutputLocFlag, "lexer", gd)
	if err != nil {
		return err
	}

	gen.ModifyPrefixPath("lexer_", "internal")

	err = gen.WriteFile()
	if err != nil {
		return err
	}

	return nil
}

func GenerateNode() error {
	nd := &internal.NodeData{}

	node_gen, err := internal.NodeGenerator.Generate(internal.OutputLocFlag, "node", nd)
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

func GenerateAst(tk_symbols []*internal.Token) error {
	gen := &internal.ASTGen{
		Ast: internal.CandidatesForAst(tk_symbols),
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
