package component

//go:generate stringer -type=ActorState

type ActorState byte
const (
	Idle ActorState = iota
	Walk
	Run
	Attack
)
