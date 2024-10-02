package ast

import (
	"errors"
	"fmt"
	"iter"
	"slices"

	gers "github.com/PlayerR9/go-errors"
	"github.com/PlayerR9/go-errors/assert"
)

var (
	// IgnoreInfo is an error that is used to ignore info. Readers must return
	// this error as is and not wrap it as callers check for this error using
	// ==.
	IgnoreInfo error
)

func init() {
	IgnoreInfo = errors.New("ignore info")
}

// NewInfoFn is a function that creates an info for a node.
//
// Parameters:
//   - node: The node to create the info for.
//
// Returns:
//   - I: The info.
//   - error: An error if the info could not be created.
//
// Return the error IgnoreInfo to not add the info to the table.
type NewInfoFn[N interface {
	Child() iter.Seq[N]

	Noder
}, I interface {
	NextInfos() ([]I, error)

	Infoer[N]
}] func(node N) (I, error)

// InfoTableOfFn is a function that creates an info table for a node.
//
// Parameters:
//   - root: The root node.
//
// Returns:
//   - map[N]I: The info table.
//   - error: An error if the info table could not be created.
type InfoTableOfFn[N interface {
	Child() iter.Seq[N]

	Noder
}, I interface {
	NextInfos() ([]I, error)

	Infoer[N]
}] func(root N) (map[N]I, error)

type InfoTableMaker[N interface {
	Child() iter.Seq[N]

	Noder
}, I interface {
	NextInfos() ([]I, error)

	Infoer[N]
}] struct {
	InitFn     func(node N, frames []string) (I, error)
	MakeInfoFn func(node N) (I, error)
}

func (itm InfoTableMaker[N, I]) Apply(root N) (map[N]I, error) {
	if root.IsNil() {
		err := gers.NewErrNilParameter("InfoTableMaker.Apply()", "root")
		return nil, err
	}

	table := make(map[N]I)

	root_info, err := itm.InitFn(root, nil)
	if err != nil {
		return nil, err
	} else if root_info.IsNil() {
		return nil, fmt.Errorf("root info is nil")
	}

	stack := []I{root_info}

	var inner_err error
	var last_top I

	for len(stack) > 0 && inner_err == nil {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		last_top = top

		if top.IsNil() {
			inner_err = errors.New("found info to be nil")
			continue
		}

		node := top.Node()
		if node.IsNil() {
			inner_err = errors.New("found node to be nil")
			continue
		}

		info, err := itm.MakeInfoFn(node)
		if err == nil {
			table[node] = info
		} else if err != IgnoreInfo {
			inner_err = err
			continue
		}

		nexts, err := top.NextInfos()
		if err != nil {
			inner_err = err
			continue
		}

		if len(nexts) > 0 {
			slices.Reverse(nexts)
			stack = append(stack, nexts...)
		}
	}

	if inner_err == nil {
		return table, nil
	}

	assert.NotNil(last_top, "last_top")

	outer_err := gers.New(BadSyntaxTree, inner_err.Error())

	for frame := range last_top.Frame() {
		outer_err.AddFrame(frame)
	}

	return nil, outer_err
}
