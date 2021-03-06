// Code generated by "stringer -type=Action"; DO NOT EDIT.

package state

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Idle-0]
	_ = x[Walk-1]
	_ = x[Run-2]
	_ = x[Attack-3]
}

const _Action_name = "IdleWalkRunAttack"

var _Action_index = [...]uint8{0, 4, 8, 11, 17}

func (i Action) String() string {
	if i >= Action(len(_Action_index)-1) {
		return "Action(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Action_name[_Action_index[i]:_Action_index[i+1]]
}
