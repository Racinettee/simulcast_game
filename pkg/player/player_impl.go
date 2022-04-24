package player

import (
	//"fmt"

	"github.com/Racinettee/simul/pkg/anim"
	comp "github.com/Racinettee/simul/pkg/component"
	ebiutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

// The players private persistent animation manager
var animationManager anim.AnimationManager = anim.NewAnimationManager()

func init() {
	animationManager.LoadAnimation("player/Chica-SpearDown", "SpearDown")
	animationManager.LoadAnimation("player/Chica-Idle-Down", "IdleDown")
	animationManager.LoadAnimation("player/Chica-WalkDown", "WalkDown")
}

type PlayerImpl struct {
	pos   f64.Vec2
	Body  *resolv.Object
	State comp.ActorState
	Dir   comp.Direction
}

// Positioned
func (player *PlayerImpl) Pos() f64.Vec2 { return player.pos }

func (player *PlayerImpl) SceneEnter() {
	player.State = comp.Idle
	player.Dir = comp.Down
	player.pos[0] = 50
	player.pos[1] = 50
	animationManager.OnAnimationExit = player.OnAnimExit()
	animationManager.SetCurrentAnimation("IdleDown")
	animationManager.PlayCurrent()

	player.Body = resolv.NewObject(player.pos[0]-(frameWidth/2), player.pos[1]-(frameHeight/2), frameWidth, frameHeight)
}

// Renderable
func (player *PlayerImpl) Render(renderer comp.Renderer) {
	// Character
	op := renderer.GetTranslation(player.pos[0], player.pos[1])
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	renderer.RenderItem(animationManager.CurrentFrame(), op)
	ebiutil.DebugPrint(renderer.ScreenSurface(), player.StateStr())
}

func (player *PlayerImpl) OnAnimExit() func(anim anim.Animation) {
	return func(anim anim.Animation) {
		switch anim.Name {
		case "SpearDown":
			player.State = comp.Idle
		case "AttackDown", "AttackUp", "AttackLeft", "AttackRight":
			break
		}
	}
}

func (p *PlayerImpl) StateStr() string {
	return p.State.String() + p.Dir.String()
}
