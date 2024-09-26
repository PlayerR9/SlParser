package ast

import (
	"iter"
)

// CheckerInfo is a struct that contains information about a node.
type CheckerInfo[N interface {
	Child() iter.Seq[N]

	Noder
}] struct {
	*Info[N]

	// depth is the depth of the node.
	depth int
}

func (ci *CheckerInfo[N]) Init(node N, frames []string) {
	if ci == nil {
		return
	}

	if ci.Info == nil {
		var new_info Info[N]
		ci.Info = &new_info
	}

	ci.Info.Init(node, frames)
}

// NewCheckerInfo creates a new CheckerInfo.
//
// Parameters:
//   - node: The node.
//   - depth: The depth of the node.
//   - frames: The frames of the node.
//
// Returns:
//   - CheckerInfo: The created CheckerInfo.
func NewCheckerInfo[N interface {
	Child() iter.Seq[N]

	Noder
}](depth int) CheckerInfo[N] {
	return CheckerInfo[N]{
		depth: depth,
	}
}

// Children returns the children of the node.
//
// Returns:
//   - []CheckerInfo[N]: The children of the node.
func (ci *CheckerInfo[N]) NextInfos() []*CheckerInfo[N] {
	if ci == nil || ci.depth == 0 {
		return nil
	}

	var new_depth int

	if ci.depth < 0 {
		new_depth = -1
	} else {
		new_depth = ci.depth - 1
	}

	new_frames := ci.AppendFrame()

	var children []*CheckerInfo[N]

	for child := range ci.node.Child() {
		sub_info := NewCheckerInfo[N](new_depth)
		sub_info.Info.Init(child, new_frames)

		children = append(children, &sub_info)
	}

	return children
}
