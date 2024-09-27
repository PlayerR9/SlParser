package ast

import (
	"errors"
	"iter"

	gers "github.com/PlayerR9/go-errors"
)

// CheckerInfo is a struct that contains information about a node.
type CheckerInfo[N interface {
	Child() iter.Seq[N]

	Noder
}] struct {
	info *Info[N]

	// depth is the depth of the node.
	depth int
}

// IsNil implements the Infoer interface.
func (ci *CheckerInfo[N]) IsNil() bool {
	return ci == nil
}

// IsSeen implements the Infoer interface.
func (ci *CheckerInfo[N]) IsSeen() bool {
	gers.AssertNotNil(ci.info, "ci.info")

	return ci.info.IsSeen()
}

// Node implements the Infoer interface.
func (ci CheckerInfo[N]) Node() N {
	gers.AssertNotNil(ci.info, "ci.info")

	return ci.info.Node()
}

// See implements the Infoer interface.
func (ci *CheckerInfo[N]) See() {
	gers.AssertNotNil(ci.info, "ci.info")

	ci.info.See()
}

// Frame implements the Infoer interface.
func (ci CheckerInfo[N]) Frame() iter.Seq[string] {
	gers.AssertNotNil(ci.info, "ci.info")

	return ci.info.Frame()
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
}](node N, frames []string, depth int) (*CheckerInfo[N], error) {
	if node.IsNil() {
		return nil, gers.NewErrNilParameter("node")
	}

	sub_info, err := NewInfo(node, frames)
	if err != nil {
		return nil, err
	}

	return &CheckerInfo[N]{
		info:  sub_info,
		depth: depth,
	}, nil
}

// Children returns the children of the node.
//
// Returns:
//   - []CheckerInfo[N]: The children of the node.
func (ci *CheckerInfo[N]) NextInfos() ([]*CheckerInfo[N], error) {
	if ci == nil {
		return nil, errors.New("receiver is nil")
	} else if ci.depth == 0 {
		return nil, nil
	}

	var new_depth int

	if ci.depth < 0 {
		new_depth = -1
	} else {
		new_depth = ci.depth - 1
	}

	new_frames := ci.info.AppendFrame()

	var children []*CheckerInfo[N]

	for child := range ci.info.Child() {
		if child.IsNil() {
			return nil, gers.NewErrAssertFail("child is nil")
		}

		next, err := NewCheckerInfo(child, new_frames, new_depth)
		if err != nil {
			return nil, err
		}

		children = append(children, next)
	}

	return children, nil
}
