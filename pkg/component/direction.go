package component

//go:generate stringer -type=Direction

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
