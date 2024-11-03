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

func CheckNode[N interface {
	~int

	String() string
}](kind string, node interface{ GetType() N }, types ...N) error {
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
	GetType() T
	Child() iter.Seq[N]

	Noder
}, T NodeType] func(n N, depth int) error

type CheckNodeFn[N interface {
	GetType() T
	Child() iter.Seq[N]

	Noder
}, T NodeType] func(n N, children []N) error

type Checker[N interface {
	GetType() T
	Child() iter.Seq[N]

	Noder
}, T NodeType] struct {
	table map[T]CheckNodeFn[N, T]
}

func NewChecker[N interface {
	GetType() T
	Child() iter.Seq[N]

	Noder
}, T NodeType]() *Checker[N, T] {
	return &Checker[N, T]{
		table: make(map[T]CheckNodeFn[N, T]),
	}
}

func (c *Checker[N, T]) Register(type_ T, fn CheckNodeFn[N, T]) error {
	if fn == nil {
		return nil
	} else if c == nil {
		return common.ErrNilReceiver
	}

	// assert.Cond(c.table != nil, "c.table must not be nil")

	c.table[type_] = fn

	return nil
}

func (c Checker[N, T]) Build() CheckFn[N, T] {
	if len(c.table) == 0 {
		return func(n N, depth int) error {
			if n.IsNil() {
				return common.ErrNilReceiver
			}

			type_ := n.GetType()

			return fmt.Errorf("node type (%q) is not supported", type_.String())
		}
	}

	table := make(map[T]CheckNodeFn[N, T], len(c.table))
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
			return fmt.Errorf("node type (%q) is not supported", type_.String())
		}

		err := fn(node, slices.Collect(node.Child()))
		if err != nil {
			return fmt.Errorf("failed to check node (%q): %w", type_.String(), err)
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

		critical := node.AreChildrenCritical()

		if inf.depth < 0 {
			for child := range node.Child() {
				p := trav.NewPair(child, CheckerInfo{
					depth: inf.depth - 1,
				}, critical)

				pairs = append(pairs, p)
			}
		} else if inf.depth != 0 {
			var pairs []trav.Pair[N]

			for child := range node.Child() {
				p := trav.NewPair(child, CheckerInfo{
					depth: inf.depth - 1,
				}, critical)

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

func (c *Checker[N, T]) Reset() {
	if c == nil {
		return
	}

	if len(c.table) > 0 {
		clear(c.table)
		c.table = nil
	}
}
