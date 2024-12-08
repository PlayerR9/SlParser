package internal

// ActionType is the type of an action.
type ActionType uint

const (
	// ActAccept is the action type for accepting the parse.
	ActAccept ActionType = iota

	// ActReduce is the action type for reducing the parse.
	ActReduce

	// ActShift is the action type for shifting the parse.
	ActShift
)
