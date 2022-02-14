package player

import (
	"fmt"
	"image"
	"log"

	"github.com/Racinettee/simul/pkg/component"
	ebi "github.com/hajimehoshi/ebiten/v2"
	ebiutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	ase "github.com/solarlune/goaseprite"
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
	pos    f64.Vec2
	Body   *resolv.Object
	Img    *ebi.Image
	Sprite *ase.File
}

// Positioned
func (player *PlayerImpl) Pos() f64.Vec2 { return player.pos }

func (player *PlayerImpl) SceneEnter() {
	player.pos[0] = 50
	player.pos[1] = 50
	player.Sprite = ase.Open("sprites/player/Chica.json")
	var err error
	player.Img, _, err = ebiutil.NewImageFromFile("sprites/player/" + player.Sprite.ImagePath)
	if err != nil {
		log.Println(err)
	}
	player.Body = resolv.NewObject(player.pos[0]-(frameWidth/2), player.pos[1]-(frameHeight/2), frameWidth, frameHeight)
	player.Sprite.Play("IdleDown")

	player.Sprite.OnTagExit = player.OnAnimExit()
}

// Renderable
func (player *PlayerImpl) Render(renderer component.Renderer) {
	// Character
	op := renderer.GetTranslation(player.pos[0], player.pos[1])
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	renderer.RenderItem(player.Img.SubImage(image.Rect(player.Sprite.CurrentFrameCoords())).(*ebi.Image), op)
}

// Behavior
func (player *PlayerImpl) Update(tick int) {
	player.Sprite.Update(float32(.5 / 60.0))

	hV := float64(0)
	vV := float64(0)

	if ebi.IsKeyPressed(ebi.KeyA) {
		hV -= 1
	}

	if ebi.IsKeyPressed(ebi.KeyD) {
		hV += 1
	}

	if ebi.IsKeyPressed(ebi.KeyW) {
		vV -= 1
	}

	if ebi.IsKeyPressed(ebi.KeyS) {
		vV += 1
		player.Sprite.Play("WalkDown")
	}
	if inpututil.IsKeyJustReleased(ebi.KeyS) {
		player.Sprite.Play("IdleDown")
	}

	if collision := player.Body.Check(hV, vV); collision != nil {
		hV, vV = 0, 0
	}

	player.pos[0] += hV
	player.pos[1] += vV
	player.Body.X, player.Body.Y = player.pos[0]-(frameWidth/2), player.pos[1]-(frameHeight/2)
	player.Body.Update()
}

func (player *PlayerImpl) OnAnimExit() func(*ase.Tag) {
	return func(tag *ase.Tag) {
		fmt.Printf("Player: %+v - %v exited\n", player, tag.Name)
	}
}
