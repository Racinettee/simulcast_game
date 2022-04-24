package player

import (
	//"fmt"
	"image"
	"log"

	comp "github.com/Racinettee/simul/pkg/component"
	ebi "github.com/hajimehoshi/ebiten/v2"
	ebiutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	State  comp.ActorState
	Dir    comp.Direction
}

// Positioned
func (player *PlayerImpl) Pos() f64.Vec2 { return player.pos }

func (player *PlayerImpl) SceneEnter() {
	player.State = comp.Idle
	player.Dir = comp.Down
	player.pos[0] = 50
	player.pos[1] = 50
	player.Sprite = ase.Open("sprites/player/Chica.json")
	var err error
	player.Img, _, err = ebiutil.NewImageFromFile("sprites/player/" + player.Sprite.ImagePath)
	if err != nil {
		log.Println(err)
	}
	player.Body = resolv.NewObject(player.pos[0]-(frameWidth/2), player.pos[1]-(frameHeight/2), frameWidth, frameHeight)
	player.Sprite.Play(player.StateStr())

	player.Sprite.OnFrameChange = player.OnAnimExit()
}

// Renderable
func (player *PlayerImpl) Render(renderer comp.Renderer) {
	// Character
	op := renderer.GetTranslation(player.pos[0], player.pos[1])
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	renderer.RenderItem(player.Img.SubImage(image.Rect(player.Sprite.CurrentFrameCoords())).(*ebi.Image), op)
	ebiutil.DebugPrint(renderer.ScreenSurface(), player.StateStr())
}

func (player *PlayerImpl) OnAnimExit() func() {
	return func() {
		currIndex := player.Sprite.FrameIndex
		prevIndex := player.Sprite.PrevFrameIndex
		currTag := player.Sprite.CurrentTag

		if currIndex == currTag.Start && prevIndex == currTag.End {
			//fmt.Printf("Player: %+v - %v exited\n", player, currTag.Name)
			switch currTag.Name {
			case "SpearDown":
				player.State = comp.Idle
			case "AttackDown", "AttackUp", "AttackLeft", "AttackRight":
				break
			}
		}
	}
}

func (p *PlayerImpl) StateStr() string {
	return p.State.String() + p.Dir.String()
}
