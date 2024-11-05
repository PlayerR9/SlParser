package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	gen "github.com/PlayerR9/SlParser/cmd/grammar/generator"
	pkg "github.com/PlayerR9/SlParser/cmd/grammar/pkg"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "[grammar]: ", log.Lshortfile)
}

var (
	InputFileLoc *string
)

func init() {
	InputFileLoc = flag.String("i", "", "input file location. Must be set and have .txt extension.")
}

func parseFlags() (string, error) {
	flag.Parse()

	if *InputFileLoc == "" {
		return "", errors.New("source file must be set")
	} else if filepath.Ext(*InputFileLoc) != ".txt" {
		return "", errors.New("source file must have extension .txt")
	}

	var builder strings.Builder

	_, err := fmt.Fprintf(&builder, "grammar -i=%s", *InputFileLoc)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func main() {
	sign, err := parseFlags()
	if err != nil {
		flag.PrintDefaults()

		logger.Fatalf("failed to parse flags: %v", err)
	}

	const (
		OutputPath string = "parsing.go"
		DirPath    string = "test"
	)

	data, err := os.ReadFile(*InputFileLoc)
	if err != nil {
		logger.Fatalf("failed to read %q: %v", *InputFileLoc, err)
	}

	rules, err := pkg.Parse(data)
	if err != nil {
		logger.Fatalf("failed to parse %q: %v", *InputFileLoc, err)
	}

	err = os.Mkdir(DirPath, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		logger.Fatalf("failed to create %q: %v", DirPath, err)
	}

	err = generateLexer(sign, DirPath)
	if err != nil {
		logger.Fatal(err)
	}

	err = generateParser(sign, DirPath, rules)
	if err != nil {
		logger.Fatal(err)
	}

	err = generateAST(sign, DirPath, rules)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("Successfully generated %q", OutputPath)
	os.Exit(0)
}

func generateLexer(sign, dir string) error {
	const (
		FileName string = "lexer.go"
	)

	data := gen.NewLexerData()

	path := filepath.Join(dir, FileName)

	err := gen.LexerGenerator.Generate(true, sign, path, data)
	if err != nil {
		return fmt.Errorf("failed to generate lexer: %w", err)
	}

	return nil
}

func generateParser(sign, dir string, rules []*pkg.Rule) error {
	const (
		FileName string = "parsing.go"
	)

	data, err := gen.NewParserData(rules)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, FileName)

	err = gen.ParserGenerator.Generate(false, sign, path, data)
	if err != nil {
		return fmt.Errorf("failed to generate parser: %w", err)
	}

	return nil
}

func generateAST(sign, dir string, rules []*pkg.Rule) error {
	const (
		FileName string = "ast.go"
	)

	data := gen.NewASTData(rules)

	path := filepath.Join(dir, FileName)

	err := gen.ASTGenerator.Generate(true, sign, path, data)
	if err != nil {
		return fmt.Errorf("failed to generate ast: %w", err)
	}

	return nil
}
