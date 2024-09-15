package generation

import "strings"

type AstCaseGen struct {
	Lengths []string
	Target  string
}

func NewAstCaseGen(target string, lengths []string) AstCaseGen {
	return AstCaseGen{
		Lengths: lengths,
		Target:  target,
	}
}

func (a AstCaseGen) Generate(indent int) string {
	var lines []string

	for _, length := range a.Lengths {
		lines = append(lines, "case "+length+":")
		lines = append(lines, "\tvar sub_nodes []ast.Noder")
		lines = append(lines, "")
		lines = append(lines, "\t// Extract here any desired sub-node...")
		lines = append(lines, "")
		lines = append(lines, "\tn := NewNode("+a.Target+", \"\", children[0].At)")
		lines = append(lines, "\ta.SetNode(&n)")
		lines = append(lines, "\t_ = a.AppendChildren(sub_nodes)")
	}

	if indent <= 0 {
		return strings.Join(lines, "\n")
	}

	indentation := strings.Repeat("\t", indent)

	for i := 0; i < len(lines); i++ {
		lines[i] = indentation + lines[i]
	}

	return strings.Join(lines, "\n")
}
