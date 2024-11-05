package generator

import (
	"io"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/PlayerR9/SlParser/cmd/grammar/pkg"
	common "github.com/PlayerR9/mygo-lib/common"
	cgen "github.com/PlayerR9/mygo-lib/generator"
)

type ParserData struct {
	PackageName string

	Symbols string
	Lines   []string
}

func (d *ParserData) SetPkgName(pkg_name string) error {
	if d == nil {
		return common.ErrNilReceiver
	}

	d.PackageName = pkg_name

	return nil
}

func NewParserData(rules []*pkg.Rule) (*ParserData, error) {
	data := &ParserData{
		Lines: make([]string, 0, len(rules)),
	}

	for _, rule := range rules {
		data.Lines = append(data.Lines, rule.Lines())
	}

	symbols := pkg.DetermineTokenTypes(rules)

	var builder strings.Builder
	w := tabwriter.NewWriter(&builder, 0, 3, 3, ' ', tabwriter.FilterHTML|tabwriter.StripEscape)
	for _, rhs := range symbols {
		data := []byte("\t" + rhs + "\tstring = " + strconv.Quote(rhs) + "\n")

		n, err := w.Write(data)
		if err != nil {
			return nil, err
		} else if n != len(data) {
			return nil, io.ErrShortWrite
		}
	}

	err := w.Flush()
	if err != nil {
		return nil, err
	}

	data.Symbols = builder.String()

	return data, nil
}

var (
	ParserGenerator *cgen.CodeGenerator[*ParserData]
)

func init() {
	ParserGenerator = cgen.Must(cgen.New[*ParserData]("parser", parser_templ))
}

const parser_templ string = `package {{ .PackageName }}

import (
	"io"
	"github.com/PlayerR9/SlParser"
	"github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/SlParser/parser"
	"github.com/PlayerR9/mygo-lib/trees"
	"github.com/PlayerR9/go-evals/result"
)

const (
{{ .Symbols }})

var (
	WriteTokenTree func(w io.Writer, tk *grammar.Token) (int, error)
	Parser parser.Parser
)

func init() {
	WriteTokenTree = trees.MakeWriteTree[*grammar.Token]()
	
	var builder parser.Builder
	{{ range $index, $line := .Lines }}
	{{ $line }}
	{{- end }}

	Parser = parser.NewParser(builder.Build())
}

var evalFn result.ApplyOnValidsFn[SlParser.Result]

func init() {
	var err error

	evalFn, err = SlParser.MakeEvaluate(Lexer, Parser, Ast)
	if err != nil {
		panic(err)
	}
}

// Result holds all the information regarding the parsing process. This is read-only.
type Result struct {
	// inner holds the inner result.
	inner *SlParser.Result
}

// Data returns the data of the result.
//
// Returns:
//   - []byte: The data of the result.
//   - bool: True if the data is set, false otherwise.
func (r Result) Data() ([]byte, bool) {
	data, err := r.inner.Data()
	return data, err == nil
}

// Tokens returns the tokens of the result.
//
// Returns:
//   - []*grammar.Token: The tokens of the result.
//   - bool: True if the tokens are set, false otherwise.
func (r Result) Tokens() ([]*grammar.Token, bool) {
	tokens, err := r.inner.Tokens()
	return tokens, err == nil
}

// ParseTree returns the parse tree of the result.
//
// Returns:
//   - *parser.Result: The parse tree of the result.
//   - bool: True if the parse tree is set, false otherwise.
func (r Result) ParseTree() (*parser.Result, bool) {
	pr, err := r.inner.ParseTree()
	return pr, err == nil
}

// Node returns the node of the result.
//
// Returns:
//   - *grammar.Node: The node of the result.
//   - bool: True if the node is set, false otherwise.
func (r Result) Node() (*grammar.Node, bool) {
	n, err := r.inner.Node()
	return n, err == nil
}

// LexerErr returns the lexer error of the result.
//
// Returns:
//   - error: The lexer error of the result.
//   - bool: True if the lexer error is set, false otherwise.
func (r Result) LexerErr() (error, bool) {
	err1, err2 := r.inner.LexerErr()
	return err1, err2 == nil
}

// Err returns the error of the result.
//
// Returns:
//   - error: The error of the result. Nil if no error is set.
func (r Result) Err() error {
	return r.inner.Err()
}

// Parse parses the given data according to the grammar.
//
// Parameters:
//   - data: The data to parse.
//
// Returns:
//   - []Result: A slice containing the result of the parsing process.
//   - error: An error if the evaluation failed.
func Parse(data []byte) ([]Result, error) {
	result := SlParser.NewResult(data)

	results, err := Evaluate([]SlParser.Result{result})

	slice := make([]Result, 0, len(results))

	for _, r := range results {
		rsl := Result{
			inner: &r,
		}
		slice = append(slice, rsl)
	}

	return slice, err
}`
