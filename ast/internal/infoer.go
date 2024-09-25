package internal

import "iter"

// Infoer is an interface for info about a node.
type Infoer[N interface {
	Child() iter.Seq[N]

	Noder
}] interface {
	// Init initializes the info.
	//
	// Parameters:
	//   - node: The node the info is about.
	//   - frames: The frames of the node. Used for stack traces.
	//
	// If length of frames is 0, then it is the first call to Init.
	Init(node N, frames []string)

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

// _Info is the internal implementation of Infoer.
type _Info[N interface {
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
func (info _Info[N]) IsNil() bool {
	return info.node.IsNil()
}

// Node implements the Infoer interface.
func (info _Info[N]) Node() N {
	return info.node
}

// Frame implements the Infoer interface.
func (info _Info[N]) Frame() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, frame := range info.frames {
			if !yield(frame) {
				return
			}
		}
	}
}

// Init implements the Infoer interface.
func (info *_Info[N]) Init(node N, frames []string) {
	if info == nil {
		return
	}

	info.node = node
	info.frames = frames
	info.is_seen = false
}

// AppendFrame appends a frame to the frames of the Info.
//
// Parameters:
//   - frame: The frame to append.
//
// Returns:
//   - []string: The new frames.
func (info _Info[N]) AppendFrame(frame string) []string {
	if len(info.frames) == 0 {
		return []string{frame}
	}

	new_frames := make([]string, len(info.frames), len(info.frames)+1)
	copy(new_frames, info.frames)

	return append(new_frames, frame)
}

// IsSeen implements the Infoer interface.
func (info _Info[N]) IsSeen() bool {
	return info.is_seen
}

// See implements the Infoer interface.
func (info *_Info[N]) See() {
	if info == nil {
		return
	}

	info.is_seen = true
}

// NextInfos returns the children of the node.
//
// Returns:
//   - []*_Info: The children of the node.
func (info *_Info[N]) NextInfos() []*_Info[N] {
	if info == nil {
		return nil
	}

	var children []*_Info[N]

	new_frames := info.AppendFrame(info.node.String())

	for child := range info.node.Child() {
		var new_info _Info[N]

		new_info.Init(child, new_frames)

		children = append(children, &new_info)
	}

	return children
}
