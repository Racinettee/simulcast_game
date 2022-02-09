package main

import (
	_ "image/png"
	"log"

	ebi "github.com/hajimehoshi/ebiten/v2"

	"github.com/Racinettee/simul/pkg/game"
)

const (
	screenWidth  = 240
	screenHeight = 240
)

func main() {
	ebi.SetWindowSize(screenWidth*2, screenHeight*2)
	ebi.SetWindowTitle("Tiles (Ebiten Demo)")
	g := &game.Game{}
	g.Init()
	if err := ebi.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
