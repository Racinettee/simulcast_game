package enemy

import (
	"github.com/Racinettee/simul/pkg/anim"
	"github.com/Racinettee/simul/pkg/component/state"
	"github.com/solarlune/resolv"
	"golang.org/x/image/math/f64"
)

var animationManager anim.AnimationManager = anim.NewAnimationManager()

func init() {
	animationManager.LoadAnimation("enemies/Gel")
}

type EnemyImpl struct {
	pos   f64.Vec2
	Body  *resolv.Object
	State state.Action
}

func (enemy *EnemyImpl) Pos() f64.Vec2 { return enemy.pos }
