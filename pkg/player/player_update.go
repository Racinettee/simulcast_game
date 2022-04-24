package player

import (
	"github.com/Racinettee/simul/pkg/component/state"
	ebi "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Behavior
func (player *PlayerImpl) Update(tick int) {
	animationManager.Update() //float32(.5 / 60.0))

	// If player is already in the attacking state then
	// we want to block movements or other actions being performed
	if player.State == state.Attack {
		return
	}
	// Place player into attacking state
	if inpututil.IsKeyJustPressed(ebi.KeySpace) {
		player.State = state.Attack
		animationManager.Play("SpearDown")
		return
	}

	currentState := state.Idle

	hV := float64(0)
	vV := float64(0)

	// Move the player, and assign walking state
	switch {
	case ebi.IsKeyPressed(ebi.KeyA):
		hV -= 1
		currentState = state.Walk
		player.Dir = state.Left
	case ebi.IsKeyPressed(ebi.KeyD):
		hV += 1
		currentState = state.Walk
		player.Dir = state.Right
	}
	switch {
	case ebi.IsKeyPressed(ebi.KeyW):
		vV -= 1
		currentState = state.Walk
		player.Dir = state.Up
	case ebi.IsKeyPressed(ebi.KeyS):
		vV += 1
		currentState = state.Walk
		player.Dir = state.Down
	}
	// Check collision with terrain objects
	if collision := player.Body.Check(hV, vV); collision != nil {
		hV, vV = 0, 0
		currentState = state.Idle
	}

	player.pos[0] += hV
	player.pos[1] += vV
	player.Body.X, player.Body.Y = player.pos[0]-(frameWidth/2), player.pos[1]-(frameHeight/2)
	player.Body.Update()

	// Play the correct animation
	player.State = currentState
	switch currentState {
	case state.Idle:
		animationManager.Play("IdleDown")
	case state.Walk:
		animationManager.Play("WalkDown")
	}
}
