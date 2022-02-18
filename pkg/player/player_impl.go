package player

import (
	"fmt"
	"image"
	"log"

	comp "github.com/Racinettee/simul/pkg/component"
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
	State comp.ActorState
}

// Positioned
func (player *PlayerImpl) Pos() f64.Vec2 { return player.pos }

func (player *PlayerImpl) SceneEnter() {
	player.State = comp.IdleState
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

	player.Sprite.OnFrameChange = player.OnAnimExit()
	//player.Sprite.OnTagExit = player.OnAnimExit()
}

// Renderable
func (player *PlayerImpl) Render(renderer comp.Renderer) {
	// Character
	op := renderer.GetTranslation(player.pos[0], player.pos[1])
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	renderer.RenderItem(player.Img.SubImage(image.Rect(player.Sprite.CurrentFrameCoords())).(*ebi.Image), op)
	ebiutil.DebugPrint(renderer.ScreenSurface(), player.State.String())
}

// Behavior
func (player *PlayerImpl) Update(tick int) {
	player.Sprite.Update((1/60.0)) //float32(.5 / 60.0))

	// If player is already in the attacking state then
	// we want to block movements or other actions being performed
	if player.State == comp.AttackState {
		return
	}
	// Place player into attacking state
	if inpututil.IsKeyJustPressed(ebi.KeySpace) {
		player.State = comp.AttackState
		player.Sprite.Play("SpearDown")
		return
	}

	currentState := comp.IdleState

	hV := float64(0)
	vV := float64(0)

	// Move the player, and assign walking state
	if ebi.IsKeyPressed(ebi.KeyA) {
		hV -= 1
		currentState = comp.WalkState
	}

	if ebi.IsKeyPressed(ebi.KeyD) {
		hV += 1
		currentState = comp.WalkState
	}

	if ebi.IsKeyPressed(ebi.KeyW) {
		vV -= 1
		currentState = comp.WalkState
	}

	if ebi.IsKeyPressed(ebi.KeyS) {
		vV += 1
		currentState = comp.WalkState
	}
	// Check collision with terrain objects
	if collision := player.Body.Check(hV, vV); collision != nil {
		hV, vV = 0, 0
		currentState = comp.IdleState
	}

	player.pos[0] += hV
	player.pos[1] += vV
	player.Body.X, player.Body.Y = player.pos[0]-(frameWidth/2), player.pos[1]-(frameHeight/2)
	player.Body.Update()

	// Play the correct animation
	player.State = currentState
	switch currentState {
	case comp.IdleState:
		player.Sprite.Play("IdleDown")
	case comp.WalkState:
		player.Sprite.Play("WalkDown")
	}
}

func (player *PlayerImpl) OnAnimExit() func() {
	return func() {
		currIndex := player.Sprite.FrameIndex
		prevIndex := player.Sprite.PrevFrameIndex
		currTag := player.Sprite.CurrentTag

		if currIndex == currTag.Start && prevIndex == currTag.End {
			fmt.Printf("Player: %+v - %v exited\n", player, currTag.Name)
			switch currTag.Name {
			case "SpearDown":
				player.State = comp.IdleState
			case "AttackDown","AttackUp","AttackLeft","AttackRight":
				break
			}
		}
	}
}
