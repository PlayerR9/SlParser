package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/PlayerR9/SlParser/make/internal"
	"github.com/PlayerR9/go-generator"
)

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	lines := bytes.Split(data, []byte{'\n'})

	var tokens []*internal.Token

	var rules []*internal.Rule

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		fields := bytes.Fields(line)
		if len(fields) == 0 {
			continue
		}

		if len(fields) <= 3 {
			panic(fmt.Errorf("invalid line: %q", line))
		}

		if !bytes.Equal(fields[1], []byte(":")) {
			panic(fmt.Errorf("missing colon: %q", line))
		}

		if !bytes.Equal(fields[len(fields)-1], []byte(";")) {
			panic(fmt.Errorf("missing semicolon: %q", line))
		}

		lhs, err := internal.MakeToken(fields[0])
		if err != nil {
			panic(fmt.Errorf("invalid lhs: %w", err))
		}

		tokens = append(tokens, lhs)

		var rhss []*internal.Token

		for i := 2; i < len(fields)-1; i++ {
			tk, err := internal.MakeToken(fields[i])
			if err != nil {
				panic(fmt.Errorf("invalid rhs at %d: %w", i, err))
			}

			rhss = append(rhss, tk)
			tokens = append(tokens, tk)
		}

		rule := internal.NewRule(lhs, rhss)
		rules = append(rules, rule)
	}

	generator.ParseFlags()

	tk_symbols, err := internal.TokenSymbols(tokens)
	if err != nil {
		panic(err)
	}

	err = GenerateParser(tk_symbols, rules)
	if err != nil {
		panic(err)
	}

	err = GenerateNode()
	if err != nil {
		panic(err)
	}

	err = GenerateAst(tk_symbols)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("go", "generate")
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully generated lexer.")
}

func GenerateParser(tk_symbols []*internal.Token, rules []*internal.Rule) error {
	ok := internal.CheckEofExists(tk_symbols)
	if !ok {
		return errors.New("missing EOF")
	}

	all_symbols := internal.ExtractSymbols(tk_symbols)

	last_terminal := internal.FindLastTerminal(tk_symbols)
	if last_terminal == nil {
		return errors.New("missing terminal")
	}

	gd := &internal.GenData{
		Symbols:      all_symbols,
		LastTerminal: last_terminal.String(),
	}

	for _, rule := range rules {
		gd.Rules = append(gd.Rules, rule.String())
	}

	gen, err := internal.Generator.Generate(internal.OutputLocFlag, "lexer", gd)
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
