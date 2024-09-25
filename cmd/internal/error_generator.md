package internal

import (
	gers "github.com/PlayerR9/go-errors"
	"github.com/PlayerR9/go-generator"
)

// ErrorGen is a generator for the errors.
type ErrorGen struct {
	// PackageName is the package name.
	PackageName string
}

// SetPackageName implements the generator.PackageNameSetter interface.
func (eg *ErrorGen) SetPackageName(pkg_name string) {
	if eg == nil {
		return
	}

	eg.PackageName = pkg_name
}

// NewErrorGen creates a new ErrorGen.
//
// Returns:
//   - *ErrorGen: The new ErrorGen. Never returns nil.
func NewErrorGen() *ErrorGen {
	return &ErrorGen{}
}

var (
	// ErrorGenerator is the error generator.
	ErrorGenerator *generator.CodeGenerator[*ErrorGen]
)

func init() {
	var err error

	ErrorGenerator, err = generator.NewCodeGeneratorFromTemplate[*ErrorGen]("error", error_templ)
	gers.AssertErr(err, "generator.NewCodeGeneratorFromTemplate[*ErrorGen](%q, error_templ)", "error")
}

// error_templ is the template for the error.
const error_templ string = `
package {{ .PackageName }}

import (
	gerr "github.com/PlayerR9/go-errors/error"
)

type ErrorCode int

const (
	// InvalidSyntax occurs when the AST is invalid or not
	// as it should be.
	InvalidSyntax ErrorCode = iota
)

// Int implements the error.ErrorCoder interface.
func (e ErrorCode) Int() int {
	return int(e)
}`
