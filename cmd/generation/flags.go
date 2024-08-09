package generation

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	prx "github.com/PlayerR9/SLParser/parser"
	ggen "github.com/PlayerR9/go-generator/generator"
)

var (
	OutputLocFlag *ggen.OutputLocVal

	// InputFileFlag is the flag used to specify the input file.
	InputFileFlag *string

	// EnableFlag is the flag used to enable debugging output.
	EnableFlag *EnableVal
)

func init() {
	OutputLocFlag = ggen.NewOutputFlag("<dir>.go", true)

	InputFileFlag = flag.String("i", "", "The input file to parse. This flag is required.")

	EnableFlag = NewEnableVal("e")
}

type FlagSet struct {
	Input  string
	Enable EnableVal
}

func ParseFlags() (*FlagSet, error) {
	ggen.ParseFlags()

	if *InputFileFlag == "" {
		return nil, errors.New("input file is required")
	}

	return &FlagSet{
		Input:  *InputFileFlag,
		Enable: *EnableFlag,
	}, nil
}

type EnableVal struct {
	l bool
	p bool
	a bool
	d bool
}

func (e *EnableVal) String() string {
	return fmt.Sprintf("l=%t, p=%t, a=%t, d=%t", e.l, e.p, e.a, e.d)
}

func (e *EnableVal) Set(value string) error {
	if value == "" {
		return nil
	}

	if strings.Contains(value, "l") {
		e.l = true
	}

	if strings.Contains(value, "p") {
		e.p = true
	}

	if strings.Contains(value, "a") {
		e.a = true
	}

	if strings.Contains(value, "d") {
		e.d = true
	}

	return nil
}

func NewEnableVal(name string) *EnableVal {
	value := &EnableVal{}

	flag.Var(value, name, "Enable debugging output: l for lexer, p for parser, d for data, and a for ast. This flag is optional.")

	return value
}

func (e *EnableVal) Get() prx.DebugSetting {
	var sum prx.DebugSetting

	if e.l {
		sum |= prx.ShowLex
	}

	if e.p {
		sum |= prx.ShowTree
	}

	if e.a {
		sum |= prx.ShowAst
	}

	if e.d {
		sum |= prx.ShowData
	}

	return sum
}
