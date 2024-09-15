package generation

import (
	"log"
	"os"
)

var (
	// Logger is the logger used to log messages.
	Logger *log.Logger
)

func init() {
	Logger = log.New(os.Stdout, "[Sl Parser]: ", log.LstdFlags)
}
