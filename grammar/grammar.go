package grammar

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// StringToEnum is a function that converts a string to an enum.
//
// Parameters:
//   - str: The string to convert.
//
// Returns:
//   - T: The converted enum.
//   - bool: True if the string was converted successfully, false otherwise.
type StringToEnum[T TokenTyper] func(str string) (T, bool)

// ReverseGrammar is a function that reverses a grammar.
//
// Parameters:
//   - grammar: The grammar to reverse.
//   - conv_func: The function that converts a string to an enum.
//
// Returns:
//   - []*Rule[T]: The reversed grammar.
//   - error: The error if any.
func ReverseGrammar[T TokenTyper](grammar string, conv_func StringToEnum[T]) ([]*Rule[T], error) {
	if grammar == "" {
		return nil, uc.NewErrInvalidParameter("grammar", uc.NewErrEmpty("grammar"))
	} else if conv_func == nil {
		return nil, uc.NewErrNilParameter("conv_func")
	}

	lines := strings.Split(grammar, "\n")

	var rules []*Rule[T]

	for i, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)

		rev, err := reverse_rule(fields, conv_func)
		if err != nil {
			return nil, uc.NewErrAt(i+1, "line", err)
		}
		uc.AssertNil(rev, "rev")

		rules = append(rules, rev)
	}

	return rules, nil
}

func reverse_rule[T TokenTyper](fields []string, conv_func StringToEnum[T]) (*Rule[T], error) {
	uc.Assert(len(fields) > 0, "expected at least one field")
	uc.AssertParam("conv_func", conv_func != nil, errors.New("value must not be nil"))

	lhs, ok := conv_func(fields[0])
	if !ok {
		return nil, fmt.Errorf("string ( %q ) is not a valid token type", fields[0])
	}

	if len(fields) == 1 {
		return nil, fmt.Errorf("expected \"equal sign\", got nothing instead")
	} else if fields[1] != "=" {
		return nil, fmt.Errorf("expected \"equal sign\", got %q instead", fields[1])
	}

	if fields[len(fields)-1] != "." {
		return nil, fmt.Errorf("expected \"dot\", got %q instead", fields[len(fields)-1])
	}

	if len(fields) == 2 {
		return nil, fmt.Errorf("expected \"rhs\", got nothing instead")
	}

	rhss := fields[2 : len(fields)-1]

	slices.Reverse(rhss)

	conv := make([]T, 0, len(rhss))

	for _, rhs := range rhss {
		token, ok := conv_func(rhs)
		if !ok {
			return nil, fmt.Errorf("string ( %q ) is not a valid token type", rhs)
		}

		conv = append(conv, token)
	}

	return NewRule(lhs, conv), nil
}
