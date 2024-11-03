package pkg

import (
	"strings"

	gslc "github.com/PlayerR9/mygo-lib/slices"
)

type Rule struct {
	Lhs  string
	Rhss []string
}

func (r Rule) String() string {
	return r.Lhs + " : " + strings.Join(r.Rhss, " ") + " ;"
}

func NewRule(lhs []byte, rhss [][]byte) (*Rule, error) {
	lhs_str, err := FixEnumName(lhs)
	if err != nil {
		return nil, err
	}

	rhss_strs := make([]string, 0, len(rhss))

	for _, rhs := range rhss {
		rhs_str, err := FixEnumName(rhs)
		if err != nil {
			return nil, err
		}

		rhss_strs = append(rhss_strs, rhs_str)
	}

	return &Rule{
		Lhs:  lhs_str,
		Rhss: rhss_strs,
	}, nil
}

func (r Rule) Symbols() []string {
	symbols := make([]string, len(r.Rhss))
	copy(symbols, r.Rhss)

	_ = gslc.Uniquefy(&symbols)

	_, _ = gslc.MayInsert(&symbols, r.Lhs)

	return symbols
}

func (r Rule) Lines() string {
	rhss := make([]string, len(r.Rhss)+1)
	rhss[0] = r.Lhs
	copy(rhss[1:], r.Rhss)

	str := strings.Join(rhss, ", ")

	return "_ = builder.AddRule(" + str + ")"
}
