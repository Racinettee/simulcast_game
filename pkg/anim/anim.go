package anim

import (
	"image"
	"log"

	ebi "github.com/hajimehoshi/ebiten/v2"
	ebiutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	ase "github.com/solarlune/goaseprite"
)

type Animation struct {
	Name   string
	Img    *ebi.Image
	Sprite *ase.File
}

type AnimationManager struct {
	Animation        map[string]Animation
	CurrentAnimation Animation
	OnAnimationExit  func(data Animation)
}

func (animManager *AnimationManager) LoadAnimation(texture, name string) {
	img, _, err := ebiutil.NewImageFromFile("sprites/" + texture + ".png")
	if err != nil {
		log.Println(err)
	}
	animation := Animation{
		Name:   name,
		Sprite: ase.Open("sprites/" + texture + ".json"),
		Img:    img,
	}
	animManager.Animation[name] = animation
	animation.Sprite.OnFrameChange = animManager.onAnimationExit()
}

func (animManager AnimationManager) GetSprite(sprite string) *ase.File {
	return animManager.Animation[sprite].Sprite
}

func (animManager AnimationManager) GetImage(sprite string) *ebi.Image {
	return animManager.Animation[sprite].Img
}

func (animManager AnimationManager) GetAnimation(sprite string) Animation {
	return animManager.Animation[sprite]
}

func (animManager *AnimationManager) SetCurrentAnimation(sprite string) bool {
	newAnim, ok := animManager.Animation[sprite]
	if ok {
		animManager.CurrentAnimation = newAnim
	}
	return ok
}

func (animManager *AnimationManager) PlayCurrent() {
	animManager.CurrentAnimation.Sprite.Play("")
}

func (animManager *AnimationManager) Play(nextAnim string) {
	animManager.SetCurrentAnimation(nextAnim)
	animManager.PlayCurrent()
}

func (animManager *AnimationManager) PlayTag(tag string) {
	animManager.CurrentAnimation.Sprite.Play(tag)
}

func (animManager *AnimationManager) CurrentFrame() *ebi.Image {
	return animManager.CurrentAnimation.Img.SubImage(image.Rect(
		animManager.CurrentAnimation.Sprite.CurrentFrameCoords())).(*ebi.Image)
}

func (animManager *AnimationManager) Update() {
	animManager.CurrentAnimation.Sprite.Update(1 / 60.0)
}

func (animManager *AnimationManager) onAnimationExit() func() {
	return func() {
		currIndex := animManager.CurrentAnimation.Sprite.FrameIndex
		prevIndex := animManager.CurrentAnimation.Sprite.PrevFrameIndex
		currTag := animManager.CurrentAnimation.Sprite.CurrentTag

		if currIndex == currTag.Start && prevIndex == currTag.End {
			if animManager.OnAnimationExit != nil {
				animManager.OnAnimationExit(animManager.CurrentAnimation)
			}
		}
	}
}

func NewAnimationManager() AnimationManager {
	return AnimationManager{
		Animation: make(map[string]Animation),
	}
}
