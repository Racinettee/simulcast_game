package player

import (
	comp "github.com/Racinettee/simul/pkg/component"
	ebi "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Behavior
func (player *PlayerImpl) Update(tick int) {
	animationManager.Update() //float32(.5 / 60.0))

	// If player is already in the attacking state then
	// we want to block movements or other actions being performed
	if player.State == comp.Attack {
		return
	}
	// Place player into attacking state
	if inpututil.IsKeyJustPressed(ebi.KeySpace) {
		player.State = comp.Attack
		animationManager.Play("SpearDown")
		return
	}

	currentState := comp.Idle

	hV := float64(0)
	vV := float64(0)

	// Move the player, and assign walking state
	if ebi.IsKeyPressed(ebi.KeyA) {
		hV -= 1
		currentState = comp.Walk
		player.Dir = comp.Left
	}

	if ebi.IsKeyPressed(ebi.KeyD) {
		hV += 1
		currentState = comp.Walk
		player.Dir = comp.Right
	}

	if ebi.IsKeyPressed(ebi.KeyW) {
		vV -= 1
		currentState = comp.Walk
		player.Dir = comp.Up
	}

	if ebi.IsKeyPressed(ebi.KeyS) {
		vV += 1
		currentState = comp.Walk
		player.Dir = comp.Down
	}
	// Check collision with terrain objects
	if collision := player.Body.Check(hV, vV); collision != nil {
		hV, vV = 0, 0
		currentState = comp.Idle
	}

	player.pos[0] += hV
	player.pos[1] += vV
	player.Body.X, player.Body.Y = player.pos[0]-(frameWidth/2), player.pos[1]-(frameHeight/2)
	player.Body.Update()

	// Play the correct animation
	player.State = currentState
	switch currentState {
	case comp.Idle:
		animationManager.Play("IdleDown")
	case comp.Walk:
		animationManager.Play("WalkDown")
	}
}
