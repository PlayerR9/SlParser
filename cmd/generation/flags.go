package generation

import (
	"errors"
	"flag"

	ggen "github.com/PlayerR9/go-generator/generator"
)

var (
	OutputLocFlag *ggen.OutputLocVal

	// InputFileFlag is the flag used to specify the input file.
	InputFileFlag *string
)

func init() {
	OutputLocFlag = ggen.NewOutputFlag("<dir>.go", true)

	InputFileFlag = flag.String("i", "", "The input file to parse. This flag is required.")
}

func ParseFlags() (string, error) {
	ggen.ParseFlags()

	if *InputFileFlag == "" {
		return "", errors.New("input file is required")
	}

	return *InputFileFlag, nil
}
