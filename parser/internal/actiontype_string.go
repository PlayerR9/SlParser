// Code generated by "stringer -type=ActionType -linecomment"; DO NOT EDIT.

package internal

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ActAccept-0]
	_ = x[ActReduce-1]
	_ = x[ActShift-2]
}

const _ActionType_name = "(ACCEPT)(REDUCE)(SHIFT)"

var _ActionType_index = [...]uint8{0, 8, 16, 23}

func (i ActionType) String() string {
	if i < 0 || i >= ActionType(len(_ActionType_index)-1) {
		return "ActionType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ActionType_name[_ActionType_index[i]:_ActionType_index[i+1]]
}
