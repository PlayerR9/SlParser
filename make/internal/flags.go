package internal

import (
	"github.com/PlayerR9/go-generator"
)

var (
	OutputLocFlag *generator.OutputLocVal
)

func init() {
	OutputLocFlag = generator.NewOutputFlag("", true)
}
