package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

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

	tk_symbols := internal.TokenSymbols(tokens)

	all_symbols := internal.ExtractSymbols(tk_symbols)
	Sort(all_symbols)

	var last_terminal string

	for _, symbol := range all_symbols {
		if strings.HasPrefix(symbol, "Ntt") {
			break
		}

		last_terminal = symbol
	}

	if last_terminal == "" {
		panic("missing terminal")
	}

	gd := &internal.GenData{
		Symbols:      all_symbols,
		LastTerminal: last_terminal,
	}

	for _, rule := range rules {
		gd.Rules = append(gd.Rules, rule.String())
	}

	generator.ParseFlags()

	gen, err := internal.Generator.Generate(internal.OutputLocFlag, "lexer", gd)
	if err != nil {
		panic(err)
	}

	err = gen.WriteFile()
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
