package internal

import (
	"errors"
	"flag"

	"github.com/PlayerR9/go-generator"
)

var (
	// OutputLocFlag is the output location flag.
	OutputLocFlag *generator.OutputLocVal

	// InputFileFlag is the flag used to specify the input file.
	InputFileFlag *string
)

func init() {
	OutputLocFlag = generator.NewOutputFlag("", true)

	InputFileFlag = flag.String("input", "", "The input file to parse. This flag is required.")
}

// ParseFlags parses the command line flags.
//
// Returns:
//   - error: if an error occurred.
func ParseFlags() error {
	generator.ParseFlags()

	if *InputFileFlag == "" {
		return errors.New("missing input file")
	}

	return nil
}
