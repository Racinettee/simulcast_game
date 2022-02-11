package game

import (
	"fmt"

	ebi "github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
	"github.com/solarlune/resolv"

	"github.com/Racinettee/simul/pkg/camera"
	"github.com/Racinettee/simul/pkg/player"
	"github.com/Racinettee/simul/pkg/tiles"
)

const (
	screenWidth  = 240
	screenHeight = 240
)

var (
	tileMap  *tiled.Map
	mapImage *ebi.Image

	// Collision stuff
	space *resolv.Space
)

func init() {
	// Decode an image from the image file's byte slice.
	// Now the byte slice is generated with //go:generate for Go 1.15 or older.
	// If you use Go 1.16 or newer, it is strongly recommended to use //go:embed to embed the image file.
	// See https://pkg.go.dev/embed for more details.
	tileMap, _ = tiled.LoadFile("tilemaps/simple_map.tmx")
	space = resolv.NewSpace(tileMap.TileWidth*tileMap.Width,
		tileMap.TileHeight*tileMap.Height, tileMap.TileWidth, tileMap.TileHeight)

	//fmt.Printf("%+v\n", tileMap.Tilesets[0])
	renderer, _ := render.NewRenderer(tileMap)
	renderer.RenderVisibleLayers()
	mapImage = ebi.NewImageFromImage(renderer.Result)
	//renderer.Clear()
	collisionObjects := tiles.CollisionObjectsOfTileLayer(tileMap.Layers[0])
	fmt.Printf("%+v\n", collisionObjects[0])
}

type Game struct {
	camera camera.FollowCam
	player player.PlayerImpl
	count  int
}

func (g *Game) Init() {
	g.camera = camera.NewFollowCam(screenWidth, screenHeight, 0, 0, 1, 0)
	g.camera.Followee = &g.player
	g.player.SceneEnter()
}

func (g *Game) Update() error {
	g.count++

	g.player.Update(g.count)

	g.camera.Update(g.count)

	return nil
}

func (g *Game) Draw(screen *ebi.Image) {
	g.camera.Surface.Clear()
	// Map
	g.camera.Surface.DrawImage(mapImage, g.camera.GetTranslation(0, 0))
	// Character
	g.player.Render(&g.camera)
	// Publish
	g.camera.Blit(screen)

	//ebiutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebi.CurrentTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
