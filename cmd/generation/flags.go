package generation

import (
	"errors"
	"flag"

	ggen "github.com/PlayerR9/lib_units/generator"
)

var (
	// InputFileFlag is the flag used to specify the input file.
	InputFileFlag *string
)

func init() {
	ggen.SetOutputFlag("<dir>.go", true)

	InputFileFlag = flag.String("i", "", "The input file to parse. This flag is required.")
}

func ParseFlags() (string, error) {
	err := ggen.ParseFlags()
	if err != nil {
		return "", err
	}

	if *InputFileFlag == "" {
		return "", errors.New("input file is required")
	}

	return *InputFileFlag, nil
}
