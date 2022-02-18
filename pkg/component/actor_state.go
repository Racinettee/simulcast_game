package component

//go:generate stringer -type=ActorState

type ActorState uint16
const (
	IdleState ActorState = iota
	WalkState
	RunState
	AttackState
)
