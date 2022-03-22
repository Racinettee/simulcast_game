package main

import (
	_ "image/png"
	"log"

	ebi "github.com/hajimehoshi/ebiten/v2"

	"github.com/Racinettee/simul/pkg/game"
)

func main() {
	g := &game.Game{}
	g.Init()
	if err := ebi.RunGame(g); err != nil {
		log.Println(err)
	}
	g.Shutdown()
}
