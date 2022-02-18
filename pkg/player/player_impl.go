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

type PlayerState = byte
const (
	IdleState PlayerState = iota
	WalkState
	RunState
	AttackState
)

type PlayerImpl struct {
	pos    f64.Vec2
	Body   *resolv.Object
	Img    *ebi.Image
	Sprite *ase.File
	State PlayerState
}

// Positioned
func (player *PlayerImpl) Pos() f64.Vec2 { return player.pos }

func (player *PlayerImpl) SceneEnter() {
	player.State = IdleState
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

	if inpututil.IsKeyJustPressed(ebi.KeySpace) {
		player.State = AttackState
		player.Sprite.Play("SpearDown")
		return
	}
	
	currentState := IdleState

	hV := float64(0)
	vV := float64(0)

	if ebi.IsKeyPressed(ebi.KeyA) {
		hV -= 1
		currentState = WalkState
	}

	if ebi.IsKeyPressed(ebi.KeyD) {
		hV += 1
		currentState = WalkState
	}

	if ebi.IsKeyPressed(ebi.KeyW) {
		vV -= 1
		currentState = WalkState
	}

	if ebi.IsKeyPressed(ebi.KeyS) {
		vV += 1
		currentState = WalkState
	}

	if collision := player.Body.Check(hV, vV); collision != nil {
		hV, vV = 0, 0
		currentState = IdleState
	}

	player.pos[0] += hV
	player.pos[1] += vV
	player.Body.X, player.Body.Y = player.pos[0]-(frameWidth/2), player.pos[1]-(frameHeight/2)
	player.Body.Update()

	player.State = currentState
	switch currentState {
	case IdleState:
		player.Sprite.Play("IdleDown")
	case WalkState:
		player.Sprite.Play("WalkDown")
	case AttackState:
		//player.Sprite.Play("AttackDown")
	}
}

func (player *PlayerImpl) OnAnimExit() func(*ase.Tag) {
	return func(tag *ase.Tag) {
		fmt.Printf("Player: %+v - %v exited\n", player, tag.Name)
		switch tag.Name {
		case "AttackDown","AttackUp","AttackLeft","AttackRight":
			break
		}
	}
}
