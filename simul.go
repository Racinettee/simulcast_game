package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"log"

	ebi "github.com/hajimehoshi/ebiten/v2"
	ebiutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	tiled "github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"

	"github.com/Racinettee/simul/pkg/camera"
)

const (
	screenWidth  = 240
	screenHeight = 240
)

const (
	frameOX     = 0
	frameOY     = 32
	frameWidth  = 32
	frameHeight = 32
	frameNum    = 8
)

var (
	runnerImage *ebi.Image
)

var (
	tileMap  *tiled.Map
	mapImage *ebi.Image
)

func init() {
	// Decode an image from the image file's byte slice.
	// Now the byte slice is generated with //go:generate for Go 1.15 or older.
	// If you use Go 1.16 or newer, it is strongly recommended to use //go:embed to embed the image file.
	// See https://pkg.go.dev/embed for more details.
	tileMap, _ = tiled.LoadFile("tilemaps/simple_map.tmx")
	renderer, _ := render.NewRenderer(tileMap)
	renderer.RenderVisibleLayers()
	mapImage = ebi.NewImageFromImage(renderer.Result)
	renderer.Clear()
}

type Game struct {
	camera camera.Camera
	count  int
}

func (g *Game) Update() error {
	g.count++

	if ebi.IsKeyPressed(ebi.KeyA) {
		g.camera.Position[0] -= 1
	}

	if ebi.IsKeyPressed(ebi.KeyD) {
		g.camera.Position[0] += 1
	}

	if ebi.IsKeyPressed(ebi.KeyW) {
		g.camera.Position[1] -= 1
	}

	if ebi.IsKeyPressed(ebi.KeyS) {
		g.camera.Position[1] += 1
	}

	return nil
}

func (g *Game) Draw(screen *ebi.Image) {
	// Map
	g.camera.Render(mapImage, screen)
	// Character
	op := ebi.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	i := (g.count / 5) % frameNum
	sx, sy := frameOX+i*frameWidth, frameOY
	screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebi.Image), &op)

	ebiutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebi.CurrentTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	runnerImage = ebi.NewImageFromImage(img)
	g := &Game{}

	ebi.SetWindowSize(screenWidth*2, screenHeight*2)
	ebi.SetWindowTitle("Tiles (Ebiten Demo)")
	if err := ebi.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
