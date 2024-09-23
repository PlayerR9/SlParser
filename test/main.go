package main

import (
	"fmt"
	"log"
	"os"

	pkg "github.com/PlayerR9/SlParser/test/parsing"
)

var (
	Debugger *log.Logger
)

func init() {
	Debugger = log.New(os.Stdout, "[DEBUG]: ", log.LstdFlags)
}

func main() {
	err := ParseCmd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// [0, 1, 4, 9, 16, 25, 36, 49, 64, 81]
}

func ParseCmd() error {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		return err
	}

	parser := pkg.NewParser()
	parser.SetMode(pkg.ShowAll)

	_, err = parser.Full(data)
	if err != nil {
		return err
	}

	return nil
}
