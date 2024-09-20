package internal

//go:generate stringer -type=ActionType -linecomment

type ActionType int

const (
	ActShift  ActionType = iota // SHIFT
	ActReduce                   // REDUCE
	ActAccept                   // ACCEPT
)
