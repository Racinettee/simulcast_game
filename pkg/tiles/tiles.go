package tiles

import (
	"github.com/lafriks/go-tiled"
	"github.com/solarlune/resolv"
)

func CollisionObjectsOfTileLayer(layer *tiled.Layer) []*resolv.Object {
	result := make([]*resolv.Object, 0)
	// Using the ID from each tile, and associating it back to its tileset, let us get the polygons
	for _, tile := range layer.Tiles {
		//fmt.Printf("%+v\n", tile)
		tileset := tile.Tileset

		// Find if this tile has any data in the tileset to speak of
		for _, tilesetTile := range tileset.Tiles {
			if tilesetTile.ID == tile.ID {
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
							newObject := resolv.NewObject(float64(tileRect.Dx()), float64(tileRect.Dy()), 1, 1)
							newObject.Shape = resolv.NewConvexPolygon(polyList...)
							result = append(result, newObject)
						}
					}
				}
			}
		}
		break
	}
	return result
}
