package state

//go:generate stringer -type=Action
//go:generate stringer -type=Direction

type Action byte

const (
	Idle Action = iota
	Walk
	Run
	Attack
)

type Direction byte

const (
	Up Direction = iota
	Right
	UpRight
	RightDown
	Down
	DownLeft
	Left
	UpLeft
)
