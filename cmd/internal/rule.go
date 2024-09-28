package internal

import (
	"fmt"
	"strings"
)

// Rule is a production of a grammar.
type Rule struct {
	// Elems is the elements of the rule.
	Elems []string
}

// String implements the fmt.Stringer interface.
func (r Rule) String() string {
	return "_ = is.AddRule(" + strings.Join(r.Elems, ", ") + ")"
}

// NewRule creates a new rule.
//
// Parameters:
//   - lhs: the left hand side of the rule.
//   - rhss: the right hand sides of the rule.
//
// Returns:
//   - *Rule: the new rule.
//   - error: the error if any.
func NewRule(lhs string, rhss []string) (*Rule, error) {
	if len(rhss) == 0 {
		return nil, fmt.Errorf("expected at least one rhss")
	}

	rule := &Rule{
		Elems: append([]string{lhs}, rhss...),
	}

	return rule, nil
}

func X() {
	// Node[0: Source]
	//     └── Node[0: Rule]
	//     │   ├── Node[0: Rhs (source)]
	//     │   ├── Node[9: Rhs (source1)]
	//     │   └── Node[17: Rhs (EOF)]

	// 		children := tk.GetChildren()
	// 		if len(children) != 2 {
	// 			return nil, fmt.Errorf("expected two children, got %d instead", len(children))
	// 		}

	// 		err := ast.CheckType(children, 1, EtEOF)
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		// source : source1 EOF ;

	// 		tmp, err := ast.LhsToAst(0, children, NtSource1, source1)
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		node := NewNode(SourceNode, "")
	// 		node.AddChildren(tmp)

	// return node, nil
}

//     └── Node[24: Rule]
//     │   ├── Node[24: Rhs (source1)]
//     │   └── Node[34: Rhs (rule)]
//     └── Node[41: Rule]
//     │   ├── Node[41: Rhs (source1)]
//     │   ├── Node[51: Rhs (rule)]
//     │   ├── Node[56: Rhs (NEWLINE)]
//     │   └── Node[64: Rhs (source1)]
//     └── Node[75: Rule]
//     │   ├── Node[75: Rhs (rule)]
//     │   ├── Node[82: Rhs (LOWERCASE_ID)]
//     │   ├── Node[95: Rhs (COLON)]
//     │   ├── Node[101: Rhs (rule1)]
//     │   └── Node[107: Rhs (SEMICOLON)]
//     └── Node[120: Rule]
//     │   ├── Node[120: Rhs (rule1)]
//     │   └── Node[128: Rhs (rhs)]
//     └── Node[134: Rule]
//     │   ├── Node[134: Rhs (rule1)]
//     │   ├── Node[142: Rhs (rhs)]
//     │   └── Node[146: Rhs (rule1)]
//     └── Node[155: Rule]
//     │   ├── Node[155: Rhs (rhs)]
//     │   └── Node[161: Rhs (UPPERCASE_ID)]
//     └── Node[176: Rule]
//         ├── Node[176: Rhs (rhs)]
//         └── Node[182: Rhs (LOWERCASE_ID)]
