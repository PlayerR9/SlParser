package ast

import (
	"errors"
	"fmt"
	"iter"
	"slices"

	slgr "github.com/PlayerR9/SlParser/grammar"
	trav "github.com/PlayerR9/SlParser/trav"
	"github.com/PlayerR9/mygo-lib/common"
	gslc "github.com/PlayerR9/mygo-lib/slices"
)

func CheckNode(kind string, node Noder, types ...string) error {
	if node == nil {
		return gslc.NewErrNotAsExpected(true, kind, nil, types...)
	}

	got := node.GetType()

	ok := slices.Contains(types, got)
	if !ok {
		return gslc.NewErrNotAsExpected(true, kind, &got, types...)
	}

	return nil
}

func CheckToken(kind string, token *slgr.Token, types ...string) error {
	if token == nil {
		return gslc.NewErrNotAsExpected(true, kind, nil, types...)
	}

	got := token.Type

	ok := slices.Contains(types, got)
	if !ok {
		return gslc.NewErrNotAsExpected(true, kind, &got, types...)
	}

	return nil
}

type CheckFn[N interface {
	Child() iter.Seq[N]

	Noder
}] func(n N, depth int) error

type CheckNodeFn[N interface {
	Child() iter.Seq[N]

	Noder
}] func(n N, children []N) error

type Checker[N interface {
	Child() iter.Seq[N]

	Noder
}] struct {
	table map[string]CheckNodeFn[N]
}

func NewChecker[N interface {
	Child() iter.Seq[N]

	Noder
}]() *Checker[N] {
	return &Checker[N]{
		table: make(map[string]CheckNodeFn[N]),
	}
}

func (c *Checker[N]) Register(type_ string, fn CheckNodeFn[N]) error {
	if fn == nil {
		return nil
	} else if c == nil {
		return common.ErrNilReceiver
	}

	// assert.Cond(c.table != nil, "c.table must not be nil")

	c.table[type_] = fn

	return nil
}

func (c Checker[N]) Build() CheckFn[N] {
	if len(c.table) == 0 {
		return func(n N, depth int) error {
			if n.IsNil() {
				return common.ErrNilReceiver
			}

			type_ := n.GetType()

			return fmt.Errorf("node type (%q) is not supported", type_)
		}
	}

	table := make(map[string]CheckNodeFn[N], len(c.table))
	for k, v := range c.table {
		table[k] = v
	}

	type CheckerInfo struct {
		depth int
	}

	do := func(node N, info trav.Info) error {
		if node.IsNil() {
			return ErrNilNode
		}

		inf, ok := info.(CheckerInfo)
		if !ok {
			return errors.New("info must be of type CheckerInfo")
		}

		if inf.depth == 0 {
			return nil
		}

		type_ := node.GetType()

		fn, ok := table[type_]
		if !ok || fn == nil {
			return fmt.Errorf("node type (%q) is not supported", type_)
		}

		err := fn(node, slices.Collect(node.Child()))
		if err != nil {
			return fmt.Errorf("failed to check node (%q): %w", type_, err)
		}

		return nil
	}

	next := func(node N, info trav.Info) ([]trav.Pair[N], error) {
		if node.IsNil() {
			return nil, ErrNilNode
		}

		inf, ok := info.(CheckerInfo)
		if !ok {
			return nil, errors.New("info must be of type CheckerInfo")
		}

		var pairs []trav.Pair[N]

		if inf.depth < 0 {
			for child := range node.Child() {
				p := trav.NewPair(child, CheckerInfo{
					depth: inf.depth - 1,
				}, false)

				pairs = append(pairs, p)
			}
		} else if inf.depth != 0 {
			var pairs []trav.Pair[N]

			for child := range node.Child() {
				p := trav.NewPair(child, CheckerInfo{
					depth: inf.depth - 1,
				}, false)

				pairs = append(pairs, p)
			}
		}

		return pairs, nil
	}

	x := trav.NewTraversor(do, next)

	return func(n N, depth int) error {
		fn := x.CFSWithInfo(func(root N) (trav.Info, error) {
			return CheckerInfo{
				depth: depth,
			}, nil
		})

		return fn(n)
	}
}

func (c *Checker[N]) Reset() {
	if c == nil {
		return
	}

	if len(c.table) > 0 {
		clear(c.table)
		c.table = nil
	}
}
