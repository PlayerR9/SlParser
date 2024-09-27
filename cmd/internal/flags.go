package internal

import (
	"errors"
	"flag"
	"os"

	"github.com/PlayerR9/go-generator"
)

var (
	// GrammarNameFlag is the flag used to specify the grammar name.
	GrammarNameFlag *string

	// InputFileFlag is the flag used to specify the input file.
	InputFileFlag *string

	// ForceFlag is the flag used to force the creation of the files even if the
	// directory already exists.
	ForceFlag *bool
)

func init() {
	GrammarNameFlag = flag.String("name", "", "The name of the grammar to generate. This flag is required.")
	InputFileFlag = flag.String("input", "", "The input file to parse. This flag is required.")
	ForceFlag = flag.Bool("y", false, "Whether to force the creation if the directory already exists.")
}

// ParseFlags parses the command line flags.
//
// Returns:
//   - string: The grammar name.
//   - error: if an error occurred.
func ParseFlags() (string, error) {
	generator.ParseFlags()

	if *InputFileFlag == "" {
		return "", errors.New("missing input file")
	}

	if *GrammarNameFlag == "" {
		return "", errors.New("missing grammar name")
	}

	dir := *GrammarNameFlag

	_, err := os.Stat(dir)
	if err == nil {
		if !*ForceFlag {
			return "", errors.New("directory already exists")
		}
	} else if !os.IsNotExist(err) {
		return "", err
	} else {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return "", err
		}
	}

	return dir, nil
}
