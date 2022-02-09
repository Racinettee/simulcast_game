package component

import "golang.org/x/image/math/f64"

type Positioned interface {
	Pos() f64.Vec2
}
