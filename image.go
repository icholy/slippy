package slippy

import (
	"image"

	"github.com/faiface/pixel"

	"github.com/icholy/slippy/tiles"
)

var (
	Placeholder = pixel.PictureDataFromImage(
		image.NewRGBA(image.Rect(0, 0, tiles.TileSize, tiles.TileSize)),
	)
)

// RectTiles returns a slice of tiles requires to fully cover the rect
func RectTiles(bounds pixel.Rect, zoom int) []tiles.Tile {
	var (
		min = tiles.FromVec(bounds.Min, zoom)
		max = tiles.FromVec(bounds.Max, zoom)
		tt  []tiles.Tile
	)
	for x := min.X; x <= max.X; x++ {
		for y := max.Y; y <= min.Y; y++ {
			t := tiles.Tile{
				X: x,
				Y: y,
				Z: zoom,
			}
			tt = append(tt, t)
		}
	}
	return tt
}
