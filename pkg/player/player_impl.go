package player

import (
	"bytes"
	"image"
	"log"

	"github.com/Racinettee/simul/pkg/component"
	ebi "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/solarlune/resolv"
	"golang.org/x/image/math/f64"
)

const (
	frameOX     = 0
	frameOY     = 32
	frameWidth  = 32
	frameHeight = 32
	frameNum    = 8
)

type PlayerImpl struct {
	pos         f64.Vec2
	Body        resolv.Shape
	runnerImage *ebi.Image
	count       int
}

// Positioned
func (player *PlayerImpl) Pos() f64.Vec2 { return player.pos }

func (player *PlayerImpl) SceneEnter() {
	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	player.runnerImage = ebi.NewImageFromImage(img)
	player.Body = resolv.NewRectangle(player.pos[0], player.pos[0], frameWidth, frameHeight)
}

// Renderable
func (player *PlayerImpl) Render(renderer component.Renderer) {
	// Character
	op := renderer.GetTranslation(player.pos[0], player.pos[1])
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	i := (player.count / 5) % frameNum
	sx, sy := frameOX+i*frameWidth, frameOY
	renderer.RenderItem(player.runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebi.Image), op)
}

// Behavior
func (player *PlayerImpl) Update(tick int) {
	player.count = tick
	if ebi.IsKeyPressed(ebi.KeyA) {
		player.pos[0] -= 1
	}

	if ebi.IsKeyPressed(ebi.KeyD) {
		player.pos[0] += 1
	}

	if ebi.IsKeyPressed(ebi.KeyW) {
		player.pos[1] -= 1
	}

	if ebi.IsKeyPressed(ebi.KeyS) {
		player.pos[1] += 1
	}
}
