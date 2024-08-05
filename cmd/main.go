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

	err = GenerateTokens(root)
	if err != nil {
		gen.Logger.Fatalf("While generating tokens: %s", err.Error())
	} else {
		gen.Logger.Printf("Successfully generated tokens.")
	}

	_, err = pkg.RenameNodes.Apply(root)
	if err != nil {
		gen.Logger.Fatalf("Error renaming nodes: %s", err.Error())
	}

	err = GenerateLexer()
	if err != nil {
		gen.Logger.Fatalf("While generating lexer: %s", err.Error())
	} else {
		gen.Logger.Printf("Successfully generated lexer.")
	}

	err = GenerateParser(root)
	if err != nil {
		gen.Logger.Fatalf("While generating parser: %s", err.Error())
	} else {
		gen.Logger.Printf("Successfully generated parser.")
	}
}

func GenerateTokens(root *ast.Node[prx.NodeType]) error {
	ee_data, err := pkg.ExtractEnums.Apply(root)
	if err != nil {
		return err
	}

	g := &gen.TokenGen{
		SpecialEnums: ee_data.GetSpecialEnums(),
		LexerEnums:   ee_data.GetLexerEnums(),
		ParserEnums:  ee_data.GetParserEnums(),
	}

	dest, err := gen.TokenGenerator.Generate("test", ".go", g)
	if err != nil {
		return err
	}

	dest.ModifyFileName("_tokens")

	err = dest.WriteFile()
	if err != nil {
		return err
	}

	return nil
}

func GenerateLexer() error {
	g := &gen.LexerGen{}

	dest, err := gen.LexerGenerator.Generate("test", ".go", g)
	if err != nil {
		return err
	}

	dest.ModifyFileName("_lexer")

	err = dest.WriteFile()
	if err != nil {
		return err
	}

	return nil
}

func GenerateParser(root *ast.Node[prx.NodeType]) error {
	rules, err := pkg.ExtractRules(root)
	if err != nil {
		return err
	}

	g := &gen.ParserGen{
		Rules: pkg.StringifyRules(rules),
	}

	dest, err := gen.ParserGenerator.Generate("test", ".go", g)
	if err != nil {
		return err
	}

	dest.ModifyFileName("_parser")

	err = dest.WriteFile()
	if err != nil {
		return err
	}

	return nil
}
