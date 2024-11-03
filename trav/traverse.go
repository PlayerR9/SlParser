package trav

import (
	"slices"

	"github.com/PlayerR9/mygo-lib/common"
)

type TravFn[N Node] func(root N) error

type DoFn[N Node] func(node N, info Info) error

type NextFn[N Node] func(node N, info Info) ([]Pair[N], error)

type InfoFn[N Node] func(node N) (Info, error)

type Traversor[N Node] struct {
	do   DoFn[N]
	next NextFn[N]
}

func NewTraversor[N Node](do DoFn[N], next NextFn[N]) Traversor[N] {
	if do == nil {
		do = func(_ N, _ Info) error {
			return nil
		}
	}

	if next == nil {
		next = func(_ N, _ Info) ([]Pair[N], error) {
			return nil, nil
		}
	}

	return Traversor[N]{
		do:   do,
		next: next,
	}
}

// Depth-First Search (DFS)
func (trav Traversor[N]) DFSWithInfo(new_info InfoFn[N]) TravFn[N] {
	if new_info == nil {
		new_info = func(_ N) (Info, error) {
			return nil, common.NewErrNilParam("new_info")
		}
	}

	return func(root N) error {
		info, err := new_info(root)
		if err != nil {
			return err
		}

		p := NewPair(root, info, true)

		stack := []Pair[N]{p}

		for len(stack) > 0 {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			err = trav.do(top.node, top.info)
			if err != nil {
				return err
			}

			next, err := trav.next(top.node, top.info)
			if err != nil {
				return err
			}

			slices.Reverse(next)
			stack = append(stack, next...)
		}

		return nil
	}
}

// Depth-First Search (DFS)
func (trav Traversor[N]) DFS() TravFn[N] {
	return func(root N) error {
		stack := []N{root}

		for len(stack) > 0 {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			err := trav.do(top, nil)
			if err != nil {
				return err
			}

			next, err := trav.next(top, nil)
			if err != nil {
				return err
			}

			nexts := make([]N, 0, len(next))

			for _, p := range next {
				nexts = append(nexts, p.node)
			}

			slices.Reverse(nexts)

			stack = append(stack, nexts...)
		}

		return nil
	}
}

// Breadth-First Search (BFS)
func (trav Traversor[N]) BFSWithInfo(new_info InfoFn[N]) TravFn[N] {
	if new_info == nil {
		new_info = func(_ N) (Info, error) {
			return nil, common.NewErrNilParam("new_info")
		}
	}

	return func(root N) error {
		info, err := new_info(root)
		if err != nil {
			return err
		}

		p := NewPair(root, info, true)

		queue := []Pair[N]{p}

		for len(queue) > 0 {
			first := queue[0]
			queue = queue[1:]

			err = trav.do(first.node, first.info)
			if err != nil {
				return err
			}

			next, err := trav.next(first.node, first.info)
			if err != nil {
				return err
			}

			queue = append(queue, next...)
		}

		return nil
	}
}

// Breadth-First Search (BFS)
func (trav Traversor[N]) BFS() TravFn[N] {
	return func(root N) error {
		queue := []N{root}

		for len(queue) > 0 {
			first := queue[0]
			queue = queue[1:]

			err := trav.do(first, nil)
			if err != nil {
				return err
			}

			next, err := trav.next(first, nil)
			if err != nil {
				return err
			}

			nexts := make([]N, 0, len(next))

			for _, p := range next {
				nexts = append(nexts, p.node)
			}

			queue = append(queue, nexts...)
		}

		return nil
	}
}

// Children-First Search (CFS)
func (trav Traversor[N]) CFSWithInfo(new_info InfoFn[N]) TravFn[N] {
	if new_info == nil {
		new_info = func(_ N) (Info, error) {
			return nil, common.NewErrNilParam("new_info")
		}
	}

	fn := func(root N) error {
		type StackElem struct {
			p    Pair[N]
			seen bool
		}

		info, err := new_info(root)
		if err != nil {
			return err
		}

		elem := &StackElem{
			p:    NewPair(root, info, true),
			seen: false,
		}

		stack := []*StackElem{elem}

		for len(stack) > 0 {
			top := stack[len(stack)-1]

			// assert.Cond(top != nil, "top must not be nil")

			node := top.p.node

			if top.seen {
				stack = stack[:len(stack)-1]

				err = trav.do(node, top.p.info)
				if err != nil {
					if top.p.is_critical {
						return err
					} else {

					}

					// We have two cases here:
					// 1. The current node's error is critical.
					// 2. The current node's error is not critical.

					// If it is not critical, we have to add to the parent's error list.
					// If it is critical, we have to return the error immediately.

					return err
				}
			} else {
				top.seen = true

				nexts, err := trav.next(top.p.node, top.p.info)
				if err != nil {
					return err
				}

				elems := make([]*StackElem, 0, len(nexts))

				for _, next := range nexts {
					elems = append(elems, &StackElem{
						p:    next,
						seen: false,
					})
				}

				slices.Reverse(elems)

				stack = append(stack, elems...)
			}

		}

		return nil
	}

	return fn
}

// Children-First Search (CFS)
func (trav Traversor[N]) CFS() TravFn[N] {
	fn := func(root N) error {
		type StackElem struct {
			node N
			seen bool
		}

		elem := &StackElem{
			node: root,
			seen: false,
		}

		stack := []*StackElem{elem}

		for len(stack) > 0 {
			top := stack[len(stack)-1]

			// assert.Cond(top != nil, "top must not be nil")

			node := top.node

			if top.seen {
				stack = stack[:len(stack)-1]

				err := trav.do(node, nil)
				if err != nil {
					return err
				}
			} else {
				top.seen = true

				nexts, err := trav.next(top.node, nil)
				if err != nil {
					return err
				}

				elems := make([]*StackElem, 0, len(nexts))

				for _, next := range nexts {
					elems = append(elems, &StackElem{
						node: next.node,
						seen: false,
					})
				}

				slices.Reverse(elems)

				stack = append(stack, elems...)
			}

		}

		return nil
	}

	return fn
}
