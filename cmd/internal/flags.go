package internal

import (
	ggen "github.com/PlayerR9/go-generator"
)

var (
	// OutputLocFlag is the output location flag.
	OutputLocFlag *ggen.OutputLocVal
)

func init() {
	OutputLocFlag = ggen.NewOutputFlag("", true)
}
