package main

import (
	"os"
	"os/exec"
	"path/filepath"

	gen "github.com/PlayerR9/SLParser/cmd/generation"
	pkg "github.com/PlayerR9/SLParser/cmd/pkg"
	ebnf "github.com/PlayerR9/SLParser/ebnf"
)

func main() {
	fs, err := gen.ParseFlags()
	if err != nil {
		gen.Logger.Fatalf("Error parsing flags: %s", err.Error())
	}

	data, err := os.ReadFile(fs.Input)
	if err != nil {
		gen.Logger.Fatalf("Error reading file: %s", err.Error())
	}

	ebnf.Parser.SetDebug(fs.Enable.Get())

	root, err := ebnf.Parser.Parse(data)
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

	rules, err := pkg.ExtractRules(root)
	if err != nil {
		gen.Logger.Fatalf("While extracting rules: %s", err.Error())
	}

	dest, err = GenerateParser(rules)
	if err != nil {
		gen.Logger.Fatalf("While generating parser: %s", err.Error())
	} else {
		gen.Logger.Printf("Successfully generated parser: %q", dest)
	}

	dest, err = GenerateAST(rules)
	if err != nil {
		gen.Logger.Fatalf("While generating ast: %s", err.Error())
	} else {
		gen.Logger.Printf("Successfully generated ast: %q", dest)
	}

	dest, err = GenerateGrammar()
	if err != nil {
		gen.Logger.Fatalf("While generating grammar: %s", err.Error())
	} else {
		gen.Logger.Printf("Successfully generated grammar: %q", dest)
	}

	dir, _ := filepath.Split(dest)

	cmd := exec.Command("go", "run", "github.com/PlayerR9/grammar/cmd", "-name=Node", "-type=NodeType", "-o="+filepath.Join(dir, "node.go"))

	err = cmd.Run()
	if err != nil {
		gen.Logger.Fatalf("While generating node: %s", err.Error())
	}

	gen.Logger.Printf("Done!")
}

func GenerateTokens(root *ebnf.Node) (string, error) {
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

	res.DestLoc = ModifyPrefixPath(res.DestLoc, "cmp_token_")

	err = res.WriteFile()
	if err != nil {
		return "", err
	}

	return res.DestLoc, nil
}

func GenerateLexer() (string, error) {
	g := &gen.LexerGen{}

	res, err := gen.LexerGenerator.Generate(gen.OutputLocFlag, "test.go", g)
	if err != nil {
		return "", err
	}

	res.DestLoc = ModifyPrefixPath(res.DestLoc, "cmp_lexer_")

	err = res.WriteFile()
	if err != nil {
		return "", err
	}

	return res.DestLoc, nil
}

func GenerateParser(rules []*pkg.Rule) (string, error) {
	dt := pkg.NewDecisionTable(rules)

	g := &gen.ParserGen{
		Table: dt,
	}

	res, err := gen.ParserGenerator.Generate(gen.OutputLocFlag, "test.go", g)
	if err != nil {
		return "", err
	}

	res.DestLoc = ModifyPrefixPath(res.DestLoc, "cmp_parser_")

	err = res.WriteFile()
	if err != nil {
		return "", err
	}

	return res.DestLoc, nil
}

func GenerateAST(rules []*pkg.Rule) (string, error) {
	table := make(map[string][]*pkg.Rule)

	for _, rule := range rules {
		prev, ok := table[rule.GetLhs()]
		if ok {
			prev = append(prev, rule)
			table[rule.GetLhs()] = prev
		} else {
			table[rule.GetLhs()] = []*pkg.Rule{rule}
		}
	}

	g := &gen.ASTGen{
		Table: table,
	}

	res, err := gen.ASTGenerator.Generate(gen.OutputLocFlag, "test.go", g)
	if err != nil {
		return "", err
	}

	res.DestLoc = ModifyPrefixPath(res.DestLoc, "cmp_ast_")

	err = res.WriteFile()
	if err != nil {
		return "", err
	}

	return res.DestLoc, nil
}

func GenerateGrammar() (string, error) {
	g := &gen.GrammarGen{}

	res, err := gen.GrammarGenerator.Generate(gen.OutputLocFlag, "test.go", g)
	if err != nil {
		return "", err
	}

	res.DestLoc = ModifyPrefixPath(res.DestLoc, "grammar_")

	err = res.WriteFile()
	if err != nil {
		return "", err
	}

	return res.DestLoc, nil
}

// ModifyPrefixPath modifies the path of the generated code.
//
// Parameters:
//   - dest_loc: The destination location of the generated code.
//   - prefix: The prefix to add to the file name. If empty, no prefix is added.
//   - sub_directories: The sub directories to create the file in.
//
// The prefix is useful for when generating multiple files as it adds a prefix without
// changing the extension.
func ModifyPrefixPath(dest_loc string, prefix string, sub_directories ...string) string {
	var loc string

	dir, file := filepath.Split(dest_loc)

	if len(sub_directories) > 0 {
		loc = filepath.Join(dir, filepath.Join(sub_directories...), prefix+file)
	} else {
		loc = filepath.Join(dir, prefix+file)
	}

	return loc
}
