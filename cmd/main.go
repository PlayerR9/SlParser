package main

import (
	"os"

	gen "github.com/PlayerR9/SLParser/cmd/generation"
	pkg "github.com/PlayerR9/SLParser/cmd/pkg"
	prx "github.com/PlayerR9/SLParser/parser"
	ast "github.com/PlayerR9/grammar/ast"
)

func main() {
	source, err := gen.ParseFlags()
	if err != nil {
		gen.Logger.Fatalf("Error parsing flags: %s", err.Error())
	}

	data, err := os.ReadFile(source)
	if err != nil {
		gen.Logger.Fatalf("Error reading file: %s", err.Error())
	}

	root, err := prx.Parse(data)
	if err != nil {
		gen.Logger.Fatalf("Error parsing file: %s", err.Error())
	}

	dest, err := GenerateTokens(root)
	if err != nil {
		gen.Logger.Fatalf("While generating tokens: %s", err.Error())
	} else {
		gen.Logger.Printf("Successfully generated tokens: %q", dest)
	}

	_, err = pkg.RenameNodes.Apply(root)
	if err != nil {
		gen.Logger.Fatalf("Error renaming nodes: %s", err.Error())
	}

	dest, err = GenerateLexer()
	if err != nil {
		gen.Logger.Fatalf("While generating lexer: %s", err.Error())
	} else {
		gen.Logger.Printf("Successfully generated lexer: %q", dest)
	}

	dest, err = GenerateParser(root)
	if err != nil {
		gen.Logger.Fatalf("While generating parser: %s", err.Error())
	} else {
		gen.Logger.Printf("Successfully generated parser: %q", dest)
	}
}

func GenerateTokens(root *ast.Node[prx.NodeType]) (string, error) {
	ee_data, err := pkg.ExtractEnums.Apply(root)
	if err != nil {
		return "", err
	}

	g := &gen.TokenGen{
		SpecialEnums: ee_data.GetSpecialEnums(),
		LexerEnums:   ee_data.GetLexerEnums(),
		ParserEnums:  ee_data.GetParserEnums(),
	}

	res, err := gen.TokenGenerator.Generate(gen.OutputLocFlag, "test.go", g)
	if err != nil {
		return "", err
	}

	dest, err := res.WriteFile("_tokens")
	if err != nil {
		return "", err
	}

	return dest, nil
}

func GenerateLexer() (string, error) {
	g := &gen.LexerGen{}

	res, err := gen.LexerGenerator.Generate(gen.OutputLocFlag, "test.go", g)
	if err != nil {
		return "", err
	}

	dest, err := res.WriteFile("_lexer")
	if err != nil {
		return "", err
	}

	return dest, nil
}

func GenerateParser(root *ast.Node[prx.NodeType]) (string, error) {
	rules, err := pkg.ExtractRules(root)
	if err != nil {
		return "", err
	}

	g := &gen.ParserGen{
		Rules: pkg.StringifyRules(rules),
	}

	res, err := gen.ParserGenerator.Generate(gen.OutputLocFlag, "test.go", g)
	if err != nil {
		return "", err
	}

	dest, err := res.WriteFile("_parser")
	if err != nil {
		return "", err
	}

	return dest, nil
}
