package camera

import (
	"github.com/Racinettee/simul/pkg/component"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/math/f64"

	camera "github.com/melonfunction/ebiten-camera"
)

type FollowCam struct {
	camera.Camera
	Followee component.Positioned
}

func NewFollowCam(w, h int, x, y, zoom, rotation float64) FollowCam {
	return FollowCam{
		Camera: *camera.NewCamera(w, h, x, y, rotation, zoom),
	}
}

func (c *FollowCam) Pos() f64.Vec2 {
	return f64.Vec2{c.X, c.Y}
}

func (c *FollowCam) Update(int) {
	pos := c.Followee.Pos()
	c.SetPosition(pos[0], pos[1])
}

func (c *FollowCam) RenderItem(item *ebiten.Image, ops *ebiten.DrawImageOptions) {
	c.Surface.DrawImage(item, ops)
}
