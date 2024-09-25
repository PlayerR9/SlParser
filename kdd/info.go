package kdd

import (
	gers "github.com/PlayerR9/go-errors"
)

type Infoer interface {
}

type DetermineInfoFn[I Infoer] func(node *Node) (I, error)

func MakeInfoTable[I Infoer](root *Node, fn DetermineInfoFn[I]) (map[*Node]I, error) {
	if root == nil || fn == nil {
		return nil, nil
	}

	stack := []*Node{root}

	table := make(map[*Node]I)

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		gers.AssertNotNil(top, "top")

		info, err := fn(top)
		if err != nil {
			return nil, err
		}

		_, ok := table[top]
		gers.Assert(!ok, "node already in the table")

		table[top] = info

		for child := range top.Child() {
			stack = append(stack, child)
		}
	}

	return table, nil
}
