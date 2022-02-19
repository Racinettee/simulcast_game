package game

import (
	"image/color"
	"log"

	"github.com/BurntSushi/toml"
	ebi "github.com/hajimehoshi/ebiten/v2"

	//ebiaudio "github.com/hajimehoshi/ebiten/v2/audio"
	//"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	ebiutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	frameWidth   = 32
	frameHeight  = 32
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
		tileMap.TileHeight*tileMap.Height, tileMap.TileWidth/4, tileMap.TileHeight/4)

	renderer, _ := render.NewRenderer(tileMap)
	renderer.RenderVisibleLayers()
	mapImage = ebi.NewImageFromImage(renderer.Result)
	collisionObjects := tiles.CollisionObjectsOfTileLayer(tileMap.Layers[0])
	space.Add(collisionObjects...)
}

type Game struct {
	camera camera.FollowCam
	player player.PlayerImpl
	count  int
	Config Config
}

func (g *Game) Init() {
	_, err := toml.DecodeFile("simul_conf.toml", &g.Config)

	if err != nil {
		log.Printf("Failed to load configuration")
	}

	g.camera = camera.NewFollowCam(screenWidth, screenHeight, 0, 0, 1, 0)
	g.camera.Followee = &g.player
	g.player.SceneEnter()
	space.Add(g.player.Body)
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

	if g.Config.DrawCollisionShapes {
		g.DrawCollisions()
	}
	// Publish
	g.camera.Blit(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) DrawCollisions() {
	pX, pY := g.camera.GetScreenCoords(g.player.Body.X, g.player.Body.Y)
	pX -= frameWidth / 2
	pY -= frameHeight / 2
	// Draw collisions
	objects := space.Objects()
	for _, obj := range objects {
		switch shape := obj.Shape.(type) {
		case *resolv.ConvexPolygon:
			points := shape.Points
			sX, sY := g.camera.GetScreenCoords(obj.X, obj.Y)

			for i := 0; i < len(points)-1; i += 1 {
				pX, pY := sX+points[i].X(), sY+points[i].Y()
				pX2, pY2 := sX+points[i+1].X(), sY+points[i+1].Y()

				ebiutil.DrawLine(g.camera.Surface, pX, pY,
					pX2, pY2, color.RGBA{255, 125, 125, 100})
			}
		default:
			sX, sY := g.camera.GetScreenCoords(obj.X, obj.Y)
			ebiutil.DrawRect(g.camera.Surface, sX, sY, obj.W, obj.H, color.RGBA{255, 125, 125, 100})
		}

	}
}
