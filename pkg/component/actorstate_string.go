// Code generated by "stringer -type=ActorState"; DO NOT EDIT.

package component

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

const _ActorState_name = "IdleWalkRunAttack"

var _ActorState_index = [...]uint8{0, 4, 8, 11, 17}

func (i ActorState) String() string {
	if i >= ActorState(len(_ActorState_index)-1) {
		return "ActorState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ActorState_name[_ActorState_index[i]:_ActorState_index[i+1]]
}
