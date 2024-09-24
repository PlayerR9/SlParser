package internal

import (
	"errors"
	"flag"

	"github.com/PlayerR9/go-generator"
)

var (
	OutputLocFlag *generator.OutputLocVal
	InputFileFlag *string
)

func init() {
	OutputLocFlag = generator.NewOutputFlag("", true)

	InputFileFlag = flag.String("input", "", "The input file to parse. This flag is required.")
}

func ParseFlags() error {
	generator.ParseFlags()

	if *InputFileFlag == "" {
		return errors.New("missing input file")
	}

	return nil
}
