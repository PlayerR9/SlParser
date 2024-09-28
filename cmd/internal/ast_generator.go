package internal

import (
	kdd "github.com/PlayerR9/SlParser/kdd"
	gers "github.com/PlayerR9/go-errors"
	"github.com/PlayerR9/go-generator"
)

// ASTGen is a generator for the AST.
type ASTGen struct {
	// PackageName is the name of the package.
	PackageName string

	// NodeTypes is the list of candidates for the node types in the AST.
	NodeTypes []string
}

// SetPackageName implements the generator.PackageNameSetter interface.
func (g *ASTGen) SetPackageName(pkg_name string) {
	if g == nil {
		return
	}

	g.PackageName = pkg_name
}

// NewASTGen creates a new ASTGen with the given tokens.
//
// Parameters:
//   - table: The information about the AST.
//
// Returns:
//   - *ASTGen: the ASTGen. Never returns nil.
func NewASTGen(table map[*kdd.Node]*Info) *ASTGen {
	candidates := CandidatesForAst(table)

	gen := &ASTGen{
		NodeTypes: candidates,
	}

	return gen
}

type RuleInfo struct {
	Lhs             string
	AllowedChildren []string
}

/*
func X(table map[*kdd.Node]*Info, root *kdd.Node) {
	//     └── Node[0: Rule]
	//     │   ├── Node[0: Rhs (source)]
	//     │   ├── Node[9: Rhs (source1)]
	//     │   └── Node[17: Rhs (EOF)]

	rule_table := make(map[string][]*RuleInfo)

	for rule := range root.Child() {
		rhss := rule.GetChildren()

		lhs := rhss[0].Data

		allowed_children := make([]string, len(rhss)-1)

		children_size := len(rhss) - 1

		// TODO: Do this in a better way.

		if rule is source, then it must have two children:
		- a rule of type source1
		- a rule of type EOF.

		Since EOF is a terminal, then do an immediate assertion.
		Since source1 is marked as a one or more, then use the lhsrhs rule.

		//     └── Node[0: Rule]
		//     │   ├── Node[0: Rhs (source)]
		//     │   ├── Node[9: Rhs (source1)]
		//     │   └── Node[17: Rhs (EOF)]

		new_rule := &RuleInfo{
			Lhs: lhs,
			ChildrenSize: children_size,
		}

		prev, ok := rule_table[lhs]
		if !ok {
			prev = []*RuleInfo{new_rule}
		} else {
			prev = append(prev, new_rule)
		}

		rule_table[lhs] = prev
	}




	// Node[0: Source]
	//     └── Node[0: Rule]
	//     │   ├── Node[0: Rhs (source)]
	//     │   ├── Node[9: Rhs (source1)]
	//     │   └── Node[17: Rhs (EOF)]
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
} */

var (
	// ASTGenerator is the generator for the AST.
	ASTGenerator *generator.CodeGenerator[*ASTGen]
)

func init() {
	ASTGenerator = gers.AssertNew(
		generator.NewCodeGeneratorFromTemplate[*ASTGen]("ast", ast_templ),
	)
}

// ast_templ is the template for the AST.
const ast_templ = `
package {{ .PackageName }}

import (
	"errors"

	"github.com/PlayerR9/SlParser/ast"
	"github.com/PlayerR9/SlParser/grammar"
	gers "github.com/PlayerR9/go-errors"
)

// NodeType is the type of a node.
type NodeType int

const (
	/*InvalidNode represents an invalid node.
	Node[InvalidNode]
	*/
	InvalidNode NodeType = iota - 1 // Invalid {{ range $index, $value := .NodeTypes }}

	/*{{ $value }}Node is [...].
	Node[{{ $value }}Node]
	*/
	{{ $value }}Node // {{ $value }}
	{{- end }}
)

var (
	ast_maker ast.AstMaker[*Node, TokenType]
)
	
func init() {
	ast_maker = make(ast.AstMaker[*Node, TokenType])

	// TODO: Add here your own custom rules...
	{{ range $index, $value := .NodeTypes }}
	ast_maker[{{ $value }}] = func(tk *grammar.ParseTree[TokenType]) (*Node, error) {
		children := tk.GetChildren()
		if len(children) == 0 {
			return nil, errors.New("expected at least one child")
		}

		// TODO: Complete this function...

		node := NewNode(tk.Pos(), {{ $value }}Node, "")
		return node, nil
	}
	{{- end }}
}
	
var (
	// What Counts as a Valid Node?
	//
	// Overall Rules:
	//  1. All children must be valid.
	//  2. A node cannot be nil.
	// {{ range $index, $node := .NodeTypes }}
	// {{ $node }} Rules:
	//
	// {{- end }}
	CheckAST ast.CheckASTWithLimit[*Node]
)

func init() {
	table := make(map[NodeType]ast.CheckNodeFn[*Node])
	{{ range $index, $node := .NodeTypes }}
	table[{{ $node }}] = func(node *Node) error {
		gers.AssertNotNil(node, "node")

		// TODO: Specify what counts as a valid node.

		return nil
	}
	{{ - end }}

	CheckAST = ast.MakeCheckFn(table)
}`
