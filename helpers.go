package SlParser

import (
	"iter"

	"github.com/PlayerR9/SlParser/ast"
	slpx "github.com/PlayerR9/SlParser/parser"
)

func HasTree[N interface {
	Child() iter.Seq[N]

	ast.Noder
}](results []Result[N], target *slpx.ParseTree) bool {
	if len(results) == 0 || target == nil {
		return false
	}

	for _, res := range results {
		r, err := res.ParseTree()
		if err != nil {
			continue
		}

		forest := r.Forest()
		if len(forest) != 1 {
			continue
		}

		tree := forest[0]

		if tree.Equals(target) {
			return true
		}
	}

	return false
}
