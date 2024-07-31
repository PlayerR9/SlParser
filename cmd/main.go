package main

import (
	"os"

	gen "github.com/PlayerR9/SLParser/cmd/generation"
	pkg "github.com/PlayerR9/SLParser/cmd/pkg"
	prx "github.com/PlayerR9/SLParser/parser"
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

	ee_data, err := pkg.ExtractEnums.Apply(root)
	if err != nil {
		gen.Logger.Fatalf("Error extracting enums: %s", err.Error())
	}

	g := &gen.Gen{
		SpecialEnums: ee_data.GetSpecialEnums(),
		LexerEnums:   ee_data.GetLexerEnums(),
		ParserEnums:  ee_data.GetParserEnums(),
	}

	_, err = pkg.RenameNodes.Apply(root)
	if err != nil {
		gen.Logger.Fatalf("Error renaming nodes: %s", err.Error())
	}

	dest, err := gen.Generator.Generate("test", ".go", g)
	if err != nil {
		gen.Logger.Fatalf("Error generating code: %s", err.Error())
	}

	gen.Logger.Printf("Successfully generated file: %q", dest)
}
