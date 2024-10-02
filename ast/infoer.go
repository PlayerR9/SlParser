package ast

import (
	"iter"

	gers "github.com/PlayerR9/go-errors"
	"github.com/PlayerR9/go-errors/assert"
)

// Infoer is an interface for info about a node.
type Infoer[N interface {
	Child() iter.Seq[N]

	Noder
}] interface {
	// IsNil checks whether the node the info is about is nil.
	//
	// Returns:
	//   - bool: true if the node is nil, false otherwise.
	IsNil() bool

	// Node returns the node the info is about.
	//
	// Returns:
	//   - N: The node the info is about.
	Node() N

	// Frame returns an iterator that yields the frames of the node.
	//
	// Returns:
	//   - iter.Seq[string]: The iterator. Never returns nil.
	Frame() iter.Seq[string]

	// IsSeen checks whether the node has been seen.
	//
	// Returns:
	//   - bool: true if the node has been seen, false otherwise.
	IsSeen() bool

	// See marks the node as seen. Does nothing if the receiver is nil.
	See()
}

// Info is the internal implementation of Infoer.
type Info[N interface {
	Child() iter.Seq[N]

	Noder
}] struct {
	// node is the node the info is about.
	node N

	// frames are the frames of the node. Used for stack traces.
	frames []string

	// is_seen is a flag that indicates whether the node has been seen.
	is_seen bool
}

// IsNil implements the Infoer interface.
func (info *Info[N]) IsNil() bool {
	return info == nil
}

// Node implements the Infoer interface.
func (info Info[N]) Node() N {
	return info.node
}

// Frame implements the Infoer interface.
func (info Info[N]) Frame() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, frame := range info.frames {
			if !yield(frame) {
				return
			}
		}
	}
}

// IsSeen implements the Infoer interface.
func (info Info[N]) IsSeen() bool {
	return info.is_seen
}

// See implements the Infoer interface.
func (info *Info[N]) See() {
	if info == nil {
		return
	}

	info.is_seen = true
}

// NewInfo creates a new info.
//
// Parameters:
//   - node: The node the info is about.
//   - frames: The frames of the node. Used for stack traces.
//
// Returns:
//   - *Info: The new info.
//   - error: An error of type error.Err with the code errors.BadParameter if the node is nil.
//
// The info is initialized as not seen. Call See to set it as seen.
func NewInfo[N interface {
	Child() iter.Seq[N]

	Noder
}](node N, frames []string) (*Info[N], error) {
	if node.IsNil() {
		return nil, gers.NewErrNilParameter("ast.NewInfo()", "node")
	}

	info := &Info[N]{
		node:    node,
		frames:  frames,
		is_seen: false,
	}

	return info, nil
}

// AppendFrame appends a frame to the frames of the Info.
//
// Returns:
//   - []string: The new frames.
func (info Info[N]) AppendFrame() []string {
	node := info.node
	assert.Cond(!node.IsNil(), "node must not be nil")

	frame := node.String()

	if len(info.frames) == 0 {
		return []string{frame}
	}

	new_frames := make([]string, len(info.frames), len(info.frames)+1)
	copy(new_frames, info.frames)

	return append(new_frames, frame)
}

// NextInfos returns the children of the node.
//
// Returns:
//   - []*Info: The children of the node.
//   - error: An error of type error.Err with the code errors.BadParameter if the node is nil.
func (info Info[N]) NextInfos() ([]*Info[N], error) {
	var children []*Info[N]

	new_frames := info.AppendFrame()

	i := 0
	for child := range info.node.Child() {
		if child.IsNil() {
			err := gers.New(assert.AssertFail, "child found to be nil")
			err.AddContext("idx", i)

			return nil, err
		}

		new_info := &Info[N]{
			node:    child,
			frames:  new_frames,
			is_seen: false,
		}

		children = append(children, new_info)
		i++
	}

	return children, nil
}

// Child returns an iterator that yields the children of the node.
//
// Returns:
//   - iter.Seq[N]: The iterator. Never returns nil.
func (info Info[N]) Child() iter.Seq[N] {
	if info.node.IsNil() {
		return func(yield func(N) bool) {}
	}

	return info.node.Child()
}
