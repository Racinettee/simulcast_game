package tiles

import (
	"fmt"

	"github.com/lafriks/go-tiled"
	"github.com/solarlune/resolv"
)

func CollisionObjectsOfTileLayer(layer *tiled.Layer) []*resolv.Object {
	result := make([]*resolv.Object, 0)
	// Using the ID from each tile, and associating it back to its tileset, let us get the polygons
	for tileNo, tile := range layer.Tiles {
		//fmt.Printf("%+v\n", tile)
		tileset := tile.Tileset
		// Find if this tile has any data in the tileset to speak of
		for _, tilesetTile := range tileset.Tiles {
			if tilesetTile.ID == tile.ID {
				x, y := layer.GetTilePosition(tileNo)
				// The tile has something in the tileset
				// Examine the object groups for collision shapes
				objectGroups := tilesetTile.ObjectGroups
				if len(objectGroups) != 0 {
					for _, object := range objectGroups[0].Objects {
						// First we will check for polygons
						if len(object.Polygons) != 0 {
							tileRect := tile.GetTileRect()
							polygon := object.Polygons[0]
							//fmt.Printf("%+v\n", polygon.Points)
							points := *polygon.Points
							// Resolv requires a list of float64 to produce polygons from
							polyList := make([]float64, len(points)*2)
							for i, j := 0, 0; i < len(points); i, j = i+1, j+2 {
								polyList[j] = points[i].X
								polyList[j+1] = points[i].Y
							}
							//fmt.Printf("%v\n", polyList)
							newObject := resolv.NewObject(float64(x), float64(y), float64(tileRect.Dx()), float64(tileRect.Dy()))
							newObject.SetShape(resolv.NewConvexPolygon(polyList...))
							result = append(result, newObject)
						} else if len(object.Ellipses) != 0 {
							newObject := resolv.NewObject(float64(x), float64(y), 8, 8)
							newObject.SetShape(resolv.NewCircle(0, 0, object.Width))
							result = append(result, newObject)
						} else {
							tileRect := tile.GetTileRect()
							newObject := resolv.NewObject(float64(x), float64(y), float64(tileRect.Dx()), float64(tileRect.Dy()))
							newObject.SetShape(resolv.NewRectangle(0, 0, float64(tileRect.Dx()), float64(tileRect.Dy())))
							fmt.Printf("Rectangle...")
							result = append(result, newObject)
						}
					}
				}
			}
		}
	}
	return result
}
