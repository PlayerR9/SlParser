package generation

import (
	"log"
	"os"

	ggen "github.com/PlayerR9/lib_units/generator"
)

var (
	// Logger is the logger used to log messages.
	Logger *log.Logger
)

func init() {
	Logger = ggen.InitLogger(os.Stdout, "SL parser")
}
